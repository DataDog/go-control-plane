package watches

import (
	"context"
	"fmt"
	"sync/atomic"

	discovery "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	"github.com/envoyproxy/go-control-plane/internal/snapshot"
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/durationpb"
)

// RawResponse is a pre-serialized xDS response containing the raw resources to
// be included in the final Discovery Response.
type RawResponse struct {
	// Request is the original request.
	Request *discovery.DiscoveryRequest

	// Version of the resources as tracked by the cache for the given type.
	// Proxy responds with this version as an acknowledgement.
	Version string

	// resources to be included in the response.
	resources []*snapshot.CachedResource

	// returnedResources tracks the resources returned for the subscription and the version when it was last returned,
	// including previously returned ones when using non-full state resources.
	// It allows the cache to know what the client knows. The server will transparently forward this
	// across requests, and the cache is responsible for its interpretation.
	returnedResources map[string]string

	// Whether this is a heartbeat response. For xDS versions that support TTL, this
	// will be converted into a response that doesn't contain the actual resource protobuf.
	// This allows for more lightweight updates that server only to update the TTL timer.
	Heartbeat bool

	// Context provided at the time of response creation. This allows associating additional
	// information with a generated response.
	Ctx context.Context

	// marshaledResponse holds an atomic reference to the serialized discovery response.
	marshaledResponse atomic.Pointer[discovery.DiscoveryResponse]
}

// RawDeltaResponse is a pre-serialized xDS response that utilizes the delta discovery request/response objects.
type RawDeltaResponse struct {
	// Request is the latest delta request on the stream.
	DeltaRequest *discovery.DeltaDiscoveryRequest

	// SystemVersionInfo holds the currently applied response system version and should be used for debugging purposes only.
	SystemVersionInfo string

	// resources to be included in the response.
	resources []*snapshot.CachedResource

	// removedResources is a list of resource aliases which should be dropped by the consuming client.
	removedResources []string

	// nextVersionMap consists of updated version mappings after this response is applied.
	nextVersionMap map[string]string

	// Context provided at the time of response creation. This allows associating additional
	// information with a generated response.
	Ctx context.Context

	// Marshaled Resources to be included in the response.
	marshaledResponse atomic.Pointer[discovery.DeltaDiscoveryResponse]
}

var (
	_ types.Response      = &RawResponse{}
	_ types.DeltaResponse = &RawDeltaResponse{}
)

// PassthroughResponse is a pre constructed xDS response that need not go through marshaling transformations.
type PassthroughResponse struct {
	// Request is the original request.
	Request *discovery.DiscoveryRequest

	// The discovery response that needs to be sent as is, without any marshaling transformations.
	DiscoveryResponse *discovery.DiscoveryResponse

	ctx context.Context

	// ReturnedResources tracks the resources returned for the subscription and the version when it was last returned,
	// including previously returned ones when using non-full state resources.
	// It allows the cache to know what the client knows. The server will transparently forward this
	// across requests, and the cache is responsible for its interpretation.
	ReturnedResources map[string]string
}

// DeltaPassthroughResponse is a pre constructed xDS response that need not go through marshaling transformations.
type DeltaPassthroughResponse struct {
	// Request is the latest delta request on the stream
	DeltaRequest *discovery.DeltaDiscoveryRequest

	// NextVersionMap consists of updated version mappings after this response is applied
	NextVersionMap map[string]string

	// This discovery response that needs to be sent as is, without any marshaling transformations
	DeltaDiscoveryResponse *discovery.DeltaDiscoveryResponse

	ctx context.Context
}

var (
	_ types.Response      = &PassthroughResponse{}
	_ types.DeltaResponse = &DeltaPassthroughResponse{}
)

// NewRawResponse allows building a response within non-yet migrated caches.
func NewRawResponse(ctx context.Context, req *discovery.DiscoveryRequest, version string, resources []*snapshot.CachedResource, returnedResources map[string]string, isHeartbeat bool) *RawResponse {
	return &RawResponse{
		Ctx:               ctx,
		Request:           req,
		Version:           version,
		resources:         resources,
		returnedResources: returnedResources,
		Heartbeat:         isHeartbeat,
	}
}

