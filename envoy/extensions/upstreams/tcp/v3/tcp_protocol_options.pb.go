// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v5.28.2
// source: envoy/extensions/upstreams/tcp/v3/tcp_protocol_options.proto

package tcpv3

import (
	_ "github.com/cncf/xds/go/udpa/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type TcpProtocolOptions struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The idle timeout for the connection. The idle timeout is defined as the period in which
	// the connection is not associated with a downstream connection. When the idle timeout is
	// reached, the connection will be closed.
	//
	// If not set, the default idle timeout is 10 minutes. To disable idle timeouts, explicitly set this to 0.
	//
	// .. warning::
	//
	//	Disabling this timeout has a highly likelihood of yielding connection leaks due to lost TCP
	//	FIN packets, etc.
	IdleTimeout *durationpb.Duration `protobuf:"bytes,1,opt,name=idle_timeout,json=idleTimeout,proto3" json:"idle_timeout,omitempty"`
}

func (x *TcpProtocolOptions) Reset() {
	*x = TcpProtocolOptions{}
	if protoimpl.UnsafeEnabled {
		mi := &file_envoy_extensions_upstreams_tcp_v3_tcp_protocol_options_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TcpProtocolOptions) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TcpProtocolOptions) ProtoMessage() {}

func (x *TcpProtocolOptions) ProtoReflect() protoreflect.Message {
	mi := &file_envoy_extensions_upstreams_tcp_v3_tcp_protocol_options_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TcpProtocolOptions.ProtoReflect.Descriptor instead.
func (*TcpProtocolOptions) Descriptor() ([]byte, []int) {
	return file_envoy_extensions_upstreams_tcp_v3_tcp_protocol_options_proto_rawDescGZIP(), []int{0}
}

func (x *TcpProtocolOptions) GetIdleTimeout() *durationpb.Duration {
	if x != nil {
		return x.IdleTimeout
	}
	return nil
}

var File_envoy_extensions_upstreams_tcp_v3_tcp_protocol_options_proto protoreflect.FileDescriptor

var file_envoy_extensions_upstreams_tcp_v3_tcp_protocol_options_proto_rawDesc = []byte{
	0x0a, 0x3c, 0x65, 0x6e, 0x76, 0x6f, 0x79, 0x2f, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f,
	0x6e, 0x73, 0x2f, 0x75, 0x70, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x73, 0x2f, 0x74, 0x63, 0x70,
	0x2f, 0x76, 0x33, 0x2f, 0x74, 0x63, 0x70, 0x5f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c,
	0x5f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x21,
	0x65, 0x6e, 0x76, 0x6f, 0x79, 0x2e, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x73,
	0x2e, 0x75, 0x70, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x73, 0x2e, 0x74, 0x63, 0x70, 0x2e, 0x76,
	0x33, 0x1a, 0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2f, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x1a, 0x1d, 0x75, 0x64, 0x70, 0x61, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x2f, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0x52, 0x0a, 0x12, 0x54, 0x63, 0x70, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x4f,
	0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x3c, 0x0a, 0x0c, 0x69, 0x64, 0x6c, 0x65, 0x5f, 0x74,
	0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44,
	0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0b, 0x69, 0x64, 0x6c, 0x65, 0x54, 0x69, 0x6d,
	0x65, 0x6f, 0x75, 0x74, 0x42, 0xa4, 0x01, 0xba, 0x80, 0xc8, 0xd1, 0x06, 0x02, 0x10, 0x02, 0x0a,
	0x2f, 0x69, 0x6f, 0x2e, 0x65, 0x6e, 0x76, 0x6f, 0x79, 0x70, 0x72, 0x6f, 0x78, 0x79, 0x2e, 0x65,
	0x6e, 0x76, 0x6f, 0x79, 0x2e, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x2e,
	0x75, 0x70, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x73, 0x2e, 0x74, 0x63, 0x70, 0x2e, 0x76, 0x33,
	0x42, 0x17, 0x54, 0x63, 0x70, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x4f, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x4e, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x65, 0x6e, 0x76, 0x6f, 0x79, 0x70, 0x72, 0x6f,
	0x78, 0x79, 0x2f, 0x67, 0x6f, 0x2d, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x2d, 0x70, 0x6c,
	0x61, 0x6e, 0x65, 0x2f, 0x65, 0x6e, 0x76, 0x6f, 0x79, 0x2f, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73,
	0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x75, 0x70, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x73, 0x2f, 0x74,
	0x63, 0x70, 0x2f, 0x76, 0x33, 0x3b, 0x74, 0x63, 0x70, 0x76, 0x33, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_envoy_extensions_upstreams_tcp_v3_tcp_protocol_options_proto_rawDescOnce sync.Once
	file_envoy_extensions_upstreams_tcp_v3_tcp_protocol_options_proto_rawDescData = file_envoy_extensions_upstreams_tcp_v3_tcp_protocol_options_proto_rawDesc
)

func file_envoy_extensions_upstreams_tcp_v3_tcp_protocol_options_proto_rawDescGZIP() []byte {
	file_envoy_extensions_upstreams_tcp_v3_tcp_protocol_options_proto_rawDescOnce.Do(func() {
		file_envoy_extensions_upstreams_tcp_v3_tcp_protocol_options_proto_rawDescData = protoimpl.X.CompressGZIP(file_envoy_extensions_upstreams_tcp_v3_tcp_protocol_options_proto_rawDescData)
	})
	return file_envoy_extensions_upstreams_tcp_v3_tcp_protocol_options_proto_rawDescData
}

var file_envoy_extensions_upstreams_tcp_v3_tcp_protocol_options_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_envoy_extensions_upstreams_tcp_v3_tcp_protocol_options_proto_goTypes = []interface{}{
	(*TcpProtocolOptions)(nil),  // 0: envoy.extensions.upstreams.tcp.v3.TcpProtocolOptions
	(*durationpb.Duration)(nil), // 1: google.protobuf.Duration
}
var file_envoy_extensions_upstreams_tcp_v3_tcp_protocol_options_proto_depIdxs = []int32{
	1, // 0: envoy.extensions.upstreams.tcp.v3.TcpProtocolOptions.idle_timeout:type_name -> google.protobuf.Duration
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_envoy_extensions_upstreams_tcp_v3_tcp_protocol_options_proto_init() }
func file_envoy_extensions_upstreams_tcp_v3_tcp_protocol_options_proto_init() {
	if File_envoy_extensions_upstreams_tcp_v3_tcp_protocol_options_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_envoy_extensions_upstreams_tcp_v3_tcp_protocol_options_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TcpProtocolOptions); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_envoy_extensions_upstreams_tcp_v3_tcp_protocol_options_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_envoy_extensions_upstreams_tcp_v3_tcp_protocol_options_proto_goTypes,
		DependencyIndexes: file_envoy_extensions_upstreams_tcp_v3_tcp_protocol_options_proto_depIdxs,
		MessageInfos:      file_envoy_extensions_upstreams_tcp_v3_tcp_protocol_options_proto_msgTypes,
	}.Build()
	File_envoy_extensions_upstreams_tcp_v3_tcp_protocol_options_proto = out.File
	file_envoy_extensions_upstreams_tcp_v3_tcp_protocol_options_proto_rawDesc = nil
	file_envoy_extensions_upstreams_tcp_v3_tcp_protocol_options_proto_goTypes = nil
	file_envoy_extensions_upstreams_tcp_v3_tcp_protocol_options_proto_depIdxs = nil
}
