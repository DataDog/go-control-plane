package on_demandv3

import (
	"bytes"
	"testing"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestDatadogOnDemandWireFormat(t *testing.T) {
	testDatadogOnDemandWireFormat(t, proto.Marshal)
}

func testDatadogOnDemandWireFormat(t *testing.T, marshal func(proto.Message) ([]byte, error)) {
	mirror := &MirrorOnDemandCds{
		DefaultPolicy: &MirrorOnDemandCds_InheritPrimary{InheritPrimary: true},
	}
	for _, tt := range []struct {
		name string
		msg  proto.Message
		want []byte
	}{
		{"filter", &OnDemand{MirrorOdcds: mirror}, []byte{0xc2, 0x3e, 0x02, 0x08, 0x01}},
		{"per route", &PerRouteConfig{MirrorOdcds: mirror}, []byte{0xc2, 0x3e, 0x02, 0x08, 0x01}},
		{"zero pending limit", &OnDemandCds{MaxPendingRequests: wrapperspb.UInt32(0)}, []byte{0x22, 0x00}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			// Golden wire bytes catch field-number drift that a same-binding round trip hides.
			got, err := marshal(tt.msg)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(got, tt.want) {
				t.Fatalf("wire encoding = %x, want %x", got, tt.want)
			}
		})
	}
}
