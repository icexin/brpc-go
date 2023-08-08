package brpc

import (
	"net"

	"google.golang.org/grpc"
)

var (
	protocols = map[string]Protocol{}
)

// type ServiceDesc = grpc.ServiceDesc

// type ClientConn interface {
// 	// Invoke performs a unary RPC and returns after the response is received
// 	// into reply.
// 	Invoke(ctx context.Context, method string, args interface{}, reply interface{}, opts ...CallOption) error
// }

// Server is the interface that must be implemented by a protocol server.
type Server interface {
	RegisterService(sd *grpc.ServiceDesc, ss interface{})
	Serve(l net.Listener) error
}

// Protocol defines how to make rpc call and serve rpc call.
type Protocol interface {
	Dial(target string, options ...DialOption) (grpc.ClientConnInterface, error)
	NewServer(options ...ServerOption) Server
}

// RegisterProtocol registers a protocol.
func RegisterProtocol(name string, p Protocol) {
	protocols[name] = p
}

// getProtocol returns a registered protocol.
func getProtocol(name string) Protocol {
	return protocols[name]
}
