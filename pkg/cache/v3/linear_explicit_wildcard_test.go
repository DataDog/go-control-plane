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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	endpoint "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	"github.com/envoyproxy/go-control-plane/pkg/log"
	"github.com/envoyproxy/go-control-plane/pkg/resource/v3"
	"github.com/envoyproxy/go-control-plane/pkg/server/stream/v3"
)

const (
	node1            = "test-node-1"
	node2            = "test-node-2"
	wildcardTestType = resource.EndpointType
)

func buildTestEndpoint(name string) types.Resource {
	return &endpoint.ClusterLoadAssignment{
		ClusterName: name,
		Endpoints: []*endpoint.LocalityLbEndpoints{{
			Locality: &core.Locality{Zone: "zone-" + name},
			LbEndpoints: []*endpoint.LbEndpoint{{
				HostIdentifier: &endpoint.LbEndpoint_Endpoint{
					Endpoint: &endpoint.Endpoint{
						Address: &core.Address{
							Address: &core.Address_SocketAddress{
								SocketAddress: &core.SocketAddress{
									Protocol: core.SocketAddress_TCP,
									Address:  "127.0.0.1",
									PortSpecifier: &core.SocketAddress_PortValue{
										PortValue: 8080,
									},
								},
							},
						},
					},
				},
			}},
		}},
	}
}

func newExplicitWildcardCache(t *testing.T) *LinearCache {
	t.Helper()
	return NewLinearCache(wildcardTestType,
		WithExplicitWildcardResources(),
		WithLogger(log.NewTestLogger(t)))
}

func setExplicitResources(t *testing.T, c *LinearCache, nodeID string, resources ...string) {
	t.Helper()
	resourceMap := make(map[string]struct{}, len(resources))
	for _, r := range resources {
		resourceMap[r] = struct{}{}
	}
	err := c.UpdateWilcardResourcesForNode(nodeID, resourceMap)
	require.NoError(t, err)
}

func createNodeWildcardWatch(t *testing.T, c *LinearCache, nodeID string, versionInfo string) (stream.Subscription, chan Response) {
	t.Helper()
	req := &Request{ResourceNames: []string{"*"}, TypeUrl: wildcardTestType, VersionInfo: versionInfo}
	sub := stream.NewSotwSubscription(req.GetResourceNames(), nodeID)
	w := make(chan Response, 1)
	_, err := c.CreateWatch(req, sub, w)
	require.NoError(t, err)
	return sub, w
}

func createNodeWildcardDeltaWatch(t *testing.T, c *LinearCache, nodeID string, versionInfo string) (stream.Subscription, chan DeltaResponse) {
	t.Helper()
	req := &DeltaRequest{
		TypeUrl:                wildcardTestType,
		ResourceNamesSubscribe: []string{"*"},
		ResponseNonce:          versionInfo,
	}
	sub := stream.NewDeltaSubscription(
		req.GetResourceNamesSubscribe(), req.GetResourceNamesUnsubscribe(), req.GetInitialResourceVersions(), nodeID,
	)
	w := make(chan DeltaResponse, 1)
	_, err := c.CreateDeltaWatch(req, sub, w)
	require.NoError(t, err)
	return sub, w
}

type protocolTest struct {
	name                 string
	createWatch          func(t *testing.T, c *LinearCache, nodeID, versionInfo string) (stream.Subscription, any)
	verifyEmpty          func(t *testing.T, respCh any, expectedVersion string) any
	verifyResources      func(t *testing.T, respCh any, expectedVersion string, resources ...string) any
	mustBlock            func(t *testing.T, respCh any)
	verifyAdded          func(t *testing.T, respCh any, resources ...string) any
	getVersion           func(resp any) string
	updateVersion        func(resp any, sub *stream.Subscription, version string)
	createAdditionalWatch func(t *testing.T, c *LinearCache, sub stream.Subscription, version string) any
}

