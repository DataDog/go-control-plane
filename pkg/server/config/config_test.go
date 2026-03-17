package config

import (
	"testing"
)

func TestShouldIgnoreWildcard(t *testing.T) {
	tests := []struct {
		name               string
		ignoreWildcardTypes []string
		typeURL            string
		expected           bool
	}{
		{
			name:               "no types configured - should not ignore",
			ignoreWildcardTypes: nil,
			typeURL:            "type.googleapis.com/envoy.config.cluster.v3.Cluster",
			expected:           false,
		},
		{
			name:               "type in ignore list - should ignore",
			ignoreWildcardTypes: []string{"type.googleapis.com/envoy.config.route.v3.VirtualHost"},
			typeURL:            "type.googleapis.com/envoy.config.route.v3.VirtualHost",
			expected:           true,
		},
		{
			name:               "type not in ignore list - should not ignore",
			ignoreWildcardTypes: []string{"type.googleapis.com/envoy.config.route.v3.VirtualHost"},
			typeURL:            "type.googleapis.com/envoy.config.cluster.v3.Cluster",
			expected:           false,
		},
		{
			name: "multiple types in ignore list - should ignore matching type",
			ignoreWildcardTypes: []string{
				"type.googleapis.com/envoy.config.route.v3.VirtualHost",
				"type.googleapis.com/envoy.config.cluster.v3.Cluster",
			},
			typeURL:  "type.googleapis.com/envoy.config.cluster.v3.Cluster",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := NewOpts()
			if tt.ignoreWildcardTypes != nil {
				opt := IgnoreWildcardForTypes(tt.ignoreWildcardTypes)
				opt(&opts)
			}

			result := opts.ShouldIgnoreWildcard(tt.typeURL)
			if result != tt.expected {
				t.Errorf("ShouldIgnoreWildcard() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestIgnoreWildcardForTypes(t *testing.T) {
	opts := NewOpts()

	// Initially should not ignore any type
	if opts.ShouldIgnoreWildcard("any.type") {
		t.Error("Expected not to ignore wildcard by default")
	}

	// Configure to ignore specific types
	types := []string{
		"type.googleapis.com/envoy.config.route.v3.VirtualHost",
		"type.googleapis.com/envoy.config.listener.v3.Listener",
	}
	opt := IgnoreWildcardForTypes(types)
	opt(&opts)

	// Should ignore configured types
	if !opts.ShouldIgnoreWildcard("type.googleapis.com/envoy.config.route.v3.VirtualHost") {
		t.Error("Expected to ignore wildcard for VirtualHost")
	}
	if !opts.ShouldIgnoreWildcard("type.googleapis.com/envoy.config.listener.v3.Listener") {
		t.Error("Expected to ignore wildcard for Listener")
	}

	// Should not ignore other types
	if opts.ShouldIgnoreWildcard("type.googleapis.com/envoy.config.cluster.v3.Cluster") {
		t.Error("Expected not to ignore wildcard for Cluster")
	}
}

func TestIgnoreWildcardAlsoDeactivatesLegacyWildcard(t *testing.T) {
	opts := NewOpts()

	vhdsType := "type.googleapis.com/envoy.config.route.v3.VirtualHost"
	clusterType := "type.googleapis.com/envoy.config.cluster.v3.Cluster"

	if !opts.IsLegacyWildcardActive(vhdsType) {
		t.Error("Expected legacy wildcard to be active for VHDS initially")
	}
	if !opts.IsLegacyWildcardActive(clusterType) {
		t.Error("Expected legacy wildcard to be active for Cluster initially")
	}

	opt := IgnoreWildcardForTypes([]string{vhdsType})
	opt(&opts)

	// Should filter wildcards and disable legacy behavior for VHDS
	if !opts.ShouldIgnoreWildcard(vhdsType) {
		t.Error("Expected to ignore wildcard for VHDS")
	}
	if opts.IsLegacyWildcardActive(vhdsType) {
		t.Error("Expected legacy wildcard to be deactivated for VHDS when ignoring wildcard")
	}

	// Other types shouldn't be affected
	if !opts.IsLegacyWildcardActive(clusterType) {
		t.Error("Expected legacy wildcard to still be active for Cluster")
	}
	if opts.ShouldIgnoreWildcard(clusterType) {
		t.Error("Expected not to ignore wildcard for Cluster")
	}
}
