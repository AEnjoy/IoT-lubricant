package version

import (
	"fmt"
	"runtime"
	"unsafe"
)

const DefaultVersion = "dev"

var (
	ServiceName       string
	Version           = DefaultVersion
	GitCommit         = DefaultVersion
	Features          = DefaultVersion
	BuildTime         string
	GoVersion         string
	BuildHostPlatform string
	PlatformVersion   string
)

func PrintVersionInfo() {
	fmt.Printf("%s-Version: %s\n", ServiceName, Version)
	fmt.Printf("Build-Time: %s\n", BuildTime)
	fmt.Printf("Go-Version: %s\n", GoVersion)
	fmt.Printf("Git-Commit: %s\n", GitCommit)
	fmt.Printf("Features: %s\n", Features)
	fmt.Printf("BuildHostPlatform: %s\n", BuildHostPlatform)
	fmt.Printf("BuildPlatformVersion: %s\n", PlatformVersion)
	fmt.Printf("Runing Platform Info: %s/%s\n", runtime.GOOS, runtime.GOARCH)
}
func MarshallJsonString() string {
	return fmt.Sprintf(`{"serviceName":"%s","version":"%s","buildTime":"%s","goVersion":"%s","gitCommit":"%s","features":"%s","buildHostPlatform":"%s","buildPlatformVersion":"%s","runingPlatformInfo":"%s/%s"}`,
		ServiceName,
		Version,
		BuildTime,
		GoVersion,
		GitCommit,
		Features,
		BuildHostPlatform,
		PlatformVersion,
		runtime.GOOS,
		runtime.GOARCH)
}
func VersionJson() []byte {
	return MarshallJson()
}
func MarshallJson() []byte {
	return StringToBytes(MarshallJsonString())
}
func StringToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&s))
}
