package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/aenjoy/iot-lubricant/pkg/logger"
	"github.com/panjf2000/ants/v2"
)

func main() {
	client, ctx, err := newCoreClient(hostAddress, userID, gatewayID)
	if err != nil {
		panic(err)
	}
	pool, err := ants.NewPool(1024, ants.WithPreAlloc(true), ants.WithNonblocking(true))
	if err != nil {
		panic(err)
	}
	for range 1024 {
		pool.Submit(func() {
			pushData2Core(client, ctx, dataCh)
		})
	}
	StartConcurrentGeneration(2)

	var sigCh = make(chan os.Signal, 50)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	logger.Info("Start to register agent")
	regAgentOnline(client, ctx)
	logger.Info("Start to send data. Press Ctrl+C abort")
	close(startSig)
	<-sigCh
	logger.Infof("Success:%d Failed:%d", sendCountSuccess, sendCountFail)
	regAgentOffline(client, ctx)
	os.Exit(0)
}