func NewTestRawResponse(req *discovery.DiscoveryRequest, version string, resources []types.ResourceWithTTL) *RawResponse {
	cachedRes := []*snapshot.CachedResource{}
	for _, res := range resources {
		newRes := snapshot.NewCachedResourceWithTTL(types.GetResourceName(res.Resource), res, version)
		cachedRes = append(cachedRes, newRes)
	}
	return &RawResponse{
		Request:   req,
		Version:   version,
		resources: cachedRes,
	}
}

// NewRawDeltaResponse allows building a response within non-yet migrated caches.
func NewRawDeltaResponse(ctx context.Context, req *discovery.DeltaDiscoveryRequest, version string, resources []*snapshot.CachedResource, removedResources []string, nextVersionMap map[string]string) *RawDeltaResponse {
	return &RawDeltaResponse{
		Ctx:               ctx,
		DeltaRequest:      req,
		SystemVersionInfo: version,
		resources:         resources,
		removedResources:  removedResources,
		nextVersionMap:    nextVersionMap,
	}
}

func NewTestRawDeltaResponse(req *discovery.DeltaDiscoveryRequest, version string, resources []types.ResourceWithTTL, removedResources []string, nextVersionMap map[string]string) *RawDeltaResponse {
	cachedRes := []*snapshot.CachedResource{}
	for _, res := range resources {
		name := types.GetResourceName(res.Resource)
		newRes := snapshot.NewCachedResourceWithTTL(name, res, nextVersionMap[name])
		cachedRes = append(cachedRes, newRes)
	}
	return &RawDeltaResponse{
		DeltaRequest:      req,
		SystemVersionInfo: version,
		resources:         cachedRes,
		removedResources:  removedResources,
		nextVersionMap:    nextVersionMap,
	}
}

// GetDiscoveryResponse performs the marshaling the first time its called and uses the cached response subsequently.
// This is necessary because the marshaled response does not change across the calls.
// This caching behavior is important in high throughput scenarios because grpc marshaling has a cost and it drives the cpu utilization under load.
func (r *RawResponse) GetDiscoveryResponse() (*discovery.DiscoveryResponse, error) {
	marshaledResponse := r.marshaledResponse.Load()

	if marshaledResponse != nil {
		return marshaledResponse, nil
	}

	marshaledResources := make([]*anypb.Any, len(r.resources))

	for i, resource := range r.resources {
		marshaledResource, err := r.marshalTTLResource(resource)
		if err != nil {
			return nil, fmt.Errorf("processing %s: %w", resource.Name, err)
		}
		marshaledResources[i] = marshaledResource
	}

	marshaledResponse = &discovery.DiscoveryResponse{
		VersionInfo: r.Version,
		Resources:   marshaledResources,
		TypeUrl:     r.GetRequest().GetTypeUrl(),
	}

	r.marshaledResponse.Store(marshaledResponse)

	return marshaledResponse, nil
}

// GetResources returned the cachedResources used in the response.
func (r *RawResponse) GetResources() []*snapshot.CachedResource {
	return r.resources
}

func (r *RawResponse) GetReturnedResources() map[string]string {
	return r.returnedResources
}

// GetRawResources is used internally within go-control-plane. Its interface and content may change
func (r *RawResponse) GetRawResources() []types.ResourceWithTTL {
	resources := make([]types.ResourceWithTTL, 0, len(r.resources))
	for _, res := range r.resources {
		resources = append(resources, types.ResourceWithTTL{Resource: res.GetResource(), TTL: res.Ttl})
	}
	return resources
}

