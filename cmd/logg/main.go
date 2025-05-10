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

func init() {
	version.ServiceName = "SvcLogger"
	version.PrintVersionInfo()
}
