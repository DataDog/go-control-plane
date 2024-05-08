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

package server_test

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
	"testing/iotest"

	"google.golang.org/protobuf/types/known/anypb"

	discovery "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	"github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"github.com/envoyproxy/go-control-plane/pkg/resource/v3"
	"github.com/envoyproxy/go-control-plane/pkg/server/v3"
)

func marshalAnyResources(resources ...types.Resource) []*anypb.Any {
	marshaled := []*anypb.Any{}
	for _, res := range resources {
		m, _ := anypb.New(res)
		marshaled = append(marshaled, m)
	}
	return marshaled
}

func marshalDiscoveryResources(resources []types.ResourceWithTTL) []*discovery.Resource {
	marshaled := []*discovery.Resource{}
	for _, res := range resources {
		m, _ := anypb.New(res.Resource)
		marshaled = append(marshaled, &discovery.Resource{
			Name:     cache.GetResourceName(res.Resource),
			Resource: m,
			Version:  cache.HashResource(m.Value),
		})
	}
	return marshaled
}

func TestGateway(t *testing.T) {
	config := makeMockConfigWatcher()
	config.responses = map[string][]cache.Response{
		resource.ClusterType: {
			&cache.PassthroughResponse{
				Request: &discovery.DiscoveryRequest{TypeUrl: resource.ClusterType},
				DiscoveryResponse: &discovery.DiscoveryResponse{
					VersionInfo: "2",
					TypeUrl:     resource.ClusterType,
					Resources:   marshalAnyResources(cluster),
				},
			},
		},
		resource.RouteType: {
			&cache.PassthroughResponse{
				Request: &discovery.DiscoveryRequest{TypeUrl: resource.RouteType},
				DiscoveryResponse: &discovery.DiscoveryResponse{
					VersionInfo: "3",
					TypeUrl:     resource.RouteType,
					Resources:   marshalAnyResources(route),
				},
			},
		},
		resource.ListenerType: {
			&cache.PassthroughResponse{
				Request: &discovery.DiscoveryRequest{TypeUrl: resource.ListenerType},
				DiscoveryResponse: &discovery.DiscoveryResponse{
					VersionInfo: "4",
					TypeUrl:     resource.RouteType,
					Resources:   marshalAnyResources(httpListener, httpScopedListener),
				},
			},
		},
	}
	gtw := server.HTTPGateway{Server: server.NewServer(context.Background(), config, nil)}

	failCases := []struct {
		path   string
		body   io.Reader
		expect int
	}{
		{
			path:   "/hello/",
			expect: http.StatusNotFound,
		},
		{
			path:   resource.FetchEndpoints,
			expect: http.StatusBadRequest,
		},
		{
			path:   resource.FetchEndpoints,
			body:   iotest.TimeoutReader(strings.NewReader("hello")),
			expect: http.StatusBadRequest,
		},
		{
			path:   resource.FetchEndpoints,
			body:   strings.NewReader("hello"),
			expect: http.StatusBadRequest,
		},
		{
			// missing response
			path:   resource.FetchEndpoints,
			body:   strings.NewReader("{\"node\": {\"id\": \"test\"}}"),
			expect: http.StatusInternalServerError,
		},
	}
	for _, cs := range failCases {
		req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, cs.path, cs.body)
		if err != nil {
			t.Fatal(err)
		}
		resp, code, err := gtw.ServeHTTP(req)
		if err == nil {
			t.Errorf("ServeHTTP succeeded, but should have failed")
		}
		if resp != nil {
			t.Errorf("handler returned wrong response")
		}
		if status := code; status != cs.expect {
			t.Errorf("handler returned wrong status: %d, want %d", status, cs.expect)
		}
	}

	for _, path := range []string{resource.FetchClusters, resource.FetchRoutes, resource.FetchListeners} {
		req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, path, strings.NewReader("{\"node\": {\"id\": \"test\"}}"))
		if err != nil {
			t.Fatal(err)
		}
		resp, code, err := gtw.ServeHTTP(req)
		if err != nil {
			t.Fatal(err)
		}
		if resp == nil {
			t.Errorf("handler returned wrong response")
		}
		if status := code; status != 200 {
			t.Errorf("handler returned wrong status: %d, want %d", status, 200)
		}
	}
}