// GetDeltaDiscoveryResponse performs the marshaling the first time its called and uses the cached response subsequently.
// We can do this because the marshaled response does not change across the calls.
// This caching behavior is important in high throughput scenarios because grpc marshaling has a cost and it drives the cpu utilization under load.
func (r *RawDeltaResponse) GetDeltaDiscoveryResponse() (*discovery.DeltaDiscoveryResponse, error) {
	marshaledResponse := r.marshaledResponse.Load()

	if marshaledResponse != nil {
		return marshaledResponse, nil
	}

	marshaledResources := make([]*discovery.Resource, len(r.resources))

	for i, resource := range r.resources {
		marshaledResource, err := resource.GetMarshaledResource()
		if err != nil {
			return nil, fmt.Errorf("processing %s: %w", resource.Name, err)
		}
		version, err := resource.GetVersion(true)
		if err != nil {
			return nil, fmt.Errorf("processing version of %s: %w", resource.Name, err)
		}
		marshaledResources[i] = &discovery.Resource{
			Name: resource.Name,
			Resource: &anypb.Any{
				TypeUrl: r.GetDeltaRequest().GetTypeUrl(),
				Value:   marshaledResource,
			},
			Version: version,
		}
	}

	marshaledResponse = &discovery.DeltaDiscoveryResponse{
		Resources:         marshaledResources,
		RemovedResources:  r.removedResources,
		TypeUrl:           r.GetDeltaRequest().GetTypeUrl(),
		SystemVersionInfo: r.SystemVersionInfo,
	}
	r.marshaledResponse.Store(marshaledResponse)

	return marshaledResponse, nil
}

// GetRawResources is used internally within go-control-plane. Its interface and content may change
func (r *RawDeltaResponse) GetRawResources() []types.ResourceWithTTL {
	resources := make([]types.ResourceWithTTL, 0, len(r.resources))
	for _, res := range r.resources {
		resources = append(resources, types.ResourceWithTTL{Resource: res.GetResource(), TTL: res.Ttl})
	}
	return resources
}

// GetRequest returns the original Discovery Request.
func (r *RawResponse) GetRequest() *discovery.DiscoveryRequest {
	return r.Request
}

func (r *RawResponse) GetContext() context.Context {
	return r.Ctx
}

// IsEmpty returns whether the response includes any resource.
// Empty responses can be expected by the client in some context.
func (r *RawResponse) IsEmpty() bool {
	return len(r.resources) == 0
}

// GetDeltaRequest returns the original DeltaRequest.
func (r *RawDeltaResponse) GetDeltaRequest() *discovery.DeltaDiscoveryRequest {
	return r.DeltaRequest
}

// GetVersion returns the response version.
// Deprecated: use GetResponseVersion instead
func (r *RawResponse) GetVersion() (string, error) {
	return r.GetResponseVersion(), nil
}

// GetResponseVersion returns the response version.
func (r *RawResponse) GetResponseVersion() string {
	return r.Version
}

// GetSystemVersion returns the raw SystemVersion.
// Deprecated: use GetResponseVersion instead
func (r *RawDeltaResponse) GetSystemVersion() (string, error) {
	return r.GetResponseVersion(), nil
}

// GetResponseVersion returns the response version.
func (r *RawDeltaResponse) GetResponseVersion() string {
	return r.SystemVersionInfo
}

// NextVersionMap returns the version map which consists of updated version mappings after this response is applied.
// Deprecated: use GetReturnedResources instead
func (r *RawDeltaResponse) GetNextVersionMap() map[string]string {
	return r.GetReturnedResources()
}

// GetReturnedResources returns the version map which consists of updated version mappings after this response is applied.
func (r *RawDeltaResponse) GetReturnedResources() map[string]string {
	return r.nextVersionMap
}

func (r *RawDeltaResponse) GetContext() context.Context {
	return r.Ctx
}

// IsEmpty returns whether the response includes any resource.
// Empty responses can be expected by the client in some context.
func (r *RawDeltaResponse) IsEmpty() bool {
	return len(r.resources)+len(r.removedResources) == 0
}

var deltaResourceTypeURL = "type.googleapis.com/" + string(proto.MessageName(&discovery.Resource{}))

