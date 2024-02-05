// Copyright 2020 Envoyproxy Authors
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

package cache

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	"github.com/envoyproxy/go-control-plane/pkg/log"
)

// cachedResource is used to track resources added by the user in the cache.
// It contains the resource itself and its associated version (currently in two different modes).
// TODO(valerian-roche): store serialized resource if already marshaled to compute the version.
type cachedResource struct {
	types.Resource

	// cacheVersion is the version of the cache at the time of last update, used in sotw.
	cacheVersion string
	// resourceVersion is the version of the resource itself (built through stable marshaling).
	// It is only set if computeResourceVersion is set to true on the cache.
	resourceVersion string
}

type watches struct {
	// sotw keeps track of current sotw watches, indexed per watch id.
	sotw map[uint64]ResponseWatch
	// delta keeps track of current delta watches, indexed per watch id.
	delta map[uint64]DeltaResponseWatch
}

func newWatches() watches {
	return watches{
		sotw:  make(map[uint64]ResponseWatch),
		delta: make(map[uint64]DeltaResponseWatch),
	}
}

func (w *watches) empty() bool {
	return len(w.sotw)+len(w.delta) == 0
}

// LinearCache supports collections of opaque resources. This cache has a
// single collection indexed by resource names and manages resource versions
// internally. It implements the cache interface for a single type URL and
// should be combined with other caches via type URL muxing. It can be used to
// supply EDS entries, for example, uniformly across a fleet of proxies.
type LinearCache struct {
	// typeURL provides the type of resources managed by the cache.
	// This information is used to reject requests watching another type, as well as to make
	// decisions based on resource type (e.g. whether sotw must return full-state).
	typeURL string

	// resources contains all resources currently set in the cache and associated versions.
	resources map[string]cachedResource

	// resourceWatches keeps track of watches currently opened specifically tracking a resource.
	// It does not contain wildcard watches.
	// It can contain resources not present in resources.
	resourceWatches map[string]watches
	// wildcardWatches keeps track of all wildcard watches currently opened.
	wildcardWatches watches
	// currentWatchID is used to index new watches.
	currentWatchID uint64

	// version is the current version of the cache. It is incremented each time resources are updated.
	version uint64
	// versionPrefix is used to modify the version returned to clients, and can be used to uniquely identify
	// cache instances and avoid issues of version reuse.
	versionPrefix string

	// areStableResourceVersionsComputed indicates whether the cache is currently computing and storing stable resource versions.
	areStableResourceVersionsComputed bool

	log log.Logger

	mu sync.RWMutex
}

var _ Cache = &LinearCache{}

// Options for modifying the behavior of the linear cache.
type LinearCacheOption func(*LinearCache)

// WithVersionPrefix sets a version prefix of the form "prefixN" in the version info.
// Version prefix can be used to distinguish replicated instances of the cache, in case
// a client re-connects to another instance.
func WithVersionPrefix(prefix string) LinearCacheOption {
	return func(cache *LinearCache) {
		cache.versionPrefix = prefix
	}
}

// WithInitialResources initializes the initial set of resources.
func WithInitialResources(resources map[string]types.Resource) LinearCacheOption {
	return func(cache *LinearCache) {
		for name, resource := range resources {
			cache.resources[name] = cachedResource{
				Resource: resource,
			}
		}
	}
}

func WithLogger(log log.Logger) LinearCacheOption {
	return func(cache *LinearCache) {
		cache.log = log
	}
}

// WithComputeStableVersions ensures the cache tracks the resources stable versions from the beginning,
// avoiding the first watch stalling until the computation has been done on all resources.
// If the cache is expected to handle delta watches, this option is fully beneficial.
// If the cache is expected to only handle sotw watches, it is currently not needed and will increase CPU usage.
func WithComputeStableVersions() LinearCacheOption {
	return func(cache *LinearCache) {
		cache.areStableResourceVersionsComputed = true
	}
}

// NewLinearCache creates a new cache. See the comments on the struct definition.
func NewLinearCache(typeURL string, opts ...LinearCacheOption) *LinearCache {
	out := &LinearCache{
		typeURL:         typeURL,
		resources:       make(map[string]cachedResource),
		resourceWatches: make(map[string]watches),
		wildcardWatches: newWatches(),
		version:         0,
		currentWatchID:  0,
		log:             log.NewDefaultLogger(),
	}
	for _, opt := range opts {
		opt(out)
	}
	for name, resource := range out.resources {
		resource.cacheVersion = out.getVersion()
		if out.areStableResourceVersionsComputed {
			version, err := computeResourceStableVersion(resource.Resource)
			if err != nil {
				out.log.Errorf("failed to build stable versions for resource %s: %s", name, err)
			} else {
				resource.resourceVersion = version
			}
		}
		out.resources[name] = resource
	}
	return out
}

