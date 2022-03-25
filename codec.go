package brpc

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"sync"

	"github.com/icexin/brpc-go/metapb"
	"github.com/keegancsmith/rpc"
	"google.golang.org/protobuf/proto"
)

// protocol spec from https://github.com/apache/incubator-brpc/blob/60159fc3f3e13490fb9806ea0a0cb0dcdbda7f7d/docs/cn/baidu_std.md

var (
	MagicStr          = [4]byte{'P', 'R', 'P', 'C'}
	ErrBadMagic       = errors.New("bad magic number")
	ErrBadMessageSize = errors.New("message size exceed")
)

const (
	maxMessageSize = 64<<20 + 24 // 64M+sizeof(rpcHeader)
)

type rpcHeader struct {
	Magic [4]byte
	X     struct {
		PacketSize int32
		MetaSize   int32
	}
}

// codec is a common rpc codec for implementing rpc.ClientCodec and rpc.ServerCodec
type codec struct {
	conn io.ReadWriteCloser
	w    *bufio.Writer
	r    *bufio.Reader

	// temporary work space
	h rpcHeader
}

func newCodec(conn io.ReadWriteCloser) *codec {
	return &codec{
		conn: conn,
		w:    bufio.NewWriter(conn),
		r:    bufio.NewReader(conn),
	}
}

// Write send rpc header and body to peer
func (c *codec) Write(meta *metapb.RpcMeta, x interface{}, cw compressWriter) error {
	buffer := new(bytes.Buffer)
	metasize := proto.Size(meta)
	h := rpcHeader{
		Magic: MagicStr,
	}
	h.X.MetaSize = int32(metasize)
	h.X.PacketSize = int32(metasize)

	// write header, for placeholder, we will change PacketSize later if no error in meta
	binary.Write(buffer, binary.BigEndian, &h)

	// write meta
	buf, err := proto.Marshal(meta)
	if err != nil {
		return err
	}
	buffer.Write(buf)

	// skip write body if error
	if meta.Response != nil && meta.Response.ErrorCode != 0 {
		buffer.WriteTo(c.w)
		return c.w.Flush()
	}

	// write body
	msg := x.(proto.Message)
	buf, err = proto.Marshal(msg)
	if err != nil {
		return err
	}
	// record the offset before we write data
	len1 := buffer.Len()
	if cw != nil {
		wc, err := cw(buffer)
		if err != nil {
			return err
		}
		wc.Write(buf)
		wc.Close()
	} else {
		buffer.Write(buf)
	}

	dataSize := buffer.Len() - len1
	h.X.PacketSize = int32(metasize + dataSize)
	// write new header
	w := bytes.NewBuffer(buffer.Bytes()[:0])
	binary.Write(w, binary.BigEndian, &h)

	buffer.WriteTo(c.w)
	return c.w.Flush()
}

func mustDecode(r io.Reader, x interface{}) {
	err := binary.Read(r, binary.BigEndian, x)
	if err != nil {
		panic(err)
	}
}

func mustReadFull(r io.Reader, buf []byte) {
	_, err := io.ReadFull(r, buf)
	if err != nil {
		panic(err)
	}
}

// ReadHeader read rpc header from peer
func (c *codec) ReadHeader(meta *metapb.RpcMeta) (err error) {
	defer func() {
		catch := recover()
		if catch != nil {
			if e, ok := catch.(error); ok {
				err = e
			} else {
				err = fmt.Errorf("%v", catch)
			}
		}
	}()

	// decode rpc header
	mustDecode(c.r, &c.h.Magic)
	if c.h.Magic != MagicStr {
		return ErrBadMagic
	}
	mustDecode(c.r, &c.h.X)

	if c.h.X.PacketSize > maxMessageSize {
		return ErrBadMessageSize
	}

	// decode rpc meta
	buf := make([]byte, c.h.X.MetaSize)
	mustReadFull(c.r, buf)
	err = proto.Unmarshal(buf, meta)
	if err != nil {
		return
	}
	if meta.GetAttachmentSize() != 0 {
		return errors.New("attachment not supported")
	}
	return
}

