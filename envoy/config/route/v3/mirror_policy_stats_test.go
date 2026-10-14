package routev3

import (
	"testing"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func TestMirrorPolicyStatPrefix(t *testing.T) {
	policy := &RouteAction_RequestMirrorPolicy{Cluster: "mirror", StatPrefix: "eu1_probe_mirror"}
	field := policy.ProtoReflect().Descriptor().Fields().ByName("stat_prefix")
	if field == nil || field.Number() != 9 {
		t.Fatal("stat_prefix must use Envoy field 9")
	}
	encoded, err := proto.Marshal(policy)
	if err != nil {
		t.Fatal(err)
	}
	var decoded RouteAction_RequestMirrorPolicy
	if err := proto.Unmarshal(encoded, &decoded); err != nil {
		t.Fatal(err)
	}
	if !proto.Equal(policy, &decoded) {
		t.Fatal("stat prefix lost during binary round trip")
	}
	if err := protojson.Unmarshal([]byte(`{"cluster":"mirror","stat_prefix":"eu1_probe_mirror"}`), &decoded); err != nil {
		t.Fatal(err)
	}
	if !proto.Equal(policy, &decoded) {
		t.Fatal("stat prefix lost during JSON decoding")
	}
}
