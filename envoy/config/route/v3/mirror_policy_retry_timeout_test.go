package routev3

import (
	"testing"
	"time"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func mirrorRetryTimeoutCases() []struct {
	name   string
	json   string
	policy *RouteAction_RequestMirrorPolicy
} {
	return []struct {
		name   string
		json   string
		policy *RouteAction_RequestMirrorPolicy
	}{
		{"unset", `{"cluster":"mirror"}`, &RouteAction_RequestMirrorPolicy{Cluster: "mirror"}},
		{"retry only", `{"cluster":"mirror","retry_policy":{"retry_on":"5xx","num_retries":1}}`, &RouteAction_RequestMirrorPolicy{
			Cluster: "mirror", RetryPolicy: &RetryPolicy{RetryOn: "5xx", NumRetries: wrapperspb.UInt32(1)},
		}},
		{"timeout disabled", `{"cluster":"mirror","retry_policy":{"retry_on":"5xx","num_retries":1},"timeout":"0s","stat_prefix":"eu1_probe_mirror"}`, &RouteAction_RequestMirrorPolicy{
			Cluster: "mirror", StatPrefix: "eu1_probe_mirror", RetryPolicy: &RetryPolicy{RetryOn: "5xx", NumRetries: wrapperspb.UInt32(1)}, Timeout: durationpb.New(0),
		}},
		{"timeout with retries", `{"cluster":"mirror","retry_policy":{"retry_on":"5xx","num_retries":1},"timeout":"0.200s"}`, &RouteAction_RequestMirrorPolicy{
			Cluster: "mirror", RetryPolicy: &RetryPolicy{RetryOn: "5xx", NumRetries: wrapperspb.UInt32(1)}, Timeout: durationpb.New(200 * time.Millisecond),
		}},
		{"timeout only", `{"cluster":"mirror","timeout":"0.200s"}`, &RouteAction_RequestMirrorPolicy{
			Cluster: "mirror", Timeout: durationpb.New(200 * time.Millisecond),
		}},
	}
}

func TestMirrorPolicyRetryTimeout(t *testing.T) {
	fields := (&RouteAction_RequestMirrorPolicy{}).ProtoReflect().Descriptor().Fields()
	for name, number := range map[protoreflect.Name]protoreflect.FieldNumber{"retry_policy": 10, "timeout": 11} {
		field := fields.ByName(name)
		if field == nil || field.Number() != number || !field.HasPresence() {
			t.Fatalf("%s must be an optional message at Envoy field %d", name, number)
		}
	}
	for _, test := range mirrorRetryTimeoutCases() {
		t.Run(test.name, func(t *testing.T) {
			encoded, err := proto.Marshal(test.policy)
			if err != nil {
				t.Fatal(err)
			}
			var decoded RouteAction_RequestMirrorPolicy
			if err := proto.Unmarshal(encoded, &decoded); err != nil {
				t.Fatal(err)
			}
			if !proto.Equal(test.policy, &decoded) {
				t.Fatal("mirror settings lost during binary round trip")
			}
			if err := protojson.Unmarshal([]byte(test.json), &decoded); err != nil {
				t.Fatal(err)
			}
			if !proto.Equal(test.policy, &decoded) {
				t.Fatal("mirror settings lost during JSON decoding")
			}
			encoded, err = protojson.Marshal(test.policy)
			if err != nil {
				t.Fatal(err)
			}
			if err := protojson.Unmarshal(encoded, &decoded); err != nil {
				t.Fatal(err)
			}
			if !proto.Equal(test.policy, &decoded) {
				t.Fatal("mirror settings lost during JSON round trip")
			}
		})
	}
}
