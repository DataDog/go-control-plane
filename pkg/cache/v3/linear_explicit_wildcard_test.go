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
	node3            = "test-node-3"
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

func createNodeWildcardWatch(t *testing.T, c *LinearCache, nodeID string) (<-chan Response, stream.Subscription) {
	t.Helper()
	req := &Request{ResourceNames: []string{"*"}, TypeUrl: wildcardTestType}
	sub := stream.NewSotwSubscription(nil, nodeID)
	w := make(chan Response, 1)
	_, err := c.CreateWatch(req, sub, w)
	require.NoError(t, err)
	return w, sub
}

//nolint:unparam
func createNodeWildcardDeltaWatch(t *testing.T, c *LinearCache, nodeID string) (<-chan DeltaResponse, stream.Subscription) {
	t.Helper()
	req := &DeltaRequest{
		TypeUrl:                wildcardTestType,
		ResourceNamesSubscribe: []string{"*"},
	}
	sub := stream.NewDeltaSubscription(
		nil, nil, nil, nodeID,
	)
	w := make(chan DeltaResponse, 1)
	_, err := c.CreateDeltaWatch(req, sub, w)
	require.NoError(t, err)
	return w, sub
}

type protocolTest struct {
	name            string
	createWatch     func(t *testing.T, c *LinearCache, nodeID, versionInfo string) (stream.Subscription, any)
	verifyEmpty     func(t *testing.T, respCh any, expectedVersion string)
	verifyResources func(t *testing.T, respCh any, expectedVersion string, resources ...string)
	mustBlock       func(t *testing.T, respCh any)
}

func getProtocolTests() []protocolTest {
	return []protocolTest{
		{
			name: "SOTW",
			createWatch: func(t *testing.T, c *LinearCache, nodeID, _ string) (stream.Subscription, any) {
				t.Helper()
				req := &Request{ResourceNames: []string{"*"}, TypeUrl: wildcardTestType, VersionInfo: ""}
				sub := stream.NewSotwSubscription(req.GetResourceNames(), nodeID)
				w := make(chan Response, 1)
				_, err := c.CreateWatch(req, sub, w)
				require.NoError(t, err)
				return sub, w
			},
			verifyEmpty: func(t *testing.T, respCh any, _ string) {
				t.Helper()
				w, ok := respCh.(chan Response)
				require.True(t, ok, "failed to cast to SOTW response channel")
				verifyResponseResources(t, w, wildcardTestType, "")
			},
			verifyResources: func(t *testing.T, respCh any, _ string, resources ...string) {
				t.Helper()
				w, ok := respCh.(chan Response)
				require.True(t, ok, "failed to cast to SOTW response channel")
				verifyResponseResources(t, w, wildcardTestType, "", resources...)
			},
			mustBlock: func(t *testing.T, respCh any) {
				t.Helper()
				w, ok := respCh.(chan Response)
				require.True(t, ok, "failed to cast to SOTW response channel")
				mustBlock(t, w)
			},
		},
		{
			name: "Delta",
			createWatch: func(t *testing.T, c *LinearCache, nodeID, _ string) (stream.Subscription, any) {
				t.Helper()
				req := &DeltaRequest{TypeUrl: wildcardTestType, ResourceNamesSubscribe: []string{"*"}}
				sub := stream.NewDeltaSubscription(req.GetResourceNamesSubscribe(), req.GetResourceNamesUnsubscribe(), req.GetInitialResourceVersions(), nodeID)
				w := make(chan DeltaResponse, 1)
				_, err := c.CreateDeltaWatch(req, sub, w)
				require.NoError(t, err)
				return sub, w
			},
			verifyEmpty: func(t *testing.T, respCh any, _ string) {
				t.Helper()
				w, ok := respCh.(chan DeltaResponse)
				require.True(t, ok, "failed to cast to Delta response channel")
				resp := <-w
				assert.Empty(t, resp.GetNextVersionMap())
			},
			verifyResources: func(t *testing.T, respCh any, _ string, resources ...string) {
				t.Helper()
				w, ok := respCh.(chan DeltaResponse)
				require.True(t, ok, "failed to cast to Delta response channel")
				resp := <-w
				versionMap := resp.GetNextVersionMap()
				assert.Len(t, versionMap, len(resources))
				for _, r := range resources {
					assert.Contains(t, versionMap, r)
				}
			},
			mustBlock: func(t *testing.T, respCh any) {
				t.Helper()
				w, ok := respCh.(chan DeltaResponse)
				require.True(t, ok, "failed to cast to Delta response channel")
				mustBlockDelta(t, w)
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
		})
	}
}

