package main

import (
	"context"
	"fmt"

	"github.com/icexin/brpc-go"
	"github.com/icexin/brpc-go/examples/echo"
)

func main() {
	clientConn, err := brpc.Dial(context.Background(), "127.0.0.1:8080")
	if err != nil {
		panic(err)
	}
	client := echo.NewEchoServerClient(clientConn)
	resp, err := client.Echo(context.Background(), &echo.EchoRequest{
		Message: "hello",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(resp.Message)
}
