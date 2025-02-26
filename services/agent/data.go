package agent

import (
	"github.com/aenjoy/iot-lubricant/pkg/edge"
	"github.com/aenjoy/iot-lubricant/pkg/utils/compress"
	"github.com/aenjoy/iot-lubricant/services/agent/data"
)

var (
	dataHandlerCh = make(chan *dataHandler, 30)
	_compressor   compress.Compressor
)

func DataHandler() {
	for ch := range dataHandlerCh {
		// todo :为了提高性能，可以池化go协程
		go func(dataIn *dataHandler) {
			d, _ := _compressor.Compress(*dataIn.dataIn)
			data.Collector.AddData(dataIn.slot, &edge.DataPacket{Data: d})
		}(ch)
	}
}

type dataHandler struct {
	slot   int
	dataIn *[]byte
}
