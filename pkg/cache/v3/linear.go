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
	"sync/atomic"

	"github.com/envoyproxy/go-control-plane/internal/snapshot"
	"github.com/envoyproxy/go-control-plane/internal/watches"
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	"github.com/envoyproxy/go-control-plane/pkg/log"
)

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

	snapshot *snapshot.TypedSnapshot

	// resourceWatches keeps track of watches currently opened specifically tracking a resource.
	// It does not contain wildcard watches.
	// It can contain resources not present in resources.
	resourceWatches map[string]cacheWatches
	// wildcardWatches keeps track of all wildcard watches currently opened.
	wildcardWatches cacheWatches
	// currentWatchID is used to index new watches.
	currentWatchID uint64

	// version is the current version of the cache. It is incremented each time resources are updated.
	version atomic.Uint64
	// versionPrefix is used to modify the version returned to clients, and can be used to uniquely identify
	// cache instances and avoid issues of version reuse.
	versionPrefix string

	// useResourceVersionsInSotw switches to a new version model for sotw watches.
	// When activated, versions are stored in subscriptions using resource versions, and the response version
	// is a hash of the returned versions to allow watch resumptions when reconnecting to the cache with a
	// new subscription.
	useResourceVersionsInSotw bool

	watchCount int

	log log.Logger

	mu sync.RWMutex
}

type cacheWatches map[uint64]snapshot.Watch

func newCacheWatches() cacheWatches {
	return make(cacheWatches)
}

var _ Cache = &LinearCache{}

// Options for modifying the behavior of the linear cache.
type LinearCacheOption func(*LinearCache)

// WithVersionPrefix sets a version prefix of the form "prefixN" in the version info.
// Version prefix can be used to distinguish replicated instances of the cache, in case
// a client re-connects to another instance.
// Deprecated: use WithSotwResourceVersions instead to avoid issues when reconnecting to other instances
// while avoiding resending resources if unchanged.
func WithVersionPrefix(prefix string) LinearCacheOption {
	return func(cache *LinearCache) {
		cache.versionPrefix = prefix
	}
}

// WithInitialResources initializes the initial set of resources.
func WithInitialResources(resources map[string]types.Resource) LinearCacheOption {
	return func(c *LinearCache) {
		c.snapshot.SetInitialResources(resources)
	}
}

func WithLogger(log log.Logger) LinearCacheOption {
	return func(cache *LinearCache) {
		cache.log = log
	}
}

// WithSotwResourceVersions changes the versions returned in sotw to encode the list of resources known
// in the subscription.
// The use of resource versions for sotw also deduplicates updates to clients if the cache updates are
// not changing the content of the resource.
// When used, the use of WithVersionPrefix is no longer needed to manage reconnection to other instances
// and should not be used.
func WithSotwResourceVersions() LinearCacheOption {
	return func(cache *LinearCache) {
		cache.useResourceVersionsInSotw = true
	}
}

// NewLinearCache creates a new cache. See the comments on the struct definition.
func NewLinearCache(typeURL string, opts ...LinearCacheOption) *LinearCache {
	out := &LinearCache{
		typeURL:         typeURL,
		resourceWatches: make(map[string]cacheWatches),
		wildcardWatches: newCacheWatches(),
		currentWatchID:  0,
		log:             log.NewDefaultLogger(),
	}

	out.snapshot = snapshot.NewTypedSnapshot(snapshot.WithUpdateCallback(out.notifyAll))
	for _, opt := range opts {
		opt(out)
	}

	return out
}

func (cache *LinearCache) notifyAll(modified []string) error {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	// Gather the list of watches impacted by the modified resources.
	resourceWatches := newCacheWatches()
	for _, name := range modified {
		for watchID, watch := range cache.resourceWatches[name] {
			resourceWatches[watchID] = watch
		}
	}

	for watchID, watch := range resourceWatches {
		response, err := cache.snapshot.ComputeResponse(watch, false, cache.getVersion())
		if err != nil {
			return err
		}

		if response != nil {
			watch.SendResponse(response)
			cache.removeWatch(watchID, watch.GetSubscription())
		} else {
			cache.log.Infof("[Linear cache] Watch %d detected as triggered but no change was found", watchID)
		}
	}

	for watchID, watch := range cache.wildcardWatches {
		response, err := cache.snapshot.ComputeResponse(watch, false, cache.getVersion())
		if err != nil {
			return err
		}

		if response != nil {
			watch.SendResponse(response)
			cache.removeWildcardWatch(watchID)
		} else {
			cache.log.Infof("[Linear cache] Wildcard watch %d detected as triggered but no change was found", watchID)
		}
	}

	return nil
}

