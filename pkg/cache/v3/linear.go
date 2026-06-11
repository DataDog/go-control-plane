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
	"maps"
	"strconv"
	"strings"
	"sync"

	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	"github.com/envoyproxy/go-control-plane/pkg/cache/internal"
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	"github.com/envoyproxy/go-control-plane/pkg/log"
)

// ResourceTransform produces the per-watch view of a stored resource. It is set
// per-watch via CreateWatchWithTransform / CreateDeltaWatchWithTransform and
// applied when the cache builds a response for that watch.
//
//   - innerVersion is the cache's own (inner) version of the stored resource,
//     handed in so the caller can embed it in the wire version if it wants.
//   - keep=false omits the resource from this watch's response.
//   - wireVersion == "" passes the stored resource through unchanged (the cache
//     uses its normal versioning on both the wire and its internal tracking).
//   - wireVersion != "" ships out's bytes and stamps wireVersion on the WIRE
//     (delta Resource.version), while the cache keeps innerVersion in the
//     internal returnedResources used for change detection. This divergence is
//     deliberate: the diff stays in inner-version space, so a wire version that
//     varies per client (e.g. a subset version) never triggers spurious re-emits.
type ResourceTransform func(name string, res types.Resource, innerVersion string) (out types.Resource, wireVersion string, keep bool)

type watch interface {
	// isDelta indicates whether the watch is a delta one.
	// It should not be used to take functional decisions, but is still currently used pending final changes.
	// It can be used to generate statistics.
	isDelta() bool
	// useResourceVersion indicates whether versions returned in the response are built using resource versions instead of cache update versions.
	useResourceVersion() bool
	// sendFullStateResponses requires that all resources matching the request, with no regards to which ones actually updated, must be provided in the response.
	// As a consequence, sending a response with no resources has a functional meaning of no matching resources available.
	sendFullStateResponses() bool

	getSubscription() Subscription
	// buildResponse computes the actual WatchResponse object to be sent on the watch.
	buildResponse(updatedResources []*internal.CachedResource, removedResources []string, returnedVersions map[string]string, version string) WatchResponse
	// sendResponse sends the response for the watch.
	// It must be called at most once.
	sendResponse(resp WatchResponse)
	// getTransform returns the per-watch resource transform, or nil.
	getTransform() ResourceTransform
	// getNode returns the node the watch was opened for, used to index the watch
	// for RefreshWatches via the cache's NodeHash.
	getNode() *core.Node
}

type watches map[uint64]watch

func newWatches() watches {
	return make(watches)
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
	resources map[string]*internal.CachedResource

	// resourceWatches keeps track of watches currently opened specifically tracking a resource.
	// It does not contain wildcard watches.
	// It can contain resources not present in resources.
	resourceWatches map[string]watches
	// prefixWatches keeps track of watches for glob collection subscriptions (e.g. "collection/*").
	// Keyed by the prefix with the trailing glob stripped but the separator kept (e.g. "collection/").
	prefixWatches map[string]watches
	// wildcardWatches keeps track of all wildcard watches currently opened.
	wildcardWatches watches
	// currentWatchID is used to index new watches.
	currentWatchID uint64

	// version is the current version of the cache. It is incremented each time resources are updated.
	version uint64
	// versionPrefix is used to modify the version returned to clients, and can be used to uniquely identify
	// cache instances and avoid issues of version reuse.
	versionPrefix string

	watchCount int

	log log.Logger

	// nodeHash, when set (via WithNodeHash), derives a per-client key from a
	// watch's node. Transform-watches are then indexed by that key in nodeIndex
	// so RefreshWatches can target a single client without scanning all watches.
	nodeHash NodeHash
	// nodeIndex maps a node key to its tracked transform-watches. Maintained in
	// trackWatch and cleaned in removeWatch/removeWildcardWatch.
	nodeIndex map[string]watches
	// nodeKeyByWatch maps a tracked transform-watch id back to its node key, so
	// any removal path can clean nodeIndex without recomputing the hash.
	nodeKeyByWatch map[uint64]string

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
			cache.resources[name] = internal.NewCachedResource(name, resource)
		}
	}
}

