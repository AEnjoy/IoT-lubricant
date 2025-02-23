/*
This utils is a tool that cover go struct to json
*/
package main

import (
	"os"

	"github.com/AEnjoy/IoT-lubricant/internal/model/form/request"
	"github.com/bytedance/sonic/encoder"
)

func main() {
	var structBody = request.AddAgentRequest{Description: "agent", GatherCycle: 1, ReportCycle: 5, Address: "127.0.0.1:5436", DataCompressAlgorithm: "default"}
	file, err := os.OpenFile("output.json", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	encoder.NewStreamEncoder(file).Encode(structBody)
}
