module github.com/icexin/brpc-go

go 1.18

require (
	github.com/golang/snappy v0.0.4
	github.com/keegancsmith/rpc v1.3.0
	github.com/pierrec/lz4 v2.6.1+incompatible
	google.golang.org/grpc v1.45.0
	google.golang.org/protobuf v1.28.0
)

require (
	github.com/frankban/quicktest v1.14.3 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/go-cmp v0.5.7 // indirect
	golang.org/x/net v0.0.0-20200822124328-c89045814202 // indirect
	golang.org/x/sys v0.0.0-20200323222414-85ca7c5b95cd // indirect
	golang.org/x/text v0.3.0 // indirect
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013 // indirect
)

replace github.com/baidu/sofa-pbrpc v1.1.4 => github.com/icexin/sofa-pbrpc v1.1.4-0.20170426051859-97df346b6e46 // indirect
