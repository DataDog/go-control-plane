package types

import (
	"context"
	"time"

	"google.golang.org/protobuf/proto"

	endpoint "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	discovery "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
)

// Request is an alias for the discovery request type.
type Request = discovery.DiscoveryRequest

// DeltaRequest is an alias for the delta discovery request type.
type DeltaRequest = discovery.DeltaDiscoveryRequest

// Response is a wrapper around Envoy's DiscoveryResponse.
type Response interface {
	// GetDiscoveryResponse returns the Constructed DiscoveryResponse.
	GetDiscoveryResponse() (*discovery.DiscoveryResponse, error)

	// GetRequest returns the request that created the watch that we're now responding to.
	// This is provided to allow the caller to correlate the response with a request.
	// Generally this will be the latest request seen on the stream for the specific type.
	GetRequest() *discovery.DiscoveryRequest

	// GetVersion returns the version in the Response.
	// The version can be a property of the resources, allowing for optimizations in subsequent calls,
	// or simply an internal property of the cache which can be used for debugging.
	// The cache implementation should be able to determine if it can provide such optimization.
	// Deprecated: use GetResponseVersion instead
	GetVersion() (string, error)

	// GetResponseVersion returns the version in the Response.
	// The version can be a property of the resources, allowing for optimizations in subsequent calls,
	// or simply an internal property of the cache which can be used for debugging.
	// The cache implementation should be able to determine if it can provide such optimization.
	GetResponseVersion() string

	// GetReturnedResources returns the map of resources and their versions returned in the subscription.
	// It may include more resources than directly set in the response to consider the full state of the client.
	// The caller is expected to provide this unchanged to the next call to CreateWatch as part of the subscription.
	GetReturnedResources() map[string]string

	// GetContext returns the context provided during response creation.
	GetContext() context.Context
}

// DeltaResponse is a wrapper around Envoy's DeltaDiscoveryResponse.
type DeltaResponse interface {
	// GetDeltaDiscoveryResponse returns the constructed DeltaDiscoveryResponse.
	GetDeltaDiscoveryResponse() (*discovery.DeltaDiscoveryResponse, error)

	// GetDeltaRequest returns the request that created the watch that we're now responding to.
	// This is provided to allow the caller to correlate the response with a request.
	// Generally this will be the latest request seen on the stream for the specific type.
	GetDeltaRequest() *discovery.DeltaDiscoveryRequest

	// GetSystemVersion returns the version in the DeltaResponse.
	// The version in delta response is not indicative of the resources included,
	// but an internal property of the cache which can be used for debugging.
	// Deprecated: use GetResponseVersion instead
	GetSystemVersion() (string, error)

	// GetResponseVersion returns the version in the DeltaResponse.
	// The version in delta response is not indicative of the resources included,
	// but an internal property of the cache which can be used for debugging.
	GetResponseVersion() string

	// GetNextVersionMap provides the version map of the internal cache.
	// The version map consists of updated version mappings after this response is applied.
	// Deprecated: use GetReturnedResources instead
	GetNextVersionMap() map[string]string

	// GetReturnedResources provides the version map of the internal cache.
	// The version map consists of updated version mappings after this response is applied.
	GetReturnedResources() map[string]string

	// GetContext returns the context provided during response creation.
	GetContext() context.Context
}

// Subscription stores the server view of the client state for a given resource type.
// This allows proper implementation of stateful aspects of the protocol (e.g. returning only some updated resources).
// Though the methods may return mutable parts of the state for performance reasons,
// the cache is expected to consider this state as immutable and thread safe between a watch creation and its cancellation.
type Subscription interface {
	// ReturnedResources returns a list of resources that were sent to the client and their associated versions.
	// The versions are:
	//  - delta protocol: version of the specific resource set in the response.
	//  - sotw protocol: version of the global response when the resource was last sent.
	ReturnedResources() map[string]string

	// SubscribedResources returns the list of resources currently subscribed to by the client for the type.
	// For delta it keeps track of subscription updates across requests
	// For sotw it is a normalized view of the last request resources
	SubscribedResources() map[string]struct{}

	// IsWildcard returns whether the client has a wildcard watch.
	// This considers subtleties related to the current migration of wildcard definitions within the protocol.
	// More details on the behavior of wildcard are present at https://www.envoyproxy.io/docs/envoy/latest/api-docs/xds_protocol#how-the-client-specifies-what-resources-to-return
	IsWildcard() bool
}

// Resource is the base interface for the xDS payload.
type Resource interface {
	proto.Message
}

// ResourceWithTTL is a Resource with an optional TTL.
type ResourceWithTTL struct {
	Resource Resource
	TTL      *time.Duration
}

// ResourceWithName provides a name for out-of-tree resources.
type ResourceWithName interface {
	proto.Message
	GetName() string
}

// GetResourceName returns the resource name for a valid xDS response type.
func GetResourceName(res Resource) string {
	switch v := res.(type) {
	case *endpoint.ClusterLoadAssignment:
		return v.GetClusterName()
	case ResourceWithName:
		return v.GetName()
	default:
		return ""
	}
}

// GetResourceName returns the resource names for a list of valid xDS response types.
func GetResourceNames(resources []ResourceWithTTL) []string {
	out := make([]string, len(resources))
	for i, r := range resources {
		out[i] = GetResourceName(r.Resource)
	}
	return out
}

type MarshaledResource = []byte

// SkipFetchError is the error returned when the cache fetch is short
// circuited due to the client's version already being up-to-date.
type SkipFetchError struct{}

// Error satisfies the error interface
func (e SkipFetchError) Error() string {
	return "skip fetch: version up to date"
}

// ResponseType enumeration of supported response types
type ResponseType int

// NOTE: The order of this enum MATTERS!
// https://www.envoyproxy.io/docs/envoy/latest/api-docs/xds_protocol#aggregated-discovery-service
// ADS expects things to be returned in a specific order.
// See the following issue for details: https://github.com/envoyproxy/go-control-plane/issues/526
const (
	Cluster ResponseType = iota
	Endpoint
	Listener
	Route
	ScopedRoute
	VirtualHost
	Secret
	Runtime
	ExtensionConfig
	RateLimitConfig
	UnknownType // token to count the total number of supported types
)