func (cache *LinearCache) computeResourceChange(sub Subscription, ignoreReturnedResources, useStableVersion bool) (updated, removed []string) {
	var changedResources []string
	var removedResources []string

	knownVersions := sub.ReturnedResources()
	if ignoreReturnedResources {
		// The response will include all resources, with no regards of resources potentially already returned.
		knownVersions = make(map[string]string)
	}

	getVersion := func(c cachedResource) string { return c.cacheVersion }
	if useStableVersion {
		getVersion = func(c cachedResource) string { return c.resourceVersion }
	}

	if sub.IsWildcard() {
		for resourceName, resource := range cache.resources {
			knownVersion, ok := knownVersions[resourceName]
			if !ok {
				// This resource is not yet known by the client (new resource added in the cache or newly subscribed).
				changedResources = append(changedResources, resourceName)
			} else if knownVersion != getVersion(resource) {
				// The client knows an outdated version.
				changedResources = append(changedResources, resourceName)
			}
		}

		// Negative check to identify resources that have been removed in the cache.
		// Sotw does not support returning "deletions", but in the case of full state resources
		// a response must then be returned.
		for resourceName := range knownVersions {
			if _, ok := cache.resources[resourceName]; !ok {
				removedResources = append(removedResources, resourceName)
			}
		}
	} else {
		for resourceName := range sub.SubscribedResources() {
			res, exists := cache.resources[resourceName]
			knownVersion, known := knownVersions[resourceName]
			if !exists {
				if known {
					// This resource was removed from the cache. If the type requires full state
					// we need to return a response.
					removedResources = append(removedResources, resourceName)
				}
				continue
			}

			if !known {
				// This resource is not yet known by the client (new resource added in the cache or newly subscribed).
				changedResources = append(changedResources, resourceName)
			} else if knownVersion != getVersion(res) {
				// The client knows an outdated version.
				changedResources = append(changedResources, resourceName)
			}
		}

		for resourceName := range knownVersions {
			// If the subscription no longer watches a resource,
			// we mark it as unknown on the client side to ensure it will be resent to the client if subscribing again later on.
			if _, ok := sub.SubscribedResources()[resourceName]; !ok {
				removedResources = append(removedResources, resourceName)
			}
		}
	}

	return changedResources, removedResources
}

func (cache *LinearCache) computeSotwResponse(watch ResponseWatch, ignoreReturnedResources bool) *RawResponse {
	changedResources, removedResources := cache.computeResourceChange(watch.subscription, ignoreReturnedResources, false)
	if len(changedResources) == 0 && len(removedResources) == 0 && !ignoreReturnedResources {
		// Nothing changed.
		return nil
	}

	returnedVersions := make(map[string]string, len(watch.subscription.ReturnedResources()))
	// Clone the current returned versions. The cache should not alter the subscription
	for resourceName, version := range watch.subscription.ReturnedResources() {
		returnedVersions[resourceName] = version
	}

	cacheVersion := cache.getVersion()
	var resources []types.ResourceWithTTL

	switch {
	// Depending on the type, the response will only include changed resources or all of them
	case !ResourceRequiresFullStateInSotw(cache.typeURL):
		// changedResources is already filtered based on the subscription.
		// TODO(valerian-roche): if the only change is a removal in the subscription,
		// or a watched resource getting deleted, this might send an empty reply.
		// While this does not violate the protocol, we might want to avoid it.
		resources = make([]types.ResourceWithTTL, 0, len(changedResources))
		for _, resourceName := range changedResources {
			cachedResource := cache.resources[resourceName]
			resources = append(resources, types.ResourceWithTTL{Resource: cachedResource.Resource})
			returnedVersions[resourceName] = cachedResource.cacheVersion
		}
	case watch.subscription.IsWildcard():
		// Include all resources for the type.
		resources = make([]types.ResourceWithTTL, 0, len(cache.resources))
		for resourceName, cachedResource := range cache.resources {
			resources = append(resources, types.ResourceWithTTL{Resource: cachedResource.Resource})
			returnedVersions[resourceName] = cachedResource.cacheVersion
		}
	default:
		// Include all resources matching the subscription, with no concern on whether
		// it has been updated or not.
		requestedResources := watch.subscription.SubscribedResources()
		// The linear cache could be very large (e.g. containing all potential CLAs)
		// Therefore drives on the subscription requested resources.
		resources = make([]types.ResourceWithTTL, 0, len(requestedResources))
		for resourceName := range requestedResources {
			cachedResource, ok := cache.resources[resourceName]
			if !ok {
				continue
			}
			resources = append(resources, types.ResourceWithTTL{Resource: cachedResource.Resource})
			returnedVersions[resourceName] = cachedResource.cacheVersion
		}
	}

	// Cleanup resources no longer existing in the cache or no longer subscribed.
	// In sotw we cannot return those if not full state,
	// but this ensures we detect unsubscription then resubscription.
	for _, resourceName := range removedResources {
		delete(returnedVersions, resourceName)
	}

	if !ignoreReturnedResources && !ResourceRequiresFullStateInSotw(cache.typeURL) && len(resources) == 0 {
		// If the request is not the initial one, and the type does not require full updates,
		// do not return if noting is to be set.
		// For full-state resources an empty response does have a semantic meaning.
		return nil
	}

	return &RawResponse{
		Request:           watch.Request,
		Resources:         resources,
		ReturnedResources: returnedVersions,
		Version:           cacheVersion,
		Ctx:               context.Background(),
	}
}

