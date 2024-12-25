package agent

import (
	"context"
	"sync"
	"time"

	"github.com/AEnjoy/IoT-lubricant/pkg/edge"
	"github.com/AEnjoy/IoT-lubricant/protobuf/meta"
)

type pool struct {
	p sync.Map // id => *agentControl
}

func newPool() *pool {
	return &pool{}
}
func (p *pool) JoinAgent(ctx context.Context, a *agentControl) error {
	a.ctx, a.cancel = context.WithCancel(ctx)

	cli, err := edge.NewAgentCli(a.agentInfo.Address)
	if err != nil {
		return err
	}

	ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err = cli.Ping(ctxTimeout, &meta.Ping{})
	if err != nil {
		return err
	}

	a.id = a.agentInfo.AgentId
	a.agentCli = cli

	p.p.Store(a.id, a)

	return nil
}
