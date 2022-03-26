package bstd

import (
	"net"

	"github.com/keegancsmith/rpc"
	"google.golang.org/grpc"
)

type server struct {
	s *rpc.Server
}

func newServer() *server {
	return &server{
		s: rpc.NewServer(),
	}
}

// ServeConn runs the server on a single connection. ServeConn blocks, serving the connection until the client hangs up. The caller typically invokes ServeConn in a go statement.
func (s *server) ServeConn(conn net.Conn) {
	codec := newServerCodec(conn)
	s.s.ServeCodec(codec)
}

// Serve accepts incoming connections on the listener l, creating a new ServerConn and service goroutine for each. The service goroutines read pbrpc requests and then call the registered handlers to reply to them. Serve returns when l.Accept fails with errors.
// TODO Handle non fatal errors
func (s *server) Serve(l net.Listener) error {
	for {
		conn, err := l.Accept()
		if err != nil {
			return err
		}
		go s.ServeConn(conn)
	}
}

func (s *server) RegisterService(sd *grpc.ServiceDesc, service interface{}) {
	s.s.RegisterName(sd.ServiceName, service)
}
