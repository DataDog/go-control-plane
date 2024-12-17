// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v5.29.1
// source: envoy/service/extension/v3/config_discovery.proto

package extensionv3

import (
	context "context"
	v3 "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	ExtensionConfigDiscoveryService_StreamExtensionConfigs_FullMethodName = "/envoy.service.extension.v3.ExtensionConfigDiscoveryService/StreamExtensionConfigs"
	ExtensionConfigDiscoveryService_DeltaExtensionConfigs_FullMethodName  = "/envoy.service.extension.v3.ExtensionConfigDiscoveryService/DeltaExtensionConfigs"
	ExtensionConfigDiscoveryService_FetchExtensionConfigs_FullMethodName  = "/envoy.service.extension.v3.ExtensionConfigDiscoveryService/FetchExtensionConfigs"
)

// ExtensionConfigDiscoveryServiceClient is the client API for ExtensionConfigDiscoveryService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ExtensionConfigDiscoveryServiceClient interface {
	StreamExtensionConfigs(ctx context.Context, opts ...grpc.CallOption) (ExtensionConfigDiscoveryService_StreamExtensionConfigsClient, error)
	DeltaExtensionConfigs(ctx context.Context, opts ...grpc.CallOption) (ExtensionConfigDiscoveryService_DeltaExtensionConfigsClient, error)
	FetchExtensionConfigs(ctx context.Context, in *v3.DiscoveryRequest, opts ...grpc.CallOption) (*v3.DiscoveryResponse, error)
}

type extensionConfigDiscoveryServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewExtensionConfigDiscoveryServiceClient(cc grpc.ClientConnInterface) ExtensionConfigDiscoveryServiceClient {
	return &extensionConfigDiscoveryServiceClient{cc}
}

