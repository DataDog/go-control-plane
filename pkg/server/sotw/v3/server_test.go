package sotw

import (
	"context"
	"testing"

	"github.com/envoyproxy/go-control-plane/pkg/server/config"
	"github.com/stretchr/testify/assert"
)

func TestFilterWildcard(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name               string
		ignoreWildcardTypes []string
		typeURL            string
		resourceNames      []string
		expectedFiltered   []string
	}{
		{
			name:               "no ignore configured - wildcard should pass through",
			ignoreWildcardTypes: nil,
			typeURL:            "type.googleapis.com/envoy.config.cluster.v3.Cluster",
			resourceNames:      []string{"*", "resource1"},
			expectedFiltered:   []string{"*", "resource1"},
		},
		{
			name:               "ignore configured for type - wildcard should be filtered",
			ignoreWildcardTypes: []string{"type.googleapis.com/envoy.config.route.v3.VirtualHost"},
			typeURL:            "type.googleapis.com/envoy.config.route.v3.VirtualHost",
			resourceNames:      []string{"*", "resource1"},
			expectedFiltered:   []string{"resource1"},
		},
		{
			name:               "ignore configured for different type - wildcard should pass through",
			ignoreWildcardTypes: []string{"type.googleapis.com/envoy.config.route.v3.VirtualHost"},
			typeURL:            "type.googleapis.com/envoy.config.cluster.v3.Cluster",
			resourceNames:      []string{"*", "resource1"},
			expectedFiltered:   []string{"*", "resource1"},
		},
		{
			name:               "only wildcard in list - should become empty when filtered",
			ignoreWildcardTypes: []string{"type.googleapis.com/envoy.config.route.v3.VirtualHost"},
			typeURL:            "type.googleapis.com/envoy.config.route.v3.VirtualHost",
			resourceNames:      []string{"*"},
			expectedFiltered:   []string{},
		},
		{
			name:               "no wildcard in list - nothing to filter",
			ignoreWildcardTypes: []string{"type.googleapis.com/envoy.config.route.v3.VirtualHost"},
			typeURL:            "type.googleapis.com/envoy.config.route.v3.VirtualHost",
			resourceNames:      []string{"resource1", "resource2"},
			expectedFiltered:   []string{"resource1", "resource2"},
		},
		{
			name:               "empty resource list - should remain empty",
			ignoreWildcardTypes: []string{"type.googleapis.com/envoy.config.route.v3.VirtualHost"},
			typeURL:            "type.googleapis.com/envoy.config.route.v3.VirtualHost",
			resourceNames:      []string{},
			expectedFiltered:   []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var opts []config.XDSOption
			if tt.ignoreWildcardTypes != nil {
				opts = append(opts, IgnoreWildcardForTypes(tt.ignoreWildcardTypes))
			}

			srv := NewServer(ctx, nil, nil, opts...).(*server)
			filtered := srv.filterWildcard(tt.typeURL, tt.resourceNames)

			assert.Equal(t, tt.expectedFiltered, filtered)
		})
	}
}
