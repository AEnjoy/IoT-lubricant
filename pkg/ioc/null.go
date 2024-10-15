package ioc

var _ Object = (*NilObject)(nil)

// NilObject is an empty object used for placeholder.
type NilObject struct {
}

func (n NilObject) Init() error {
	return nil
}

func (n NilObject) Weight() uint16 {
	return 1
}

func (n NilObject) Version() string {
	return ""
}