func (c *extensionConfigDiscoveryServiceClient) StreamExtensionConfigs(ctx context.Context, opts ...grpc.CallOption) (ExtensionConfigDiscoveryService_StreamExtensionConfigsClient, error) {
	stream, err := c.cc.NewStream(ctx, &ExtensionConfigDiscoveryService_ServiceDesc.Streams[0], ExtensionConfigDiscoveryService_StreamExtensionConfigs_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &extensionConfigDiscoveryServiceStreamExtensionConfigsClient{stream}
	return x, nil
}

type ExtensionConfigDiscoveryService_StreamExtensionConfigsClient interface {
	Send(*v3.DiscoveryRequest) error
	Recv() (*v3.DiscoveryResponse, error)
	grpc.ClientStream
}

type extensionConfigDiscoveryServiceStreamExtensionConfigsClient struct {
	grpc.ClientStream
}

func (x *extensionConfigDiscoveryServiceStreamExtensionConfigsClient) Send(m *v3.DiscoveryRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *extensionConfigDiscoveryServiceStreamExtensionConfigsClient) Recv() (*v3.DiscoveryResponse, error) {
	m := new(v3.DiscoveryResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *extensionConfigDiscoveryServiceClient) DeltaExtensionConfigs(ctx context.Context, opts ...grpc.CallOption) (ExtensionConfigDiscoveryService_DeltaExtensionConfigsClient, error) {
	stream, err := c.cc.NewStream(ctx, &ExtensionConfigDiscoveryService_ServiceDesc.Streams[1], ExtensionConfigDiscoveryService_DeltaExtensionConfigs_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &extensionConfigDiscoveryServiceDeltaExtensionConfigsClient{stream}
	return x, nil
}

type ExtensionConfigDiscoveryService_DeltaExtensionConfigsClient interface {
	Send(*v3.DeltaDiscoveryRequest) error
	Recv() (*v3.DeltaDiscoveryResponse, error)
	grpc.ClientStream
}

type extensionConfigDiscoveryServiceDeltaExtensionConfigsClient struct {
	grpc.ClientStream
}

func (x *extensionConfigDiscoveryServiceDeltaExtensionConfigsClient) Send(m *v3.DeltaDiscoveryRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *extensionConfigDiscoveryServiceDeltaExtensionConfigsClient) Recv() (*v3.DeltaDiscoveryResponse, error) {
	m := new(v3.DeltaDiscoveryResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *extensionConfigDiscoveryServiceClient) FetchExtensionConfigs(ctx context.Context, in *v3.DiscoveryRequest, opts ...grpc.CallOption) (*v3.DiscoveryResponse, error) {
	out := new(v3.DiscoveryResponse)
	err := c.cc.Invoke(ctx, ExtensionConfigDiscoveryService_FetchExtensionConfigs_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ExtensionConfigDiscoveryServiceServer is the server API for ExtensionConfigDiscoveryService service.
// All implementations should embed UnimplementedExtensionConfigDiscoveryServiceServer
// for forward compatibility
type ExtensionConfigDiscoveryServiceServer interface {
	StreamExtensionConfigs(ExtensionConfigDiscoveryService_StreamExtensionConfigsServer) error
	DeltaExtensionConfigs(ExtensionConfigDiscoveryService_DeltaExtensionConfigsServer) error
	FetchExtensionConfigs(context.Context, *v3.DiscoveryRequest) (*v3.DiscoveryResponse, error)
}

// UnimplementedExtensionConfigDiscoveryServiceServer should be embedded to have forward compatible implementations.
type UnimplementedExtensionConfigDiscoveryServiceServer struct {
}

func (UnimplementedExtensionConfigDiscoveryServiceServer) StreamExtensionConfigs(ExtensionConfigDiscoveryService_StreamExtensionConfigsServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamExtensionConfigs not implemented")
}
func (UnimplementedExtensionConfigDiscoveryServiceServer) DeltaExtensionConfigs(ExtensionConfigDiscoveryService_DeltaExtensionConfigsServer) error {
	return status.Errorf(codes.Unimplemented, "method DeltaExtensionConfigs not implemented")
}
func (UnimplementedExtensionConfigDiscoveryServiceServer) FetchExtensionConfigs(context.Context, *v3.DiscoveryRequest) (*v3.DiscoveryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FetchExtensionConfigs not implemented")
}

// UnsafeExtensionConfigDiscoveryServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ExtensionConfigDiscoveryServiceServer will
// result in compilation errors.
type UnsafeExtensionConfigDiscoveryServiceServer interface {
	mustEmbedUnimplementedExtensionConfigDiscoveryServiceServer()
}

func RegisterExtensionConfigDiscoveryServiceServer(s grpc.ServiceRegistrar, srv ExtensionConfigDiscoveryServiceServer) {
	s.RegisterService(&ExtensionConfigDiscoveryService_ServiceDesc, srv)
}

func _ExtensionConfigDiscoveryService_StreamExtensionConfigs_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ExtensionConfigDiscoveryServiceServer).StreamExtensionConfigs(&extensionConfigDiscoveryServiceStreamExtensionConfigsServer{stream})
}

type ExtensionConfigDiscoveryService_StreamExtensionConfigsServer interface {
	Send(*v3.DiscoveryResponse) error
	Recv() (*v3.DiscoveryRequest, error)
	grpc.ServerStream
}

type extensionConfigDiscoveryServiceStreamExtensionConfigsServer struct {
	grpc.ServerStream
}

func (x *extensionConfigDiscoveryServiceStreamExtensionConfigsServer) Send(m *v3.DiscoveryResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *extensionConfigDiscoveryServiceStreamExtensionConfigsServer) Recv() (*v3.DiscoveryRequest, error) {
	m := new(v3.DiscoveryRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _ExtensionConfigDiscoveryService_DeltaExtensionConfigs_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ExtensionConfigDiscoveryServiceServer).DeltaExtensionConfigs(&extensionConfigDiscoveryServiceDeltaExtensionConfigsServer{stream})
}

type ExtensionConfigDiscoveryService_DeltaExtensionConfigsServer interface {
	Send(*v3.DeltaDiscoveryResponse) error
	Recv() (*v3.DeltaDiscoveryRequest, error)
	grpc.ServerStream
}

type extensionConfigDiscoveryServiceDeltaExtensionConfigsServer struct {
	grpc.ServerStream
}

func (x *extensionConfigDiscoveryServiceDeltaExtensionConfigsServer) Send(m *v3.DeltaDiscoveryResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *extensionConfigDiscoveryServiceDeltaExtensionConfigsServer) Recv() (*v3.DeltaDiscoveryRequest, error) {
	m := new(v3.DeltaDiscoveryRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _ExtensionConfigDiscoveryService_FetchExtensionConfigs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(v3.DiscoveryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExtensionConfigDiscoveryServiceServer).FetchExtensionConfigs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ExtensionConfigDiscoveryService_FetchExtensionConfigs_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExtensionConfigDiscoveryServiceServer).FetchExtensionConfigs(ctx, req.(*v3.DiscoveryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ExtensionConfigDiscoveryService_ServiceDesc is the grpc.ServiceDesc for ExtensionConfigDiscoveryService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ExtensionConfigDiscoveryService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "envoy.service.extension.v3.ExtensionConfigDiscoveryService",
	HandlerType: (*ExtensionConfigDiscoveryServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "FetchExtensionConfigs",
			Handler:    _ExtensionConfigDiscoveryService_FetchExtensionConfigs_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "StreamExtensionConfigs",
			Handler:       _ExtensionConfigDiscoveryService_StreamExtensionConfigs_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "DeltaExtensionConfigs",
			Handler:       _ExtensionConfigDiscoveryService_DeltaExtensionConfigs_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "envoy/service/extension/v3/config_discovery.proto",
}
