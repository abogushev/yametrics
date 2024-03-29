// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.11
// source: internal/protocol/proto/metrics.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// MetricsClient is the client API for Metrics service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MetricsClient interface {
	SaveMetrics(ctx context.Context, opts ...grpc.CallOption) (Metrics_SaveMetricsClient, error)
}

type metricsClient struct {
	cc grpc.ClientConnInterface
}

func NewMetricsClient(cc grpc.ClientConnInterface) MetricsClient {
	return &metricsClient{cc}
}

func (c *metricsClient) SaveMetrics(ctx context.Context, opts ...grpc.CallOption) (Metrics_SaveMetricsClient, error) {
	stream, err := c.cc.NewStream(ctx, &Metrics_ServiceDesc.Streams[0], "/yametrics.Metrics/SaveMetrics", opts...)
	if err != nil {
		return nil, err
	}
	x := &metricsSaveMetricsClient{stream}
	return x, nil
}

type Metrics_SaveMetricsClient interface {
	Send(*Metric) error
	CloseAndRecv() (*emptypb.Empty, error)
	grpc.ClientStream
}

type metricsSaveMetricsClient struct {
	grpc.ClientStream
}

func (x *metricsSaveMetricsClient) Send(m *Metric) error {
	return x.ClientStream.SendMsg(m)
}

func (x *metricsSaveMetricsClient) CloseAndRecv() (*emptypb.Empty, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(emptypb.Empty)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// MetricsServer is the server API for Metrics service.
// All implementations must embed UnimplementedMetricsServer
// for forward compatibility
type MetricsServer interface {
	SaveMetrics(Metrics_SaveMetricsServer) error
	mustEmbedUnimplementedMetricsServer()
}

// UnimplementedMetricsServer must be embedded to have forward compatible implementations.
type UnimplementedMetricsServer struct {
}

func (UnimplementedMetricsServer) SaveMetrics(Metrics_SaveMetricsServer) error {
	return status.Errorf(codes.Unimplemented, "method SaveMetrics not implemented")
}
func (UnimplementedMetricsServer) mustEmbedUnimplementedMetricsServer() {}

// UnsafeMetricsServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MetricsServer will
// result in compilation errors.
type UnsafeMetricsServer interface {
	mustEmbedUnimplementedMetricsServer()
}

func RegisterMetricsServer(s grpc.ServiceRegistrar, srv MetricsServer) {
	s.RegisterService(&Metrics_ServiceDesc, srv)
}

func _Metrics_SaveMetrics_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(MetricsServer).SaveMetrics(&metricsSaveMetricsServer{stream})
}

type Metrics_SaveMetricsServer interface {
	SendAndClose(*emptypb.Empty) error
	Recv() (*Metric, error)
	grpc.ServerStream
}

type metricsSaveMetricsServer struct {
	grpc.ServerStream
}

func (x *metricsSaveMetricsServer) SendAndClose(m *emptypb.Empty) error {
	return x.ServerStream.SendMsg(m)
}

func (x *metricsSaveMetricsServer) Recv() (*Metric, error) {
	m := new(Metric)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Metrics_ServiceDesc is the grpc.ServiceDesc for Metrics service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Metrics_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "yametrics.Metrics",
	HandlerType: (*MetricsServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "SaveMetrics",
			Handler:       _Metrics_SaveMetrics_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "internal/protocol/proto/metrics.proto",
}
