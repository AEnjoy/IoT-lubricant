package api

import (
	"encoding/base64"

	"github.com/aenjoy/iot-lubricant/pkg/model"
	"github.com/aenjoy/iot-lubricant/services/datastoreAssistant/driver"
	logg "github.com/aenjoy/iot-lubricant/services/logg/api"
	"github.com/bytedance/sonic"
)

func DsnTest(dsnType, dsn, userId string) string {
	switch dsnType {
	case "mysql":
		dsn, err := base64.StdEncoding.DecodeString(dsn)
		if err != nil {
			logg.L.Errorf("failed to decode dsn:%s", dsn)
			return "Failed"
		}
		_, f, err := driver.NewMySQLDriver(string(dsn), "", userId)
		if err != nil {
			logg.L.Errorf("test failed:failed to create mysql driver:%s userId:%s", dsn, userId)
			return "Failed"
		}
		_ = f()
		return "Success"
	case "TDEngine":
		dsn, err := base64.StdEncoding.DecodeString(dsn)
		if err != nil {
			logg.L.Errorf("failed to decode dsn:%s", dsn)
			return "Failed"
		}
		var info model.TDEngineLinkerInfo
		err = sonic.Unmarshal(dsn, &info)
		if err != nil {
			logg.L.Errorf("failed to unmarshal dsn:%s", dsn)
			return "Failed"
		}
		_, f, err := driver.NewTDEngineDriver(userId,
			info.Host,
			info.User,
			info.Pass, info.Db, info.Port,
			nil,
			nil)
		if err != nil {
			logg.L.Errorf("test failed:failed to create TDEngine driver:%s userId:%s", dsn, userId)
			return "Failed"
		}
		_ = f()
		return "Success"
	default:
		logg.L.Errorf("test failed:unsupported dsn type:%s", dsnType)
		return "Unsupported"
	}
}
