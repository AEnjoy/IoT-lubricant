package logCollect

import (
	"context"

	"github.com/aenjoy/iot-lubricant/pkg/constant"
	"github.com/aenjoy/iot-lubricant/pkg/logger"
	mqV2 "github.com/aenjoy/iot-lubricant/pkg/utils/mq/v2"
	"github.com/aenjoy/iot-lubricant/pkg/version"
	svcpb "github.com/aenjoy/iot-lubricant/protobuf/svc"
	"github.com/aenjoy/iot-lubricant/services/corepkg/datastore"
	"github.com/aenjoy/iot-lubricant/services/corepkg/ioc"
	logg "github.com/aenjoy/iot-lubricant/services/logg/api"

	"google.golang.org/protobuf/proto"
)

var _ ioc.Object = (*Log)(nil)

type Log struct {
	mqV2.Mq
}

func (l *Log) Init() error {
	ds := ioc.Controller.Get(ioc.APP_NAME_CORE_DATABASE_STORE).(*datastore.DataStore)
	l.Mq = ds.V2mq
	logg.L, _ = logg.NewLogger(l, false)
	logg.SetServiceName(version.ServiceName)
	logg.L = logg.L.
		WithVersionJson(version.VersionJson()).
		WithPrintToStdout().
		AsRoot()
	logger.Debugf("Success Init Logger Compent<internal>: %s", version.VersionJson())
	logg.L.Debugf("Success Init Logger Compent<Logg>: %s", version.VersionJson())
	logger.Debugf("Logger Info: %s", logg.L.String())

	go l.loggerStore()
	return nil
}

func (Log) Weight() uint16 {
	return ioc.SvcLoggerService
}

func (Log) Version() string {
	return "dev"
}

var logDataCh = make(chan *svcpb.Logs, 300)

func (*Log) Transfer(_ context.Context, data *svcpb.Logs, _ bool) (retval error) {
	logDataCh <- data
	return
}
func (l *Log) loggerStore() {
	for dataCh := range logDataCh {
		logdata, err := proto.Marshal(dataCh)
		if err != nil {
			return
		}
		_ = l.QueuePublish(constant.MESSAGE_SVC_LOGGER, logdata)
	}
}
