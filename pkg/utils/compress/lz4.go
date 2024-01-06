package compress

import (
	"bytes"

	"github.com/pierrec/lz4"
)

var _ Compressor = (*Lz4)(nil)

type Lz4 struct {
}

func (l *Lz4) GetName() string {
	return "lz4"
}

func (l *Lz4) Compress(in []byte) ([]byte, error) {
	var buf bytes.Buffer

	writer := lz4.NewWriter(&buf)
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

func (l *Lz4) Decompress(in []byte) ([]byte, error) {
	reader := lz4.NewReader(bytes.NewReader(in))
	var out bytes.Buffer
	_, err := out.ReadFrom(reader)
	if err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}
