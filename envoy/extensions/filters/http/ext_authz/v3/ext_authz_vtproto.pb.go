//go:build vtprotobuf
// +build vtprotobuf
// Code generated by protoc-gen-go-vtproto. DO NOT EDIT.
// source: envoy/extensions/filters/http/ext_authz/v3/ext_authz.proto

package ext_authzv3

import (
	wrapperspb "github.com/planetscale/vtprotobuf/types/known/wrapperspb"
	proto "google.golang.org/protobuf/proto"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	bits "math/bits"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

func (m *ExtAuthz) MarshalVTStrict() (dAtA []byte, err error) {
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

func (m *ExtAuthz) MarshalToVTStrict(dAtA []byte) (int, error) {
	size := m.SizeVT()
	return m.MarshalToSizedBufferVTStrict(dAtA[:size])
}

func (m *ExtAuthz) MarshalToSizedBufferVTStrict(dAtA []byte) (int, error) {
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
	if m.EncodeRawHeaders {
		i--
		if m.EncodeRawHeaders {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x1
		i--
		dAtA[i] = 0xb8
	}
	if len(m.RouteTypedMetadataContextNamespaces) > 0 {
		for iNdEx := len(m.RouteTypedMetadataContextNamespaces) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.RouteTypedMetadataContextNamespaces[iNdEx])
			copy(dAtA[i:], m.RouteTypedMetadataContextNamespaces[iNdEx])
			i = encodeVarint(dAtA, i, uint64(len(m.RouteTypedMetadataContextNamespaces[iNdEx])))
			i--
			dAtA[i] = 0x1
			i--
			dAtA[i] = 0xb2
		}
	}
	if len(m.RouteMetadataContextNamespaces) > 0 {
		for iNdEx := len(m.RouteMetadataContextNamespaces) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.RouteMetadataContextNamespaces[iNdEx])
			copy(dAtA[i:], m.RouteMetadataContextNamespaces[iNdEx])
			i = encodeVarint(dAtA, i, uint64(len(m.RouteMetadataContextNamespaces[iNdEx])))
			i--
			dAtA[i] = 0x1
			i--
			dAtA[i] = 0xaa
		}
	}
	if m.ChargeClusterResponseStats != nil {
		size, err := (*wrapperspb.BoolValue)(m.ChargeClusterResponseStats).MarshalToSizedBufferVTStrict(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarint(dAtA, i, uint64(size))
		i--
		dAtA[i] = 0x1
		i--
		dAtA[i] = 0xa2
	}
	if m.FailureModeAllowHeaderAdd {
		i--
		if m.FailureModeAllowHeaderAdd {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x1
		i--
		dAtA[i] = 0x98
	}
	if m.IncludeTlsSession {
		i--
		if m.IncludeTlsSession {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x1
		i--
		dAtA[i] = 0x90
	}
	if m.AllowedHeaders != nil {
		if vtmsg, ok := interface{}(m.AllowedHeaders).(interface {
			MarshalToSizedBufferVTStrict([]byte) (int, error)
		}); ok {
			size, err := vtmsg.MarshalToSizedBufferVTStrict(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarint(dAtA, i, uint64(size))
		} else {
			encoded, err := proto.Marshal(m.AllowedHeaders)
			if err != nil {
				return 0, err
			}
			i -= len(encoded)
			copy(dAtA[i:], encoded)
			i = encodeVarint(dAtA, i, uint64(len(encoded)))
		}
		i--
		dAtA[i] = 0x1
		i--
		dAtA[i] = 0x8a
	}
	if len(m.TypedMetadataContextNamespaces) > 0 {
		for iNdEx := len(m.TypedMetadataContextNamespaces) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.TypedMetadataContextNamespaces[iNdEx])
			copy(dAtA[i:], m.TypedMetadataContextNamespaces[iNdEx])
			i = encodeVarint(dAtA, i, uint64(len(m.TypedMetadataContextNamespaces[iNdEx])))
			i--
			dAtA[i] = 0x1
			i--
			dAtA[i] = 0x82
		}
	}
	if len(m.BootstrapMetadataLabelsKey) > 0 {
		i -= len(m.BootstrapMetadataLabelsKey)
		copy(dAtA[i:], m.BootstrapMetadataLabelsKey)
		i = encodeVarint(dAtA, i, uint64(len(m.BootstrapMetadataLabelsKey)))
		i--
		dAtA[i] = 0x7a
	}
	if m.FilterEnabledMetadata != nil {
		if vtmsg, ok := interface{}(m.FilterEnabledMetadata).(interface {
			MarshalToSizedBufferVTStrict([]byte) (int, error)
		}); ok {
			size, err := vtmsg.MarshalToSizedBufferVTStrict(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarint(dAtA, i, uint64(size))
		} else {
			encoded, err := proto.Marshal(m.FilterEnabledMetadata)
			if err != nil {
				return 0, err
			}
			i -= len(encoded)
			copy(dAtA[i:], encoded)
			i = encodeVarint(dAtA, i, uint64(len(encoded)))
		}
		i--
		dAtA[i] = 0x72
	}
	if len(m.StatPrefix) > 0 {
		i -= len(m.StatPrefix)
		copy(dAtA[i:], m.StatPrefix)
		i = encodeVarint(dAtA, i, uint64(len(m.StatPrefix)))
		i--
		dAtA[i] = 0x6a
	}
	if m.TransportApiVersion != 0 {
		i = encodeVarint(dAtA, i, uint64(m.TransportApiVersion))
		i--
		dAtA[i] = 0x60
	}
	if m.DenyAtDisable != nil {
		if vtmsg, ok := interface{}(m.DenyAtDisable).(interface {
			MarshalToSizedBufferVTStrict([]byte) (int, error)
		}); ok {
			size, err := vtmsg.MarshalToSizedBufferVTStrict(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarint(dAtA, i, uint64(size))
		} else {
			encoded, err := proto.Marshal(m.DenyAtDisable)
			if err != nil {
				return 0, err
			}
			i -= len(encoded)
			copy(dAtA[i:], encoded)
			i = encodeVarint(dAtA, i, uint64(len(encoded)))
		}
		i--
		dAtA[i] = 0x5a
	}
	if m.IncludePeerCertificate {
		i--
		if m.IncludePeerCertificate {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x50
	}
	if m.FilterEnabled != nil {
		if vtmsg, ok := interface{}(m.FilterEnabled).(interface {
			MarshalToSizedBufferVTStrict([]byte) (int, error)
		}); ok {
			size, err := vtmsg.MarshalToSizedBufferVTStrict(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarint(dAtA, i, uint64(size))
		} else {
			encoded, err := proto.Marshal(m.FilterEnabled)
			if err != nil {
				return 0, err
			}
			i -= len(encoded)
			copy(dAtA[i:], encoded)
			i = encodeVarint(dAtA, i, uint64(len(encoded)))
		}
		i--
		dAtA[i] = 0x4a
	}
	if len(m.MetadataContextNamespaces) > 0 {
		for iNdEx := len(m.MetadataContextNamespaces) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.MetadataContextNamespaces[iNdEx])
			copy(dAtA[i:], m.MetadataContextNamespaces[iNdEx])
			i = encodeVarint(dAtA, i, uint64(len(m.MetadataContextNamespaces[iNdEx])))
			i--
			dAtA[i] = 0x42
		}
	}
	if m.StatusOnError != nil {
		if vtmsg, ok := interface{}(m.StatusOnError).(interface {
			MarshalToSizedBufferVTStrict([]byte) (int, error)
		}); ok {
			size, err := vtmsg.MarshalToSizedBufferVTStrict(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarint(dAtA, i, uint64(size))
		} else {
			encoded, err := proto.Marshal(m.StatusOnError)
			if err != nil {
				return 0, err
			}
			i -= len(encoded)
			copy(dAtA[i:], encoded)
			i = encodeVarint(dAtA, i, uint64(len(encoded)))
		}
		i--
		dAtA[i] = 0x3a
	}
	if m.ClearRouteCache {
		i--
		if m.ClearRouteCache {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x30
	}
	if m.WithRequestBody != nil {
		size, err := m.WithRequestBody.MarshalToSizedBufferVTStrict(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarint(dAtA, i, uint64(size))
		i--
		dAtA[i] = 0x2a
	}
	if msg, ok := m.Services.(*ExtAuthz_HttpService); ok {
		size, err := msg.MarshalToSizedBufferVTStrict(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
	}
	if m.FailureModeAllow {
		i--
		if m.FailureModeAllow {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x10
	}
	if msg, ok := m.Services.(*ExtAuthz_GrpcService); ok {
		size, err := msg.MarshalToSizedBufferVTStrict(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
	}
	return len(dAtA) - i, nil
}

func (m *ExtAuthz_GrpcService) MarshalToVTStrict(dAtA []byte) (int, error) {
	size := m.SizeVT()
	return m.MarshalToSizedBufferVTStrict(dAtA[:size])
}

func (m *ExtAuthz_GrpcService) MarshalToSizedBufferVTStrict(dAtA []byte) (int, error) {
	i := len(dAtA)
	if m.GrpcService != nil {
		if vtmsg, ok := interface{}(m.GrpcService).(interface {
			MarshalToSizedBufferVTStrict([]byte) (int, error)
		}); ok {
			size, err := vtmsg.MarshalToSizedBufferVTStrict(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarint(dAtA, i, uint64(size))
		} else {
			encoded, err := proto.Marshal(m.GrpcService)
			if err != nil {
				return 0, err
			}
			i -= len(encoded)
			copy(dAtA[i:], encoded)
			i = encodeVarint(dAtA, i, uint64(len(encoded)))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}
func (m *ExtAuthz_HttpService) MarshalToVTStrict(dAtA []byte) (int, error) {
	size := m.SizeVT()
	return m.MarshalToSizedBufferVTStrict(dAtA[:size])
}

func (m *ExtAuthz_HttpService) MarshalToSizedBufferVTStrict(dAtA []byte) (int, error) {
	i := len(dAtA)
	if m.HttpService != nil {
		size, err := m.HttpService.MarshalToSizedBufferVTStrict(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarint(dAtA, i, uint64(size))
		i--
		dAtA[i] = 0x1a
	}
	return len(dAtA) - i, nil
}
func (m *BufferSettings) MarshalVTStrict() (dAtA []byte, err error) {
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

func (m *BufferSettings) MarshalToVTStrict(dAtA []byte) (int, error) {
	size := m.SizeVT()
	return m.MarshalToSizedBufferVTStrict(dAtA[:size])
}

func (m *BufferSettings) MarshalToSizedBufferVTStrict(dAtA []byte) (int, error) {
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
	if m.PackAsBytes {
		i--
		if m.PackAsBytes {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x18
	}
	if m.AllowPartialMessage {
		i--
		if m.AllowPartialMessage {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x10
	}
	if m.MaxRequestBytes != 0 {
		i = encodeVarint(dAtA, i, uint64(m.MaxRequestBytes))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *HttpService) MarshalVTStrict() (dAtA []byte, err error) {
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

func (m *HttpService) MarshalToVTStrict(dAtA []byte) (int, error) {
	size := m.SizeVT()
	return m.MarshalToSizedBufferVTStrict(dAtA[:size])
}

func (m *HttpService) MarshalToSizedBufferVTStrict(dAtA []byte) (int, error) {
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
	if m.AuthorizationResponse != nil {
		size, err := m.AuthorizationResponse.MarshalToSizedBufferVTStrict(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarint(dAtA, i, uint64(size))
		i--
		dAtA[i] = 0x42
	}
	if m.AuthorizationRequest != nil {
		size, err := m.AuthorizationRequest.MarshalToSizedBufferVTStrict(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarint(dAtA, i, uint64(size))
		i--
		dAtA[i] = 0x3a
	}
	if len(m.PathPrefix) > 0 {
		i -= len(m.PathPrefix)
		copy(dAtA[i:], m.PathPrefix)
		i = encodeVarint(dAtA, i, uint64(len(m.PathPrefix)))
		i--
		dAtA[i] = 0x12
	}
	if m.ServerUri != nil {
		if vtmsg, ok := interface{}(m.ServerUri).(interface {
			MarshalToSizedBufferVTStrict([]byte) (int, error)
		}); ok {
			size, err := vtmsg.MarshalToSizedBufferVTStrict(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarint(dAtA, i, uint64(size))
		} else {
			encoded, err := proto.Marshal(m.ServerUri)
			if err != nil {
				return 0, err
			}
			i -= len(encoded)
			copy(dAtA[i:], encoded)
			i = encodeVarint(dAtA, i, uint64(len(encoded)))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *AuthorizationRequest) MarshalVTStrict() (dAtA []byte, err error) {
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

func (m *AuthorizationRequest) MarshalToVTStrict(dAtA []byte) (int, error) {
	size := m.SizeVT()
	return m.MarshalToSizedBufferVTStrict(dAtA[:size])
}

func (m *AuthorizationRequest) MarshalToSizedBufferVTStrict(dAtA []byte) (int, error) {
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
	if len(m.HeadersToAdd) > 0 {
		for iNdEx := len(m.HeadersToAdd) - 1; iNdEx >= 0; iNdEx-- {
			if vtmsg, ok := interface{}(m.HeadersToAdd[iNdEx]).(interface {
				MarshalToSizedBufferVTStrict([]byte) (int, error)
			}); ok {
				size, err := vtmsg.MarshalToSizedBufferVTStrict(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarint(dAtA, i, uint64(size))
			} else {
				encoded, err := proto.Marshal(m.HeadersToAdd[iNdEx])
				if err != nil {
					return 0, err
				}
				i -= len(encoded)
				copy(dAtA[i:], encoded)
				i = encodeVarint(dAtA, i, uint64(len(encoded)))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	if m.AllowedHeaders != nil {
		if vtmsg, ok := interface{}(m.AllowedHeaders).(interface {
			MarshalToSizedBufferVTStrict([]byte) (int, error)
		}); ok {
			size, err := vtmsg.MarshalToSizedBufferVTStrict(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarint(dAtA, i, uint64(size))
		} else {
			encoded, err := proto.Marshal(m.AllowedHeaders)
			if err != nil {
				return 0, err
			}
			i -= len(encoded)
			copy(dAtA[i:], encoded)
			i = encodeVarint(dAtA, i, uint64(len(encoded)))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *AuthorizationResponse) MarshalVTStrict() (dAtA []byte, err error) {
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

func (m *AuthorizationResponse) MarshalToVTStrict(dAtA []byte) (int, error) {
	size := m.SizeVT()
	return m.MarshalToSizedBufferVTStrict(dAtA[:size])
}

func (m *AuthorizationResponse) MarshalToSizedBufferVTStrict(dAtA []byte) (int, error) {
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
	if m.DynamicMetadataFromHeaders != nil {
		if vtmsg, ok := interface{}(m.DynamicMetadataFromHeaders).(interface {
			MarshalToSizedBufferVTStrict([]byte) (int, error)
		}); ok {
			size, err := vtmsg.MarshalToSizedBufferVTStrict(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarint(dAtA, i, uint64(size))
		} else {
			encoded, err := proto.Marshal(m.DynamicMetadataFromHeaders)
			if err != nil {
				return 0, err
			}
			i -= len(encoded)
			copy(dAtA[i:], encoded)
			i = encodeVarint(dAtA, i, uint64(len(encoded)))
		}
		i--
		dAtA[i] = 0x2a
	}
	if m.AllowedClientHeadersOnSuccess != nil {
		if vtmsg, ok := interface{}(m.AllowedClientHeadersOnSuccess).(interface {
			MarshalToSizedBufferVTStrict([]byte) (int, error)
		}); ok {
			size, err := vtmsg.MarshalToSizedBufferVTStrict(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarint(dAtA, i, uint64(size))
		} else {
			encoded, err := proto.Marshal(m.AllowedClientHeadersOnSuccess)
			if err != nil {
				return 0, err
			}
			i -= len(encoded)
			copy(dAtA[i:], encoded)
			i = encodeVarint(dAtA, i, uint64(len(encoded)))
		}
		i--
		dAtA[i] = 0x22
	}
	if m.AllowedUpstreamHeadersToAppend != nil {
		if vtmsg, ok := interface{}(m.AllowedUpstreamHeadersToAppend).(interface {
			MarshalToSizedBufferVTStrict([]byte) (int, error)
		}); ok {
			size, err := vtmsg.MarshalToSizedBufferVTStrict(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarint(dAtA, i, uint64(size))
		} else {
			encoded, err := proto.Marshal(m.AllowedUpstreamHeadersToAppend)
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
	if m.AllowedClientHeaders != nil {
		if vtmsg, ok := interface{}(m.AllowedClientHeaders).(interface {
			MarshalToSizedBufferVTStrict([]byte) (int, error)
		}); ok {
			size, err := vtmsg.MarshalToSizedBufferVTStrict(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarint(dAtA, i, uint64(size))
		} else {
			encoded, err := proto.Marshal(m.AllowedClientHeaders)
			if err != nil {
				return 0, err
			}
			i -= len(encoded)
			copy(dAtA[i:], encoded)
			i = encodeVarint(dAtA, i, uint64(len(encoded)))
		}
		i--
		dAtA[i] = 0x12
	}
	if m.AllowedUpstreamHeaders != nil {
		if vtmsg, ok := interface{}(m.AllowedUpstreamHeaders).(interface {
			MarshalToSizedBufferVTStrict([]byte) (int, error)
		}); ok {
			size, err := vtmsg.MarshalToSizedBufferVTStrict(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarint(dAtA, i, uint64(size))
		} else {
			encoded, err := proto.Marshal(m.AllowedUpstreamHeaders)
			if err != nil {
				return 0, err
			}
			i -= len(encoded)
			copy(dAtA[i:], encoded)
			i = encodeVarint(dAtA, i, uint64(len(encoded)))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *ExtAuthzPerRoute) MarshalVTStrict() (dAtA []byte, err error) {
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

func (m *ExtAuthzPerRoute) MarshalToVTStrict(dAtA []byte) (int, error) {
	size := m.SizeVT()
	return m.MarshalToSizedBufferVTStrict(dAtA[:size])
}

func (m *ExtAuthzPerRoute) MarshalToSizedBufferVTStrict(dAtA []byte) (int, error) {
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
	if msg, ok := m.Override.(*ExtAuthzPerRoute_CheckSettings); ok {
		size, err := msg.MarshalToSizedBufferVTStrict(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
	}
	if msg, ok := m.Override.(*ExtAuthzPerRoute_Disabled); ok {
		size, err := msg.MarshalToSizedBufferVTStrict(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
	}
	return len(dAtA) - i, nil
}

func (m *ExtAuthzPerRoute_Disabled) MarshalToVTStrict(dAtA []byte) (int, error) {
	size := m.SizeVT()
	return m.MarshalToSizedBufferVTStrict(dAtA[:size])
}

func (m *ExtAuthzPerRoute_Disabled) MarshalToSizedBufferVTStrict(dAtA []byte) (int, error) {
	i := len(dAtA)
	i--
	if m.Disabled {
		dAtA[i] = 1
	} else {
		dAtA[i] = 0
	}
	i--
	dAtA[i] = 0x8
	return len(dAtA) - i, nil
}
func (m *ExtAuthzPerRoute_CheckSettings) MarshalToVTStrict(dAtA []byte) (int, error) {
	size := m.SizeVT()
	return m.MarshalToSizedBufferVTStrict(dAtA[:size])
}

func (m *ExtAuthzPerRoute_CheckSettings) MarshalToSizedBufferVTStrict(dAtA []byte) (int, error) {
	i := len(dAtA)
	if m.CheckSettings != nil {
		size, err := m.CheckSettings.MarshalToSizedBufferVTStrict(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarint(dAtA, i, uint64(size))
		i--
		dAtA[i] = 0x12
	}
	return len(dAtA) - i, nil
}
func (m *CheckSettings) MarshalVTStrict() (dAtA []byte, err error) {
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

func (m *CheckSettings) MarshalToVTStrict(dAtA []byte) (int, error) {
	size := m.SizeVT()
	return m.MarshalToSizedBufferVTStrict(dAtA[:size])
}

func (m *CheckSettings) MarshalToSizedBufferVTStrict(dAtA []byte) (int, error) {
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
	if m.WithRequestBody != nil {
		size, err := m.WithRequestBody.MarshalToSizedBufferVTStrict(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarint(dAtA, i, uint64(size))
		i--
		dAtA[i] = 0x1a
	}
	if m.DisableRequestBodyBuffering {
		i--
		if m.DisableRequestBodyBuffering {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x10
	}
	if len(m.ContextExtensions) > 0 {
		for k := range m.ContextExtensions {
			v := m.ContextExtensions[k]
			baseI := i
			i -= len(v)
			copy(dAtA[i:], v)
			i = encodeVarint(dAtA, i, uint64(len(v)))
			i--
			dAtA[i] = 0x12
			i -= len(k)
			copy(dAtA[i:], k)
			i = encodeVarint(dAtA, i, uint64(len(k)))
			i--
			dAtA[i] = 0xa
			i = encodeVarint(dAtA, i, uint64(baseI-i))
			i--
			dAtA[i] = 0xa
		}
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
func (m *ExtAuthz) SizeVT() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if vtmsg, ok := m.Services.(interface{ SizeVT() int }); ok {
		n += vtmsg.SizeVT()
	}
	if m.FailureModeAllow {
		n += 2
	}
	if m.WithRequestBody != nil {
		l = m.WithRequestBody.SizeVT()
		n += 1 + l + sov(uint64(l))
	}
	if m.ClearRouteCache {
		n += 2
	}
	if m.StatusOnError != nil {
		if size, ok := interface{}(m.StatusOnError).(interface {
			SizeVT() int
		}); ok {
			l = size.SizeVT()
		} else {
			l = proto.Size(m.StatusOnError)
		}
		n += 1 + l + sov(uint64(l))
	}
	if len(m.MetadataContextNamespaces) > 0 {
		for _, s := range m.MetadataContextNamespaces {
			l = len(s)
			n += 1 + l + sov(uint64(l))
		}
	}
	if m.FilterEnabled != nil {
		if size, ok := interface{}(m.FilterEnabled).(interface {
			SizeVT() int
		}); ok {
			l = size.SizeVT()
		} else {
			l = proto.Size(m.FilterEnabled)
		}
		n += 1 + l + sov(uint64(l))
	}
	if m.IncludePeerCertificate {
		n += 2
	}
	if m.DenyAtDisable != nil {
		if size, ok := interface{}(m.DenyAtDisable).(interface {
			SizeVT() int
		}); ok {
			l = size.SizeVT()
		} else {
			l = proto.Size(m.DenyAtDisable)
		}
		n += 1 + l + sov(uint64(l))
	}
	if m.TransportApiVersion != 0 {
		n += 1 + sov(uint64(m.TransportApiVersion))
	}
	l = len(m.StatPrefix)
	if l > 0 {
		n += 1 + l + sov(uint64(l))
	}
	if m.FilterEnabledMetadata != nil {
		if size, ok := interface{}(m.FilterEnabledMetadata).(interface {
			SizeVT() int
		}); ok {
			l = size.SizeVT()
		} else {
			l = proto.Size(m.FilterEnabledMetadata)
		}
		n += 1 + l + sov(uint64(l))
	}
	l = len(m.BootstrapMetadataLabelsKey)
	if l > 0 {
		n += 1 + l + sov(uint64(l))
	}
	if len(m.TypedMetadataContextNamespaces) > 0 {
		for _, s := range m.TypedMetadataContextNamespaces {
			l = len(s)
			n += 2 + l + sov(uint64(l))
		}
	}
	if m.AllowedHeaders != nil {
		if size, ok := interface{}(m.AllowedHeaders).(interface {
			SizeVT() int
		}); ok {
			l = size.SizeVT()
		} else {
			l = proto.Size(m.AllowedHeaders)
		}
		n += 2 + l + sov(uint64(l))
	}
	if m.IncludeTlsSession {
		n += 3
	}
	if m.FailureModeAllowHeaderAdd {
		n += 3
	}
	if m.ChargeClusterResponseStats != nil {
		l = (*wrapperspb.BoolValue)(m.ChargeClusterResponseStats).SizeVT()
		n += 2 + l + sov(uint64(l))
	}
	if len(m.RouteMetadataContextNamespaces) > 0 {
		for _, s := range m.RouteMetadataContextNamespaces {
			l = len(s)
			n += 2 + l + sov(uint64(l))
		}
	}
	if len(m.RouteTypedMetadataContextNamespaces) > 0 {
		for _, s := range m.RouteTypedMetadataContextNamespaces {
			l = len(s)
			n += 2 + l + sov(uint64(l))
		}
	}
	if m.EncodeRawHeaders {
		n += 3
	}
	n += len(m.unknownFields)
	return n
}

func (m *ExtAuthz_GrpcService) SizeVT() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.GrpcService != nil {
		if size, ok := interface{}(m.GrpcService).(interface {
			SizeVT() int
		}); ok {
			l = size.SizeVT()
		} else {
			l = proto.Size(m.GrpcService)
		}
		n += 1 + l + sov(uint64(l))
	}
	return n
}
func (m *ExtAuthz_HttpService) SizeVT() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.HttpService != nil {
		l = m.HttpService.SizeVT()
		n += 1 + l + sov(uint64(l))
	}
	return n
}
func (m *BufferSettings) SizeVT() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.MaxRequestBytes != 0 {
		n += 1 + sov(uint64(m.MaxRequestBytes))
	}
	if m.AllowPartialMessage {
		n += 2
	}
	if m.PackAsBytes {
		n += 2
	}
	n += len(m.unknownFields)
	return n
}

func (m *HttpService) SizeVT() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.ServerUri != nil {
		if size, ok := interface{}(m.ServerUri).(interface {
			SizeVT() int
		}); ok {
			l = size.SizeVT()
		} else {
			l = proto.Size(m.ServerUri)
		}
		n += 1 + l + sov(uint64(l))
	}
	l = len(m.PathPrefix)
	if l > 0 {
		n += 1 + l + sov(uint64(l))
	}
	if m.AuthorizationRequest != nil {
		l = m.AuthorizationRequest.SizeVT()
		n += 1 + l + sov(uint64(l))
	}
	if m.AuthorizationResponse != nil {
		l = m.AuthorizationResponse.SizeVT()
		n += 1 + l + sov(uint64(l))
	}
	n += len(m.unknownFields)
	return n
}

func (m *AuthorizationRequest) SizeVT() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.AllowedHeaders != nil {
		if size, ok := interface{}(m.AllowedHeaders).(interface {
			SizeVT() int
		}); ok {
			l = size.SizeVT()
		} else {
			l = proto.Size(m.AllowedHeaders)
		}
		n += 1 + l + sov(uint64(l))
	}
	if len(m.HeadersToAdd) > 0 {
		for _, e := range m.HeadersToAdd {
			if size, ok := interface{}(e).(interface {
				SizeVT() int
			}); ok {
				l = size.SizeVT()
			} else {
				l = proto.Size(e)
			}
			n += 1 + l + sov(uint64(l))
		}
	}
	n += len(m.unknownFields)
	return n
}

func (m *AuthorizationResponse) SizeVT() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.AllowedUpstreamHeaders != nil {
		if size, ok := interface{}(m.AllowedUpstreamHeaders).(interface {
			SizeVT() int
		}); ok {
			l = size.SizeVT()
		} else {
			l = proto.Size(m.AllowedUpstreamHeaders)
		}
		n += 1 + l + sov(uint64(l))
	}
	if m.AllowedClientHeaders != nil {
		if size, ok := interface{}(m.AllowedClientHeaders).(interface {
			SizeVT() int
		}); ok {
			l = size.SizeVT()
		} else {
			l = proto.Size(m.AllowedClientHeaders)
		}
		n += 1 + l + sov(uint64(l))
	}
	if m.AllowedUpstreamHeadersToAppend != nil {
		if size, ok := interface{}(m.AllowedUpstreamHeadersToAppend).(interface {
			SizeVT() int
		}); ok {
			l = size.SizeVT()
		} else {
			l = proto.Size(m.AllowedUpstreamHeadersToAppend)
		}
		n += 1 + l + sov(uint64(l))
	}
	if m.AllowedClientHeadersOnSuccess != nil {
		if size, ok := interface{}(m.AllowedClientHeadersOnSuccess).(interface {
			SizeVT() int
		}); ok {
			l = size.SizeVT()
		} else {
			l = proto.Size(m.AllowedClientHeadersOnSuccess)
		}
		n += 1 + l + sov(uint64(l))
	}
	if m.DynamicMetadataFromHeaders != nil {
		if size, ok := interface{}(m.DynamicMetadataFromHeaders).(interface {
			SizeVT() int
		}); ok {
			l = size.SizeVT()
		} else {
			l = proto.Size(m.DynamicMetadataFromHeaders)
		}
		n += 1 + l + sov(uint64(l))
	}
	n += len(m.unknownFields)
	return n
}

func (m *ExtAuthzPerRoute) SizeVT() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if vtmsg, ok := m.Override.(interface{ SizeVT() int }); ok {
		n += vtmsg.SizeVT()
	}
	n += len(m.unknownFields)
	return n
}

func (m *ExtAuthzPerRoute_Disabled) SizeVT() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	n += 2
	return n
}
func (m *ExtAuthzPerRoute_CheckSettings) SizeVT() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.CheckSettings != nil {
		l = m.CheckSettings.SizeVT()
		n += 1 + l + sov(uint64(l))
	}
	return n
}
func (m *CheckSettings) SizeVT() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.ContextExtensions) > 0 {
		for k, v := range m.ContextExtensions {
			_ = k
			_ = v
			mapEntrySize := 1 + len(k) + sov(uint64(len(k))) + 1 + len(v) + sov(uint64(len(v)))
			n += mapEntrySize + 1 + sov(uint64(mapEntrySize))
		}
	}
	if m.DisableRequestBodyBuffering {
		n += 2
	}
	if m.WithRequestBody != nil {
		l = m.WithRequestBody.SizeVT()
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
