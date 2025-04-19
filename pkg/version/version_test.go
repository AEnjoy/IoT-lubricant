package version

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	ServiceName = "unitTest"
	BuildTime = time.Now().String()
	GoVersion = "1.17"
	Version = "1.0.0"
	BuildHostPlatform = "linux/amd64"
	PlatformVersion = "debian 12"
	assert.Equal(t, string(MarshallJson()), MarshallJsonString())
}
