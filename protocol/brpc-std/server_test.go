package bstd

import (
	"context"
	"io"
	"net"
	"testing"

	"github.com/icexin/brpc-go"
	"github.com/icexin/brpc-go/examples/echo"
)

type echoService struct {
	echo.UnimplementedEchoServerServer
}

func (s *echoService) Echo(ctx context.Context, req *echo.EchoRequest) (*echo.EchoResponse, error) {
	return &echo.EchoResponse{
		Message: "reply: " + req.Message,
	}, nil
}

func newTestServer() (func(), string) {
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}
	server := brpc.NewServer(ProtocolName)
	echo.RegisterEchoServerServer(server, &echoService{})
	go server.Serve(l)

	addr := l.Addr().String()
	return func() {
		l.Close()
	}, addr
}

func newTestClient(addr string) (func(), echo.EchoServerClient) {
	conn, err := brpc.Dial(ProtocolName, addr)
	if err != nil {
		panic(err)
	}
	closefunc := func() {
		conn.(io.Closer).Close()
	}
	return closefunc, echo.NewEchoServerClient(conn)
}

func TestServer(t *testing.T) {
	closefunc, addr := newTestServer()
	defer closefunc()

	closefunc, client := newTestClient(addr)
	defer closefunc()

	resp, err := client.Echo(context.Background(), &echo.EchoRequest{
		Message: "hello",
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Message != "reply: hello" {
		t.Fatalf("unexpected response: %s", resp.Message)
	}
}

func BenchmarkServer(b *testing.B) {
	closefunc, addr := newTestServer()
	defer closefunc()

	closefunc, client := newTestClient(addr)
	defer closefunc()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			resp, err := client.Echo(context.Background(), &echo.EchoRequest{
				Message: "hello",
			})
			if err != nil {
				b.Fatal(err)
			}
			if resp.Message != "reply: hello" {
				b.Fatalf("unexpected response: %s", resp.Message)
			}
		}
	})
}
