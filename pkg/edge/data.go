package edge

import (
	"bytes"
	"time"

	"github.com/AEnjoy/IoT-lubricant/pkg/utils/compress"
)

type DataPacket struct {
	Data      []byte
	Timestamp time.Time
}

var (
	triggerChan    = make(chan struct{})
	compressedChan = make(chan *DataPacket)
	dataSetCh      = make(chan []byte, 3)
	dataChan2      = make(chan *DataPacket) // send to gateway
)

func compressor(method string, dataChan <-chan []byte, compressedChan chan<- *DataPacket) {
	compressor, _ := compress.NewCompressor(method)
	for dataPacket := range dataChan {
		timeNow := time.Now()
		data, _ := compressor.Compress(dataPacket)
		compressedChan <- &DataPacket{
			Data:      data,
			Timestamp: timeNow,
		}
	}
}
func transmitter(cycle int, compressedChan <-chan *DataPacket, triggerChan chan struct{}, dataChan chan<- *DataPacket) {
	var buffer [][]byte
	var firstPacketTime *time.Time
	for {
		select {
		case compressedData := <-compressedChan:
			if firstPacketTime == nil {
				firstPacketTime = &compressedData.Timestamp
			}
			buffer = append(buffer, compressedData.Data)
			if len(buffer) >= cycle {
				dataChan <- &DataPacket{
					Data:      bytes.Join(buffer, compress.Sepa),
					Timestamp: *firstPacketTime,
				}
				buffer = nil
				firstPacketTime = nil
			}
		case <-triggerChan:
			if len(buffer) > 0 {
				dataChan <- &DataPacket{
					Data:      bytes.Join(buffer, compress.Sepa),
					Timestamp: *firstPacketTime,
				}
				buffer = nil
				firstPacketTime = nil
			}
		}
	}
}
