package agent

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/AEnjoy/IoT-lubricant/internal/app/gateway/internal/data"
	"github.com/AEnjoy/IoT-lubricant/internal/model"
	"github.com/AEnjoy/IoT-lubricant/pkg/logger"
	level "github.com/AEnjoy/IoT-lubricant/pkg/types/code"
	"github.com/AEnjoy/IoT-lubricant/pkg/types/exception"
	"github.com/AEnjoy/IoT-lubricant/pkg/types/exception/code"
	"github.com/AEnjoy/IoT-lubricant/protobuf/agent"
	"github.com/AEnjoy/IoT-lubricant/protobuf/meta"
)

const exceptionSigMaxSize = 10

type agentControl struct {
	id   string
	slot []int // for api paths

	agentCli    agent.EdgeServiceClient
	agentInfo   *model.Agent
	dataCollect data.Apis

	ctx    context.Context
	cancel context.CancelFunc

	exitSig   chan struct{} // 销毁信号
	exceptSig chan *exception.Exception

	gatherLock sync.Mutex
	start      bool // online/offline

	once sync.Once
}

func (c *agentControl) init() {
	c.once.Do(func() {
		c.exitSig = make(chan struct{})
		c.exceptSig = make(chan *exception.Exception, exceptionSigMaxSize)
		c.dataCollect = data.NewDataStoreApis(c.id)

		if len(c.slot) == 0 {
			c.slot = []int{0}
		}
		go func() {
			for {
				select {
				case <-c.exitSig:
					c._checkOut()
					return
				default:
				}
				_, err := c.agentCli.Ping(c.ctx, &meta.Ping{})
				if err != nil {
					if c.start {
						c.start = false
						c._offlineWarn()
					}
				} else {
					c.start = true
				}
				time.Sleep(time.Second * time.Duration(rand.Int31n(5)+1))
			}
		}()
	})
}
func (c *agentControl) IsStarted() bool {
	return c.start
}

// 销毁
func (c *agentControl) _checkOut() {
	logger.Info("agent control checkout", c.id)
	_ = c._stop()
	c.cancel()
}
func (c *agentControl) _stop() error {
	_, err := c.agentCli.StopGather(c.ctx, &agent.StopGatherRequest{})
	return err
}
func (c *agentControl) _start() error {
	_, err := c.agentCli.StartGather(c.ctx, &agent.StartGatherRequest{})
	return err
}
func (c *agentControl) _offlineWarn() {
	c.exceptSig <- exception.ErrNewException(nil, code.WarnAgentOffline,
		exception.WithLevel(level.Warn),
		exception.WithMsg(fmt.Sprintf("AgentID: %s", c.id)),
	)
}
func (c *agentControl) _gather(wg *sync.WaitGroup) {
	defer wg.Done()
	for _, i := range c.slot {
		resp, err := c.agentCli.GetGatherData(c.ctx, &agent.GetDataRequest{
			AgentID: c.id,
			Slot:    int32(i),
		})
		if err != nil {
			c.exceptSig <- exception.ErrNewException(err, code.ErrGaterDataReqFailed,
				exception.WithLevel(level.Error),
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
	if c.start {
		return errors.New("agent is already stared")
	}

	c.ctx, c.cancel = context.WithCancel(ctx)
	//c.start = true
	return nil
}

func (c *agentControl) StartGather() error {
	if !c.start {
		return errors.New("agent is not started or has been offline")
	}
	if !c.gatherLock.TryLock() {
		return errors.New("agent is already gathering")
	}

	go func() {
		defer c.gatherLock.Unlock()

		ticker := time.NewTicker(time.Duration(c.agentInfo.GatherCycle) * time.Second)
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
func (c *agentControl) StopGather() {
	c._checkOut()
}
func (c *agentControl) GetDataApi() data.Apis {
	return c.dataCollect
}