func (cache *LinearCache) computeDeltaResponse(watch DeltaResponseWatch) *RawDeltaResponse {
	changedResources, removedResources := cache.computeResourceChange(watch.subscription, false, true)
	if len(changedResources) == 0 && len(removedResources) == 0 {
		// Nothing changed.
		return nil
	}

	returnedVersions := make(map[string]string, len(watch.subscription.ReturnedResources()))
	// Clone the current returned versions. The cache should not alter the subscription
	for resourceName, version := range watch.subscription.ReturnedResources() {
		returnedVersions[resourceName] = version
	}

	cacheVersion := cache.getVersion()
	resources := make([]types.Resource, 0, len(changedResources))
	for _, resourceName := range changedResources {
		resource := cache.resources[resourceName]
		resources = append(resources, resource.Resource)
		returnedVersions[resourceName] = resource.resourceVersion
	}
	// Cleanup resources no longer existing in the cache or no longer subscribed.
	for _, resourceName := range removedResources {
		delete(returnedVersions, resourceName)
	}

	return &RawDeltaResponse{
		DeltaRequest:      watch.Request,
		Resources:         resources,
		RemovedResources:  removedResources,
		NextVersionMap:    returnedVersions,
		SystemVersionInfo: cacheVersion,
		Ctx:               context.Background(),
	}
}

func (cache *LinearCache) notifyAll(modified []string) {
	// Gather the list of watches impacted by the modified resources.
	sotwWatches := make(map[uint64]ResponseWatch)
	deltaWatches := make(map[uint64]DeltaResponseWatch)
	for _, name := range modified {
		for watchID, watch := range cache.resourceWatches[name].sotw {
			sotwWatches[watchID] = watch
		}
		for watchID, watch := range cache.resourceWatches[name].delta {
			deltaWatches[watchID] = watch
		}
	}

	// sotw watches
	for watchID, watch := range sotwWatches {
		response := cache.computeSotwResponse(watch, false)
		if response != nil {
			watch.Response <- response
			cache.removeWatch(watchID, watch.subscription)
		} else {
			cache.log.Warnf("[Linear cache] Watch %d detected as triggered but no change was found", watchID)
		}
	}

	for watchID, watch := range cache.wildcardWatches.sotw {
		response := cache.computeSotwResponse(watch, false)
		if response != nil {
			watch.Response <- response
			delete(cache.wildcardWatches.sotw, watchID)
		} else {
			cache.log.Warnf("[Linear cache] Wildcard watch %d detected as triggered but no change was found", watchID)
		}
	}

	// delta watches
	for watchID, watch := range deltaWatches {
		response := cache.computeDeltaResponse(watch)
		if response != nil {
			watch.Response <- response
			cache.removeDeltaWatch(watchID, watch.subscription)
		} else {
			cache.log.Warnf("[Linear cache] Delta watch %d detected as triggered but no change was found", watchID)
		}
	}

	for watchID, watch := range cache.wildcardWatches.delta {
		response := cache.computeDeltaResponse(watch)
		if response != nil {
			watch.Response <- response
			delete(cache.wildcardWatches.delta, watchID)
		} else {
			cache.log.Warnf("[Linear cache] Wildcard delta watch %d detected as triggered but no change was found", watchID)
		}
	}
}

