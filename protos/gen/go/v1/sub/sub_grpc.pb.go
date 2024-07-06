// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.4.0
// - protoc             (unknown)
// source: v1/sub/sub.proto

package sub

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.62.0 or later.
const _ = grpc.SupportPackageIsVersion8

const (
	SubService_Unsubscribe_FullMethodName = "/sub.v1.SubService/Unsubscribe"
	SubService_Subscribe_FullMethodName   = "/sub.v1.SubService/Subscribe"
)

// SubServiceClient is the client API for SubService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SubServiceClient interface {
	Unsubscribe(ctx context.Context, in *UnsubscribeRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	Subscribe(ctx context.Context, in *SubscribeRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type subServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewSubServiceClient(cc grpc.ClientConnInterface) SubServiceClient {
	return &subServiceClient{cc}
}

func (c *subServiceClient) Unsubscribe(ctx context.Context, in *UnsubscribeRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, SubService_Unsubscribe_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *subServiceClient) Subscribe(ctx context.Context, in *SubscribeRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, SubService_Subscribe_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SubServiceServer is the server API for SubService service.
// All implementations must embed UnimplementedSubServiceServer
// for forward compatibility
type SubServiceServer interface {
	Unsubscribe(context.Context, *UnsubscribeRequest) (*emptypb.Empty, error)
	Subscribe(context.Context, *SubscribeRequest) (*emptypb.Empty, error)
	mustEmbedUnimplementedSubServiceServer()
}

// UnimplementedSubServiceServer must be embedded to have forward compatible implementations.
type UnimplementedSubServiceServer struct {
}

func (UnimplementedSubServiceServer) Unsubscribe(context.Context, *UnsubscribeRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Unsubscribe not implemented")
}
func (UnimplementedSubServiceServer) Subscribe(context.Context, *SubscribeRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Subscribe not implemented")
}
func (UnimplementedSubServiceServer) mustEmbedUnimplementedSubServiceServer() {}

// UnsafeSubServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SubServiceServer will
// result in compilation errors.
type UnsafeSubServiceServer interface {
	mustEmbedUnimplementedSubServiceServer()
}

func RegisterSubServiceServer(s grpc.ServiceRegistrar, srv SubServiceServer) {
	s.RegisterService(&SubService_ServiceDesc, srv)
}

func _SubService_Unsubscribe_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UnsubscribeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SubServiceServer).Unsubscribe(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SubService_Unsubscribe_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SubServiceServer).Unsubscribe(ctx, req.(*UnsubscribeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SubService_Subscribe_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SubscribeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SubServiceServer).Subscribe(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SubService_Subscribe_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SubServiceServer).Subscribe(ctx, req.(*SubscribeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// SubService_ServiceDesc is the grpc.ServiceDesc for SubService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SubService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "sub.v1.SubService",
	HandlerType: (*SubServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Unsubscribe",
			Handler:    _SubService_Unsubscribe_Handler,
		},
		{
			MethodName: "Subscribe",
			Handler:    _SubService_Subscribe_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "v1/sub/sub.proto",
}
