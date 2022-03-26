# Usage

- Install the protobuf compiler https://grpc.io/docs/protoc-installation/
- Install the Go protobuf plugin and brpc plugin

``` bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1
```

- Generate protobuf files

``` bash
protoc --go_out=. --go_opt=paths=source_relative --go-brpc_out=. --go-brpc_opt=paths=source_relative  *.proto
```

For how to build server and client code, see [server.go](./server/server.go) and [client.go](./client/client.go)