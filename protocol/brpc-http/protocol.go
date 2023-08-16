package bhttp

import (
	"github.com/icexin/brpc-go"
	"google.golang.org/grpc"
)

const (
	ProtocolName = "brpc-http"
)

type protocol struct {
}

func (p *protocol) Dial(target string, options ...brpc.DialOption) (grpc.ClientConnInterface, error) {
	return dial(target, options...)
}

func (p *protocol) NewServer(options ...brpc.ServerOption) brpc.Server {
	return newServer(options...)
}

func init() {
	brpc.RegisterProtocol(ProtocolName, &protocol{})
}
