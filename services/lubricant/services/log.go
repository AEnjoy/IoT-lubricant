package services

import (
	"context"

	"github.com/aenjoy/iot-lubricant/pkg/constant"
	"github.com/aenjoy/iot-lubricant/pkg/utils/mq"
	"github.com/aenjoy/iot-lubricant/pkg/version"
	svcpb "github.com/aenjoy/iot-lubricant/protobuf/svc"
	logg "github.com/aenjoy/iot-lubricant/services/logg/api"
	"github.com/aenjoy/iot-lubricant/services/lubricant/ioc"
	"google.golang.org/protobuf/proto"
)

var _ ioc.Object = (*Log)(nil)

type Log struct {
	mq.Mq
}

func (l *Log) Init() error {
	l.Mq = ioc.Controller.Get(ioc.APP_NAME_CORE_Internal_MQ_SERVICE).(mq.Mq)
	logg.L, _ = logg.NewLogger(l, false)
	logg.SetServiceName(version.ServiceName)
	logg.L = logg.L.
		WithVersionJson(version.VersionJson()).
		WithPrintToStdout().
		AsRoot()

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
		_ = l.PublishBytes(constant.MESSAGE_SVC_LOGGER, logdata)
	}
}
