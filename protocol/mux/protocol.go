package mux

import (
	"github.com/icexin/brpc-go"
	"google.golang.org/grpc"
)

const (
	ProtocolName = "brpc-mux"
)

type protocol struct {
}

func (p *protocol) Dial(target string, options ...brpc.DialOption) (grpc.ClientConnInterface, error) {
	panic("not implemented")
}

func (p *protocol) NewServer(options ...brpc.ServerOption) brpc.Server {
	return newServer()
}

func init() {
	brpc.RegisterProtocol(ProtocolName, &protocol{})
}
