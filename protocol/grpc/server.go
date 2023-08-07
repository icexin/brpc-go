package grpc

import (
	"net"

	"google.golang.org/grpc"
)

type server struct {
	grpcServer *grpc.Server
}

func newServer() *server {
	return &server{
		grpcServer: grpc.NewServer(),
	}
}

// Serve accepts incoming connections on the listener l, creating a new ServerConn and service goroutine for each. The service goroutines read pbrpc requests and then call the registered handlers to reply to them. Serve returns when l.Accept fails with errors.
func (s *server) Serve(l net.Listener) error {
	return s.grpcServer.Serve(l)
}

func (s *server) RegisterService(sd *grpc.ServiceDesc, ss interface{}) {
	s.grpcServer.RegisterService(sd, ss)
}
