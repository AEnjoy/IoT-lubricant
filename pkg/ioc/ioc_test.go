package ioc

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var _ Object = (*testObject)(nil)

type testObject struct {
	w uint16
}

var testContainer = &MapContainer{
	name:   "test",
	storge: make(map[string]Object),
}

func (t *testObject) Init() error {
	return nil
}

func (t *testObject) Weight() uint16 {
	return t.w
}

func (t testObject) Version() string {
	return ""
}

func TestIoCThreadSafe(t *testing.T) {
	for i := 0; i < 1000; i++ {
		go testContainer.Registry(strconv.Itoa(i), &testObject{uint16(i)})
	}
	<-time.After(time.Microsecond * 100)
	assert.NoError(t, testContainer.Init())
}
