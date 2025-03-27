package main

import (
	"context"

	"github.com/aenjoy/iot-lubricant/cmd/logg/internal"
	"github.com/aenjoy/iot-lubricant/pkg/version"
	"github.com/aenjoy/iot-lubricant/services/logg"
)

func main() {
	app := logg.NewApp(
		logg.UseContext(context.Background()),
		logg.UseDb(internal.GetDb()),
		logg.UseMq(internal.GetMq()),
	)
	panic(app.Run())
}

var (
	ServiceName       = "SvcLogger"
	Version           string
	BuildTime         string
	GoVersion         string
	GitCommit         string
	Features          string
	BuildHostPlatform string
	PlatformVersion   string
)

func init() {
	version.ServiceName = ServiceName
	version.Version = Version
	version.BuildTime = BuildTime
	version.GoVersion = GoVersion
	version.GitCommit = GitCommit
	version.Features = Features
	version.BuildHostPlatform = BuildHostPlatform
	version.PlatformVersion = PlatformVersion
	version.PrintVersionInfo()
}
