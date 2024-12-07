package edge

import (
	"bytes"
	"time"

	"github.com/AEnjoy/IoT-lubricant/internal/edge/data"
	"github.com/AEnjoy/IoT-lubricant/pkg/edge"
	"github.com/AEnjoy/IoT-lubricant/pkg/utils/compress"
)

var (
	triggerChan    = make(chan struct{})
	compressedChan = make(chan *edge.DataPacket) // 压缩后的数据
	dataSetCh      = make(chan []byte, 3)        //采集到的原始数据
	dataChan2      = make(chan *edge.DataPacket) // send to gateway

)

func compressor(method string, dataChan <-chan []byte, compressedChan chan<- *edge.DataPacket) {
	compressor, _ := compress.NewCompressor(method)
	for dataPacket := range dataChan {
		timeNow := time.Now()
		data, _ := compressor.Compress(dataPacket)
		compressedChan <- &edge.DataPacket{
			Data:      data,
			Timestamp: timeNow,
		}
	}
}
func transmitter(cycle int, compressedChan <-chan *edge.DataPacket, triggerChan chan struct{}, dataChan chan<- *edge.DataPacket) {
	var buffer [][]byte
	var firstPacketTime *time.Time
	for {
		select {
		case compressedData := <-compressedChan:
			data.DCL.Lock()
			if firstPacketTime == nil {
				firstPacketTime = &compressedData.Timestamp
			}
			buffer = append(buffer, compressedData.Data)
			if len(buffer) >= cycle {
				d := &edge.DataPacket{
					Data:      bytes.Join(buffer, compress.Sepa),
					Timestamp: *firstPacketTime,
				}
				dataChan <- d
				data.DataCollect = append(data.DataCollect, d)
				buffer = nil
				firstPacketTime = nil
			}
			data.DCL.Unlock()
		case <-triggerChan:
			if len(buffer) > 0 {
				dataChan <- &edge.DataPacket{
					Data:      bytes.Join(buffer, compress.Sepa),
					Timestamp: *firstPacketTime,
				}
				buffer = nil
				firstPacketTime = nil
			}
		}
	}
}

var _compressor compress.Compressor

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

var dataHandlerCh = make(chan *dataHandler, 30)