func WithLogger(log log.Logger) LinearCacheOption {
	return func(cache *LinearCache) {
		cache.log = log
	}
}

// WithNodeHash enables per-node indexing of transform-watches so RefreshWatches
// can target a single client. The hash must produce, for a watch's node, the
// same key the caller passes to RefreshWatches. Only watches created via
// CreateWatchWithTransform / CreateDeltaWatchWithTransform are indexed.
func WithNodeHash(hash NodeHash) LinearCacheOption {
	return func(cache *LinearCache) {
		cache.nodeHash = hash
	}
}

// NewLinearCache creates a new cache. See the comments on the struct definition.
func NewLinearCache(typeURL string, opts ...LinearCacheOption) *LinearCache {
	out := &LinearCache{
		typeURL:         typeURL,
		resources:       make(map[string]*internal.CachedResource),
		resourceWatches: make(map[string]watches),
		prefixWatches:   make(map[string]watches),
		wildcardWatches: newWatches(),
		version:         0,
		currentWatchID:  0,
		log:             log.NewDefaultLogger(),
		nodeIndex:       make(map[string]watches),
		nodeKeyByWatch:  make(map[uint64]string),
	}
	for _, opt := range opts {
		opt(out)
	}
	version := out.getVersion()
	for _, resource := range out.resources {
		resource.SetCacheVersion(version)
	}
	return out
}

// isSubscribedTo checks whether resourceName is covered by the subscription,
// either as an exact match in SubscribedResources or via a SubscribedPrefixes entry.
// When there are no prefix subscriptions this is a single O(1) map lookup.
func isSubscribedTo(sub Subscription, resourceName string) bool {
	if _, ok := sub.SubscribedResources()[resourceName]; ok {
		return true
	}
	for prefix := range sub.SubscribedPrefixes() {
		if strings.HasPrefix(resourceName, prefix) {
			return true
		}
	}
	return false
}

