package compress

import (
	"bytes"
	"compress/gzip"
	"io"
)

var _ Compressor = (*Gzip)(nil)

type Gzip struct {
}

func (g *Gzip) GetName() string {
	return "gzip"
}

func (g *Gzip) Compress(in []byte) ([]byte, error) {
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	_, err := w.Write(in)
	if err != nil {
		return nil, err
	}
	defer w.Close()
	return b.Bytes(), nil
}

func (g *Gzip) Decompress(in []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(in))
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	return io.ReadAll(reader)
}