func getProtocolTests() []protocolTest {
	return []protocolTest{
		{
			name: "SOTW",
			createWatch: func(t *testing.T, c *LinearCache, nodeID, versionInfo string) (stream.Subscription, any) {
				sub, w := createNodeWildcardWatch(t, c, nodeID, versionInfo)
				return sub, w
			},
			verifyEmpty: func(t *testing.T, respCh any, _ string) any {
				w, ok := respCh.(chan Response)
				require.True(t, ok, "failed to cast to SOTW response channel")
				return verifyResponseResources(t, w, wildcardTestType, "")
			},
			verifyResources: func(t *testing.T, respCh any, _ string, resources ...string) any {
				w, ok := respCh.(chan Response)
				require.True(t, ok, "failed to cast to SOTW response channel")
				return verifyResponseResources(t, w, wildcardTestType, "", resources...)
			},
			mustBlock: func(t *testing.T, respCh any) {
				w, ok := respCh.(chan Response)
				require.True(t, ok, "failed to cast to SOTW response channel")
				mustBlock(t, w)
			},
			verifyAdded: func(t *testing.T, respCh any, resources ...string) any {
				w, ok := respCh.(chan Response)
				require.True(t, ok, "failed to cast to SOTW response channel")
				resp := <-w
				returnedResources := resp.GetReturnedResources()
				for _, r := range resources {
					assert.Contains(t, returnedResources, r, "resource %s should have been added", r)
				}
				return resp
			},
			// Returns the cache version from the response (e.g., "3", "4").
			// The cache version increments with each cache mutation.
			getVersion: func(resp any) string {
				sotwResp, ok := resp.(Response)
				if !ok {
					panic("failed to cast to SOTW response")
				}
				version, _ := sotwResp.GetVersion()
				return version
			},
			updateVersion: func(resp any, sub *stream.Subscription, version string) {
				sotwResp, ok := resp.(Response)
				if !ok {
					panic("failed to cast to SOTW response")
				}
				req := &Request{TypeUrl: wildcardTestType, VersionInfo: version}
				updateFromSotwResponse(sotwResp, sub, req)
			},
			createAdditionalWatch: func(t *testing.T, c *LinearCache, sub stream.Subscription, version string) any {
				req := &Request{
					ResourceNames: []string{"*"},
					TypeUrl:       wildcardTestType,
					VersionInfo:   version,
				}
				w := make(chan Response, 1)
				_, err := c.CreateWatch(req, sub, w)
				require.NoError(t, err)
				return w
			},
		},
		{
			name: "Delta",
			createWatch: func(t *testing.T, c *LinearCache, nodeID, versionInfo string) (stream.Subscription, any) {
				sub, w := createNodeWildcardDeltaWatch(t, c, nodeID, versionInfo)
				return sub, w
			},
			verifyEmpty: func(t *testing.T, respCh any, _ string) any {
				t.Helper()
				w, ok := respCh.(chan DeltaResponse)
				require.True(t, ok, "failed to cast to Delta response channel")
				resp := <-w
				assert.Empty(t, resp.GetNextVersionMap())
				return resp
			},
			verifyResources: func(t *testing.T, respCh any, _ string, resources ...string) any {
				t.Helper()
				w, ok := respCh.(chan DeltaResponse)
				require.True(t, ok, "failed to cast to Delta response channel")
				resp := <-w
				versionMap := resp.GetNextVersionMap()
				assert.Len(t, versionMap, len(resources))
				for _, r := range resources {
					assert.Contains(t, versionMap, r)
				}
				return resp
			},
			mustBlock: func(t *testing.T, respCh any) {
				t.Helper()
				w, ok := respCh.(chan DeltaResponse)
				require.True(t, ok, "failed to cast to Delta response channel")
				mustBlockDelta(t, w)
			},
			verifyAdded: func(t *testing.T, respCh any, resources ...string) any {
				w, ok := respCh.(chan DeltaResponse)
				require.True(t, ok, "failed to cast to Delta response channel")
				resp := <-w
				versionMap := resp.GetNextVersionMap()
				for _, r := range resources {
					assert.Contains(t, versionMap, r, "resource %s should have been added", r)
				}
				return resp
			},
			// Any non-empty value works - the cache only checks
			// if it's empty (first request) or non-empty (subsequent request).
			getVersion: func(resp any) string {
				return "non-empty"
			},
			updateVersion: func(resp any, sub *stream.Subscription, version string) {
				deltaResp, ok := resp.(DeltaResponse)
				if !ok {
					panic("failed to cast to Delta response")
				}
				sub.SetReturnedResources(deltaResp.GetNextVersionMap())
			},
			createAdditionalWatch: func(t *testing.T, c *LinearCache, sub stream.Subscription, version string) any {
				req := &DeltaRequest{
					TypeUrl:       wildcardTestType,
					ResponseNonce: version,
				}
				w := make(chan DeltaResponse, 1)
				_, err := c.CreateDeltaWatch(req, sub, w)
				require.NoError(t, err)
				return w
			},
		},
	}
}