// ReadBody read rpc body from peer which corresponding to last rpc header
func (c *codec) ReadBody(x interface{}, cr compressReader) error {
	dataSize := c.h.X.PacketSize - c.h.X.MetaSize
	if x == nil {
		_, err := c.r.Discard(int(dataSize))
		return err
	}

	msg := x.(proto.Message)

	var buf []byte
	if cr != nil {
		// construct a compress reader
		rc, err := cr(io.LimitReader(c.r, int64(dataSize)))
		if err != nil {
			return err
		}
		// read all uncompressed data to buf
		buf, err = ioutil.ReadAll(rc)
		if err != nil {
			rc.Close()
			return err
		}
		rc.Close()
	} else {
		buf = make([]byte, dataSize)
		_, err := io.ReadFull(c.r, buf)
		if err != nil {
			return err
		}
	}

	return proto.Unmarshal(buf, msg)
}

func (c *codec) Close() error {
	return c.conn.Close()
}

type clientCodec struct {
	c *codec

	// temporary work space
	m metapb.RpcMeta
}

// newClientCodec returns a new rpc.ClientCodec using sofa-pbrpc on conn.
func newClientCodec(conn io.ReadWriteCloser) rpc.ClientCodec {
	return &clientCodec{c: newCodec(conn)}
}

func splitServiceMethod(serviceMethod string) (string, string) {
	if i := strings.LastIndex(serviceMethod, "."); i >= 0 {
		return serviceMethod[:i], serviceMethod[i+1:]
	}
	return "", ""
}

func (c *clientCodec) WriteRequest(req *rpc.Request, x interface{}) error {
	serviceName, methodName := splitServiceMethod(req.ServiceMethod)
	m := &metapb.RpcMeta{
		CorrelationId: int64(req.Seq),
		Request: &metapb.RpcRequestMeta{
			ServiceName: serviceName,
			MethodName:  methodName,
		},
	}
	// TODO support request compress type
	return c.c.Write(m, x, nil)
}

func (c *clientCodec) ReadResponseHeader(resp *rpc.Response) error {
	err := c.c.ReadHeader(&c.m)
	if err != nil {
		return err
	}
	resp.Seq = uint64(c.m.GetCorrelationId())
	if c.m.GetResponse().GetErrorCode() != 0 {
		resp.Error = fmt.Sprintf("code:%d, reason:%s", c.m.Response.ErrorCode, c.m.Response.ErrorText)
	}

	return nil
}

func (c *clientCodec) ReadResponseBody(x interface{}) error {
	return c.c.ReadBody(x, newCompressReader(metapb.CompressType(c.m.CompressType)))
}

func (c *clientCodec) Close() error {
	return c.c.Close()
}

type serverCodec struct {
	c *codec

	// since ReadRequestHeader and ReadRequestBody are called in pairs,
	// reqmeta only shared between them
	reqmeta *metapb.RpcMeta

	mutex   sync.Mutex
	pending map[uint64]*metapb.RpcMeta
}

// newServerCodec returns a new rpc.ServerCodec using sofa-pbrpc on conn.
func newServerCodec(conn io.ReadWriteCloser) rpc.ServerCodec {
	return &serverCodec{
		c:       newCodec(conn),
		pending: make(map[uint64]*metapb.RpcMeta),
	}
}

func (s *serverCodec) ReadRequestHeader(req *rpc.Request) error {
	meta := new(metapb.RpcMeta)
	err := s.c.ReadHeader(meta)
	if err != nil {
		return err
	}
	req.Seq = uint64(meta.CorrelationId)
	req.ServiceMethod = meta.Request.ServiceName + "." + meta.Request.MethodName
	s.reqmeta = meta

	s.mutex.Lock()
	s.pending[req.Seq] = meta
	s.mutex.Unlock()
	return nil
}

func (s *serverCodec) ReadRequestBody(x interface{}) error {
	cr := newCompressReader(metapb.CompressType(s.reqmeta.CompressType))
	return s.c.ReadBody(x, cr)
}

func (s *serverCodec) WriteResponse(resp *rpc.Response, x interface{}) error {
	s.mutex.Lock()
	reqmeta := s.pending[resp.Seq]
	delete(s.pending, resp.Seq)
	s.mutex.Unlock()

	meta := &metapb.RpcMeta{
		CorrelationId: int64(resp.Seq),
		CompressType:  reqmeta.CompressType,
		Response:      &metapb.RpcResponseMeta{},
	}

	if resp.Error != "" {
		meta.Response.ErrorCode = 500
		meta.Response.ErrorText = resp.Error
	}

	return s.c.Write(meta, x, newCompressWriter(metapb.CompressType(meta.CompressType)))
}

func (s *serverCodec) Close() error {
	return s.c.Close()
}