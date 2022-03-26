# Usage

- Install the protobuf compiler https://grpc.io/docs/protoc-installation/
- Install the Go protobuf plugin and brpc plugin

``` bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26
go install github.com/icexin/brpc-go/protoc-gen-go-brpc@latest
```

- Generate protobuf files

``` bash
protoc --go_out=. --go_opt=paths=source_relative --go-brpc_out=. --go-brpc_opt=paths=source_relative  *.proto
```

For how to build server and client code, see [server.go](./server/server.go) and [client.go](./client/client.go)