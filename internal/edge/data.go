package edge

import (
	"github.com/AEnjoy/IoT-lubricant/internal/edge/data"
	"github.com/AEnjoy/IoT-lubricant/pkg/edge"
	"github.com/AEnjoy/IoT-lubricant/pkg/utils/compress"
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
