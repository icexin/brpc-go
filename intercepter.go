package brpc

import (
	"context"

	"google.golang.org/grpc"
)

func toUnaryHandler(info *grpc.UnaryServerInfo, handler grpc.UnaryHandler, inter grpc.UnaryServerInterceptor) grpc.UnaryHandler {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		return inter(ctx, req, info, handler)
	}
}

func ChainIntercepter(ints ...grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if len(ints) == 0 {
			return handler(ctx, req)
		}
		return ints[0](ctx, req, info, toUnaryHandler(info, handler, ChainIntercepter(ints[1:]...)))
	}
}

func WithInterceptor(interceptor grpc.UnaryServerInterceptor) ServerOption {
	return BServerOption(func(o *ServerOptions) {
		o.Interceptor = interceptor
	})
}
