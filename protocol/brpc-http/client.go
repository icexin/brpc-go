package bhttp

import (
	"bytes"
	"context"
	"io"
	"net/http"

	"github.com/icexin/brpc-go"
	"google.golang.org/protobuf/proto"
)

type clientConn struct {
	target string
}

// Invoke performs a unary RPC and returns after the response is received
// into reply.
func (c *clientConn) Invoke(ctx context.Context, method string, args interface{}, reply interface{}, opts ...brpc.CallOption) error {
	argmsg := args.(proto.Message)
	replymsg := reply.(proto.Message)

	buf, err := proto.Marshal(argmsg)
	if err != nil {
		return err
	}
	url := c.target + "/" + method

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

	err = proto.Unmarshal(buf, replymsg)
	if err != nil {
		return err
	}
	return nil
}

func dial(target string, options ...brpc.DialOption) (brpc.ClientConn, error) {
	return &clientConn{
		target: target,
	}, nil
}
