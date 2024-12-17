// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v5.29.1
// source: envoy/config/retry/previous_priorities/previous_priorities_config.proto

package previous_priorities

import (
	_ "github.com/cncf/xds/go/udpa/annotations"
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// A retry host selector that attempts to spread retries between priorities, even if certain
// priorities would not normally be attempted due to higher priorities being available.
//
// As priorities get excluded, load will be distributed amongst the remaining healthy priorities
// based on the relative health of the priorities, matching how load is distributed during regular
// host selection. For example, given priority healths of {100, 50, 50}, the original load will be
// {100, 0, 0} (since P0 has capacity to handle 100% of the traffic). If P0 is excluded, the load
// changes to {0, 50, 50}, because P1 is only able to handle 50% of the traffic, causing the
// remaining to spill over to P2.
//
// Each priority attempted will be excluded until there are no healthy priorities left, at which
// point the list of attempted priorities will be reset, essentially starting from the beginning.
// For example, given three priorities P0, P1, P2 with healthy % of 100, 0 and 50 respectively, the
// following sequence of priorities would be selected (assuming update_frequency = 1):
// Attempt 1: P0 (P0 is 100% healthy)
// Attempt 2: P2 (P0 already attempted, P2 only healthy priority)
// Attempt 3: P0 (no healthy priorities, reset)
// Attempt 4: P2
//
// In the case of all upstream hosts being unhealthy, no adjustments will be made to the original
// priority load, so behavior should be identical to not using this plugin.
//
// Using this PriorityFilter requires rebuilding the priority load, which runs in O(# of
// priorities), which might incur significant overhead for clusters with many priorities.
// [#extension: envoy.retry_priorities.previous_priorities]
type PreviousPrioritiesConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// How often the priority load should be updated based on previously attempted priorities. Useful
	// to allow each priorities to receive more than one request before being excluded or to reduce
	// the number of times that the priority load has to be recomputed.
	//
	// For example, by setting this to 2, then the first two attempts (initial attempt and first
	// retry) will use the unmodified priority load. The third and fourth attempt will use priority
	// load which excludes the priorities routed to with the first two attempts, and the fifth and
	// sixth attempt will use the priority load excluding the priorities used for the first four
	// attempts.
	//
	// Must be greater than 0.
	UpdateFrequency int32 `protobuf:"varint,1,opt,name=update_frequency,json=updateFrequency,proto3" json:"update_frequency,omitempty"`
}