// computeResourceChange compares the subscription known resources and the cache current state to compute the list of resources
// which have changed and should be notified to the user.
//
// The useResourceVersion argument defines what version type to use for resources:
//   - if set to false versions are based on when resources were updated in the cache.
//   - if set to true versions are a stable property of the resource, with no regard to when it was added to the cache.
func (cache *LinearCache) computeResourceChange(sub Subscription, useResourceVersion bool) (updated, removed []string, err error) {
	var changedResources []string
	var removedResources []string

	knownVersions := sub.ReturnedResources()

	if sub.IsWildcard() {
		for resourceName, resource := range cache.resources {
			knownVersion, ok := knownVersions[resourceName]
			if !ok {
				// This resource is not yet known by the client (new resource added in the cache or newly subscribed).
				changedResources = append(changedResources, resourceName)
			} else {
				resourceVersion, err := resource.GetVersion(useResourceVersion)
				if err != nil {
					return nil, nil, fmt.Errorf("failed to compute version of %s: %w", resourceName, err)
				}
				if knownVersion != resourceVersion {
					// The client knows an outdated version.
					changedResources = append(changedResources, resourceName)
				}
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
		// Exact subscriptions.
		for resourceName := range sub.SubscribedResources() {
			res, exists := cache.resources[resourceName]
			if !exists {
				continue
			}
			knownVersion, known := knownVersions[resourceName]
			if !known {
				changedResources = append(changedResources, resourceName)
			} else {
				resourceVersion, err := res.GetVersion(useResourceVersion)
				if err != nil {
					return nil, nil, fmt.Errorf("failed to compute version of %s: %w", resourceName, err)
				}
				if knownVersion != resourceVersion {
					changedResources = append(changedResources, resourceName)
				}
			}
		}

		// Prefix subscriptions: single pass over the cache, checking all prefixes per resource.
		subscribedPrefixes := sub.SubscribedPrefixes()
		if len(subscribedPrefixes) > 0 {
			exactSubs := sub.SubscribedResources()
			for cachedName, res := range cache.resources {
				if _, exactSub := exactSubs[cachedName]; exactSub {
					continue
				}
				for prefix := range subscribedPrefixes {
					if !strings.HasPrefix(cachedName, prefix) {
						continue
					}
					knownVersion, known := knownVersions[cachedName]
					if !known {
						changedResources = append(changedResources, cachedName)
					} else {
						resourceVersion, err := res.GetVersion(useResourceVersion)
						if err != nil {
							return nil, nil, fmt.Errorf("failed to compute version of %s: %w", cachedName, err)
						}
						if knownVersion != resourceVersion {
							changedResources = append(changedResources, cachedName)
						}
					}
					break
				}
			}
		}

		// Consolidated removal: detect both cache deletion and unsubscription.
		// Cache deletion is checked first as it is O(1) and more common than
		// subscription changes; isSubscribedTo may iterate prefixes.
		for resourceName := range knownVersions {
			if _, inCache := cache.resources[resourceName]; !inCache {
				removedResources = append(removedResources, resourceName)
			} else if !isSubscribedTo(sub, resourceName) {
				removedResources = append(removedResources, resourceName)
			}
		}
	}

	return changedResources, removedResources, nil
}

func (cache *LinearCache) computeResponse(watch watch, replyEvenIfEmpty bool) (WatchResponse, error) {
	sub := watch.getSubscription()
	changedResources, removedResources, err := cache.computeResourceChange(sub, watch.useResourceVersion())
	if err != nil {
		return nil, err
	}

	if len(changedResources) == 0 && len(removedResources) == 0 && !replyEvenIfEmpty {
		// Nothing changed.
		return nil, nil
	}

	// In sotw the list of resources to actually return depends on:
	//  - whether the type requires full-state in each reply (lds and cds).
	//  - whether the request is wildcard.
	// resourcesToReturn will include all the resource names to reply based on the changes detected.
	var resourcesToReturn []string
	switch {
	// For lds and cds, answers will always include all existing subscribed resources, with no regard to which resource was changed or removed.
	// For other types, the response only includes updated resources (sotw cannot notify for deletion).
	case !watch.sendFullStateResponses():
		// TODO(valerian-roche): remove this leak of delta/sotw behavior here.
		if !watch.isDelta() && !replyEvenIfEmpty && len(changedResources) == 0 {
			// If the request is not the initial one, and the type does not require full updates,
			// do not return if nothing is to be set.
			// For full-state resources an empty response does have a semantic meaning.
			return nil, nil
		}
		// changedResources is already filtered based on the subscription.
		resourcesToReturn = changedResources

	case sub.IsWildcard():
		// Include all resources for the type.
		resourcesToReturn = make([]string, 0, len(cache.resources))
		for resourceName := range cache.resources {
			resourcesToReturn = append(resourcesToReturn, resourceName)
		}

	default:
		// Include all resources matching the subscription, with no concern on whether it has been updated or not.
		requestedResources := sub.SubscribedResources()
		// The linear cache could be very large (e.g. containing all potential CLAs)
		// Therefore drives on the subscription requested resources.
		resourcesToReturn = make([]string, 0, len(requestedResources))
		for resourceName := range requestedResources {
			if _, ok := cache.resources[resourceName]; ok {
				resourcesToReturn = append(resourcesToReturn, resourceName)
			}
		}
		subscribedPrefixes := sub.SubscribedPrefixes()
		if len(subscribedPrefixes) > 0 {
			for cachedName := range cache.resources {
				if _, exactSub := requestedResources[cachedName]; exactSub {
					continue
				}
				for prefix := range subscribedPrefixes {
					if strings.HasPrefix(cachedName, prefix) {
						resourcesToReturn = append(resourcesToReturn, cachedName)
						break
					}
				}
			}
		}
	}

	return cache.buildResponseForResources(watch, resourcesToReturn, removedResources)
}

// buildResponseForResources builds a WatchResponse for the given resource names,
// applying the watch's transform (if any). The internal returnedVersions always
// carry the stored resource's inner version, even when the transform stamps a
// different version on the wire — keeping change detection in inner-version space.
// Must be called under lock.
func (cache *LinearCache) buildResponseForResources(watch watch, resourcesToReturn []string, removedResources []string) (WatchResponse, error) {
	sub := watch.getSubscription()
	transform := watch.getTransform()

	// returnedVersions includes all resources currently known to the subscription and their version.
	// Clone the current returned versions. The cache should not alter the subscription
	returnedVersions := maps.Clone(sub.ReturnedResources())

	resources := make([]*internal.CachedResource, 0, len(resourcesToReturn))
	for _, resourceName := range resourcesToReturn {
		cachedResource := cache.resources[resourceName]
		version, err := cachedResource.GetVersion(watch.useResourceVersion())
		if err != nil {
			return nil, fmt.Errorf("failed to compute version of %s: %w", resourceName, err)
		}

		emit := cachedResource
		if transform != nil {
			out, wireVersion, keep := transform(resourceName, cachedResource.GetRawResource().Resource, version)
			if !keep {
				continue
			}
			if wireVersion != "" {
				// Ship the transformed bytes stamped with wireVersion on the wire,
				// but keep the inner version in returnedVersions below.
				emit = internal.NewCachedResource(resourceName, out, internal.WithResourceVersion(wireVersion))
			}
		}

		resources = append(resources, emit)
		returnedVersions[resourceName] = version
	}

	// Cleanup resources no longer existing in the cache or no longer subscribed.
	// In sotw we cannot return those if not full state,
	// but this ensures we detect unsubscription then resubscription.
	for _, resourceName := range removedResources {
		delete(returnedVersions, resourceName)
	}

	return watch.buildResponse(resources, removedResources, returnedVersions, cache.getVersion()), nil
}

// RefreshWatches forces a resend of the resources matching matchResource to the
// transform-watches of a single node (keyed by the cache's NodeHash). It bypasses
// the version diff — used for out-of-band per-client updates where the stored
// resource is unchanged but the per-client view (and thus the wire version) is not.
// matchResource is called per subscribed resource name and must not block.
func (cache *LinearCache) RefreshWatches(nodeKey string, matchResource func(name string) bool) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	type tracked struct {
		id uint64
		w  watch
	}
	// Snapshot first: sending removes watches (one-shot), mutating nodeIndex.
	toRefresh := make([]tracked, 0, len(cache.nodeIndex[nodeKey]))
	for id, w := range cache.nodeIndex[nodeKey] {
		toRefresh = append(toRefresh, tracked{id, w})
	}

	for _, t := range toRefresh {
		sub := t.w.getSubscription()
		var matched []string
		for name := range sub.SubscribedResources() {
			if !matchResource(name) {
				continue
			}
			if _, ok := cache.resources[name]; ok {
				matched = append(matched, name)
			}
		}
		if len(matched) == 0 {
			continue
		}

		response, err := cache.buildResponseForResources(t.w, matched, nil)
		if err != nil {
			cache.log.Errorf("[linear cache] failed to build refresh response for node %q: %v", nodeKey, err)
			continue
		}
		if response == nil {
			continue
		}
		t.w.sendResponse(response)
		// One-shot, like notifyAll: the server re-creates the watch on the next request.
		if t.w.getSubscription().IsWildcard() {
			cache.removeWildcardWatch(t.id)
		} else {
			cache.removeWatch(t.id, t.w.getSubscription())
		}
	}
}

func (cache *LinearCache) notifyAll(modified []string) error {
	// Gather the list of watches impacted by the modified resources.
	triggeredWatches := newWatches()
	for _, name := range modified {
		maps.Copy(triggeredWatches, cache.resourceWatches[name])
		for prefix, pw := range cache.prefixWatches {
			if strings.HasPrefix(name, prefix) {
				maps.Copy(triggeredWatches, pw)
			}
		}
	}

	// non-wildcard and prefix watches
	for watchID, watch := range triggeredWatches {
		response, err := cache.computeResponse(watch, false)
		if err != nil {
			return err
		}

		if response != nil {
			watch.sendResponse(response)
			cache.removeWatch(watchID, watch.getSubscription())
		} else {
			cache.log.Infof("[Linear cache] Watch %d detected as triggered but no change was found", watchID)
		}
	}

	for watchID, watch := range cache.wildcardWatches {
		response, err := cache.computeResponse(watch, false)
		if err != nil {
			return err
		}

		if response != nil {
			watch.sendResponse(response)
			cache.removeWildcardWatch(watchID)
		} else {
			cache.log.Infof("[Linear cache] Wildcard watch %d detected as triggered but no change was found", watchID)
		}
	}

	return nil
}

func (cache *LinearCache) isResourceAllowed(res types.SnapshotResource) error {
	if res.TypeURL() != cache.typeURL {
		return fmt.Errorf("resource %s has wrong type: expected %s and received %s", res.Name, cache.typeURL, res.Serialized.GetTypeUrl())
	}
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
	cache.resources[name] = internal.NewCachedResource(name, res, internal.WithCacheVersion(cache.getVersion()))

	return cache.notifyAll([]string{name})
}

// UpdateSnapshotResource updates a resource in the collection.
func (cache *LinearCache) UpdateSnapshotResource(res types.SnapshotResource) error {
	if err := cache.isResourceAllowed(res); err != nil {
		return err
	}
	cache.mu.Lock()
	defer cache.mu.Unlock()

	cache.version++
	cache.resources[res.Name] = res.AsCachedResource(cache.getVersion())

	return cache.notifyAll([]string{res.Name})
}

// DeleteResource removes a resource in the collection.
func (cache *LinearCache) DeleteResource(name string) error {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	cache.version++
	delete(cache.resources, name)

	return cache.notifyAll([]string{name})
}

// UpdateResources updates/deletes a list of resources in the cache.
// Calling UpdateResources instead of iterating on UpdateResource and DeleteResource
// is significantly more efficient when using delta or wildcard watches.
func (cache *LinearCache) UpdateResources(toUpdate map[string]types.Resource, toDelete []string) error {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	cache.version++
	version := cache.getVersion()
	modified := make([]string, 0, len(toUpdate)+len(toDelete))
	for name, resource := range toUpdate {
		cache.resources[name] = internal.NewCachedResource(name, resource, internal.WithCacheVersion(version))
		modified = append(modified, name)
	}
	for _, name := range toDelete {
		delete(cache.resources, name)
		modified = append(modified, name)
	}

	return cache.notifyAll(modified)
}

// UpdateSnapshotResources updates/deletes a list of resources in the cache.
// Calling UpdateSnapshotResources instead of iterating on UpdateSnapshotResource and DeleteResource
// is significantly more efficient when using delta or wildcard watches.
func (cache *LinearCache) UpdateSnapshotResources(toUpdate []types.SnapshotResource, toDelete []string) error {
	// Do a validation pass beforehand to not leave the cache in a half state.
	for _, res := range toUpdate {
		if err := cache.isResourceAllowed(res); err != nil {
			return err
		}
	}

	cache.mu.Lock()
	defer cache.mu.Unlock()

	cache.version++
	version := cache.getVersion()
	modified := make([]string, 0, len(toUpdate)+len(toDelete))
	for _, res := range toUpdate {
		cache.resources[res.Name] = res.AsCachedResource(version)
		modified = append(modified, res.Name)
	}
	for _, name := range toDelete {
		delete(cache.resources, name)
		modified = append(modified, name)
	}

	return cache.notifyAll(modified)
}

// SetResources replaces current resources with a new set of resources.
// Given the use of lazy serialization, if most resources are actually the same
// using UpdateResources instead will be much more efficient.
func (cache *LinearCache) SetResources(resources map[string]types.Resource) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	cache.version++
	version := cache.getVersion()

	modified := make([]string, 0, len(resources))
	// Collect deleted resource names.
	for name := range cache.resources {
		if _, found := resources[name]; !found {
			delete(cache.resources, name)
			modified = append(modified, name)
		}
	}

	// We assume all resources passed to SetResources are changed.
	// Otherwise we would have to do proto.Equal on resources which is pretty expensive operation
	for name, resource := range resources {
		cache.resources[name] = internal.NewCachedResource(name, resource, internal.WithCacheVersion(version))
		modified = append(modified, name)
	}

	if err := cache.notifyAll(modified); err != nil {
		cache.log.Errorf("Failed to notify watches: %s", err.Error())
	}
}