// UpdateResource updates a resource in the collection.
func (cache *LinearCache) UpdateResource(name string, res types.Resource) error {
	if res == nil {
		return errors.New("nil resource")
	}

	cache.incrementVersion()
	return cache.snapshot.UpdateResources(map[string]types.Resource{name: res}, nil)
}

// DeleteResource removes a resource in the collection.
func (cache *LinearCache) DeleteResource(name string) error {
	cache.incrementVersion()
	return cache.snapshot.UpdateResources(nil, []string{name})
}

// UpdateResources updates/deletes a list of resources in the cache.
// Calling UpdateResources instead of iterating on UpdateResource and DeleteResource
// is significantly more efficient when using delta or wildcard watches.
func (cache *LinearCache) UpdateResources(toUpdate map[string]types.Resource, toDelete []string) error {
	cache.incrementVersion()
	return cache.snapshot.UpdateResources(toUpdate, toDelete)
}

// SetResources replaces current resources with a new set of resources.
// If only some resources are to be updated, UpdateResources is more efficient.
func (cache *LinearCache) SetResources(resources map[string]types.Resource) {
	cache.incrementVersion()
	if err := cache.snapshot.SetResources(resources); err != nil {
		cache.log.Errorf("Failed to notify watches: %s", err.Error())
	}
}

