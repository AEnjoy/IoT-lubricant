package main

import (
	"fmt"
	"runtime"
)

var (
	Version           string
	BuildTime         string
	GoVersion         string
	GitCommit         string
	Features          string
	BuildHostPlatform string
	PlatformVersion   string
)

func printBuildInfo() {
	fmt.Printf("IoT-lubricant-Version: %s\n", Version)
	fmt.Printf("Build-Time: %s\n", BuildTime)
	fmt.Printf("Go-Version: %s\n", GoVersion)
	fmt.Printf("Git-Commit: %s\n", GitCommit)
	fmt.Printf("Features: %s\n", Features)
	fmt.Printf("BuildHostPlatform: %s\n", BuildHostPlatform)
	fmt.Printf("BuildPlatformVersion: %s\n", PlatformVersion)
	fmt.Printf("Runing Platform Info: %s/%s\n", runtime.GOOS, runtime.GOARCH)
}
