// Copyright 2018 Envoyproxy Authors
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

// Package cache defines a configuration cache for the server.
package cache

import (
	"context"

	"github.com/envoyproxy/go-control-plane/internal/watches"
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
)

// ConfigWatcher requests watches for configuration resources by a node, last
// applied version identifier, and resource names hint. The watch should send
// the responses when they are ready. The watch can be canceled by the
// consumer, in effect terminating the watch for the request.
// ConfigWatcher implementation must be thread-safe.
type ConfigWatcher interface {
	// CreateWatch returns a new open watch from a non-empty request.
	// This is the entrypoint to propagate configuration changes the
	// provided Response channel. State from the gRPC server is utilized
	// to make sure consuming cache implementations can see what the server has sent to clients.
	//
	// An individual consumer normally issues a single open watch by each type URL.
	//
	// The provided channel produces requested resources as responses, once they are available.
	//
	// Cancel is an optional function to release resources in the producer. If
	// provided, the consumer may call this function multiple times.
	CreateWatch(*Request, Subscription, chan Response) (cancel func(), err error)

	// CreateDeltaWatch returns a new open incremental xDS watch.
	// This is the entrypoint to propagate configuration changes the
	// provided DeltaResponse channel. State from the gRPC server is utilized
	// to make sure consuming cache implementations can see what the server has sent to clients.
	//
	// The provided channel produces requested resources as responses, or spontaneous updates in accordance
	// with the incremental xDS specification.
	//
	// Cancel is an optional function to release resources in the producer. If
	// provided, the consumer may call this function multiple times.
	CreateDeltaWatch(*DeltaRequest, Subscription, chan DeltaResponse) (cancel func(), err error)
}

// ConfigFetcher fetches configuration resources from cache
type ConfigFetcher interface {
	// Fetch implements the polling method of the config cache using a non-empty request.
	Fetch(context.Context, *Request) (Response, error)
}

// Cache is a generic config cache with a watcher.
type Cache interface {
	ConfigWatcher
	ConfigFetcher
}

// For backward compatibility, reimport those interfaces within the Cache package
type Request = types.Request
type DeltaRequest = types.DeltaRequest
type Response = types.Response
type DeltaResponse = types.DeltaResponse
type Subscription = types.Subscription

type RawResponse = watches.RawResponse
type RawDeltaResponse = watches.RawDeltaResponse
type PassthroughResponse = watches.PassthroughResponse
type DeltaPassthroughResponse = watches.DeltaPassthroughResponse
