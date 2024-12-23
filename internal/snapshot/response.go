package snapshot

import (
	"fmt"
	"maps"

	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
)

type WatchResponse interface {
	GetReturnedResources() map[string]string
	GetResponseVersion() string
}

type Watch interface {
	// IsDelta indicates whether the watch is a delta one.
	// It should not be used to take functional decisions, but is still currently used pending final changes.
	// It can be used to generate statistics.
	IsDelta() bool
	// UseResourceVersion indicates whether versions returned in the response are built using resource versions instead of cache update versions.
	UseResourceVersion() bool
	// SendFullStateResponses requires that all resources matching the request, with no regards to which ones actually updated, must be provided in the response.
	// As a consequence, sending a response with no resources has a functional meaning of no matching resources available.
	SendFullStateResponses() bool

	GetSubscription() types.Subscription
	// SendResponse sends the response for the watch.
	// It must be called at most once.
	SendResponse(resp WatchResponse)

	// BuildResponse computes the actual WatchResponse object to be sent on the watch.
	// cachedResource instances are passed by copy as their lazy accessors are not thread-safe.
	BuildResponse(updatedResources []*CachedResource, removedResources []string, returnedVersions map[string]string, version string) WatchResponse
}

// computeResourceChangeLocked compares the subscription known resources and the cache current state to compute the list of resources
// which have changed and should be notified to the user.
//
// The useResourceVersion argument defines what version type to use for resources:
//   - if set to false versions are based on when resources were updated in the cache.
//   - if set to true versions are a stable property of the resource, with no regard to when it was added to the cache.
func computeResourceChangeLocked(sub types.Subscription, useResourceVersion bool, cachedResources map[string]*CachedResource) (updated, removed []string, err error) {
	var changedResources []string
	var removedResources []string

	knownVersions := sub.ReturnedResources()
	if sub.IsWildcard() {
		for resourceName, resource := range cachedResources {
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
			if _, ok := cachedResources[resourceName]; !ok {
				removedResources = append(removedResources, resourceName)
			}
		}
	} else {
		for resourceName := range sub.SubscribedResources() {
			res, exists := cachedResources[resourceName]
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
			} else {
				resourceVersion, err := res.GetVersion(useResourceVersion)
				if err != nil {
					return nil, nil, fmt.Errorf("failed to compute version of %s: %w", resourceName, err)
				}
				if knownVersion != resourceVersion {
					// The client knows an outdated version.
					changedResources = append(changedResources, resourceName)
				}
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

	return changedResources, removedResources, nil
}

func (s *TypedSnapshot) ComputeResponse(watch Watch, replyEvenIfEmpty bool, cacheVersion string) (WatchResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	sub := watch.GetSubscription()
	changedResources, removedResources, err := computeResourceChangeLocked(sub, watch.UseResourceVersion(), s.resources)
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
	case !watch.SendFullStateResponses():
		// TODO(valerian-roche): remove this leak of delta/sotw behavior here.
		if !watch.IsDelta() && !replyEvenIfEmpty && len(changedResources) == 0 {
			// If the request is not the initial one, and the type does not require full updates,
			// do not return if nothing is to be set.
			// For full-state resources an empty response does have a semantic meaning.
			return nil, nil
		}

		// changedResources is already filtered based on the subscription.
		resourcesToReturn = changedResources
	case sub.IsWildcard():
		// Include all resources for the type.
		resourcesToReturn = make([]string, 0, len(s.resources))
		for resourceName := range s.resources {
			resourcesToReturn = append(resourcesToReturn, resourceName)
		}
	default:
		// Include all resources matching the subscription, with no concern on whether it has been updated or not.
		requestedResources := sub.SubscribedResources()
		// The linear cache could be very large (e.g. containing all potential CLAs)
		// Therefore drives on the subscription requested resources.
		resourcesToReturn = make([]string, 0, len(requestedResources))
		for resourceName := range requestedResources {
			if _, ok := s.resources[resourceName]; ok {
				resourcesToReturn = append(resourcesToReturn, resourceName)
			}
		}
	}

	// returnedVersions includes all resources currently known to the subscription and their version.
	// Clone the current returned versions. The cache should not alter the subscription.
	returnedVersions := maps.Clone(sub.ReturnedResources())

	resources := make([]*CachedResource, 0, len(resourcesToReturn))
	for _, resourceName := range resourcesToReturn {
		cachedResource := s.resources[resourceName]
		resources = append(resources, cachedResource)
		version, err := cachedResource.GetVersion(watch.UseResourceVersion())
		if err != nil {
			return nil, fmt.Errorf("failed to compute version of %s: %w", resourceName, err)
		}
		returnedVersions[resourceName] = version
	}
	// Cleanup resources no longer existing in the cache or no longer subscribed.
	// In sotw we cannot return those if not full state,
	// but this ensures we detect unsubscription then resubscription.
	for _, resourceName := range removedResources {
		delete(returnedVersions, resourceName)
	}

	return watch.BuildResponse(resources, removedResources, returnedVersions, cacheVersion), nil
}
