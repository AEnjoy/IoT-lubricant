package reporterHandler

import (
	"context"
	"strings"

	logg "github.com/aenjoy/iot-lubricant/services/logg/api"
)

func (a app) gatewayOfflineGuardPayload(payload any) {
	str := string(payload.([]byte))
	// str is `"%s<!SPLIT!>%s", userid, gatewayid`
	var userid, gatewayid string
	result := strings.Split(str, "<!SPLIT!>")
	if len(result) == 2 {
		userid = result[0]
		gatewayid = result[1]
	} else {
		logg.L.Errorf("internalError: failed to split gateway id: %s", str)
		return
	}

	txn := a.ICoreDb.Begin()
	err := a.ICoreDb.SetGatewayStatus(context.Background(), txn, userid, gatewayid, "offline")
	a.ICoreDb.Commit(txn)
	if err != nil {
		logg.L.Errorf("failed to set gateway status: %v", err)
	}
}