func computeResourceStableVersion(res types.Resource) (string, error) {
	// hash our version in here and build the version map
	marshaledResource, err := MarshalResource(res)
	if err != nil {
		return "", err
	}
	v := HashResource(marshaledResource)
	if v == "" {
		return "", errors.New("failed to build resource version")
	}
	return v, nil
}

func (cache *LinearCache) addResourceToCache(name string, res types.Resource) error {
	update := cachedResource{
		Resource:     res,
		cacheVersion: cache.getVersion(),
	}
	if cache.areStableResourceVersionsComputed {
		version, err := computeResourceStableVersion(res)
		if err != nil {
			return err
		}
		update.resourceVersion = version
	}
	cache.resources[name] = update
	return nil
}

// UpdateResource updates a resource in the collection.
func (cache *LinearCache) UpdateResource(name string, res types.Resource) error {
	if res == nil {
		return errors.New("nil resource")
	}
	cache.mu.Lock()
	defer cache.mu.Unlock()

	cache.version++
	if err := cache.addResourceToCache(name, res); err != nil {
		return err
	}

	cache.notifyAll([]string{name})

	return nil
}

// DeleteResource removes a resource in the collection.
func (cache *LinearCache) DeleteResource(name string) error {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	cache.version++
	delete(cache.resources, name)

	cache.notifyAll([]string{name})
	return nil
}

// UpdateResources updates/deletes a list of resources in the cache.
// Calling UpdateResources instead of iterating on UpdateResource and DeleteResource
// is significantly more efficient when using delta or wildcard watches.
func (cache *LinearCache) UpdateResources(toUpdate map[string]types.Resource, toDelete []string) error {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	cache.version++
	modified := make([]string, 0, len(toUpdate)+len(toDelete))
	for name, resource := range toUpdate {
		if err := cache.addResourceToCache(name, resource); err != nil {
			return err
		}
		modified = append(modified, name)
	}
	for _, name := range toDelete {
		delete(cache.resources, name)
		modified = append(modified, name)
	}

	cache.notifyAll(modified)

	return nil
}

// SetResources replaces current resources with a new set of resources.
// This function is useful for wildcard xDS subscriptions.
// This way watches that are subscribed to all resources are triggered only once regardless of how many resources are changed.
func (cache *LinearCache) SetResources(resources map[string]types.Resource) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	cache.version++

	modified := make([]string, 0, len(resources))
	// Collect deleted resource names.
	for name := range cache.resources {
		if _, found := resources[name]; !found {
			delete(cache.resources, name)
			modified = append(modified, name)
		}
	}

	// Collect changed resource names.
	// We assume all resources passed to SetResources are changed.
	// Otherwise we would have to do proto.Equal on resources which is pretty expensive operation
	for name, resource := range resources {
		if err := cache.addResourceToCache(name, resource); err != nil {
			cache.log.Errorf("Failed to add resources to the cache: %s", err)
		}
		modified = append(modified, name)
	}

	cache.notifyAll(modified)
}

// GetResources returns current resources stored in the cache
func (cache *LinearCache) GetResources() map[string]types.Resource {
	cache.mu.RLock()
	defer cache.mu.RUnlock()

	// create a copy of our internal storage to avoid data races
	// involving mutations of our backing map
	resources := make(map[string]types.Resource, len(cache.resources))
	for k, v := range cache.resources {
		resources[k] = v.Resource
	}
	return resources
}

