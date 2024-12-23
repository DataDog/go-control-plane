package snapshot

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	"google.golang.org/protobuf/proto"
)

// cachedResource is used to track resources added by the user in the cache.
// It contains the resource itself and its associated version (currently in two different modes).
type CachedResource struct {
	Name string

	resource types.Resource
	Ttl      *time.Duration

	// snapshotVersion is the version of the snapshot at the time of last update, used in sotw.
	snapshotVersion string

	marshalFunc                func() ([]byte, error)
	computeResourceVersionFunc func() (string, error)
}

func NewCachedResource(name string, res types.Resource, cacheVersion string) *CachedResource {
	marshalFunc := sync.OnceValues(func() ([]byte, error) {
		return marshalResource(res)
	})
	return &CachedResource{
		Name:            name,
		resource:        res,
		snapshotVersion: cacheVersion,
		marshalFunc:     marshalFunc,
		computeResourceVersionFunc: sync.OnceValues(func() (string, error) {
			marshaled, err := marshalFunc()
			if err != nil {
				return "", fmt.Errorf("marshaling resource: %w", err)
			}
			return hashResource(marshaled), nil
		}),
	}
}

func NewCachedResourceWithTTL(name string, res types.ResourceWithTTL, cacheVersion string) *CachedResource {
	cachedRes := NewCachedResource(name, res.Resource, cacheVersion)
	cachedRes.Ttl = res.TTL
	return cachedRes
}

func (c *CachedResource) GetResource() types.Resource {
	return c.resource
}

// GetMarshaledResource lazily marshals the resource and returns the bytes.
func (c *CachedResource) GetMarshaledResource() ([]byte, error) {
	return c.marshalFunc()
}

// getResourceVersion lazily hashes the resource and returns the stable hash used to track version changes.
func (c *CachedResource) getResourceVersion() (string, error) {
	return c.computeResourceVersionFunc()
}

// getVersion returns the requested version.
func (c *CachedResource) GetVersion(useResourceVersion bool) (string, error) {
	if !useResourceVersion {
		return c.snapshotVersion, nil
	}

	return c.getResourceVersion()
}

// marshalResource converts the Resource to MarshaledResource.
func marshalResource(resource types.Resource) (types.MarshaledResource, error) {
	return proto.MarshalOptions{Deterministic: true}.Marshal(resource)
}

// hashResource will take a resource and create a SHA256 hash sum out of the marshaled bytes
func hashResource(resource []byte) string {
	hasher := sha256.New()
	hasher.Write(resource)

	return hex.EncodeToString(hasher.Sum(nil))
}
