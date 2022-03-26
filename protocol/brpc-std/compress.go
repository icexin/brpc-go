package bstd

import (
	"compress/gzip"
	"compress/zlib"
	"io"
	"io/ioutil"

	"github.com/golang/snappy"
	"github.com/icexin/brpc-go/protocol/brpc-std/metapb"
	"github.com/pierrec/lz4"
)

type (
	compressReader func(r io.Reader) (io.ReadCloser, error)
	compressWriter func(w io.Writer) (io.WriteCloser, error)
)

// gzip compress
func gzipCompressReader(r io.Reader) (io.ReadCloser, error) {
	return gzip.NewReader(r)
}

func gzipCompressWriter(w io.Writer) (io.WriteCloser, error) {
	return gzip.NewWriter(w), nil
}

// zlib compress
func zlibCompressReader(r io.Reader) (io.ReadCloser, error) {
	return zlib.NewReader(r)
}

func zlibCompressWriter(w io.Writer) (io.WriteCloser, error) {
	return zlib.NewWriter(w), nil
}

// snappy compress
func snappyCompressReader(r io.Reader) (io.ReadCloser, error) {
	return ioutil.NopCloser(snappy.NewReader(r)), nil
}

func snappyCompressWriter(w io.Writer) (io.WriteCloser, error) {
	return snappy.NewWriter(w), nil
}

// lz4 compress
func lz4CompressReader(r io.Reader) (io.ReadCloser, error) {
	return ioutil.NopCloser(lz4.NewReader(r)), nil
}

func lz4CompressWriter(w io.Writer) (io.WriteCloser, error) {
	return lz4.NewWriter(w), nil
}

func newCompressReader(tp metapb.CompressType) compressReader {
	switch tp {
	case metapb.CompressType_COMPRESS_TYPE_GZIP:
		return compressReader(gzipCompressReader)
	case metapb.CompressType_COMPRESS_TYPE_ZLIB:
		return compressReader(zlibCompressReader)
	case metapb.CompressType_COMPRESS_TYPE_SNAPPY:
		return compressReader(snappyCompressReader)
	case metapb.CompressType_COMPRESS_TYPE_LZ4:
		return compressReader(lz4CompressReader)
	default:
		return nil
	}
}

func newCompressWriter(tp metapb.CompressType) compressWriter {
	switch tp {
	case metapb.CompressType_COMPRESS_TYPE_GZIP:
		return compressWriter(gzipCompressWriter)
	case metapb.CompressType_COMPRESS_TYPE_ZLIB:
		return compressWriter(zlibCompressWriter)
	case metapb.CompressType_COMPRESS_TYPE_SNAPPY:
		return compressWriter(snappyCompressWriter)
	case metapb.CompressType_COMPRESS_TYPE_LZ4:
		return compressWriter(lz4CompressWriter)
	default:
		return nil
	}
}
