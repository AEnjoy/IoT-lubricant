package compress

var Sepa = []byte("#^#")

type Compressor interface {
	Compress([]byte) ([]byte, error)
	Decompress([]byte) ([]byte, error)
	GetName() string
}

func NewCompressor(method string) (Compressor, error) {
	switch method {
	case "lz4":
		return &Lz4{}, nil
	case "gzip":
		return &Gzip{}, nil
	case "zstd": // Need cgo supported with libzstd
		return &Zstd{}, nil
	default: // default '-','default', use default
		return &Default{}, nil
	}
}
