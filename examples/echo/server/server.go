package main

import (
	"context"
	"net"

	"github.com/icexin/brpc-go"
	"github.com/icexin/brpc-go/examples/echo"
)

type echoService struct {
}

func (s *echoService) Echo(ctx context.Context, req *echo.EchoRequest) (*echo.EchoResponse, error) {
	return &echo.EchoResponse{
		Message: "reply: " + req.Message,
	}, nil
}

func main() {
	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	server := brpc.NewServer()
	echo.RegisterEchoServerServer(server, &echoService{})
	server.Serve(l)
}
