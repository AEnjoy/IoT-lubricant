//go:build !cgo

package compress

var _ Compressor = (*Zstd)(nil)

type Zstd struct{ Default }
