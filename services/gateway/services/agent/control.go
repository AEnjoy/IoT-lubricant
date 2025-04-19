package agent

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/aenjoy/iot-lubricant/pkg/edge"
	"github.com/aenjoy/iot-lubricant/pkg/logger"
	"github.com/aenjoy/iot-lubricant/pkg/model"
	"github.com/aenjoy/iot-lubricant/pkg/types/action"
	errLevel "github.com/aenjoy/iot-lubricant/pkg/types/code"
	"github.com/aenjoy/iot-lubricant/pkg/types/exception"
	exceptionCode "github.com/aenjoy/iot-lubricant/pkg/types/exception/code"
	object "github.com/aenjoy/iot-lubricant/pkg/types/task"
	"github.com/aenjoy/iot-lubricant/pkg/utils/crontab"
	agentpb "github.com/aenjoy/iot-lubricant/protobuf/agent"
	corepb "github.com/aenjoy/iot-lubricant/protobuf/core"
	metapb "github.com/aenjoy/iot-lubricant/protobuf/meta"
	svcpb "github.com/aenjoy/iot-lubricant/protobuf/svc"
	"github.com/aenjoy/iot-lubricant/services/gateway/services/data"
	logg "github.com/aenjoy/iot-lubricant/services/logg/api"

	"google.golang.org/genproto/googleapis/rpc/status"
)

const exceptionSigMaxSize = 10

type agentControl struct {
	id   string
	Slot []int // for api paths

	AgentCli    agentpb.EdgeServiceClient
	AgentInfo   *model.Agent
	dataCollect data.Apis

	rootCtx context.Context
	ctx     context.Context
	cancel  context.CancelFunc

	exitSig   chan struct{} // 销毁信号
	exceptSig chan *exception.Exception
	reporter  chan *corepb.ReportRequest

	gatherLock sync.Mutex
	online     bool // online/offline

	_initOnce sync.Once
	_exitOnce sync.Once
}

func (c *agentControl) setCtx(ctx context.Context) {
	if ctx == nil || ctx == context.TODO() {
		ctx = c.rootCtx
	}
	c.ctx, c.cancel = context.WithCancel(ctx)
}
func (c *agentControl) tryConnect() (cli agentpb.EdgeServiceClient, closeCallBack func(), err error) {
	cli, closeCallBack, err = edge.NewAgentCliWithClose(c.AgentInfo.Address)
	if err != nil {
		return nil, nil, err
	}

	ctxTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = cli.Ping(ctxTimeout, &metapb.Ping{})
	if err != nil {
		return nil, nil, err
	}
	return
}
func (c *agentControl) init(ctx context.Context) {
	c._initOnce.Do(func() {
		c.rootCtx = ctx
		c.id = c.AgentInfo.AgentId
		c.setCtx(context.TODO())
		c.exitSig = make(chan struct{})
		c.exceptSig = make(chan *exception.Exception, exceptionSigMaxSize)
		c.dataCollect = data.NewDataStoreApis(c.id)

		_ = crontab.RegisterCron(c.getAgentLogs, "@every 5s")
		if handelFunc != nil {
			go handelFunc(c.exceptSig)
		}
		if len(c.Slot) == 0 {
			c.Slot = []int{0}
		}
		go func() {
			for {
				select {
				case <-c.exitSig:
					return
				default:
				}
				c._checkOnline()
				time.Sleep(time.Second * time.Duration(rand.Int31n(5)+1))
			}
		}()
	})
}

var _ping = &metapb.Ping{}

func (c *agentControl) _checkOnline() bool {
	_, err := c.AgentCli.Ping(c.ctx, _ping)
	if err != nil {
		if c.online {
			c.reporter <- &corepb.ReportRequest{
				AgentId: c.id,
				Req: &corepb.ReportRequest_AgentStatus{
					AgentStatus: &corepb.AgentStatusRequest{
						Req: &status.Status{Message: "offline"},
					},
				},
			}
			c.online = false
			c._offlineWarn()
		}
	} else {
		if !c.online {
			logg.L.Debug("agent online")
			c.reporter <- &corepb.ReportRequest{
				AgentId: c.id,
				Req: &corepb.ReportRequest_AgentStatus{
					AgentStatus: &corepb.AgentStatusRequest{
						Req: &status.Status{Message: "online"},
					},
				},
			}
		}
		c.online = true
	}
	return c.online
}
func (c *agentControl) IsStarted() bool {
	return c.online
}
func (c *agentControl) getAgentLogs() {
	begin := time.Now().Unix()
	log := logg.L.
		WithProtocol("grpc").
		WithAction(action.CollectAgentLogs).
		WithNotPrintToStdout()
	logs, err := c.AgentCli.CollectLogs(c.ctx, &agentpb.CollectLogsRequest{})
	end := time.Now().Unix()
	if err != nil {
		if !c.online {
			return
		}
		log.WithLoglevel(svcpb.Level_ERROR).
			WithCost(time.Duration(end-begin)).
			Errorf("failed to get agent logs:%v", err)
		return
	}
	if len(logs.GetLogs()) > 0 {
		c.reporter <- &corepb.ReportRequest{
			AgentId: c.id,
			Req: &corepb.ReportRequest_ReportLog{
				ReportLog: &corepb.ReportLogRequest{
					AgentId: c.id,
					Logs:    logs.GetLogs(),
				},
			},
		}
	}
}

