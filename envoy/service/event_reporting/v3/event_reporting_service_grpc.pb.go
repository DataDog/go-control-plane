// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v5.29.1
// source: envoy/service/event_reporting/v3/event_reporting_service.proto

package event_reportingv3

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	EventReportingService_StreamEvents_FullMethodName = "/envoy.service.event_reporting.v3.EventReportingService/StreamEvents"
)

// EventReportingServiceClient is the client API for EventReportingService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type EventReportingServiceClient interface {
	// Envoy will connect and send StreamEventsRequest messages forever.
	// The management server may send StreamEventsResponse to configure event stream. See below.
	// This API is designed for high throughput with the expectation that it might be lossy.
	StreamEvents(ctx context.Context, opts ...grpc.CallOption) (EventReportingService_StreamEventsClient, error)
}

type eventReportingServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewEventReportingServiceClient(cc grpc.ClientConnInterface) EventReportingServiceClient {
	return &eventReportingServiceClient{cc}
}

func (c *eventReportingServiceClient) StreamEvents(ctx context.Context, opts ...grpc.CallOption) (EventReportingService_StreamEventsClient, error) {
	stream, err := c.cc.NewStream(ctx, &EventReportingService_ServiceDesc.Streams[0], EventReportingService_StreamEvents_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &eventReportingServiceStreamEventsClient{stream}
	return x, nil
}

type EventReportingService_StreamEventsClient interface {
	Send(*StreamEventsRequest) error
	Recv() (*StreamEventsResponse, error)
	grpc.ClientStream
}

type eventReportingServiceStreamEventsClient struct {
	grpc.ClientStream
}

func (x *eventReportingServiceStreamEventsClient) Send(m *StreamEventsRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *eventReportingServiceStreamEventsClient) Recv() (*StreamEventsResponse, error) {
	m := new(StreamEventsResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// EventReportingServiceServer is the server API for EventReportingService service.
// All implementations should embed UnimplementedEventReportingServiceServer
// for forward compatibility
type EventReportingServiceServer interface {
	// Envoy will connect and send StreamEventsRequest messages forever.
	// The management server may send StreamEventsResponse to configure event stream. See below.
	// This API is designed for high throughput with the expectation that it might be lossy.
	StreamEvents(EventReportingService_StreamEventsServer) error
}

// UnimplementedEventReportingServiceServer should be embedded to have forward compatible implementations.
type UnimplementedEventReportingServiceServer struct {
}

func (UnimplementedEventReportingServiceServer) StreamEvents(EventReportingService_StreamEventsServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamEvents not implemented")
}

// UnsafeEventReportingServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to EventReportingServiceServer will
// result in compilation errors.
type UnsafeEventReportingServiceServer interface {
	mustEmbedUnimplementedEventReportingServiceServer()
}

func RegisterEventReportingServiceServer(s grpc.ServiceRegistrar, srv EventReportingServiceServer) {
	s.RegisterService(&EventReportingService_ServiceDesc, srv)
}

func _EventReportingService_StreamEvents_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(EventReportingServiceServer).StreamEvents(&eventReportingServiceStreamEventsServer{stream})
}

type EventReportingService_StreamEventsServer interface {
	Send(*StreamEventsResponse) error
	Recv() (*StreamEventsRequest, error)
	grpc.ServerStream
}

type eventReportingServiceStreamEventsServer struct {
	grpc.ServerStream
}

func (x *eventReportingServiceStreamEventsServer) Send(m *StreamEventsResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *eventReportingServiceStreamEventsServer) Recv() (*StreamEventsRequest, error) {
	m := new(StreamEventsRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// EventReportingService_ServiceDesc is the grpc.ServiceDesc for EventReportingService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var EventReportingService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "envoy.service.event_reporting.v3.EventReportingService",
	HandlerType: (*EventReportingServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "StreamEvents",
			Handler:       _EventReportingService_StreamEvents_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "envoy/service/event_reporting/v3/event_reporting_service.proto",
}
