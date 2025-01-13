package agent

import (
	"context"
	"sync"
)

type pool struct {
	p sync.Map // id => *agentControl
}

func newPool() *pool {
	return &pool{}
}
func (p *pool) JoinAgent(ctx context.Context, a *agentControl) error {
	cli, _, err := a.tryConnect()
	if err != nil {
		return err
	}
	a.AgentCli = cli

	a.init(ctx)
	p.p.Store(a.id, a)

	return nil
}
func (p *pool) GetAgentControl(id string) *agentControl {
	if v, ok := p.p.Load(id); ok {
		return v.(*agentControl)
	}
	return nil
}
