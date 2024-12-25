package agent

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/AEnjoy/IoT-lubricant/internal/app/gateway/internal/data"
	"github.com/AEnjoy/IoT-lubricant/internal/model"
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

	gather bool
	start  bool

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
	c._stop()
	c.cancel()
}
func (c *agentControl) _stop() {
	_, _ = c.agentCli.StopGather(c.ctx, &agent.StopGatherRequest{})
}
func (c *agentControl) _offlineWarn() {
	c.exceptSig <- exception.ErrNewException(nil, code.WarnAgentOffline,
		exception.WithLevel(level.Warn),
		exception.WithMsg(fmt.Sprintf("AgentID: %s", c.id)),
	)
}
