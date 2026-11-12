//go:build vtprotobuf

package on_demandv3

import (
	"testing"

	"google.golang.org/protobuf/proto"
)

func TestDatadogOnDemandWireFormatVT(t *testing.T) {
	testDatadogOnDemandWireFormat(t, func(m proto.Message) ([]byte, error) {
		return m.(interface{ MarshalVTStrict() ([]byte, error) }).MarshalVTStrict()
	})
}
