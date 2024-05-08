package cache

import (
	"fmt"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	"github.com/envoyproxy/go-control-plane/pkg/resource/v3"
)

// Resources is a versioned group of resources.
type Resources struct {
	// Version information.
	Version string

	// Items in the group indexed by name.
	Items map[string]types.ResourceWithTTL
}

// cachedResource is used to track resources added by the user in the cache.
// It contains the resource itself and its associated version (currently in two different modes).
type cachedResource struct {
	// name is the resource name, to avoid having to rely on type-specific logic to retrieve
	// the resource name.
	name string

	// ResourceWithTTL is the resource provided by the user. In the future this would likely also
	// allow to provide opaque, serialized resources.
	types.ResourceWithTTL

	// cacheVersion is the version of the cache at the time of last update, used in sotw.
	cacheVersion string

	// resourceVersion is the version of the resource itself (a hash of its content after deterministic marshaling).
	// It is lazy initialized and should be accessed through getStableVersion.
	resourceVersion string

	// marshaledContent is the serialized version of the resource, built through deterministic marshaling.
	// It is lazy initialized on getStableVersion calls.
	marshaledContent []byte
}

func newCachedResource(name string, res types.Resource, cacheVersion string) *cachedResource {
	return &cachedResource{
		name:            name,
		ResourceWithTTL: types.ResourceWithTTL{Resource: res},
		cacheVersion:    cacheVersion,
	}
}

func newCachedResources(res map[string]types.ResourceWithTTL, cacheVersion string) map[string]*cachedResource {
	resources := make(map[string]*cachedResource, len(res))
	for name, r := range res {
		resources[name] = &cachedResource{
			name:            name,
			ResourceWithTTL: r,
			cacheVersion:    cacheVersion,
		}
	}
	return resources
}

func (c *cachedResource) getStableVersion() (string, error) {
	if c.resourceVersion != "" {
		return c.resourceVersion, nil
	}

	marshaledResource, err := MarshalResource(c.Resource)
	if err != nil {
		return "", err
	}
	c.marshaledContent = marshaledResource
	c.resourceVersion = HashResource(marshaledResource)
	return c.resourceVersion, nil
}

func (c *cachedResource) getVersion(useStableVersion bool) (string, error) {
	if !useStableVersion {
		return c.cacheVersion, nil
	}

	return c.getStableVersion()
}

type returnedResource struct {
	types.ResourceWithTTL

	// name of the resource. This avoids the need of introspection, which cannot necessarily be used with opaque types.
	name string
	// content and version can be set on resource creation to avoid reserializing the resource in every single response to all clients.
	// If not set the computation will occur on access and will not be latched to avoid concurrency issues.
	serialized *anypb.Any
	version    string
}

func newReturnedResource(res types.ResourceWithTTL) returnedResource {
	return returnedResource{
		ResourceWithTTL: res,
		name:            GetResourceName(res.Resource),
	}
}

func newReturnedResourceFromCache(res *cachedResource) returnedResource {
	r := returnedResource{
		name:            res.name,
		ResourceWithTTL: res.ResourceWithTTL,
		version:         res.resourceVersion,
	}
	if res.marshaledContent != nil {
		r.serialized = anyFromContent(res.marshaledContent, res.Resource)
	}
	return r
}

func anyFromContent(content []byte, res proto.Message) *anypb.Any {
	a := new(anypb.Any)
	a.Value = content
	a.TypeUrl = resource.APITypePrefix + string(proto.MessageName(res))
	return a
}

func (c returnedResource) buildAnyEntryWithVersion() (*anypb.Any, string, error) {
	a, err := c.buildAnyEntry()
	if err != nil {
		return nil, "", err
	}

	if c.version != "" {
		return a, c.version, nil
	}

	return a, HashResource(a.Value), nil
}

func (c returnedResource) buildAnyEntry() (*anypb.Any, error) {
	if c.Resource == nil {
		return nil, nil
	}

	if c.serialized != nil {
		return c.serialized, nil
	}

	m, err := MarshalResource(c.Resource)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize resource: %w", err)
	}

	return anyFromContent(m, c.Resource), nil
}

// getReturnedResourceNames returns the resource names for a list of valid xDS response types.
func getReturnedResourceNames(resources []returnedResource) []string {
	out := make([]string, len(resources))
	for i, r := range resources {
		out[i] = r.name
	}
	return out
}

// IndexResourcesByName creates a map from the resource name to the resource.
func IndexResourcesByName(items []types.ResourceWithTTL) map[string]types.ResourceWithTTL {
	indexed := make(map[string]types.ResourceWithTTL, len(items))
	for _, item := range items {
		indexed[GetResourceName(item.Resource)] = item
	}
	return indexed
}

// IndexReturnedResourcesByName creates a map from the resource name to the resource.
func IndexReturnedResourcesByName(items []returnedResource) map[string]types.ResourceWithTTL {
	indexed := make(map[string]types.ResourceWithTTL, len(items))
	for _, item := range items {
		indexed[item.name] = item.ResourceWithTTL
	}
	return indexed
}

// IndexRawResourcesByName creates a map from the resource name to the resource.
func IndexRawResourcesByName(items []types.Resource) map[string]types.Resource {
	indexed := make(map[string]types.Resource, len(items))
	for _, item := range items {
		indexed[GetResourceName(item)] = item
	}
	return indexed
}

// NewResources creates a new resource group.
func NewResources(version string, items []types.Resource) Resources {
	itemsWithTTL := make([]types.ResourceWithTTL, 0, len(items))
	for _, item := range items {
		itemsWithTTL = append(itemsWithTTL, types.ResourceWithTTL{Resource: item})
	}
	return NewResourcesWithTTL(version, itemsWithTTL)
}

// NewResourcesWithTTL creates a new resource group.
func NewResourcesWithTTL(version string, items []types.ResourceWithTTL) Resources {
	return Resources{
		Version: version,
		Items:   IndexResourcesByName(items),
	}
}