// SetSnapshotResources replaces current resources with a new set of resources.
// Given the use of lazy serialization, if most resources are actually the same
// using UpdateResources instead will be much more efficient.
func (cache *LinearCache) SetSnapshotResources(resources []types.SnapshotResource) error {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	cache.version++
	version := cache.getVersion()

	modified := make([]string, 0, len(resources))
	newResources := make(map[string]*internal.CachedResource, len(resources))
	// We assume all resources passed to SetSnapshotResources are changed.
	// Otherwise we would have to do proto.Equal on resources which is pretty expensive operation
	for _, res := range resources {
		if err := cache.isResourceAllowed(res); err != nil {
			return err
		}
		newResources[res.Name] = res.AsCachedResource(version)
		modified = append(modified, res.Name)
	}

	// Collect deleted resource names.
	for name := range cache.resources {
		if _, found := newResources[name]; !found {
			modified = append(modified, name)
		}
	}

	cache.resources = newResources
	return cache.notifyAll(modified)
}

// GetResources returns current resources stored in the cache.
func (cache *LinearCache) GetResources() map[string]types.Resource {
	cache.mu.RLock()
	defer cache.mu.RUnlock()

	// create a copy of our internal storage to avoid data races
	// involving mutations of our backing map
	resources := make(map[string]types.Resource, len(cache.resources))
	for k, v := range cache.resources {
		resources[k] = v.GetRawResource().Resource
	}
	return resources
}

