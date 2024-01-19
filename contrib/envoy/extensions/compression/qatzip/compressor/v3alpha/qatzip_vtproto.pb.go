//go:build vtprotobuf
// +build vtprotobuf
// Code generated by protoc-gen-go-vtproto. DO NOT EDIT.
// source: contrib/envoy/extensions/compression/qatzip/compressor/v3alpha/qatzip.proto

package v3alpha

import (
	wrapperspb "github.com/planetscale/vtprotobuf/types/known/wrapperspb"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	bits "math/bits"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

func (m *Qatzip) MarshalVTStrict() (dAtA []byte, err error) {
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

func (m *Qatzip) MarshalToVTStrict(dAtA []byte) (int, error) {
	size := m.SizeVT()
	return m.MarshalToSizedBufferVTStrict(dAtA[:size])
}

func (m *Qatzip) MarshalToSizedBufferVTStrict(dAtA []byte) (int, error) {
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
	if m.ChunkSize != nil {
		size, err := (*wrapperspb.UInt32Value)(m.ChunkSize).MarshalToSizedBufferVTStrict(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarint(dAtA, i, uint64(size))
		i--
		dAtA[i] = 0x2a
	}
	if m.StreamBufferSize != nil {
		size, err := (*wrapperspb.UInt32Value)(m.StreamBufferSize).MarshalToSizedBufferVTStrict(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarint(dAtA, i, uint64(size))
		i--
		dAtA[i] = 0x22
	}
	if m.InputSizeThreshold != nil {
		size, err := (*wrapperspb.UInt32Value)(m.InputSizeThreshold).MarshalToSizedBufferVTStrict(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarint(dAtA, i, uint64(size))
		i--
		dAtA[i] = 0x1a
	}
	if m.HardwareBufferSize != 0 {
		i = encodeVarint(dAtA, i, uint64(m.HardwareBufferSize))
		i--
		dAtA[i] = 0x10
	}
	if m.CompressionLevel != nil {
		size, err := (*wrapperspb.UInt32Value)(m.CompressionLevel).MarshalToSizedBufferVTStrict(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarint(dAtA, i, uint64(size))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarint(dAtA []byte, offset int, v uint64) int {
	offset -= sov(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Qatzip) SizeVT() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.CompressionLevel != nil {
		l = (*wrapperspb.UInt32Value)(m.CompressionLevel).SizeVT()
		n += 1 + l + sov(uint64(l))
	}
	if m.HardwareBufferSize != 0 {
		n += 1 + sov(uint64(m.HardwareBufferSize))
	}
	if m.InputSizeThreshold != nil {
		l = (*wrapperspb.UInt32Value)(m.InputSizeThreshold).SizeVT()
		n += 1 + l + sov(uint64(l))
	}
	if m.StreamBufferSize != nil {
		l = (*wrapperspb.UInt32Value)(m.StreamBufferSize).SizeVT()
		n += 1 + l + sov(uint64(l))
	}
	if m.ChunkSize != nil {
		l = (*wrapperspb.UInt32Value)(m.ChunkSize).SizeVT()
		n += 1 + l + sov(uint64(l))
	}
	n += len(m.unknownFields)
	return n
}

func sov(x uint64) (n int) {
	return (bits.Len64(x|1) + 6) / 7
}
func soz(x uint64) (n int) {
	return sov(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}