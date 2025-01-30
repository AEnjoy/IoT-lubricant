package main

import (
	"fmt"
	"runtime"
)

var (
	Version         string
	BuildTime       string
	GoVersion       string
	GitTag          string
	Features        string
	Platform        string
	PlatformVersion string
)

func printBuildInfo() {
	fmt.Printf("IoT-lubricant-Version: %s\n", Version)
	fmt.Printf("Build-Time: %s\n", BuildTime)
	fmt.Printf("Go-Version: %s\n", GoVersion)
	fmt.Printf("Git-Tag: %s\n", GitTag)
	fmt.Printf("Features: %s\n", Features)
	fmt.Printf("Platform: %s\n", Platform)
	fmt.Printf("Platform-Version: %s\n", PlatformVersion)
	fmt.Printf("Runing Platform Info: %s/%s\n", runtime.GOOS, runtime.GOARCH)
}