// 销毁
func (c *agentControl) _checkOut() {
	logger.Info("agent control checkout", c.id)
	c._exitOnce.Do(func() {
		close(c.exitSig)
		close(c.exceptSig)
		data.ManualPushAgentData(c.id)
	})
}
func (c *agentControl) _stopGather() error {
	resp, err := c.AgentCli.StopGather(c.ctx, &agentpb.StopGatherRequest{})
	if err != nil {
		return err
	}
	//c.cancel()

	if resp.GetMessage() != "success" && resp.GetMessage() != "" && resp.GetMessage() != "Gather is not working" {
		return errors.New(resp.GetMessage())
	}
	return nil
}
func (c *agentControl) _start() error {
	_, err := c.AgentCli.StartGather(c.ctx, &agentpb.StartGatherRequest{})
	return err
}
func (c *agentControl) _offlineWarn() {
	c.exceptSig <- exception.ErrNewException(nil, exceptionCode.WarnAgentOffline,
		exception.WithLevel(errLevel.Warn),
		exception.WithContext(string(object.TargetAgent), c.id),
	)
}
func (c *agentControl) _gather(wg *sync.WaitGroup) {
	defer wg.Done()
	for _, i := range c.Slot {
		resp, err := c.AgentCli.GetGatherData(c.ctx, &agentpb.GetDataRequest{
			AgentID: c.id,
			Slot:    int32(i),
		})
		if err != nil {
			c.exceptSig <- exception.ErrNewException(err, exceptionCode.ErrGaterDataReqFailed,
				exception.WithLevel(errLevel.Error),
				exception.WithMsg(fmt.Sprintf("AgentID: %s", c.id)),
				exception.WithMsg(fmt.Sprintf("Slot: %d", i)),
			)
			continue
		}
		_ = c.dataCollect.Push(resp)
	}
}

// Start 启动 Agent 对于是否能启动成功，请稍后通过 IsStarted 判断 (未完全实现)
func (c *agentControl) Start(ctx context.Context) error {
	// todo:
	//if !c.online {
	//	return errors.New("agent is not started or has been offline")
	//}
	if c.gatherLock.TryLock() {
		defer c.gatherLock.Unlock()
		c.ctx, c.cancel = context.WithCancel(ctx)
	}
	return nil
}
func (c *agentControl) IsGathering() bool {
	if c.gatherLock.TryLock() {
		defer c.gatherLock.Unlock()
		return false
	}
	return true
}
func (c *agentControl) StartGather() error {
	c.setCtx(context.TODO())
	if !c._checkOnline() {
		return errors.New("agent is not started or has been offline")
	}
	if !c.gatherLock.TryLock() {
		return errors.New("agent is already gathering")
	}

	go func() {
		defer c.gatherLock.Unlock()
		if err := c._start(); err != nil {
			if !strings.Contains(err.Error(), "Gather is working now") {
				logg.L.Error("agent: Gather start failed", err)
				c.exceptSig <- exception.ErrNewException(err, exceptionCode.ErrGaterStartFailed,
					exception.WithLevel(errLevel.Error),
					exception.WithMsg(fmt.Sprintf("AgentID: %s", c.id)),
				)
				return
			}
			logg.L.Warn("agent: Gather is working now", err)
		}

		ticker := time.NewTicker(time.Duration(c.AgentInfo.GatherCycle) * time.Second)
		defer ticker.Stop()

		var wg sync.WaitGroup
		for {
			select {
			case <-c.ctx.Done():
				wg.Wait()
				return
			case <-ticker.C:
				wg.Add(1)
				go c._gather(&wg)
			}
		}
	}()
	return nil
}
func (c *agentControl) StopGather() error {
	err := c._stopGather()
	c.cancel()
	return err
}

// Exit 销毁
func (c *agentControl) Exit() {
	_ = c._stopGather()
	c.exitSig <- struct{}{}
	c._checkOut()
}
func (c *agentControl) GetDataApi() data.Apis {
	return c.dataCollect
}
func newAgentControl(a *model.Agent, reporter chan *corepb.ReportRequest) *agentControl {
	return &agentControl{
		AgentInfo: a,
		reporter:  reporter,
	}
}