// GetResource returns the resource stored under name and its version.
// A nil resource with a nil error means the name is not in the cache.
//
// useResourceVersion controls the returned version:
//   - if set to false versions reflect the cache version when the entry was added
//   - if set to true versions are a stable property of the resource, with no regard to when it was added to the cache.
func (cache *LinearCache) GetResource(name string, useResourceVersion bool) (types.Resource, string, error) {
	cache.mu.RLock()
	defer cache.mu.RUnlock()
	r, ok := cache.resources[name]
	if !ok {
		return nil, "", nil
	}
	version, err := r.GetVersion(useResourceVersion)
	if err != nil {
		return nil, "", fmt.Errorf("linear cache: get version for %q: %w", name, err)
	}
	return r.GetRawResource().Resource, version, nil
}

// The implementations of sotw and delta watches handling is nearly identical. The main distinctions are:
//   - handling of version in sotw when the request is the first of a subscription. Delta has a proper handling based on the request providing known versions.
//   - building the initial resource versions in delta if they've not been computed yet.
//   - computeSotwResponse and computeDeltaResponse has slightly different implementations due to sotw requirements to return full state for certain resources only.
func (cache *LinearCache) CreateWatch(request *Request, sub Subscription, value chan Response) (func(), error) {
	return cache.createWatch(request, sub, value, nil)
}

