package main

import (
	"context"
	"net"

	"github.com/icexin/brpc-go"
	"github.com/icexin/brpc-go/examples/echo"
)

type echoService struct {
	echo.UnimplementedEchoServerServer
}

func (s *echoService) Echo(context.Context, *echo.EchoRequest) (*echo.EchoResponse, error) {
	return &echo.EchoResponse{
		Message: "hello",
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
