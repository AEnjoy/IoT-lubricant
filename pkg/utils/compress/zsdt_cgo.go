//go:build cgo

// we need `go build -tags external_libzstd`

package compress

import (
	"bytes"
	"runtime"

	"github.com/DataDog/zstd"
)

var _ Compressor = (*Zstd)(nil)

type Zstd struct {
}

func (z *Zstd) GetName() string {
	return "zstd"
}

func (z *Zstd) Compress(in []byte) ([]byte, error) {
	var buf bytes.Buffer
	writer := zstd.NewWriter(&buf)
	_ = writer.SetNbWorkers(runtime.NumCPU())
	_, err := writer.Write(in)
	if err != nil {
		return nil, err
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (z *Zstd) Decompress(in []byte) ([]byte, error) {
	reader := zstd.NewReader(bytes.NewReader(in))
	var out bytes.Buffer
	_, err := out.ReadFrom(reader)
	if err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}
