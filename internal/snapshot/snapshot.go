package snapshot

import (
	"strconv"
	"sync"

	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
)

type TypedSnapshot struct {
	// resources contains all resources currently set in the cache and associated versions.
	resources map[string]*CachedResource

	version uint64

	// updateCallback is called on any resource update.
	// It can be safely set to nil if unneeded.
	updateCallback func(resourceNames []string) error

	mu sync.RWMutex
}

type SnapshotOptions = func(*TypedSnapshot)

func NewTypedSnapshot(options ...func(*TypedSnapshot)) *TypedSnapshot {
	s := &TypedSnapshot{
		resources: make(map[string]*CachedResource),
	}
	for _, opt := range options {
		opt(s)
	}
	return s
}

func WithUpdateCallback(cb func(resourceNames []string) error) SnapshotOptions {
	return func(ts *TypedSnapshot) {
		ts.updateCallback = cb
	}
}

func WithInitialResources(resources map[string]types.Resource) SnapshotOptions {
	return func(ts *TypedSnapshot) {
		ts.SetInitialResources(resources)
	}
}

func (s *TypedSnapshot) SetInitialResources(resources map[string]types.Resource) {
	s.resources, _ = setResources(resources, nil, s.versionLocked(), cachedResourceBuilder)
}

func (s *TypedSnapshot) incrementVersionLocked() string {
	s.version++
	return s.versionLocked()
}

func (s *TypedSnapshot) versionLocked() string {
	return strconv.FormatUint(s.version, 10)
}

func (s *TypedSnapshot) UpdateResources(resources map[string]types.Resource, removed []string) error {
	s.mu.Lock()
	updated := updateResources(resources, removed, s.resources, s.incrementVersionLocked(), cachedResourceBuilder)
	s.mu.Unlock()
	if s.updateCallback != nil {
		return s.updateCallback(updated)
	}
	return nil
}

func (s *TypedSnapshot) UpdateResourcesWithTTL(resources map[string]types.ResourceWithTTL, removed []string) error {
	s.mu.Lock()
	updated := updateResources(resources, removed, s.resources, s.incrementVersionLocked(), cachedResourceWithTTLBuilder)
	s.mu.Unlock()
	if s.updateCallback != nil {
		return s.updateCallback(updated)
	}
	return nil
}

func (s *TypedSnapshot) SetResources(resources map[string]types.Resource) error {
	s.mu.Lock()
	newResources, impactedResources := setResources(resources, s.resources, s.incrementVersionLocked(), cachedResourceBuilder)
	s.resources = newResources
	s.mu.Unlock()

	if s.updateCallback != nil {
		return s.updateCallback(impactedResources)
	}
	return nil
}

func (s *TypedSnapshot) SetResourcesWithTTL(resources map[string]types.ResourceWithTTL) error {
	s.mu.Lock()
	newResources, impactedResources := setResources(resources, s.resources, s.incrementVersionLocked(), cachedResourceWithTTLBuilder)
	s.resources = newResources
	s.mu.Unlock()

	if s.updateCallback != nil {
		return s.updateCallback(impactedResources)
	}
	return nil
}

func (s *TypedSnapshot) NumResources() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.resources)
}

// GetResources returns current resources stored in the snapshot
func (s *TypedSnapshot) GetResources() map[string]types.Resource {
	s.mu.RLock()
	defer s.mu.RUnlock()

	resources := make(map[string]types.Resource, len(s.resources))
	for k, v := range s.resources {
		resources[k] = v.GetResource()
	}
	return resources
}

func cachedResourceBuilder(name string, res types.Resource, cacheVersion string) *CachedResource {
	return NewCachedResource(name, res, cacheVersion)
}

func cachedResourceWithTTLBuilder(name string, res types.ResourceWithTTL, cacheVersion string) *CachedResource {
	return NewCachedResourceWithTTL(name, res, cacheVersion)
}

func updateResources[T any](newResources map[string]T, removed []string, resources map[string]*CachedResource, cacheVersion string, builder func(string, T, string) *CachedResource) []string {
	impactedResources := make([]string, 0, len(resources)+len(removed))
	for name, res := range newResources {
		resources[name] = builder(name, res, cacheVersion)
		impactedResources = append(impactedResources, name)
	}
	for _, resName := range removed {
		delete(resources, resName)
	}
	impactedResources = append(impactedResources, removed...)
	return impactedResources
}

func setResources[T any](newResources map[string]T, existingResources map[string]*CachedResource, cacheVersion string, builder func(string, T, string) *CachedResource) (map[string]*CachedResource, []string) {
	impactedResources := make([]string, 0, len(newResources))

	previousResources := existingResources
	resources := make(map[string]*CachedResource, len(newResources))
	for name, res := range newResources {
		resources[name] = builder(name, res, cacheVersion)
		impactedResources = append(impactedResources, name)
	}

	for name := range previousResources {
		if _, ok := resources[name]; !ok {
			impactedResources = append(impactedResources, name)
		}
	}
	return resources, impactedResources
}
