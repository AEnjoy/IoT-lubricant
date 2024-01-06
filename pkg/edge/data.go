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
	compressedChan = make(chan []byte)
	dataSetCh      = make(chan []byte, 3)
	dataChan2      = make(chan []byte) // send to gateway
)

func (a *app) dataService(method string) {
	go compressor(method, dataSetCh, compressedChan)
	go transmitter(a.config.ReportCycle, compressedChan, triggerChan, dataChan2)
}

func compressor(method string, dataChan <-chan []byte, compressedChan chan<- []byte) {
	compressor, _ := compress.NewCompressor(method)
	for dataPacket := range dataChan {
		data, _ := compressor.Compress(dataPacket)
		compressedChan <- data
	}
}
func transmitter(cycle int, compressedChan <-chan []byte, triggerChan chan struct{}, dataChan chan<- []byte) {
	var buffer [][]byte
	for {
		select {
		case compressedData := <-compressedChan:
			buffer = append(buffer, compressedData)
			if len(buffer) >= cycle {
				dataChan <- bytes.Join(buffer, compress.Sepa)
				buffer = nil
			}
		case <-triggerChan:
			if len(buffer) > 0 {
				dataChan <- bytes.Join(buffer, compress.Sepa)
				buffer = nil
			}
		}
	}
}
