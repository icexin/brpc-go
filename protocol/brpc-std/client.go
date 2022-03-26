package bstd

import (
	"bytes"
	"context"
	"net"

	"github.com/icexin/brpc-go"
	"github.com/keegancsmith/rpc"
)

// clientConn represents a client connection to an RPC server.
type clientConn struct {
	c *rpc.Client
}

// Close tears down the ClientConn and all underlying connections.
func (c *clientConn) Close() error {
	return c.c.Close()
}

func grpcMethodToBrpcMethod(method string) string {
	methodbuf := []byte(method[1:])
	i := bytes.Index(methodbuf, []byte{'/'})
	methodbuf[i] = '.'
	return string(methodbuf)
}

// Invoke sends the RPC request on the wire and returns after response is received. Invoke is called by generated code. Also users can call Invoke directly when it is really needed in their use cases.
func (c *clientConn) Invoke(ctx context.Context, method string, args interface{}, reply interface{}, opts ...brpc.CallOption) error {
	brpcMethod := grpcMethodToBrpcMethod(method)
	return c.c.Call(ctx, brpcMethod, args, reply)
}

func dial(target string, options ...brpc.DialOption) (brpc.ClientConn, error) {
	conn, err := net.Dial("tcp", target)
	if err != nil {
		return nil, err
	}
	c := rpc.NewClientWithCodec(newClientCodec(conn))
	return &clientConn{c: c}, nil
}
