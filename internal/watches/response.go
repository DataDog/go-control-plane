package watches

import (
	"context"
	"encoding/hex"
	"hash/fnv"
	"sort"

	"github.com/envoyproxy/go-control-plane/internal/snapshot"
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
)

// ResponseWatch is a watch record keeping both the request and an open channel for the response.
type ResponseWatch struct {
	// Request is the original request for the watch.
	Request *types.Request

	// Response is the channel to push responses to.
	Response chan types.Response

	// Subscription stores the current client subscription state.
	Subscription types.Subscription

	// enableResourceVersion indicates whether versions returned in the response are built using stable versions instead of cache update versions.
	enableResourceVersion bool

	// fullStateResponses requires that all resources matching the request, with no regards to which ones actually updated, must be provided in the response.
	fullStateResponses bool

	// versionPrefix is prepended to computed versions when resource versions are used.
	versionPrefix string
}

// NewResponseWatch builds a new watch to build responses for.
func NewResponseWatch(req *types.Request, resp chan types.Response, sub types.Subscription, useResourceVersion bool, sendFullStateResponses bool, versionPrefix string) ResponseWatch {
	return ResponseWatch{
		Request:               req,
		Response:              resp,
		Subscription:          sub,
		enableResourceVersion: useResourceVersion,
		fullStateResponses:    sendFullStateResponses,
		versionPrefix:         versionPrefix,
	}
}

func (w ResponseWatch) IsDelta() bool {
	return false
}

func (w ResponseWatch) UseResourceVersion() bool {
	return w.enableResourceVersion
}

func (w ResponseWatch) SendFullStateResponses() bool {
	return w.fullStateResponses
}

func (w ResponseWatch) GetSubscription() types.Subscription {
	return w.Subscription
}

func (w ResponseWatch) SendResponse(resp snapshot.WatchResponse) {
	w.Response <- resp.(*RawResponse)
}

func (w ResponseWatch) BuildResponse(updatedResources []*snapshot.CachedResource, _ []string, returnedVersions map[string]string, version string) snapshot.WatchResponse {
	responseVersion := version
	if w.enableResourceVersion {
		responseVersion = w.versionPrefix + computeSotwResourceVersion(returnedVersions)
	}

	return &RawResponse{
		Request:           w.Request,
		resources:         updatedResources,
		returnedResources: returnedVersions,
		Version:           responseVersion,
		Ctx:               context.Background(),
	}
}

// DeltaResponseWatch is a watch record keeping both the delta request and an open channel for the delta response.
type DeltaResponseWatch struct {
	// Request is the most recent delta request for the watch
	Request *types.DeltaRequest

	// Response is the channel to push the delta responses to
	Response chan types.DeltaResponse

	// Subscription stores the current client subscription state.
	Subscription types.Subscription
}

// NewDeltaResponseWatch builds a new watch to build responses for.
func NewDeltaResponseWatch(req *types.DeltaRequest, resp chan types.DeltaResponse, sub types.Subscription) DeltaResponseWatch {
	return DeltaResponseWatch{
		Request:      req,
		Response:     resp,
		Subscription: sub,
	}
}

func (w DeltaResponseWatch) IsDelta() bool {
	return true
}

func (w DeltaResponseWatch) UseResourceVersion() bool {
	return true
}

func (w DeltaResponseWatch) SendFullStateResponses() bool {
	return false
}

func (w DeltaResponseWatch) GetSubscription() types.Subscription {
	return w.Subscription
}

func (w DeltaResponseWatch) SendResponse(resp snapshot.WatchResponse) {
	w.Response <- resp.(*RawDeltaResponse)
}

func (w DeltaResponseWatch) BuildResponse(updatedResources []*snapshot.CachedResource, removedResources []string, returnedVersions map[string]string, version string) snapshot.WatchResponse {
	return &RawDeltaResponse{
		DeltaRequest:      w.Request,
		resources:         updatedResources,
		removedResources:  removedResources,
		nextVersionMap:    returnedVersions,
		SystemVersionInfo: version,
		Ctx:               context.Background(),
	}
}

func computeSotwResourceVersion(versionMap map[string]string) string {
	// To enforce a stable hash we need to have an ordered vision of the map.
	keys := make([]string, 0, len(versionMap))
	for key := range versionMap {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	mapHasher := fnv.New64()

	itemHasher := fnv.New64()
	for _, key := range keys {
		itemHasher.Reset()
		itemHasher.Write([]byte(key))
		mapHasher.Write(itemHasher.Sum(nil))
		itemHasher.Reset()
		itemHasher.Write([]byte(versionMap[key]))
		mapHasher.Write(itemHasher.Sum(nil))
	}
	return hex.EncodeToString(mapHasher.Sum(nil))
}
