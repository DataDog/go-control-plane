//go:build !disable_pgv

package routev3

import "testing"

func TestMirrorPolicyRetryValidation(t *testing.T) {
	for _, test := range mirrorRetryTimeoutCases() {
		if err := test.policy.ValidateAll(); err != nil {
			t.Fatalf("%s: %v", test.name, err)
		}
	}
	policy := &RouteAction_RequestMirrorPolicy{
		Cluster:     "mirror",
		RetryPolicy: &RetryPolicy{RetryBackOff: &RetryPolicy_RetryBackOff{}},
	}
	if err := policy.Validate(); err == nil {
		t.Fatal("missing retry backoff base interval must fail validation")
	}
	if err := policy.ValidateAll(); err == nil {
		t.Fatal("missing retry backoff base interval must fail aggregate validation")
	}
}
