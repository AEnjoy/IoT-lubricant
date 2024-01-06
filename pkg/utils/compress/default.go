package compress

var _ Compressor = (*Default)(nil)

type Default struct {
	// Do not Compress
}

func (d Default) Compress(in []byte) ([]byte, error) {
	return in, nil
}

func (d Default) Decompress(in []byte) ([]byte, error) {
	return in, nil
}

func (d Default) GetName() string {
	return "-"
}
