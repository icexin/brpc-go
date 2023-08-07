package bhttp

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/icexin/brpc-go"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

type clientConn struct {
	target string
}

// Invoke performs a unary RPC and returns after the response is received
// into reply.
func (c *clientConn) Invoke(ctx context.Context, method string, args interface{}, reply interface{}, opts ...grpc.CallOption) error {
	argmsg := args.(proto.Message)
	replymsg := reply.(proto.Message)

	buf, err := proto.Marshal(argmsg)
	if err != nil {
		return err
	}
	url := c.target + method

	request, err := http.NewRequest("POST", url, bytes.NewReader(buf))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/proto")
	request = request.WithContext(ctx)
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	buf, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http error, status:%d, body:%s", resp.StatusCode, buf)
	}

	err = proto.Unmarshal(buf, replymsg)
	if err != nil {
		return err
	}
	return nil
}

func (c *clientConn) NewStream(ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	panic("not implemented")
}

func dial(target string, options ...brpc.DialOption) (grpc.ClientConnInterface, error) {
	return &clientConn{
		target: target,
	}, nil
}
