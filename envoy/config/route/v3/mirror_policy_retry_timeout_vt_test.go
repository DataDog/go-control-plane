//go:build vtprotobuf

package routev3

import (
	"testing"

	"google.golang.org/protobuf/proto"
)

func TestMirrorPolicyRetryTimeoutVT(t *testing.T) {
	for _, test := range mirrorRetryTimeoutCases() {
		t.Run(test.name, func(t *testing.T) {
			encoded, err := test.policy.MarshalVTStrict()
			if err != nil {
				t.Fatal(err)
			}
			if len(encoded) != test.policy.SizeVT() {
				t.Fatal("vtproto size differs from serialized size")
			}
			var decoded RouteAction_RequestMirrorPolicy
			if err := proto.Unmarshal(encoded, &decoded); err != nil {
				t.Fatal(err)
			}
			if !proto.Equal(test.policy, &decoded) {
				t.Fatal("mirror settings lost during vtproto serialization")
			}
		})
	}
}