// The implementations of sotw and delta watches handling is nearly identical. The main distinctions are:
//   - handling of version in sotw when the request is the first of a subscription. Delta has a proper handling based on the request providing known versions.
//   - building the initial resource versions in delta if they've not been computed yet.
//   - computeSotwResponse and computeDeltaResponse has slightly different implementations due to sotw requirements to return full state for certain resources only.
func (cache *LinearCache) CreateWatch(request *Request, sub Subscription, value chan Response) (func(), error) {
	if request.GetTypeUrl() != cache.typeURL {
		return nil, fmt.Errorf("request type %s does not match cache type %s", request.GetTypeUrl(), cache.typeURL)
	}

	// If the request does not include a version the client considers it has no current state.
	// In this case we will always reply to allow proper initialization of dependencies in the client.
	ignoreCurrentSubscriptionResources := request.GetVersionInfo() == ""
	if !strings.HasPrefix(request.GetVersionInfo(), cache.versionPrefix) {
		// If the version of the request does not match the cache prefix, we will send a response in all cases to match the legacy behavior.
		ignoreCurrentSubscriptionResources = true
		cache.log.Debugf("[linear cache] received watch with version %s not matching the cache prefix %s. Will return all known resources", request.GetVersionInfo(), cache.versionPrefix)
	}

	// A major difference between delta and sotw is the ability to not resend everything when connecting to a new control-plane
	// In delta the request provides the version of the resources it does know, even if the request is wildcard or does request more resources
	// In sotw the request only provides the global version of the control-plane, and there is no way for the control-plane to know if resources have
	// been added since in the requested resources. In the context of generalized wildcard, even wildcard could be new, and taking the assumption
	// that wildcard implies that the client already knows all resources at the given version is no longer true.
	// We could optimize the reconnection case here if:
	//  - we take the assumption that clients will not start requesting wildcard while providing a version. We could then ignore requests providing the resources.
	//  - we use the version as some form of hash of resources known, and we can then consider it as a way to correctly verify whether all resources are unchanged.
	// For now it is not done as:
	//  - for the first case, while the protocol documentation does not explicitly mention the case, it does not mark it impossible and explicitly references unsubscribing from wildcard.
	//  - for the second one we could likely do it with little difficulty if need be, but if users rely on the current monotonic version it could impact their callbacks implementations.
	watch := ResponseWatch{Request: request, Response: value, subscription: sub}

	cache.mu.Lock()
	defer cache.mu.Unlock()

	response := cache.computeSotwResponse(watch, ignoreCurrentSubscriptionResources)
	if response != nil {
		cache.log.Debugf("[linear cache] replying to the watch with resources %v (subscription values %v, known %v)", response.GetReturnedResources(), sub.SubscribedResources(), sub.ReturnedResources())
		watch.Response <- response
		return func() {}, nil
	}

	watchID := cache.nextWatchID()
	// Create open watches since versions are up to date.
	if sub.IsWildcard() {
		cache.log.Infof("[linear cache] open watch %d for %s all resources, system version %q", watchID, cache.typeURL, cache.getVersion())
		cache.wildcardWatches.sotw[watchID] = watch
		return func() {
			cache.mu.Lock()
			defer cache.mu.Unlock()
			delete(cache.wildcardWatches.sotw, watchID)
		}, nil
	}

	cache.log.Infof("[linear cache] open watch %d for %s resources %v, system version %q", watchID, cache.typeURL, sub.SubscribedResources(), cache.getVersion())
	for name := range sub.SubscribedResources() {
		watches, exists := cache.resourceWatches[name]
		if !exists {
			watches = newWatches()
			cache.resourceWatches[name] = watches
		}
		watches.sotw[watchID] = watch
	}
	return func() {
		cache.mu.Lock()
		defer cache.mu.Unlock()
		cache.removeWatch(watchID, watch.subscription)
	}, nil
}

// Must be called under lock
func (cache *LinearCache) removeWatch(watchID uint64, sub Subscription) {
	// Make sure we clean the watch for ALL resources it might be associated with,
	// as the channel will no longer be listened to
	for resource := range sub.SubscribedResources() {
		resourceWatches := cache.resourceWatches[resource]
		delete(resourceWatches.sotw, watchID)
		if resourceWatches.empty() {
			delete(cache.resourceWatches, resource)
		}
	}
}