// CreateWatchWithTransform is CreateWatch with a per-watch ResourceTransform
// applied when building responses, and the watch indexed by node (if WithNodeHash
// is set) for RefreshWatches.
func (cache *LinearCache) CreateWatchWithTransform(request *Request, sub Subscription, value chan Response, transform ResourceTransform) (func(), error) {
	return cache.createWatch(request, sub, value, transform)
}

func (cache *LinearCache) createWatch(request *Request, sub Subscription, value chan Response, transform ResourceTransform) (func(), error) {
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
	// For now it is not done as:
	//  - for the first case, while the protocol documentation does not explicitly mention the case, it does not mark it impossible and explicitly references unsubscribing from wildcard.
	//  - for the second one we could likely do it with little difficulty if need be, but if users rely on the current monotonic version it could impact their callbacks implementations.
	watch := ResponseWatch{
		Request:            request,
		Response:           value,
		subscription:       sub,
		fullStateResponses: ResourceRequiresFullStateInSotw(cache.typeURL),
		transform:          transform,
	}

	cache.mu.Lock()
	defer cache.mu.Unlock()

	response, err := cache.computeResponse(watch, replyEvenIfEmpty)
	if err != nil {
		return nil, fmt.Errorf("failed to compute the watch respnse: %w", err)
	}
	if response != nil {
		cache.log.Debugf("[linear cache] replying to the watch with resources %v (subscription values %v, known %v)", response.GetReturnedResources(), sub.SubscribedResources(), sub.ReturnedResources())
		watch.sendResponse(response)
		return func() {}, nil
	}

	return cache.trackWatch(watch), nil
}