func TestLinearExplicitWildcardBothProtocols(t *testing.T) {
	protocols := getProtocolTests()

	for _, protocol := range protocols {
		t.Run(protocol.name, func(t *testing.T) {
			t.Run("handles empty resource list", func(t *testing.T) {
				c := newExplicitWildcardCache(t)

				require.NoError(t, c.UpdateResource("a", buildTestEndpoint("a")))
				setExplicitResources(t, c, node1)

				_, w := protocol.createWatch(t, c, node1, "")
				protocol.verifyEmpty(t, w, "")
			})

			t.Run("multiple nodes with different explicit resources", func(t *testing.T) {
				c := newExplicitWildcardCache(t)

				require.NoError(t, c.UpdateResource("a", buildTestEndpoint("a")))
				require.NoError(t, c.UpdateResource("b", buildTestEndpoint("b")))
				require.NoError(t, c.UpdateResource("c", buildTestEndpoint("c")))

				setExplicitResources(t, c, node1, "a", "b")
				setExplicitResources(t, c, node2, "b", "c")

				_, w1 := protocol.createWatch(t, c, node1, "")
				_, w2 := protocol.createWatch(t, c, node2, "")

				protocol.verifyResources(t, w1, "", "a", "b")
				protocol.verifyResources(t, w2, "", "b", "c")
			})

			t.Run("explicit resource not in cache", func(t *testing.T) {
				c := newExplicitWildcardCache(t)

				require.NoError(t, c.UpdateResource("a", buildTestEndpoint("a")))
				require.NoError(t, c.UpdateResource("b", buildTestEndpoint("b")))

				setExplicitResources(t, c, node1, "a", "b", "nonexistent")

				_, w := protocol.createWatch(t, c, node1, "")
				protocol.verifyResources(t, w, "", "a", "b")
			})

			t.Run("don't receive entire cache with explicit resources", func(t *testing.T) {
				c := newExplicitWildcardCache(t)

				require.NoError(t, c.UpdateResource("a", buildTestEndpoint("a")))
				require.NoError(t, c.UpdateResource("b", buildTestEndpoint("b")))
				require.NoError(t, c.UpdateResource("c", buildTestEndpoint("c")))

				setExplicitResources(t, c, node1, "a", "b")

				sub, w := protocol.createWatch(t, c, node1, "")
				resp := protocol.verifyResources(t, w, "", "a", "b")
				version := protocol.getVersion(resp)
				protocol.updateVersion(resp, &sub, version)

				checkTotalWatchCount(t, c, 0)
			})

			t.Run("explicit resources updated while watch active", func(t *testing.T) {
				c := newExplicitWildcardCache(t)

				require.NoError(t, c.UpdateResource("a", buildTestEndpoint("a")))
				require.NoError(t, c.UpdateResource("b", buildTestEndpoint("b")))
				require.NoError(t, c.UpdateResource("c", buildTestEndpoint("c")))

				setExplicitResources(t, c, node1, "a", "b")

				sub, w := protocol.createWatch(t, c, node1, "")
				resp := protocol.verifyResources(t, w, "", "a", "b")
				version := protocol.getVersion(resp)
				protocol.updateVersion(resp, &sub, version)

				w2 := protocol.createAdditionalWatch(t, c, sub, version)
				protocol.mustBlock(t, w2)

				setExplicitResources(t, c, node1, "b", "c")
				protocol.verifyAdded(t, w2, "c")
			})

			t.Run("resource versions tracked correctly on update", func(t *testing.T) {
				c := newExplicitWildcardCache(t)

				require.NoError(t, c.UpdateResource("a", buildTestEndpoint("a")))
				require.NoError(t, c.UpdateResource("b", buildTestEndpoint("b")))

				setExplicitResources(t, c, node1, "a", "b")

				sub, w := protocol.createWatch(t, c, node1, "")
				resp := protocol.verifyResources(t, w, "", "a", "b")
				version := protocol.getVersion(resp)
				protocol.updateVersion(resp, &sub, version)

				w2 := protocol.createAdditionalWatch(t, c, sub, version)
				protocol.mustBlock(t, w2)

				require.NoError(t, c.UpdateResource("a", buildTestEndpoint("a-updated")))
				protocol.verifyAdded(t, w2, "a")
			})

			t.Run("wildcard and non-wildcard watches coexist", func(t *testing.T) {
				c := newExplicitWildcardCache(t)

				require.NoError(t, c.UpdateResource("a", buildTestEndpoint("a")))
				require.NoError(t, c.UpdateResource("b", buildTestEndpoint("b")))

				setExplicitResources(t, c, node1, "a")

				sub, wWildcard := protocol.createWatch(t, c, node1, "")
				respWildcard := protocol.verifyResources(t, wWildcard, "", "a")
				version := protocol.getVersion(respWildcard)
				protocol.updateVersion(respWildcard, &sub, version)

				reqNonWildcard := &Request{
					ResourceNames: []string{"a", "b"},
					TypeUrl:       wildcardTestType,
					VersionInfo:   "",
				}
				wNonWildcard := make(chan Response, 1)
				subNonWildcard := stream.NewSotwSubscription(reqNonWildcard.GetResourceNames(), node1)
				_, err := c.CreateWatch(reqNonWildcard, subNonWildcard, wNonWildcard)
				require.NoError(t, err)
				verifyResponseResources(t, wNonWildcard, wildcardTestType, "", "a", "b")
			})

			t.Run("explicit wildcard only affects wildcard subscriptions", func(t *testing.T) {
				c := newExplicitWildcardCache(t)

				require.NoError(t, c.UpdateResource("a", buildTestEndpoint("a")))
				require.NoError(t, c.UpdateResource("b", buildTestEndpoint("b")))

				setExplicitResources(t, c, node1, "a")

				req := &Request{
					ResourceNames: []string{"b"},
					TypeUrl:       wildcardTestType,
					VersionInfo:   "",
				}
				w := make(chan Response, 1)
				sub := stream.NewSotwSubscription(req.GetResourceNames(), node1)
				_, err := c.CreateWatch(req, sub, w)
				require.NoError(t, err)

				verifyResponseResources(t, w, wildcardTestType, "", "b")
			})

			t.Run("cache update triggers correct watches", func(t *testing.T) {
				c := newExplicitWildcardCache(t)

				require.NoError(t, c.UpdateResource("a", buildTestEndpoint("a")))
				require.NoError(t, c.UpdateResource("b", buildTestEndpoint("b")))

				setExplicitResources(t, c, node1, "a")
				setExplicitResources(t, c, node2, "b")

				sub1, w1 := protocol.createWatch(t, c, node1, "")
				resp1 := protocol.verifyResources(t, w1, "", "a")
				version1 := protocol.getVersion(resp1)
				protocol.updateVersion(resp1, &sub1, version1)

				sub2, w2 := protocol.createWatch(t, c, node2, "")
				resp2 := protocol.verifyResources(t, w2, "", "b")
				version2 := protocol.getVersion(resp2)
				protocol.updateVersion(resp2, &sub2, version2)

				w1Next := protocol.createAdditionalWatch(t, c, sub1, version1)
				protocol.mustBlock(t, w1Next)

				w2Next := protocol.createAdditionalWatch(t, c, sub2, version2)
				protocol.mustBlock(t, w2Next)

				require.NoError(t, c.UpdateResource("a", buildTestEndpoint("a-updated")))

				protocol.verifyAdded(t, w1Next, "a")
				protocol.mustBlock(t, w2Next)
			})

			t.Run("triggers pending wildcard watches", func(t *testing.T) {
				c := newExplicitWildcardCache(t)

				require.NoError(t, c.UpdateResource("a", buildTestEndpoint("a")))
				require.NoError(t, c.UpdateResource("b", buildTestEndpoint("b")))

				sub, w := protocol.createWatch(t, c, node1, "")
				resp := protocol.verifyEmpty(t, w, "")
				version := protocol.getVersion(resp)
				protocol.updateVersion(resp, &sub, version)

				w2 := protocol.createAdditionalWatch(t, c, sub, version)
				protocol.mustBlock(t, w2)

				setExplicitResources(t, c, node1, "a", "b")
				protocol.verifyResources(t, w2, "", "a", "b")
			})

			t.Run("only notifies watches for specified node", func(t *testing.T) {
				c := newExplicitWildcardCache(t)

				require.NoError(t, c.UpdateResource("a", buildTestEndpoint("a")))

				sub1, w1 := protocol.createWatch(t, c, node1, "")
				resp1 := protocol.verifyEmpty(t, w1, "")
				version1 := protocol.getVersion(resp1)
				protocol.updateVersion(resp1, &sub1, version1)

				sub2, w2 := protocol.createWatch(t, c, node2, "")
				resp2 := protocol.verifyEmpty(t, w2, "")
				version2 := protocol.getVersion(resp2)
				protocol.updateVersion(resp2, &sub2, version2)

				w1Next := protocol.createAdditionalWatch(t, c, sub1, version1)
				protocol.mustBlock(t, w1Next)

				w2Next := protocol.createAdditionalWatch(t, c, sub2, version2)
				protocol.mustBlock(t, w2Next)

				setExplicitResources(t, c, node1, "a")
				protocol.verifyResources(t, w1Next, "", "a")
				protocol.mustBlock(t, w2Next)
			})
		})
	}
}

