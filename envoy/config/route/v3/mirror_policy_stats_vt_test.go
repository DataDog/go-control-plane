//go:build vtprotobuf

package routev3

import (
	"testing"

	"google.golang.org/protobuf/proto"
)

func TestMirrorPolicyStatPrefixVT(t *testing.T) {
	policy := &RouteAction_RequestMirrorPolicy{Cluster: "mirror", StatPrefix: "eu1_probe_mirror"}
	encoded, err := policy.MarshalVTStrict()
	if err != nil {
		t.Fatal(err)
	}
	if len(encoded) != policy.SizeVT() {
		t.Fatal("vtproto size differs from serialized size")
	}
	var decoded RouteAction_RequestMirrorPolicy
	if err := proto.Unmarshal(encoded, &decoded); err != nil {
		t.Fatal(err)
	}
	if !proto.Equal(policy, &decoded) {
		t.Fatal("stat prefix lost during vtproto serialization")
	}
}