func (r *RawResponse) marshalTTLResource(resource *snapshot.CachedResource) (*anypb.Any, error) {
	if resource.Ttl == nil {
		marshaled, err := resource.GetMarshaledResource()
		if err != nil {
			return nil, fmt.Errorf("marshaling: %w", err)
		}
		return &anypb.Any{
			TypeUrl: r.GetRequest().GetTypeUrl(),
			Value:   marshaled,
		}, nil
	}

	wrappedResource := &discovery.Resource{
		Name: resource.Name,
		Ttl:  durationpb.New(*resource.Ttl),
	}

	if !r.Heartbeat {
		marshaled, err := resource.GetMarshaledResource()
		if err != nil {
			return nil, fmt.Errorf("marshaling: %w", err)
		}
		rsrc := new(anypb.Any)
		rsrc.TypeUrl = r.GetRequest().GetTypeUrl()
		rsrc.Value = marshaled

		wrappedResource.Resource = rsrc
	}

	marshaled, err := proto.MarshalOptions{Deterministic: true}.Marshal(wrappedResource)
	if err != nil {
		return nil, fmt.Errorf("marshaling discovery resource: %w", err)
	}

	return &anypb.Any{
		TypeUrl: deltaResourceTypeURL,
		Value:   marshaled,
	}, nil
}

// GetDiscoveryResponse returns the final passthrough Discovery Response.
func (r *PassthroughResponse) GetDiscoveryResponse() (*discovery.DiscoveryResponse, error) {
	return r.DiscoveryResponse, nil
}

func (r *PassthroughResponse) GetReturnedResources() map[string]string {
	return r.ReturnedResources
}

// GetDeltaDiscoveryResponse returns the final passthrough Delta Discovery Response.
func (r *DeltaPassthroughResponse) GetDeltaDiscoveryResponse() (*discovery.DeltaDiscoveryResponse, error) {
	return r.DeltaDiscoveryResponse, nil
}

// GetRequest returns the original Discovery Request.
func (r *PassthroughResponse) GetRequest() *discovery.DiscoveryRequest {
	return r.Request
}

// GetDeltaRequest returns the original Delta Discovery Request.
func (r *DeltaPassthroughResponse) GetDeltaRequest() *discovery.DeltaDiscoveryRequest {
	return r.DeltaRequest
}

// GetVersion returns the response version.
// Deprecated: use GetResponseVersion instead
func (r *PassthroughResponse) GetVersion() (string, error) {
	discoveryResponse, _ := r.GetDiscoveryResponse()
	if discoveryResponse != nil {
		return discoveryResponse.GetVersionInfo(), nil
	}
	return "", fmt.Errorf("DiscoveryResponse is nil")
}

// GetResponseVersion returns the response version, or empty if not set.
func (r *PassthroughResponse) GetResponseVersion() string {
	discoveryResponse, _ := r.GetDiscoveryResponse()
	if discoveryResponse != nil {
		return discoveryResponse.GetVersionInfo()
	}
	return ""
}

func (r *PassthroughResponse) GetContext() context.Context {
	return r.ctx
}

// GetSystemVersion returns the response version.
// Deprecated: use GetResponseVersion instead
func (r *DeltaPassthroughResponse) GetSystemVersion() (string, error) {
	deltaDiscoveryResponse, _ := r.GetDeltaDiscoveryResponse()
	if deltaDiscoveryResponse != nil {
		return deltaDiscoveryResponse.GetSystemVersionInfo(), nil
	}
	return "", fmt.Errorf("DeltaDiscoveryResponse is nil")
}

// GetResponseVersion returns the response version, or empty if not set.
func (r *DeltaPassthroughResponse) GetResponseVersion() string {
	deltaDiscoveryResponse, _ := r.GetDeltaDiscoveryResponse()
	if deltaDiscoveryResponse != nil {
		return deltaDiscoveryResponse.GetSystemVersionInfo()
	}
	return ""
}

// NextVersionMap returns the version map from a DeltaPassthroughResponse.
// Deprecated: use GetReturnedResources instead
func (r *DeltaPassthroughResponse) GetNextVersionMap() map[string]string {
	return r.GetReturnedResources()
}

// GetReturnedResources returns the version map from a DeltaPassthroughResponse.
func (r *DeltaPassthroughResponse) GetReturnedResources() map[string]string {
	return r.NextVersionMap
}

func (r *DeltaPassthroughResponse) GetContext() context.Context {
	return r.ctx
}