func (cache *LinearCache) CreateDeltaWatch(request *DeltaRequest, sub Subscription, value chan DeltaResponse) (func(), error) {
	return cache.createDeltaWatch(request, sub, value, nil)
}

// CreateDeltaWatchWithTransform is CreateDeltaWatch with a per-watch
// ResourceTransform applied when building responses, and the watch indexed by
// node (if WithNodeHash is set) for RefreshWatches.
func (cache *LinearCache) CreateDeltaWatchWithTransform(request *DeltaRequest, sub Subscription, value chan DeltaResponse, transform ResourceTransform) (func(), error) {
	return cache.createDeltaWatch(request, sub, value, transform)
}

func (cache *LinearCache) createDeltaWatch(request *DeltaRequest, sub Subscription, value chan DeltaResponse, transform ResourceTransform) (func(), error) {
	if request.GetTypeUrl() != cache.typeURL {
		return nil, fmt.Errorf("request type %s does not match cache type %s", request.GetTypeUrl(), cache.typeURL)
	}

	watch := DeltaResponseWatch{Request: request, Response: value, subscription: sub, transform: transform}

	// On first request on a wildcard subscription, envoy does expect a response to come in to
	// conclude initialization.
	replyEvenIfEmpty := sub.IsWildcard() && request.GetResponseNonce() == ""

	cache.mu.Lock()
	defer cache.mu.Unlock()

	response, err := cache.computeResponse(watch, replyEvenIfEmpty)
	if err != nil {
		return nil, fmt.Errorf("failed to compute the watch respnse: %w", err)
	}

	if response != nil {
		cache.log.Debugf("[linear cache] replying to the delta watch (subscription values %v, known %v)", sub.SubscribedResources(), sub.ReturnedResources())
		watch.sendResponse(response)
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

// Must be called under lock.
func (cache *LinearCache) trackWatch(watch watch) func() {
	cache.watchCount++

	watchID := cache.nextWatchID()
	sub := watch.getSubscription()

	// Index transform-watches by node so RefreshWatches can target a single
	// client. Cleaned by removeWatch/removeWildcardWatch via nodeKeyByWatch.
	if cache.nodeHash != nil && watch.getTransform() != nil {
		cache.addNodeIndex(cache.nodeHash.ID(watch.getNode()), watchID, watch)
	}

	if sub.IsWildcard() {
		cache.log.Infof("[linear cache] open watch %d (delta: %t) for %s all resources", watchID, watch.isDelta(), cache.typeURL)
		cache.log.Debugf("[linear cache] subscription details for watch %d: known versions %v, system version %q", watchID, sub.ReturnedResources(), cache.getVersion())
		cache.wildcardWatches[watchID] = watch
		return func() {
			cache.mu.Lock()
			defer cache.mu.Unlock()
			cache.removeWildcardWatch(watchID)
		}
	}

	cache.log.Infof("[linear cache] open watch %d (delta: %t) for %s resources %v prefixes %v", watchID, watch.isDelta(), cache.typeURL, sub.SubscribedResources(), sub.SubscribedPrefixes())
	cache.log.Debugf("[linear cache] subscription details for watch %d: known versions %v, system version %q", watchID, sub.ReturnedResources(), cache.getVersion())
	for name := range sub.SubscribedResources() {
		watches, exists := cache.resourceWatches[name]
		if !exists {
			watches = newWatches()
			cache.resourceWatches[name] = watches
		}
		watches[watchID] = watch
	}
	for prefix := range sub.SubscribedPrefixes() {
		pw, exists := cache.prefixWatches[prefix]
		if !exists {
			pw = newWatches()
			cache.prefixWatches[prefix] = pw
		}
		pw[watchID] = watch
	}
	return func() {
		cache.mu.Lock()
		defer cache.mu.Unlock()
		cache.removeWatch(watchID, sub)
	}
}

// Must be called under lock.
func (cache *LinearCache) removeWildcardWatch(watchID uint64) {
	cache.watchCount--
	delete(cache.wildcardWatches, watchID)
	cache.removeFromNodeIndex(watchID)
}

// addNodeIndex records a transform-watch under its node key. Must be called under lock.
func (cache *LinearCache) addNodeIndex(nodeKey string, watchID uint64, watch watch) {
	ws, ok := cache.nodeIndex[nodeKey]
	if !ok {
		ws = newWatches()
		cache.nodeIndex[nodeKey] = ws
	}
	ws[watchID] = watch
	cache.nodeKeyByWatch[watchID] = nodeKey
}

// removeFromNodeIndex drops a watch from the node index, if present. Must be
// called under lock. Safe to call for watches that were never indexed.
func (cache *LinearCache) removeFromNodeIndex(watchID uint64) {
	nodeKey, ok := cache.nodeKeyByWatch[watchID]
	if !ok {
		return
	}
	delete(cache.nodeKeyByWatch, watchID)
	ws := cache.nodeIndex[nodeKey]
	delete(ws, watchID)
	if len(ws) == 0 {
		delete(cache.nodeIndex, nodeKey)
	}
}

// Must be called under lock.
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
	for prefix := range sub.SubscribedPrefixes() {
		pw := cache.prefixWatches[prefix]
		delete(pw, watchID)
		if len(pw) == 0 {
			delete(cache.prefixWatches, prefix)
		}
	}
	cache.watchCount--
	cache.removeFromNodeIndex(watchID)
}

func (cache *LinearCache) getVersion() string {
	return cache.versionPrefix + strconv.FormatUint(cache.version, 10)
}

func (cache *LinearCache) Fetch(context.Context, *Request) (Response, error) {
	return nil, errors.New("fetch is not implemented by LinearCache")
}

// NumResources returns the number of resources currently in the cache.
// As GetResources is building a clone it is expensive to get metrics otherwise.
func (cache *LinearCache) NumResources() int {
	cache.mu.RLock()
	defer cache.mu.RUnlock()
	return len(cache.resources)
}

// NumWatches returns the number of active watches for a resource name, including wildcard and prefix ones.
func (cache *LinearCache) NumWatches(name string) int {
	cache.mu.RLock()
	defer cache.mu.RUnlock()
	count := len(cache.resourceWatches[name]) + len(cache.wildcardWatches)
	for prefix, pw := range cache.prefixWatches {
		if strings.HasPrefix(name, prefix) {
			count += len(pw)
		}
	}
	return count
}

// NumPrefixWatches returns the total number of active prefix watches.
func (cache *LinearCache) NumPrefixWatches() int {
	cache.mu.RLock()
	defer cache.mu.RUnlock()
	count := 0
	for _, pw := range cache.prefixWatches {
		count += len(pw)
	}
	return count
}

// NumWildcardWatches returns the number of wildcard watches.
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