func (cache *LinearCache) CreateDeltaWatch(request *DeltaRequest, sub Subscription, value chan DeltaResponse) (func(), error) {
	if request.GetTypeUrl() != cache.typeURL {
		return nil, fmt.Errorf("request type %s does not match cache type %s", request.GetTypeUrl(), cache.typeURL)
	}

	watch := DeltaResponseWatch{Request: request, Response: value, subscription: sub}

	cache.mu.Lock()
	defer cache.mu.Unlock()

	if !cache.areStableResourceVersionsComputed {
		// If we had no previously open delta watches, we need to build the version map for the first time.
		// The version map will not be destroyed when the last delta watch is removed.
		// This avoids constantly rebuilding when only a few delta watches are open.
		cache.areStableResourceVersionsComputed = true
		cache.log.Infof("[linear cache] activating stable resource version computation for %s", cache.typeURL)

		for name, cachedResource := range cache.resources {
			version, err := computeResourceStableVersion(cachedResource.Resource)
			if err != nil {
				return func() {}, fmt.Errorf("failed to build stable versions for resource %s: %w", name, err)
			}
			cachedResource.resourceVersion = version
			cache.resources[name] = cachedResource
		}
	}

	response := cache.computeDeltaResponse(watch)
	if response != nil {
		cache.log.Debugf("[linear cache] replying to the delta watch (subscription values %v, known %v)", sub.SubscribedResources(), sub.ReturnedResources())
		watch.Response <- response
		return nil, nil
	}

	watchID := cache.nextWatchID()
	// Create open watches since versions are up to date.
	if sub.IsWildcard() {
		cache.log.Infof("[linear cache] open delta watch %d for all %s resources, system version %q", watchID, cache.typeURL, cache.getVersion())
		cache.wildcardWatches.delta[watchID] = watch
		return func() {
			cache.mu.Lock()
			defer cache.mu.Unlock()
			delete(cache.wildcardWatches.delta, watchID)
		}, nil
	}

	cache.log.Infof("[linear cache] open delta watch %d for %s resources %v, system version %q", watchID, cache.typeURL, sub.SubscribedResources(), cache.getVersion())
	for name := range sub.SubscribedResources() {
		watches, exists := cache.resourceWatches[name]
		if !exists {
			watches = newWatches()
			cache.resourceWatches[name] = watches
		}
		watches.delta[watchID] = watch
	}
	return func() {
		cache.mu.Lock()
		defer cache.mu.Unlock()
		cache.removeDeltaWatch(watchID, watch.subscription)
	}, nil
}

func (cache *LinearCache) getVersion() string {
	return cache.versionPrefix + strconv.FormatUint(cache.version, 10)
}

// cancellation function for cleaning stale watches
func (cache *LinearCache) removeDeltaWatch(watchID uint64, sub Subscription) {
	// Make sure we clean the watch for ALL resources it might be associated with,
	// as the channel will no longer be listened to
	for resource := range sub.SubscribedResources() {
		resourceWatches := cache.resourceWatches[resource]
		delete(resourceWatches.delta, watchID)
		if resourceWatches.empty() {
			delete(cache.resourceWatches, resource)
		}
	}
}

func (cache *LinearCache) nextWatchID() uint64 {
	cache.currentWatchID++
	if cache.currentWatchID == 0 {
		panic("watch id count overflow")
	}
	return cache.currentWatchID
}

func (cache *LinearCache) Fetch(context.Context, *Request) (Response, error) {
	return nil, errors.New("not implemented")
}

// Number of resources currently on the cache.
// As GetResources is building a clone it is expensive to get metrics otherwise.
func (cache *LinearCache) NumResources() int {
	cache.mu.RLock()
	defer cache.mu.RUnlock()
	return len(cache.resources)
}

// NumDeltaWatches returns the number of active sotw watches for a resource name.
func (cache *LinearCache) NumWatches(name string) int {
	cache.mu.RLock()
	defer cache.mu.RUnlock()
	return len(cache.resourceWatches[name].sotw) + len(cache.wildcardWatches.sotw)
}

// NumDeltaWatches returns the number of active delta watches for a resource name.
func (cache *LinearCache) NumDeltaWatchesForResource(name string) int {
	cache.mu.RLock()
	defer cache.mu.RUnlock()
	return len(cache.resourceWatches[name].delta) + len(cache.wildcardWatches.delta)
}

// NumDeltaWatches returns the total number of active delta watches.
// Warning: it is quite inefficient, and NumDeltaWatchesForResource should be preferred.
func (cache *LinearCache) NumDeltaWatches() int {
	cache.mu.RLock()
	defer cache.mu.RUnlock()
	uniqueWatches := map[uint64]struct{}{}
	for _, watches := range cache.resourceWatches {
		for id := range watches.delta {
			uniqueWatches[id] = struct{}{}
		}
	}
	return len(uniqueWatches) + len(cache.wildcardWatches.delta)
}
