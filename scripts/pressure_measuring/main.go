package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/aenjoy/iot-lubricant/pkg/logger"
)

func main() {
	go writeLog()
	client, ctx, err := newCoreClient(hostAddress, userID, gatewayID)
	if err != nil {
		panic(err)
	}

	go pushData2Core(client, ctx, dataCh)
	StartConcurrentGeneration(2)

	var sigCh = make(chan os.Signal, 50)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	logger.Info("Start to register agent")
	regAgentOnline(client, ctx)
	handleTaskResp(client, ctx)
	logger.Info("Start to send data. Press Ctrl+C abort")
	close(startSig)
	<-sigCh
	logger.Infof("Success:%d Failed:%d", sendCountSuccess, sendCountFail)
	regAgentOffline(client, ctx)
	os.Exit(0)
}
