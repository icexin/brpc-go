package mux

import "github.com/icexin/brpc-go"

const (
	ProtocolName = "brpc-mux"
)

type protocol struct {
}

func (p *protocol) Dial(target string, options ...brpc.DialOption) (brpc.ClientConn, error) {
	panic("not implemented")
}

func (p *protocol) NewServer(options ...brpc.ServerOption) brpc.Server {
	return newServer()
}

func init() {
	brpc.RegisterProtocol(ProtocolName, &protocol{})
}
