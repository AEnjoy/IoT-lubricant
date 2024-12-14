package edge

import "github.com/AEnjoy/IoT-lubricant/pkg/edge/config"

func SaveConfig() {
	if err := config.SaveConfig(config.SaveType_ALL); err != nil {
		panic(err)
	}
}
