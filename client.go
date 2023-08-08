package brpc

import (
	"errors"

	"google.golang.org/grpc"
)

var (
	ErrProtocolNotFound = errors.New("protocol not found")
)

// Dial creates a client connection to the given target according to the protocol.
func Dial(protocol, target string, options ...DialOption) (grpc.ClientConnInterface, error) {
	proto := getProtocol(protocol)
	if proto == nil {
		return nil, ErrProtocolNotFound
	}
	return proto.Dial(target, options...)
}
