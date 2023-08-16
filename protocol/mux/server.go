package mux

import (
	"log"
	"net"
	"sync"

	"github.com/icexin/brpc-go"
	bhttp "github.com/icexin/brpc-go/protocol/brpc-http"
	bstd "github.com/icexin/brpc-go/protocol/brpc-std"
	bgrpc "github.com/icexin/brpc-go/protocol/grpc"
	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
)

type server struct {
	mux        cmux.CMux
	grpcServer brpc.Server
	brpcServer brpc.Server
	httpServer brpc.Server
}

func newServer(options ...brpc.ServerOption) *server {
	return &server{
		grpcServer: brpc.NewServer(bgrpc.ProtocolName, options...),
		brpcServer: brpc.NewServer(bstd.ProtocolName, options...),
		httpServer: brpc.NewServer(bhttp.ProtocolName, options...),
	}
}

func (s *server) Serve(l net.Listener) error {
	s.mux = cmux.New(l)
	// grpcL := s.mux.Match(cmux.HTTP2HeaderField("content-type", "application/grpc"))
	grpcL := s.mux.Match(cmux.HTTP2())
	httpL := s.mux.Match(cmux.HTTP1Fast())
	brpcL := s.mux.Match(cmux.PrefixMatcher(string(bstd.MagicStr[:])))
	wg := new(sync.WaitGroup)
	wg.Add(3)
	go func() {
		err := s.grpcServer.Serve(grpcL)
		if err != nil {
			log.Printf("grpc server error:%v", err)
		}
		wg.Done()
	}()
	go func() {
		err := s.httpServer.Serve(httpL)
		if err != nil {
			log.Printf("http server error:%v", err)
		}
		wg.Done()
	}()
	go func() {
		err := s.brpcServer.Serve(brpcL)
		if err != nil {
			log.Printf("brpc server error:%v", err)
		}
		wg.Done()
	}()
	return s.mux.Serve()
}

func (s *server) RegisterService(sd *grpc.ServiceDesc, srv interface{}) {
	s.grpcServer.RegisterService(sd, srv)
	s.brpcServer.RegisterService(sd, srv)
	s.httpServer.RegisterService(sd, srv)
}
