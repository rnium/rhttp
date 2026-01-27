package dynamicdata

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/rnium/rhttp/internal/build"
	"github.com/rnium/rhttp/pkg/rhttp"
)

func nRandomBytes(n int) []byte {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return b
}

func buildUUID() string {
	b := nRandomBytes(16)

	// Set version (4) and variant (RFC 4122)
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80

	// Format as string
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

type dripReader struct {
	interval time.Duration
	numBytes int
	count    int
}

func newDripReader(numbytes int, duration float64) *dripReader {
	interval := duration / float64(numbytes) * 1000
	return &dripReader{
		numBytes: numbytes,
		interval: time.Millisecond * time.Duration(interval),
	}
}

func (dr *dripReader) Read(p []byte) (int, error) {
	if dr.count >= dr.numBytes {
		return 0, io.EOF
	}
	n := copy(p, []byte{'*'})
	dr.count += n
	time.Sleep(dr.interval)
	return n, nil
}

type chunkedReader struct {
	numStreams int
	count      int
	rd         *build.ReadResponseData
}

func newChunkedReader(num int, req *rhttp.Request) *chunkedReader {
	return &chunkedReader{
		numStreams: num,
		count: 0,
		rd: build.BuildReadData(req),
	}
}

func (cr *chunkedReader) Read(p []byte) (int, error) {
	if cr.count >= cr.numStreams {
		return 0, io.EOF
	}
	data, _ := json.Marshal(cr.rd)
	n := copy(p, data)
	cr.count += 1
	time.Sleep(time.Millisecond * 100)
	return n, nil
}
