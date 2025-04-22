package api

import (
	"time"

	logg "github.com/aenjoy/iot-lubricant/services/logg/api"
	"github.com/bytedance/sonic"
)

func CovertData(in any, nodeId string, ts time.Time) map[string]any {
	data, ok := in.([]byte)
	if !ok {
		logg.L.
			WithAction("CovertData").
			Errorf("failed to convert data to []byte")
		return nil
	}
	var retVal = map[string]any{}
	err := sonic.Unmarshal(data, &retVal)
	if err != nil {
		logg.L.
			WithAction("CovertData").
			Errorf("failed to unmarshal data:%v", err)
		return nil
	}
	retVal["node_id"] = nodeId
	retVal["ts"] = ts.Unix()
	return retVal
}
