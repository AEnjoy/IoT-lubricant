package agent

import (
	"net/http"

	metapb "github.com/aenjoy/iot-lubricant/protobuf/meta"
)

func errIsTargetNotEqual(info *metapb.CommonResponse) bool {
	if info == nil {
		return false
	}
	return info.GetMessage() == "target agentID error" && info.GetCode() == http.StatusInternalServerError
}
