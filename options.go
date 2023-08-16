package brpc

import "google.golang.org/grpc"

// DialOption configures how we set up the connection.
type DialOption interface{}

// CallOption configures how we call the server.
type CallOption interface{}

// ServerOption configures how we set up the server.
type ServerOption interface{}

type ServerOptions struct {
	Interceptor grpc.UnaryServerInterceptor
}

type BServerOption func(*ServerOptions)