func (x *PreviousPrioritiesConfig) Reset() {
	*x = PreviousPrioritiesConfig{}
	if protoimpl.UnsafeEnabled {
		mi := &file_envoy_config_retry_previous_priorities_previous_priorities_config_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PreviousPrioritiesConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PreviousPrioritiesConfig) ProtoMessage() {}

func (x *PreviousPrioritiesConfig) ProtoReflect() protoreflect.Message {
	mi := &file_envoy_config_retry_previous_priorities_previous_priorities_config_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PreviousPrioritiesConfig.ProtoReflect.Descriptor instead.
func (*PreviousPrioritiesConfig) Descriptor() ([]byte, []int) {
	return file_envoy_config_retry_previous_priorities_previous_priorities_config_proto_rawDescGZIP(), []int{0}
}

func (x *PreviousPrioritiesConfig) GetUpdateFrequency() int32 {
	if x != nil {
		return x.UpdateFrequency
	}
	return 0
}

var File_envoy_config_retry_previous_priorities_previous_priorities_config_proto protoreflect.FileDescriptor

var file_envoy_config_retry_previous_priorities_previous_priorities_config_proto_rawDesc = []byte{
	0x0a, 0x47, 0x65, 0x6e, 0x76, 0x6f, 0x79, 0x2f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2f, 0x72,
	0x65, 0x74, 0x72, 0x79, 0x2f, 0x70, 0x72, 0x65, 0x76, 0x69, 0x6f, 0x75, 0x73, 0x5f, 0x70, 0x72,
	0x69, 0x6f, 0x72, 0x69, 0x74, 0x69, 0x65, 0x73, 0x2f, 0x70, 0x72, 0x65, 0x76, 0x69, 0x6f, 0x75,
	0x73, 0x5f, 0x70, 0x72, 0x69, 0x6f, 0x72, 0x69, 0x74, 0x69, 0x65, 0x73, 0x5f, 0x63, 0x6f, 0x6e,
	0x66, 0x69, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x26, 0x65, 0x6e, 0x76, 0x6f, 0x79,
	0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x72, 0x65, 0x74, 0x72, 0x79, 0x2e, 0x70, 0x72,
	0x65, 0x76, 0x69, 0x6f, 0x75, 0x73, 0x5f, 0x70, 0x72, 0x69, 0x6f, 0x72, 0x69, 0x74, 0x69, 0x65,
	0x73, 0x1a, 0x1e, 0x75, 0x64, 0x70, 0x61, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x2f, 0x6d, 0x69, 0x67, 0x72, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x1a, 0x1d, 0x75, 0x64, 0x70, 0x61, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x2f, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x1a, 0x17, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64,
	0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x4e, 0x0a, 0x18, 0x50, 0x72, 0x65,
	0x76, 0x69, 0x6f, 0x75, 0x73, 0x50, 0x72, 0x69, 0x6f, 0x72, 0x69, 0x74, 0x69, 0x65, 0x73, 0x43,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x32, 0x0a, 0x10, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x5f,
	0x66, 0x72, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x42,
	0x07, 0xfa, 0x42, 0x04, 0x1a, 0x02, 0x20, 0x00, 0x52, 0x0f, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65,
	0x46, 0x72, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x79, 0x42, 0xec, 0x01, 0xf2, 0x98, 0xfe, 0x8f,
	0x05, 0x38, 0x12, 0x36, 0x65, 0x6e, 0x76, 0x6f, 0x79, 0x2e, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73,
	0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x72, 0x65, 0x74, 0x72, 0x79, 0x2e, 0x70, 0x72, 0x69, 0x6f, 0x72,
	0x69, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x65, 0x76, 0x69, 0x6f, 0x75, 0x73, 0x5f, 0x70, 0x72, 0x69,
	0x6f, 0x72, 0x69, 0x74, 0x69, 0x65, 0x73, 0x2e, 0x76, 0x33, 0xba, 0x80, 0xc8, 0xd1, 0x06, 0x02,
	0x10, 0x01, 0x0a, 0x34, 0x69, 0x6f, 0x2e, 0x65, 0x6e, 0x76, 0x6f, 0x79, 0x70, 0x72, 0x6f, 0x78,
	0x79, 0x2e, 0x65, 0x6e, 0x76, 0x6f, 0x79, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x72,
	0x65, 0x74, 0x72, 0x79, 0x2e, 0x70, 0x72, 0x65, 0x76, 0x69, 0x6f, 0x75, 0x73, 0x5f, 0x70, 0x72,
	0x69, 0x6f, 0x72, 0x69, 0x74, 0x69, 0x65, 0x73, 0x42, 0x1d, 0x50, 0x72, 0x65, 0x76, 0x69, 0x6f,
	0x75, 0x73, 0x50, 0x72, 0x69, 0x6f, 0x72, 0x69, 0x74, 0x69, 0x65, 0x73, 0x43, 0x6f, 0x6e, 0x66,
	0x69, 0x67, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x4d, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x65, 0x6e, 0x76, 0x6f, 0x79, 0x70, 0x72, 0x6f, 0x78, 0x79,
	0x2f, 0x67, 0x6f, 0x2d, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x2d, 0x70, 0x6c, 0x61, 0x6e,
	0x65, 0x2f, 0x65, 0x6e, 0x76, 0x6f, 0x79, 0x2f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2f, 0x72,
	0x65, 0x74, 0x72, 0x79, 0x2f, 0x70, 0x72, 0x65, 0x76, 0x69, 0x6f, 0x75, 0x73, 0x5f, 0x70, 0x72,
	0x69, 0x6f, 0x72, 0x69, 0x74, 0x69, 0x65, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_envoy_config_retry_previous_priorities_previous_priorities_config_proto_rawDescOnce sync.Once
	file_envoy_config_retry_previous_priorities_previous_priorities_config_proto_rawDescData = file_envoy_config_retry_previous_priorities_previous_priorities_config_proto_rawDesc
)

func file_envoy_config_retry_previous_priorities_previous_priorities_config_proto_rawDescGZIP() []byte {
	file_envoy_config_retry_previous_priorities_previous_priorities_config_proto_rawDescOnce.Do(func() {
		file_envoy_config_retry_previous_priorities_previous_priorities_config_proto_rawDescData = protoimpl.X.CompressGZIP(file_envoy_config_retry_previous_priorities_previous_priorities_config_proto_rawDescData)
	})
	return file_envoy_config_retry_previous_priorities_previous_priorities_config_proto_rawDescData
}

var file_envoy_config_retry_previous_priorities_previous_priorities_config_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_envoy_config_retry_previous_priorities_previous_priorities_config_proto_goTypes = []interface{}{
	(*PreviousPrioritiesConfig)(nil), // 0: envoy.config.retry.previous_priorities.PreviousPrioritiesConfig
}
var file_envoy_config_retry_previous_priorities_previous_priorities_config_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_envoy_config_retry_previous_priorities_previous_priorities_config_proto_init() }
func file_envoy_config_retry_previous_priorities_previous_priorities_config_proto_init() {
	if File_envoy_config_retry_previous_priorities_previous_priorities_config_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_envoy_config_retry_previous_priorities_previous_priorities_config_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PreviousPrioritiesConfig); i {
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
			RawDescriptor: file_envoy_config_retry_previous_priorities_previous_priorities_config_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_envoy_config_retry_previous_priorities_previous_priorities_config_proto_goTypes,
		DependencyIndexes: file_envoy_config_retry_previous_priorities_previous_priorities_config_proto_depIdxs,
		MessageInfos:      file_envoy_config_retry_previous_priorities_previous_priorities_config_proto_msgTypes,
	}.Build()
	File_envoy_config_retry_previous_priorities_previous_priorities_config_proto = out.File
	file_envoy_config_retry_previous_priorities_previous_priorities_config_proto_rawDesc = nil
	file_envoy_config_retry_previous_priorities_previous_priorities_config_proto_goTypes = nil
	file_envoy_config_retry_previous_priorities_previous_priorities_config_proto_depIdxs = nil
}
