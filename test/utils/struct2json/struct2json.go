/*
This utils is a tool that cover go struct to json
*/
package main

import (
	"os"

	"github.com/bytedance/sonic/encoder"
)

func main() {
	var structBody any
	file, err := os.OpenFile("output.json", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	_ = encoder.NewStreamEncoder(file).Encode(structBody)
}