// GetResources returns current resources stored in the cache
func (cache *LinearCache) GetResources() map[string]types.Resource {
	return cache.snapshot.GetResources()
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
	replyEvenIfEmpty := request.GetVersionInfo() == ""
	if !strings.HasPrefix(request.GetVersionInfo(), cache.versionPrefix) {
		// If the version of the request does not match the cache prefix, we will send a response in all cases to match the legacy behavior.
		replyEvenIfEmpty = true
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
	// When using the `WithSotwResourceVersions` option, this optimization is activated and avoids resending all the dataset on wildcard watch resumption if no change has occurred.
	watch := watches.NewResponseWatch(request, value, sub, cache.useResourceVersionsInSotw, ResourceRequiresFullStateInSotw(cache.typeURL), cache.versionPrefix)

	// Lock to ensure a change in the cache will not be missed on the watch.
	// The call to notifyAll must considered watches from the version of the cache when it is called.
	// This lock could be removed by keeping a record of impacted objects per cache revision.
	cache.mu.Lock()
	defer cache.mu.Unlock()

	response, err := cache.snapshot.ComputeResponse(watch, replyEvenIfEmpty, cache.getVersion())
	if err != nil {
		return nil, fmt.Errorf("failed to compute the watch respnse: %w", err)
	}
	shouldReply := false
	if response != nil {
		// If the request
		//  - is the first
		//  - is wildcard
		//  - provides a non-empty version, matching the version prefix
		// and the cache uses resource versions, if the generated versions are the same as the previous one, we do not return the response.
		// This avoids resending all data if the new subscription is just a resumption of the previous one.
		// This optimization is only done on wildcard as we cannot track if a subscription is "new" at this stage and needs to be returned.
		// In the context of wildcard it could be incorrect if the subscription is newly wildcard, and we already returned all objects,
		// but as of Q1-2024 there are no known usecases of a subscription becoming wildcard (in envoy of xds-grpc).
		if cache.useResourceVersionsInSotw && sub.IsWildcard() && request.GetResponseNonce() == "" && !replyEvenIfEmpty {
			if request.GetVersionInfo() != response.GetResponseVersion() {
				// The response has a different returned version map as the request
				shouldReply = true
			} else {
				// We confirmed the content of the known resources, store them in the watch we create.
				subscription := newWatchSubscription(sub)
				subscription.returnedResources = response.GetReturnedResources()
				watch.Subscription = subscription
				sub = subscription
			}
		} else {
			shouldReply = true
		}
	}

	if shouldReply {
		cache.log.Debugf("[linear cache] replying to the watch with resources %v (subscription values %v, known %v)", response.GetReturnedResources(), sub.SubscribedResources(), sub.ReturnedResources())
		watch.SendResponse(response)
		return func() {}, nil
	}

	return cache.trackWatch(watch), nil
}

func (cache *LinearCache) CreateDeltaWatch(request *DeltaRequest, sub Subscription, value chan DeltaResponse) (func(), error) {
	if request.GetTypeUrl() != cache.typeURL {
		return nil, fmt.Errorf("request type %s does not match cache type %s", request.GetTypeUrl(), cache.typeURL)
	}

	watch := watches.NewDeltaResponseWatch(request, value, sub)

	// On first request on a wildcard subscription, envoy does expect a response to come in to
	// conclude initialization.
	replyEvenIfEmpty := false
	if sub.IsWildcard() && request.GetResponseNonce() == "" {
		replyEvenIfEmpty = true
	}

	cache.mu.Lock()
	defer cache.mu.Unlock()

	response, err := cache.snapshot.ComputeResponse(watch, replyEvenIfEmpty, cache.getVersion())
	if err != nil {
		return nil, fmt.Errorf("failed to compute the watch respnse: %w", err)
	}
	if response != nil {
		cache.log.Debugf("[linear cache] replying to the delta watch (subscription values %v, known %v)", sub.SubscribedResources(), sub.ReturnedResources())
		watch.SendResponse(response)
		return nil, nil
	}

	return cache.trackWatch(watch), nil
}

func (cache *LinearCache) nextWatchID() uint64 {
	cache.currentWatchID++
	if cache.currentWatchID == 0 {
		panic("watch id count overflow")
	}
	return cache.currentWatchID
}

// Must be called under lock
func (cache *LinearCache) trackWatch(watch snapshot.Watch) func() {
	cache.watchCount++

	watchID := cache.nextWatchID()
	sub := watch.GetSubscription()
	// Create open watches since versions are up to date.
	if sub.IsWildcard() {
		cache.log.Infof("[linear cache] open watch %d for %s all resources", watchID, cache.typeURL)
		cache.log.Debugf("[linear cache] subscription details for watch %d: known versions %v, system version %q", watchID, sub.ReturnedResources(), cache.getVersion())
		cache.wildcardWatches[watchID] = watch
		return func() {
			cache.mu.Lock()
			defer cache.mu.Unlock()
			cache.removeWildcardWatch(watchID)
		}
	}

	cache.log.Infof("[linear cache] open watch %d for %s resources %v", watchID, cache.typeURL, sub.SubscribedResources())
	cache.log.Debugf("[linear cache] subscription details for watch %d: known versions %v, system version %q", watchID, sub.ReturnedResources(), cache.getVersion())
	for name := range sub.SubscribedResources() {
		watches, exists := cache.resourceWatches[name]
		if !exists {
			watches = newCacheWatches()
			cache.resourceWatches[name] = watches
		}
		watches[watchID] = watch
	}
	return func() {
		cache.mu.Lock()
		defer cache.mu.Unlock()
		cache.removeWatch(watchID, sub)
	}
}

// Must be called under lock
func (cache *LinearCache) removeWatch(watchID uint64, sub Subscription) {
	// Make sure we clean the watch for ALL resources it might be associated with,
	// as the channel will no longer be listened to
	for resource := range sub.SubscribedResources() {
		resourceWatches := cache.resourceWatches[resource]
		delete(resourceWatches, watchID)
		if len(resourceWatches) == 0 {
			delete(cache.resourceWatches, resource)
		}
	}
	cache.watchCount--
}

// Must be called under lock
func (cache *LinearCache) removeWildcardWatch(watchID uint64) {
	cache.watchCount--
	delete(cache.wildcardWatches, watchID)
}

func (cache *LinearCache) Fetch(context.Context, *Request) (Response, error) {
	return nil, errors.New("not implemented")
}

func (cache *LinearCache) getVersion() string {
	version := cache.version.Load()
	return cache.versionPrefix + strconv.FormatUint(version, 10)
}

func (cache *LinearCache) incrementVersion() {
	cache.version.Add(1)
}

// NumResources returns the number of resources currently in the cache.
// As GetResources is building a clone it is expensive to get metrics otherwise.
func (cache *LinearCache) NumResources() int {
	return cache.snapshot.NumResources()
}

// NumWatches returns the number of active watches for a resource name.
func (cache *LinearCache) NumWatches(name string) int {
	cache.mu.RLock()
	defer cache.mu.RUnlock()
	return len(cache.resourceWatches[name]) + len(cache.wildcardWatches)
}

// TotalWatches returns the number of active watches on the cache in general.
func (cache *LinearCache) NumWildcardWatches() int {
	cache.mu.RLock()
	defer cache.mu.RUnlock()
	return len(cache.wildcardWatches)
}

// NumCacheWatches returns the number of active watches on the cache in general.
func (cache *LinearCache) NumCacheWatches() int {
	cache.mu.RLock()
	defer cache.mu.RUnlock()
	return cache.watchCount
}
