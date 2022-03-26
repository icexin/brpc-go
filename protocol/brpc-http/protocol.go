package bhttp

import "github.com/icexin/brpc-go"

const (
	ProtocolName = "brpc-http"
)

type protocol struct {
}

func (p *protocol) Dial(target string, options ...brpc.DialOption) (brpc.ClientConn, error) {
	return dial(target, options...)
}

func (p *protocol) NewServer(options ...brpc.ServerOption) brpc.Server {
	panic("not implemented")
}

func init() {
	brpc.RegisterProtocol(ProtocolName, &protocol{})
}
