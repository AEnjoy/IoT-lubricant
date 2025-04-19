package agent

import "github.com/aenjoy/iot-lubricant/pkg/edge/config"

func SaveConfig() {
	if err := config.SaveConfig(config.SaveType_ALL); err != nil {
		panic(err)
	}
}
