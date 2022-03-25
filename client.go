package brpc

import (
	"bytes"
	"context"
	"errors"
	"net"

	"github.com/keegancsmith/rpc"
	"google.golang.org/grpc"
)

// ClientConn represents a client connection to an RPC server.
type ClientConn struct {
	c *rpc.Client
}

// Close tears down the ClientConn and all underlying connections.
func (c *ClientConn) Close() error {
	return c.c.Close()
}

func grpcMethodToBrpcMethod(method string) string {
	methodbuf := []byte(method[1:])
	i := bytes.Index(methodbuf, []byte{'/'})
	methodbuf[i] = '.'
	return string(methodbuf)
}

// Invoke sends the RPC request on the wire and returns after response is received. Invoke is called by generated code. Also users can call Invoke directly when it is really needed in their use cases.
func (c *ClientConn) Invoke(ctx context.Context, method string, args interface{}, reply interface{}, opts ...grpc.CallOption) error {
	brpcMethod := grpcMethodToBrpcMethod(method)
	return c.c.Call(ctx, brpcMethod, args, reply)
}

// NewStream begins a streaming RPC.
func (c *ClientConn) NewStream(ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	return nil, errors.New("unimplemented")
}

// Dial creates a client connection to the given target.
// The provided Context must be non-nil. If the context expires before the connection is complete, an error is returned. Once successfully connected, any expiration of the context will not affect the connection.
func Dial(ctx context.Context, addr string) (*ClientConn, error) {
	dialer := new(net.Dialer)
	conn, err := dialer.DialContext(ctx, "tcp", addr)
	if err != nil {
		return nil, err
	}

	return NewClientConn(conn), nil
}

// NewClientConn creates a ClientConn on a given connection
func NewClientConn(conn net.Conn) *ClientConn {
	c := rpc.NewClientWithCodec(newClientCodec(conn))
	return &ClientConn{c: c}
}
