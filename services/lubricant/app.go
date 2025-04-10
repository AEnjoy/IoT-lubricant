package lubricant

import (
	"fmt"

	"github.com/aenjoy/iot-lubricant/pkg/logger"
	"github.com/aenjoy/iot-lubricant/services/lubricant/ioc"
	"github.com/aenjoy/iot-lubricant/services/lubricant/repo"
	"github.com/aenjoy/iot-lubricant/services/lubricant/router"
	"github.com/gin-gonic/gin"
)

type app struct {
	hostName string
	port     string

	httpServer *gin.Engine
	dbClient   *repo.CoreDb
}

func (s *app) Run() error {
	err := s.httpServer.Run(fmt.Sprintf(":%s", s.port))
	if err != nil {
		logger.Errorln("Web Server start error, error Info is: ", err)
		logger.Info("Web Server will not start, please check the configuration or error Info.")
		return err
	}
	return nil
}
func NewApp(opts ...func(*app) error) *app {
	var server = new(app)
	for _, opt := range opts {
		if err := opt(server); err != nil {
			logger.Fatalf("Failed to apply option: %v", err)
		}
	}
	return server
}

func SetHostName(hostName string) func(*app) error {
	return func(s *app) error {
		s.hostName = hostName
		return nil
	}
}

func UseDB(dbClient *repo.CoreDb) func(*app) error {
	return func(s *app) error {
		s.dbClient = dbClient
		return nil
	}
}

func UseGinEngine() func(*app) error {
	return func(s *app) error {
		s.httpServer = ioc.Controller.Get(ioc.APP_NAME_CORE_WEB_SERVER).(*router.WebService).Engine
		return nil
	}
}

func SetPort(port string) func(*app) error {
	return func(s *app) error {
		s.port = port
		return nil
	}
}

func UseServerKey() func(*app) error {
	return func(*app) error {
		return initKeys()
	}
}
func UseCasdoor() func(*app) error {
	return func(s *app) error {
		return initCasdoor()
	}
}
func UseSignalHandler(handler func()) func(*app) error {
	return func(a *app) error {
		go handler()
		return nil
	}
}
