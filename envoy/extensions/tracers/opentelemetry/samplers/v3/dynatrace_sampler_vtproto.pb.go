//go:build vtprotobuf
// +build vtprotobuf
// Code generated by protoc-gen-go-vtproto. DO NOT EDIT.
// source: envoy/extensions/tracers/opentelemetry/samplers/v3/dynatrace_sampler.proto

package samplersv3

import (
	proto "google.golang.org/protobuf/proto"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

func (m *DynatraceSamplerConfig) MarshalVTStrict() (dAtA []byte, err error) {
	if m == nil {
		return nil, nil
	}
	size := m.SizeVT()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBufferVTStrict(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *DynatraceSamplerConfig) MarshalToVTStrict(dAtA []byte) (int, error) {
	size := m.SizeVT()
	return m.MarshalToSizedBufferVTStrict(dAtA[:size])
}

func (m *DynatraceSamplerConfig) MarshalToSizedBufferVTStrict(dAtA []byte) (int, error) {
	if m == nil {
		return 0, nil
	}
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.unknownFields != nil {
		i -= len(m.unknownFields)
		copy(dAtA[i:], m.unknownFields)
	}
	if m.RootSpansPerMinute != 0 {
		i = encodeVarint(dAtA, i, uint64(m.RootSpansPerMinute))
		i--
		dAtA[i] = 0x20
	}
	if m.HttpService != nil {
		if vtmsg, ok := interface{}(m.HttpService).(interface {
			MarshalToSizedBufferVTStrict([]byte) (int, error)
		}); ok {
			size, err := vtmsg.MarshalToSizedBufferVTStrict(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarint(dAtA, i, uint64(size))
		} else {
			encoded, err := proto.Marshal(m.HttpService)
			if err != nil {
				return 0, err
			}
			i -= len(encoded)
			copy(dAtA[i:], encoded)
			i = encodeVarint(dAtA, i, uint64(len(encoded)))
		}
		i--
		dAtA[i] = 0x1a
	}
	if m.ClusterId != 0 {
		i = encodeVarint(dAtA, i, uint64(m.ClusterId))
		i--
		dAtA[i] = 0x10
	}
	if len(m.Tenant) > 0 {
		i -= len(m.Tenant)
		copy(dAtA[i:], m.Tenant)
		i = encodeVarint(dAtA, i, uint64(len(m.Tenant)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *DynatraceSamplerConfig) SizeVT() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Tenant)
	if l > 0 {
		n += 1 + l + sov(uint64(l))
	}
	if m.ClusterId != 0 {
		n += 1 + sov(uint64(m.ClusterId))
	}
	if m.HttpService != nil {
		if size, ok := interface{}(m.HttpService).(interface {
			SizeVT() int
		}); ok {
			l = size.SizeVT()
		} else {
			l = proto.Size(m.HttpService)
		}
		n += 1 + l + sov(uint64(l))
	}
	if m.RootSpansPerMinute != 0 {
		n += 1 + sov(uint64(m.RootSpansPerMinute))
	}
	n += len(m.unknownFields)
	return n
}
