package agent

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/AEnjoy/IoT-lubricant/pkg/edge"
	"github.com/AEnjoy/IoT-lubricant/pkg/logger"
	"github.com/AEnjoy/IoT-lubricant/pkg/model"
	errLevel "github.com/AEnjoy/IoT-lubricant/pkg/types/code"
	"github.com/AEnjoy/IoT-lubricant/pkg/types/exception"
	exceptCode "github.com/AEnjoy/IoT-lubricant/pkg/types/exception/code"
	agentpb "github.com/AEnjoy/IoT-lubricant/protobuf/agent"
	metapb "github.com/AEnjoy/IoT-lubricant/protobuf/meta"
	data2 "github.com/AEnjoy/IoT-lubricant/services/gateway/services/data"
)

const exceptionSigMaxSize = 10

type agentControl struct {
	id   string
	Slot []int // for api paths

	AgentCli    agentpb.EdgeServiceClient
	AgentInfo   *model.Agent
	dataCollect data2.Apis

	ctx    context.Context
	cancel context.CancelFunc

	exitSig   chan struct{} // 销毁信号
	exceptSig chan *exception.Exception

	gatherLock sync.Mutex
	online     bool // online/offline

	_initOnce sync.Once
	_exitOnce sync.Once
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
		c.id = c.AgentInfo.AgentId
		c.ctx, c.cancel = context.WithCancel(ctx)
		c.exitSig = make(chan struct{})
		c.exceptSig = make(chan *exception.Exception, exceptionSigMaxSize)
		c.dataCollect = data2.NewDataStoreApis(c.id)

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
				_, err := c.AgentCli.Ping(c.ctx, &metapb.Ping{})
				if err != nil {
					if c.online {
						c.online = false
						c._offlineWarn()
					}
				} else {
					c.online = true
				}
				time.Sleep(time.Second * time.Duration(rand.Int31n(5)+1))
			}
		}()
	})
}
func (c *agentControl) IsStarted() bool {
	return c.online
}

// 销毁
func (c *agentControl) _checkOut() {
	logger.Info("agent control checkout", c.id)
	c._exitOnce.Do(func() {
		close(c.exitSig)
		close(c.exceptSig)
		data2.ManualPushAgentData(c.id)
	})
}
func (c *agentControl) _stopGather() error {
	resp, err := c.AgentCli.StopGather(c.ctx, &agentpb.StopGatherRequest{})
	if err != nil {
		return err
	}
	c.cancel()

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
	c.exceptSig <- exception.ErrNewException(nil, exceptCode.WarnAgentOffline,
		exception.WithLevel(errLevel.Warn),
		exception.WithMsg(fmt.Sprintf("AgentID: %s", c.id)),
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
			c.exceptSig <- exception.ErrNewException(err, exceptCode.ErrGaterDataReqFailed,
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

func (c *agentControl) StartGather() error {
	if !c.online {
		return errors.New("agent is not started or has been offline")
	}
	if !c.gatherLock.TryLock() {
		return errors.New("agent is already gathering")
	}

	go func() {
		if err := c._start(); err != nil {
			if !strings.Contains(err.Error(), "Gather is working now") {
				logger.Errorln("agent: Gather start failed", err)
				c.exceptSig <- exception.ErrNewException(err, exceptCode.ErrGaterStartFailed,
					exception.WithLevel(errLevel.Error),
					exception.WithMsg(fmt.Sprintf("AgentID: %s", c.id)),
				)
				return
			}
			logger.Warnln("agent: Gather is working now", err)
		}
		defer c.gatherLock.Unlock()

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
	return c._stopGather()
}

// Exit 销毁
func (c *agentControl) Exit() {
	_ = c._stopGather()
	c.exitSig <- struct{}{}
	c._checkOut()
}
func (c *agentControl) GetDataApi() data2.Apis {
	return c.dataCollect
}
func newAgentControl(a *model.Agent) *agentControl {
	return &agentControl{
		AgentInfo: a,
	}
}
