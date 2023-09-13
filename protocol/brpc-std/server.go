package bstd

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"runtime/debug"
	"sync"

	"github.com/icexin/brpc-go"
	"github.com/icexin/brpc-go/protocol/brpc-std/metapb"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

type service struct {
	desc    *grpc.ServiceDesc
	methods map[string]*grpc.MethodDesc
	srv     interface{}
}

type server struct {
	services map[string]*service
	opts     *brpc.ServerOptions
}

func newServer(options ...brpc.ServerOption) *server {
	var opts brpc.ServerOptions
	for _, opt := range options {
		if o, ok := opt.(brpc.BServerOption); ok {
			o(&opts)
		}
	}
	return &server{
		services: make(map[string]*service),
		opts:     &opts,
	}
}

// ServeConn runs the server on a single connection. ServeConn blocks, serving the connection until the client hangs up. The caller typically invokes ServeConn in a go statement.
func (s *server) ServeConn(conn net.Conn) {
	defer conn.Close()
	codec := newCodec(conn)
	writeLock := new(sync.Mutex)

	wg := new(sync.WaitGroup)
	for {
		req, srv, method, dec, keepReading, err := s.readRequest(codec)
		if err != nil {
			log.Printf("read request error:%v", err)
			if !keepReading {
				break
			}
			// send a response if we actually managed to read a header.
			if req != nil {
				s.sendResponse(writeLock, codec, err, req, nil, nil)
			}
			continue
		}
		wg.Add(1)
		go func() {
			defer func() {
				r := recover()
				if r != nil {
					debug.PrintStack()
					s.sendResponse(writeLock, codec, fmt.Errorf("panic:%v", r), req, nil, nil)
				}
				wg.Done()
			}()
			resp, err := method.Handler(srv, context.Background(), dec, s.opts.Interceptor)
			if err != nil {
				log.Printf("call method error:%v", err)
			}
			s.sendResponse(writeLock, codec, err, req, resp, nil)
		}()
	}
	// We've seen that there are no more requests.
	// Wait for responses to be sent before closing connection
	wg.Wait()
}

func (s *server) sendResponse(lock *sync.Mutex, codec *codec, errReq error, req *metapb.RpcMeta, resp interface{}, cw compressWriter) {
	meta := &metapb.RpcMeta{
		CorrelationId: int64(req.GetCorrelationId()),
		Response:      &metapb.RpcResponseMeta{},
	}
	if errReq != nil {
		meta.Response.ErrorCode = 500
		meta.Response.ErrorText = errReq.Error()
	}
	lock.Lock()
	err := codec.Write(meta, resp, cw)
	lock.Unlock()
	if err != nil {
		log.Printf("write response error:%v", err)
	}
}

func (s *server) readRequest(codec *codec) (reqMeta *metapb.RpcMeta, srv interface{}, method *grpc.MethodDesc, dec func(interface{}) error, keepReading bool, err error) {
	var service *service
	reqMeta, service, method, keepReading, err = s.readRequestHeader(codec)
	if err != nil {
		if !keepReading {
			return
		}
		// discard body
		codec.ReadBodyBytes(true, nil)
		return
	}

	var body []byte
	body, err = codec.ReadBodyBytes(false, newCompressReader(metapb.CompressType(reqMeta.GetCompressType())))
	if err != nil {
		return
	}
	srv = service.srv
	dec = func(v interface{}) error {
		msg, ok := v.(proto.Message)
		if !ok {
			return fmt.Errorf("request type is not proto.Message")
		}
		return proto.Unmarshal(body, msg)
	}
	return
}

func (s *server) readRequestHeader(codec *codec) (reqMeta *metapb.RpcMeta, srv *service, method *grpc.MethodDesc, keepReading bool, err error) {
	reqMeta = new(metapb.RpcMeta)
	err = codec.ReadHeader(reqMeta)
	if err != nil {
		reqMeta = nil
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return
		}
		err = fmt.Errorf("read rpc header error:%w", err)
		return
	}

	// We read the header successfully. If we see an error now,
	// we can still recover and read the next request.
	keepReading = true
	req := reqMeta.GetRequest()
	srv, ok := s.services[req.ServiceName]
	if !ok {
		err = fmt.Errorf("can't find service:%s", req.ServiceName)
		return
	}
	method, ok = srv.methods[req.MethodName]
	if !ok {
		err = fmt.Errorf("can't find method:%s", req.MethodName)
	}
	return
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

func (s *server) RegisterService(sd *grpc.ServiceDesc, srv interface{}) {
	methods := make(map[string]*grpc.MethodDesc)
	for i := range sd.Methods {
		m := &sd.Methods[i]
		methods[m.MethodName] = m
	}
	s.services[sd.ServiceName] = &service{
		desc:    sd,
		methods: methods,
		srv:     srv,
	}
}
