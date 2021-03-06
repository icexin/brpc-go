// Code generated by protoc-gen-go-brpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-brpc v1.2.0
// - protoc             v3.19.4
// source: echo.proto

package echo

import (
	context "context"
	brpc_go "github.com/icexin/brpc-go"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// EchoServerClient is the client API for EchoServer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type EchoServerClient interface {
	Echo(ctx context.Context, in *EchoRequest, opts ...brpc_go.CallOption) (*EchoResponse, error)
}

type echoServerClient struct {
	cc brpc_go.ClientConn
}

func NewEchoServerClient(cc brpc_go.ClientConn) EchoServerClient {
	return &echoServerClient{cc}
}

func (c *echoServerClient) Echo(ctx context.Context, in *EchoRequest, opts ...brpc_go.CallOption) (*EchoResponse, error) {
	out := new(EchoResponse)
	err := c.cc.Invoke(ctx, "/brpc.test.EchoServer/Echo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EchoServerServer is the server API for EchoServer service.
// All implementations should embed UnimplementedEchoServerServer
// for forward compatibility
type EchoServerServer interface {
	Echo(context.Context, *EchoRequest) (*EchoResponse, error)
}

// UnimplementedEchoServerServer should be embedded to have forward compatible implementations.
type UnimplementedEchoServerServer struct {
}

func (UnimplementedEchoServerServer) Echo(context.Context, *EchoRequest) (*EchoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Echo not implemented")
}

// wrapperEchoServerServer is a wrapper for EchoServerServer that implements the interface of net/rpc Service
type wrapperEchoServerServer struct {
	svr EchoServerServer
}

func (w *wrapperEchoServerServer) Echo(ctx context.Context, req *EchoRequest, resp *EchoResponse) error {
	res, err := w.svr.Echo(ctx, req)
	if err != nil {
		return err
	}
	*resp = *res
	return nil
}

// UnsafeEchoServerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to EchoServerServer will
// result in compilation errors.
type UnsafeEchoServerServer interface {
	mustEmbedUnimplementedEchoServerServer()
}

func RegisterEchoServerServer(s grpc.ServiceRegistrar, srv EchoServerServer) {
	s.RegisterService(&EchoServer_ServiceDesc, &wrapperEchoServerServer{srv})
}

func _EchoServer_Echo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EchoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EchoServerServer).Echo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/brpc.test.EchoServer/Echo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EchoServerServer).Echo(ctx, req.(*EchoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// EchoServer_ServiceDesc is the grpc.ServiceDesc for EchoServer service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var EchoServer_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "brpc.test.EchoServer",
	HandlerType: (*EchoServerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Echo",
			Handler:    _EchoServer_Echo_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "echo.proto",
}
