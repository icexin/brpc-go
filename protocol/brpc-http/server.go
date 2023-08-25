package bhttp

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"path"
	"strings"

	"github.com/icexin/brpc-go"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type service struct {
	desc    *grpc.ServiceDesc
	methods map[string]*grpc.MethodDesc
	srv     interface{}
}

type server struct {
	services   map[string]*service
	httpServer *http.Server
	opts       *brpc.ServerOptions
}

func newServer(options ...brpc.ServerOption) *server {
	var opts brpc.ServerOptions
	for _, opt := range options {
		if o, ok := opt.(brpc.BServerOption); ok {
			o(&opts)
		}
	}
	return &server{
		services:   make(map[string]*service),
		httpServer: &http.Server{},
		opts:       &opts,
	}
}

// Serve accepts incoming connections on the listener l, creating a new ServerConn and service goroutine for each. The service goroutines read pbrpc requests and then call the registered handlers to reply to them. Serve returns when l.Accept fails with errors.
// TODO Handle non fatal errors
func (s *server) Serve(l net.Listener) error {
	return s.httpServer.Serve(l)
}

func (s *server) Shutdown() {
	s.httpServer.Shutdown(context.Background())
}

func (s *server) handleService(srv *service, w http.ResponseWriter, r *http.Request) {
	methodName := path.Base(r.URL.Path)
	method, ok := srv.methods[methodName]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "method %s not found", methodName)
		return
	}

	buf, err := s.serveRequest(srv, method, r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, err.Error())
		return
	}
	w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
	w.Write(buf)
}

func (s *server) serveRequest(srv *service, method *grpc.MethodDesc, r *http.Request) ([]byte, error) {
	var (
		enc func(interface{}) ([]byte, error)
		dec func([]byte, interface{}) error
	)

	switch r.Header.Get("Content-Type") {
	case "application/json":
		enc = jsonEncode
		dec = jsonDecode
	case "application/proto", "application/protobuf":
		enc = protoEncode
		dec = protoDecode
	default:
		return nil, fmt.Errorf("invalid content type:%s", r.Header.Get("Content-Type"))
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("read body error:%w", err)
	}
	decFunc := func(v interface{}) error {
		err := dec(body, v)
		if err != nil {
			return err
		}
		return nil
	}
	resp, err := method.Handler(srv.srv, r.Context(), decFunc, s.opts.Interceptor)
	if err != nil {
		return nil, err
	}
	buf, err := enc(resp)
	if err != nil {
		return nil, fmt.Errorf("encode response error:%w", err)
	}
	return buf, nil
}

func (s *server) RegisterService(sd *grpc.ServiceDesc, srv interface{}) {
	methods := make(map[string]*grpc.MethodDesc)
	for i := range sd.Methods {
		m := &sd.Methods[i]
		methods[m.MethodName] = m
	}
	service := &service{
		desc:    sd,
		methods: methods,
		srv:     srv,
	}
	fs := strings.Split(sd.ServiceName, ".")
	shortName := fs[len(fs)-1]
	s.services[sd.ServiceName] = service
	s.services[shortName] = service
	handler := func(w http.ResponseWriter, r *http.Request) {}
	http.HandleFunc("/"+shortName+"/", handler)
	http.HandleFunc("/"+sd.ServiceName+"/", handler)
}

func jsonEncode(v interface{}) ([]byte, error) {
	msg, ok := v.(proto.Message)
	if !ok {
		return nil, fmt.Errorf("invalid type:%T", v)
	}
	opts := protojson.MarshalOptions{
		UseProtoNames: true,
	}
	return opts.Marshal(msg)
}

func jsonDecode(buf []byte, v interface{}) error {
	msg, ok := v.(proto.Message)
	if !ok {
		return fmt.Errorf("invalid type:%T", v)
	}
	return protojson.Unmarshal(buf, msg)
}

func protoEncode(v interface{}) ([]byte, error) {
	msg, ok := v.(proto.Message)
	if !ok {
		return nil, fmt.Errorf("invalid type:%T", v)
	}
	return proto.Marshal(msg)
}

func protoDecode(buf []byte, v interface{}) error {
	msg, ok := v.(proto.Message)
	if !ok {
		return fmt.Errorf("invalid type:%T", v)
	}
	return proto.Unmarshal(buf, msg)
}