func TestLinearExplicitWildcardDeltaOnly(t *testing.T) {
	t.Run("resource removal from explicit list", func(t *testing.T) {
		// Note: SOTW doesn't properly handle resource removal from explicit list yet
		c := newExplicitWildcardCache(t)

		require.NoError(t, c.UpdateResource("a", buildTestEndpoint("a")))
		require.NoError(t, c.UpdateResource("b", buildTestEndpoint("b")))
		require.NoError(t, c.UpdateResource("c", buildTestEndpoint("c")))

		setExplicitResources(t, c, node1, "a", "b", "c")

		sub, w := createNodeWildcardDeltaWatch(t, c, node1, "")
		resp := <-w
		sub.SetReturnedResources(resp.GetNextVersionMap())
		assert.Len(t, resp.GetNextVersionMap(), 3)

		w2 := make(chan DeltaResponse, 1)
		req := &DeltaRequest{
			TypeUrl:       wildcardTestType,
			ResponseNonce: "1",
		}
		_, err := c.CreateDeltaWatch(req, sub, w2)
		require.NoError(t, err)
		mustBlockDelta(t, w2)

		setExplicitResources(t, c, node1, "a", "b")
		resp2 := <-w2
		deltaResp, err := resp2.GetDeltaDiscoveryResponse()
		require.NoError(t, err)
		removedResources := deltaResp.GetRemovedResources()
		assert.Contains(t, removedResources, "c", "resource c should have been removed")
	})

	t.Run("update explicit list triggers watch with add and remove", func(t *testing.T) {
		// Note: SOTW doesn't properly handle simultaneous add/remove from explicit list yet
		c := newExplicitWildcardCache(t)

		require.NoError(t, c.UpdateResource("a", buildTestEndpoint("a")))
		require.NoError(t, c.UpdateResource("b", buildTestEndpoint("b")))
		require.NoError(t, c.UpdateResource("c", buildTestEndpoint("c")))

		setExplicitResources(t, c, node1, "a", "b")

		sub, w := createNodeWildcardDeltaWatch(t, c, node1, "")
		resp := <-w
		sub.SetReturnedResources(resp.GetNextVersionMap())
		assert.Len(t, resp.GetNextVersionMap(), 2)

		w2 := make(chan DeltaResponse, 1)
		req := &DeltaRequest{
			TypeUrl:       wildcardTestType,
			ResponseNonce: "1",
		}
		_, err := c.CreateDeltaWatch(req, sub, w2)
		require.NoError(t, err)
		mustBlockDelta(t, w2)

		setExplicitResources(t, c, node1, "a", "c")
		resp2 := <-w2

		versionMap := resp2.GetNextVersionMap()
		assert.Contains(t, versionMap, "c", "resource c should have been added")

		deltaResp, err := resp2.GetDeltaDiscoveryResponse()
		require.NoError(t, err)
		removedResources := deltaResp.GetRemovedResources()
		assert.Contains(t, removedResources, "b", "resource b should have been removed")
	})
}

func TestLinearExplicitWildcardErrors(t *testing.T) {
	t.Run("UpdateWilcardResourcesForNode without mode enabled", func(t *testing.T) {
		c := NewLinearCache(wildcardTestType, WithLogger(log.NewTestLogger(t)))

		resourceMap := map[string]struct{}{"a": {}}
		err := c.UpdateWilcardResourcesForNode(node1, resourceMap)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "WithExplicitWildcardResources")
	})
}
