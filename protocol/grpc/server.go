package grpc

import (
	"net"

	"github.com/icexin/brpc-go"
	"google.golang.org/grpc"
)

type server struct {
	grpcServer *grpc.Server
}

func newServer(boptions ...brpc.ServerOption) *server {
	var options []grpc.ServerOption
	for _, o := range boptions {
		if opt, ok := o.(grpc.ServerOption); ok {
			options = append(options, opt)
		}
	}
	return &server{
		grpcServer: grpc.NewServer(options...),
	}
}

// Serve accepts incoming connections on the listener l, creating a new ServerConn and service goroutine for each. The service goroutines read pbrpc requests and then call the registered handlers to reply to them. Serve returns when l.Accept fails with errors.
func (s *server) Serve(l net.Listener) error {
	return s.grpcServer.Serve(l)
}

func (s *server) RegisterService(sd *grpc.ServiceDesc, ss interface{}) {
	s.grpcServer.RegisterService(sd, ss)
}