func TestLinearUpdateWildcardResourcesForNode(t *testing.T) {
	t.Run("triggers pending wildcard watches", func(t *testing.T) {
		c := newExplicitWildcardCache(t)

		require.NoError(t, c.UpdateResource("a", buildTestEndpoint("a")))
		require.NoError(t, c.UpdateResource("b", buildTestEndpoint("b")))

		req := &Request{
			ResourceNames: []string{"*"},
			TypeUrl:       wildcardTestType,
			VersionInfo:   "",
		}
		sub := stream.NewSotwSubscription(req.GetResourceNames(), node1)
		w := make(chan Response, 1)
		_, err := c.CreateWatch(req, sub, w)
		require.NoError(t, err)

		resp, r := verifyResponseContent(t, w, wildcardTestType, "2")
		assert.Empty(t, r.GetResources())
		updateFromSotwResponse(resp, &sub, req)

		w2 := make(chan Response, 1)
		_, err = c.CreateWatch(req, sub, w2)
		require.NoError(t, err)
		mustBlock(t, w2)

		setExplicitResources(t, c, node1, "a", "b")
		verifyResponseResources(t, w2, wildcardTestType, "3", "a", "b")
	})

	t.Run("only notifies watches for specified node", func(t *testing.T) {
		c := newExplicitWildcardCache(t)

		require.NoError(t, c.UpdateResource("a", buildTestEndpoint("a")))

		req1 := &Request{
			ResourceNames: []string{"*"},
			TypeUrl:       wildcardTestType,
			VersionInfo:   "",
		}
		req2 := &Request{
			ResourceNames: []string{"*"},
			TypeUrl:       wildcardTestType,
			VersionInfo:   "",
		}
		sub1 := stream.NewSotwSubscription(req1.GetResourceNames(), node1)
		sub2 := stream.NewSotwSubscription(req2.GetResourceNames(), node2)
		w1 := make(chan Response, 1)
		w2 := make(chan Response, 1)

		_, err := c.CreateWatch(req1, sub1, w1)
		require.NoError(t, err)
		_, err = c.CreateWatch(req2, sub2, w2)
		require.NoError(t, err)

		resp1, r1 := verifyResponseContent(t, w1, wildcardTestType, "1")
		assert.Empty(t, r1.GetResources())
		resp2, r2 := verifyResponseContent(t, w2, wildcardTestType, "1")
		assert.Empty(t, r2.GetResources())
		updateFromSotwResponse(resp1, &sub1, req1)
		updateFromSotwResponse(resp2, &sub2, req2)

		req1.VersionInfo = "1"
		req2.VersionInfo = "1"

		w1Next := make(chan Response, 1)
		_, err = c.CreateWatch(req1, sub1, w1Next)
		require.NoError(t, err)

		w2Next := make(chan Response, 1)
		_, err = c.CreateWatch(req2, sub2, w2Next)
		require.NoError(t, err)

		mustBlock(t, w1Next)
		mustBlock(t, w2Next)

		setExplicitResources(t, c, node1, "a")
		verifyResponseResources(t, w1Next, wildcardTestType, "2", "a")
		mustBlock(t, w2Next)
	})

	t.Run("handles empty resource list", func(t *testing.T) {
		c := newExplicitWildcardCache(t)

		require.NoError(t, c.UpdateResource("a", buildTestEndpoint("a")))

		setExplicitResources(t, c, node1)

		w, _ := createNodeWildcardWatch(t, c, node1)
		verifyResponseResources(t, w, wildcardTestType, "")
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

func TestLinearExplicitWildcardSotw(t *testing.T) {
	t.Run("don't receive entire cache with explicit resources", func(t *testing.T) {
		c := newExplicitWildcardCache(t)

		require.NoError(t, c.UpdateResource("a", buildTestEndpoint("a")))
		require.NoError(t, c.UpdateResource("b", buildTestEndpoint("b")))
		require.NoError(t, c.UpdateResource("c", buildTestEndpoint("c")))

		setExplicitResources(t, c, node1, "a", "b")

		w, sub := createNodeWildcardWatch(t, c, node1)
		resp := verifyResponseResources(t, w, wildcardTestType, "", "a", "b")

		updateFromSotwResponse(resp, &sub, &Request{TypeUrl: wildcardTestType})

		checkTotalWatchCount(t, c, 0)
	})

	t.Run("explicit resources updated while watch active", func(t *testing.T) {
		c := newExplicitWildcardCache(t)

		require.NoError(t, c.UpdateResource("a", buildTestEndpoint("a")))
		require.NoError(t, c.UpdateResource("b", buildTestEndpoint("b")))
		require.NoError(t, c.UpdateResource("c", buildTestEndpoint("c")))

		setExplicitResources(t, c, node1, "a", "b")

		w, sub := createNodeWildcardWatch(t, c, node1)
		resp := verifyResponseResources(t, w, wildcardTestType, "", "a", "b")
		version, _ := resp.GetVersion()
		updateFromSotwResponse(resp, &sub, &Request{TypeUrl: wildcardTestType, VersionInfo: version})

		w2 := make(chan Response, 1)
		req := &Request{
			ResourceNames: []string{"*"},
			TypeUrl:       wildcardTestType,
			VersionInfo:   version,
		}
		_, err := c.CreateWatch(req, sub, w2)
		require.NoError(t, err)
		mustBlock(t, w2)

		setExplicitResources(t, c, node1, "b", "c")
		resp2 := <-w2
		returnedResources := resp2.GetReturnedResources()
		assert.Contains(t, returnedResources, "b")
		assert.Contains(t, returnedResources, "c")
	})
}

func TestLinearExplicitWildcardDelta(t *testing.T) {
	t.Run("delta wildcard with explicit resources", func(t *testing.T) {
		c := newExplicitWildcardCache(t)

		require.NoError(t, c.UpdateResource("a", buildTestEndpoint("a")))
		require.NoError(t, c.UpdateResource("b", buildTestEndpoint("b")))
		require.NoError(t, c.UpdateResource("c", buildTestEndpoint("c")))

		setExplicitResources(t, c, node1, "a", "b")

		w, sub := createNodeWildcardDeltaWatch(t, c, node1)
		resp := <-w
		nextVersionMap := resp.GetNextVersionMap()
		sub.SetReturnedResources(nextVersionMap)

		assert.Len(t, nextVersionMap, 2)
		assert.Contains(t, nextVersionMap, "a")
		assert.Contains(t, nextVersionMap, "b")

		req := &DeltaRequest{
			TypeUrl:       wildcardTestType,
			ResponseNonce: "1",
		}
		w2 := make(chan DeltaResponse, 1)
		_, err := c.CreateDeltaWatch(req, sub, w2)
		require.NoError(t, err)
		mustBlockDelta(t, w2)
	})

	t.Run("delta wildcard resource versions tracked correctly", func(t *testing.T) {
		c := newExplicitWildcardCache(t)

		require.NoError(t, c.UpdateResource("a", buildTestEndpoint("a")))
		require.NoError(t, c.UpdateResource("b", buildTestEndpoint("b")))

		setExplicitResources(t, c, node1, "a", "b")

		w, sub := createNodeWildcardDeltaWatch(t, c, node1)
		resp := <-w
		sub.SetReturnedResources(resp.GetNextVersionMap())

		req := &DeltaRequest{
			TypeUrl:       wildcardTestType,
			ResponseNonce: "1",
		}
		w2 := make(chan DeltaResponse, 1)
		_, err := c.CreateDeltaWatch(req, sub, w2)
		require.NoError(t, err)
		mustBlockDelta(t, w2)

		require.NoError(t, c.UpdateResource("a", buildTestEndpoint("a-updated")))
		resp2 := <-w2
		nextVersionMap2 := resp2.GetNextVersionMap()
		assert.Len(t, nextVersionMap2, 2)
		assert.Contains(t, nextVersionMap2, "a")
	})

	t.Run("delta wildcard with resource removal from explicit list", func(t *testing.T) {
		c := newExplicitWildcardCache(t)

		require.NoError(t, c.UpdateResource("a", buildTestEndpoint("a")))
		require.NoError(t, c.UpdateResource("b", buildTestEndpoint("b")))
		require.NoError(t, c.UpdateResource("c", buildTestEndpoint("c")))

		setExplicitResources(t, c, node1, "a", "b", "c")

		w, sub := createNodeWildcardDeltaWatch(t, c, node1)
		resp := <-w
		sub.SetReturnedResources(resp.GetNextVersionMap())

		req := &DeltaRequest{
			TypeUrl:       wildcardTestType,
			ResponseNonce: "1",
		}
		w2 := make(chan DeltaResponse, 1)
		_, err := c.CreateDeltaWatch(req, sub, w2)
		require.NoError(t, err)
		mustBlockDelta(t, w2)

		setExplicitResources(t, c, node1, "a", "b")
		resp2 := <-w2
		deltaResp, err := resp2.GetDeltaDiscoveryResponse()
		require.NoError(t, err)
		removedResources := deltaResp.GetRemovedResources()
		assert.Contains(t, removedResources, "c")
	})

	t.Run("delta wildcard update explicit list triggers watch", func(t *testing.T) {
		c := newExplicitWildcardCache(t)

		require.NoError(t, c.UpdateResource("a", buildTestEndpoint("a")))
		require.NoError(t, c.UpdateResource("b", buildTestEndpoint("b")))
		require.NoError(t, c.UpdateResource("c", buildTestEndpoint("c")))

		setExplicitResources(t, c, node1, "a", "b")

		w, sub := createNodeWildcardDeltaWatch(t, c, node1)
		resp := <-w
		sub.SetReturnedResources(resp.GetNextVersionMap())

		req := &DeltaRequest{
			TypeUrl:       wildcardTestType,
			ResponseNonce: "1",
		}
		w2 := make(chan DeltaResponse, 1)
		_, err := c.CreateDeltaWatch(req, sub, w2)
		require.NoError(t, err)
		mustBlockDelta(t, w2)

		setExplicitResources(t, c, node1, "a", "c")
		resp2 := <-w2
		nextVersionMap2 := resp2.GetNextVersionMap()
		deltaResp, err := resp2.GetDeltaDiscoveryResponse()
		require.NoError(t, err)
		removedResources := deltaResp.GetRemovedResources()
		assert.Contains(t, nextVersionMap2, "c")
		assert.Contains(t, removedResources, "b")
	})
}

func TestLinearExplicitWildcardMixedSubscriptions(t *testing.T) {
	t.Run("wildcard and non-wildcard watches coexist", func(t *testing.T) {
		c := newExplicitWildcardCache(t)

		require.NoError(t, c.UpdateResource("a", buildTestEndpoint("a")))
		require.NoError(t, c.UpdateResource("b", buildTestEndpoint("b")))

		setExplicitResources(t, c, node1, "a")

		wWildcard, subWildcard := createNodeWildcardWatch(t, c, node1)
		respWildcard := verifyResponseResources(t, wWildcard, wildcardTestType, "", "a")
		updateFromSotwResponse(respWildcard, &subWildcard, &Request{TypeUrl: wildcardTestType})

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

		w1, sub1 := createNodeWildcardWatch(t, c, node1)
		resp1 := verifyResponseResources(t, w1, wildcardTestType, "", "a")
		version1, _ := resp1.GetVersion()
		updateFromSotwResponse(resp1, &sub1, &Request{TypeUrl: wildcardTestType, VersionInfo: version1})

		w2, sub2 := createNodeWildcardWatch(t, c, node2)
		resp2 := verifyResponseResources(t, w2, wildcardTestType, "", "b")
		version2, _ := resp2.GetVersion()
		updateFromSotwResponse(resp2, &sub2, &Request{TypeUrl: wildcardTestType, VersionInfo: version2})

		w1Next := make(chan Response, 1)
		_, err := c.CreateWatch(&Request{
			ResourceNames: []string{"*"},
			TypeUrl:       wildcardTestType,
			VersionInfo:   version1,
		}, sub1, w1Next)
		require.NoError(t, err)
		mustBlock(t, w1Next)

		w2Next := make(chan Response, 1)
		_, err = c.CreateWatch(&Request{
			ResourceNames: []string{"*"},
			TypeUrl:       wildcardTestType,
			VersionInfo:   version2,
		}, sub2, w2Next)
		require.NoError(t, err)
		mustBlock(t, w2Next)

		require.NoError(t, c.UpdateResource("a", buildTestEndpoint("a-updated")))

		resp1Next := <-w1Next
		assert.Contains(t, resp1Next.GetReturnedResources(), "a")
		mustBlock(t, w2Next)
	})
}
