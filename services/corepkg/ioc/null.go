package ioc

var _ Object = (*NilObject)(nil)

// NilObject is an empty object used for placeholder.
type NilObject struct {
}

func (NilObject) Init() error {
	return nil
}

func (NilObject) Weight() uint16 {
	return 1
}

func (NilObject) Version() string {
	return ""
}
