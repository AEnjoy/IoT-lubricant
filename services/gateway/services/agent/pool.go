package agent

import (
	"context"
	"sync"
)

type pool struct {
	p sync.Map // id => *agentControl
}

var _pool *pool

func newPool() *pool {
	if _pool == nil {
		_pool = new(pool)
	}
	return _pool
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
func (p *pool) RemoveAgent(id string) {
	p.p.Delete(id)
}
func (p *pool) GetAgentControl(id string) *agentControl {
	if v, ok := p.p.Load(id); ok {
		return v.(*agentControl)
	}
	// 对于不存在的情况，先判断是否存在于数据库，如果是，重新加载后不存在则返回nil
	return p._tryReloadAgentControl(id)
}
func (p *pool) _tryReloadAgentControl(id string) *agentControl {
	agent, err := _apis.(*agentApis).db.GetAgent(id)
	if err != nil {
		return nil
	}

	ctrl := newAgentControl(&agent, _apis.(*agentApis).reporter)
	err = p.JoinAgent(context.Background(), ctrl)
	if err != nil {
		return nil
	}
	return ctrl
}
