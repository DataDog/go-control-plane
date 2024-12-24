package types_test

import (
	"testing"

	"github.com/envoyproxy/go-control-plane/envoy/api/v2/cluster"
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	"github.com/envoyproxy/go-control-plane/pkg/test/resource/v3"
	"github.com/stretchr/testify/assert"
)

const (
	clusterName         = "cluster0"
	routeName           = "route0"
	embeddedRouteName   = "embeddedRoute0"
	scopedRouteName     = "scopedRoute0"
	listenerName        = "listener0"
	scopedListenerName  = "scopedListener0"
	virtualHostName     = "virtualHost0"
	runtimeName         = "runtime0"
	tlsName             = "secret0"
	rootName            = "root0"
	extensionConfigName = "extensionConfig0"
)

var (
	testEndpoint        = resource.MakeEndpoint(clusterName, 8080)
	testCluster         = resource.MakeCluster(resource.Ads, clusterName)
	testRoute           = resource.MakeRouteConfig(routeName, clusterName)
	testEmbeddedRoute   = resource.MakeRouteConfig(embeddedRouteName, clusterName)
	testScopedRoute     = resource.MakeScopedRouteConfig(scopedRouteName, routeName, []string{"1.2.3.4"})
	testVirtualHost     = resource.MakeVirtualHost(virtualHostName, clusterName)
	testListener        = resource.MakeRouteHTTPListener(resource.Ads, listenerName, 80, routeName)
	testListenerDefault = resource.MakeRouteHTTPListenerDefaultFilterChain(resource.Ads, listenerName, 80, routeName)
	testScopedListener  = resource.MakeScopedRouteHTTPListenerForRoute(resource.Ads, scopedListenerName, 80, embeddedRouteName)
	testRuntime         = resource.MakeRuntime(runtimeName)
	testSecret          = resource.MakeSecrets(tlsName, rootName)
	testExtensionConfig = resource.MakeExtensionConfig(resource.Ads, extensionConfigName, routeName)
)

type customResource struct {
	cluster.Filter // Any proto would work here.
}

const customName = "test-name"

func (cs *customResource) GetName() string { return customName }

var _ types.ResourceWithName = &customResource{}

func TestGetResourceName(t *testing.T) {
	assert.Equal(t, clusterName, types.GetResourceName(testEndpoint))
	assert.Equal(t, clusterName, types.GetResourceName(testCluster))
	assert.Equal(t, routeName, types.GetResourceName(testRoute))
	assert.Equal(t, scopedRouteName, types.GetResourceName(testScopedRoute))
	assert.Equal(t, virtualHostName, types.GetResourceName(testVirtualHost))
	assert.Equal(t, listenerName, types.GetResourceName(testListener))
	assert.Equal(t, runtimeName, types.GetResourceName(testRuntime))
	assert.Equal(t, customName, types.GetResourceName(&customResource{}))
	assert.Equal(t, "", types.GetResourceName(nil))
}

func TestGetResourceNames(t *testing.T) {
	tests := []struct {
		name  string
		input []types.ResourceWithTTL
		want  []string
	}{
		{
			name:  "empty",
			input: []types.ResourceWithTTL{},
			want:  []string{},
		},
		{
			name:  "many",
			input: []types.ResourceWithTTL{{Resource: testRuntime}, {Resource: testListener}, {Resource: testListenerDefault}, {Resource: testVirtualHost}},
			want:  []string{runtimeName, listenerName, listenerName, virtualHostName},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			got := types.GetResourceNames(test.input)
			assert.ElementsMatch(t, test.want, got)
		})
	}
}
