package agent

import (
	"context"
	"fmt"
	"time"

	"github.com/AEnjoy/IoT-lubricant/internal/model"
	"github.com/AEnjoy/IoT-lubricant/internal/model/repo"
	"github.com/AEnjoy/IoT-lubricant/pkg/docker"
	errCh "github.com/AEnjoy/IoT-lubricant/pkg/error"
	"github.com/AEnjoy/IoT-lubricant/pkg/logger"
	errLevel "github.com/AEnjoy/IoT-lubricant/pkg/types/code"
	"github.com/AEnjoy/IoT-lubricant/pkg/types/exception"
	exceptCode "github.com/AEnjoy/IoT-lubricant/pkg/types/exception/code"
	agentpb "github.com/AEnjoy/IoT-lubricant/protobuf/agent"
	proxypb "github.com/AEnjoy/IoT-lubricant/protobuf/proxy"
	"github.com/bytedance/sonic"
	grpcCode "google.golang.org/genproto/googleapis/rpc/code"
	"gorm.io/gorm"
)

var _ Apis = (*agentApis)(nil)

type agentApis struct {
	db repo.GatewayDbOperator

	*pool
}

func (a *agentApis) init() {
	ctx := context.Background()
	agents, err := a.db.GetAllAgents(nil)
	if err != nil {
		panic(err)
	}
	//todo: 这里可以加个并发限速防止goroutine合并/丢弃
	for _, agent := range agents {
		go func(agent *model.Agent) {
			ins := a.db.GetAgentInstance(nil, agent.AgentId)
			control := newAgentControl(agent)

			if ins.Local && !docker.IsContainerRunning(ins.ContainerID) {
				if err := docker.StartContainer(ins.ContainerID); err != nil {
					logger.Error("online container failed", ins.ContainerID, err)
					return
				}
			}

			if err := a.pool.JoinAgent(ctx, control); err != nil {
				logger.Error("agent join to handel pool failed", agent.AgentId, err)
			}
		}(&agent)
		time.Sleep(100 * time.Millisecond) //避免Goroutine启动过快失败
	}
}
func (a *agentApis) StartAgent(id string) error {
	ctrl := a.pool.GetAgentControl(id)
	if ctrl == nil {
		return exception.New(exceptCode.ErrorGatewayAgentNotFound, exception.WithLevel(errLevel.Error), exception.WithMsg(fmt.Sprintf("agentID:%s", id)))
	}
	if !a.isLocalAgentDevice(id) {
		return exception.New(exceptCode.OperationOnlyAtLocal, exception.WithLevel(errLevel.Error),
			exception.WithMsg(fmt.Sprintf("agentID:%s", id)),
			exception.WithReason("Cannot control manually added agent instances to online, please manage manually."))
	}

	ins := a.db.GetAgentInstance(nil, id)
	if !docker.IsContainerRunning(ins.ContainerID) {
		err := bootAgentInstance(ins.ContainerID)
		if err != nil {
			return exception.ErrNewException(err, exceptCode.ErrorAgentStartFailed,
				exception.WithLevel(errLevel.Error),
				exception.WithMsg(fmt.Sprintf("agentID:%s", id)))
		}
		return nil
	}
	logger.Warnln("agent already started", id)
	return nil
}

// StopAgent 对于来自Core的操作，认为StopAgent操作，实际上只是停止采集，不会真正删除agent实例
func (a *agentApis) StopAgent(id string) error {
	ctrl := a.pool.GetAgentControl(id)
	return ctrl.StopGather()
}
func (a *agentApis) KillAgent(id string) error {
	panic("implement me")
}
func (a *agentApis) RemoveAgent(id string) error {
	txn := a.db.Begin()
	errorCh := errCh.NewErrorChan()
	defer errCh.HandleErrorCh(errorCh).ErrorWillDo(func(error) {
		a.db.Rollback(txn)
	}).SuccessWillDo(func() {
		a.db.Commit(txn)
	}).Do()

	var (
		ctx  = context.Background()
		ctrl = a.pool.GetAgentControl(id)
		ins  = a.db.GetAgentInstance(txn, id)
	)

	err := ctrl.StopGather()
	if err != nil {
		errorCh.Report(err, exceptCode.StopAgentFailed, "stop gather failed", true)
		return err
	}
	if ins.Local {
		if err = docker.Stop(ctx, ins.ContainerID); err != nil {
			errorCh.Report(err, exceptCode.StopAgentFailed, "stop agent container failed", true)
			return err
		}
		if err = docker.Remove(ctx, ins.ContainerID); err != nil {
			errorCh.Report(err, exceptCode.RemoveAgentFailed, "remove agent container failed", true)
			return err
		}
	} else {
		logger.Warnf("AgentID:[%s] Not supporting the removal of manually added agents, "+
			"which will result in the deletion of database records", id)
	}

	ctrl.Exit()
	a.pool.RemoveAgent(id)
	return nil
}

// todo: update 和 edit 需要做下功能区分
func (a *agentApis) UpdateAgent(id string, conf *model.CreateAgentRequest) error {
	txn := a.db.Begin()
	errorCh := errCh.NewErrorChan()
	defer errCh.HandleErrorCh(errorCh).ErrorWillDo(func(error) {
		a.db.Rollback(txn)
	}).SuccessWillDo(func() {
		a.db.Commit(txn)
	}).Do()

	// todo:可以校验请求是否与数据库的Agent最初的配置一致
	ins := a.db.GetAgentInstance(txn, id)
	defer func(db repo.GatewayDbOperator, txn *gorm.DB, id string, ins model.AgentInstance) {
		if err := db.UpdateAgentInstance(txn, id, ins); err != nil {
			err = exception.ErrNewException(err, exceptCode.ErrorAgentUpdateFailed,
				exception.WithLevel(errLevel.Error),
				exception.WithMsg(fmt.Sprintf("agentID:%s", id)),
				exception.WithMsg("at update agent instance database commit"))
			errorCh.Report(err, exceptCode.ErrorAgentUpdateFailed, "update agent instance failed", true)
		}
	}(a.db, txn, id, ins)
	originCreateConf := &model.CreateAgentRequest{CreateAgentConf: new(model.CreateAgentConf)}
	if err := sonic.Unmarshal([]byte(ins.CreateConf), originCreateConf.CreateAgentConf); err != nil {
		errorCh.Report(err, exceptCode.ErrorDecodeJSON, "unmarshal agent conf failed", true)
		return err
	}

	// 四种情况: 1.local-> local  2.remote -> remote 3.local -> remote 4.remote -> local

	if ins.Local {
		// 1.local-> local
		if conf == nil {
			conf = originCreateConf
		}
		if conf.AgentContainerInfo != nil {
			if conf.AgentContainerInfo != nil {
				id, err := docker.UpdateContainer(context.Background(), conf.AgentContainerInfo, ins.ContainerID)
				if err != nil {
					errorCh.Report(err, exceptCode.ErrorAgentUpdateFailed, "update agent container failed", true)
					return err
				}
				ins.ContainerID = id
			}
			if conf.DriverContainerInfo != nil {
				id, err := docker.UpdateContainer(context.Background(), conf.DriverContainerInfo, ins.DriverID)
				if err != nil {
					errorCh.Report(err, exceptCode.ErrorAgentUpdateFailed, "update driver container failed", true)
					return err
				}
				ins.DriverID = id
			}
		}
	} else {
		err := exception.New(exceptCode.ErrorAgentUpdateNotSupportRemote, exception.WithLevel(errLevel.Error),
			exception.WithMsg("agent container conf is needed"))
		errorCh.Report(err, exceptCode.ErrorAgentUpdateNotSupportRemote, "please manually update your remote agent or recreate agent", true)
		return err
	}
	return nil
}
func (a *agentApis) EditAgent(_ string, req *proxypb.EditAgentRequest) error {
	txn := a.db.Begin()
	errorCh := errCh.NewErrorChan()
	defer errCh.HandleErrorCh(errorCh).ErrorWillDo(func(error) {
		a.db.Rollback(txn)
	}).SuccessWillDo(func() {
		a.db.Commit(txn)
	}).Do()

	ins := a.db.GetAgentInstance(txn, req.GetAgentId())
	defer func(db repo.GatewayDbOperator, txn *gorm.DB, id string, ins model.AgentInstance) {
		if err := db.UpdateAgentInstance(txn, id, ins); err != nil {
			err = exception.ErrNewException(err, exceptCode.ErrorAgentUpdateFailed,
				exception.WithLevel(errLevel.Error),
				exception.WithMsg(fmt.Sprintf("agentID:%s", id)),
				exception.WithMsg("at update agent instance database commit"))
			errorCh.Report(err, exceptCode.ErrorAgentUpdateFailed, "update agent instance failed", true)
		}
	}(a.db, txn, req.GetAgentId(), ins)

	// 四种情况: 1.local-> local  2.remote -> remote 3.local -> remote 4.remote -> local
	ctrl := a.pool.GetAgentControl(req.GetAgentId())
	m := model.ProxypbEditAgentRequest2Agent(req)
	if len(req.GetConf()) > 0 {
		createConf := model.CreateAgentRequest{}
		if err := sonic.Unmarshal(req.GetConf(), &createConf); err != nil {
			errorCh.Report(err, exceptCode.ErrorDecodeJSON, "unmarshal agent conf failed", true)
			return err
		}
		if err := a.RemoveAgent(req.GetAgentId()); err != nil {
			errorCh.Report(err, exceptCode.RemoveAgentFailed, "update agent container failed(at remove old agent)", true)
			return err
		}
		if err := a.CreateAgent(&createConf); err != nil {
			errorCh.Report(err, exceptCode.ErrorAgentUpdateFailed, "update agent setting failed(at create new agent)", true)
			return err
		}
		m.Cycle = createConf.AgentInfo.Cycle
		if createConf.DriverContainerInfo == nil && createConf.AgentContainerInfo == nil && createConf.AgentInfo.Address != "" {
			// -> remote
			m.Address = createConf.AgentInfo.Address
		}
	}

	defer func(db repo.GatewayDbOperator, txn *gorm.DB, id string, agent *model.Agent) {
		err := db.UpdateAgent(txn, id, agent)
		if err != nil {
			errorCh.Report(err, exceptCode.UpdateAgentFailed, "update db agent info failed due to:%v", true)
		}
	}(a.db, txn, req.GetAgentId(), m)

	if err := ctrl.StopGather(); err != nil {
		errorCh.Report(err, exceptCode.StopAgentFailed, "stop gather failed due to:%v", true)
		return err
	}

	setResp, err := ctrl.AgentCli.SetAgent(context.Background(),
		&agentpb.SetAgentRequest{AgentID: req.GetAgentId(), AgentInfo: req.Info})
	if err != nil {
		errorCh.Report(err, exceptCode.SetAgentFailed, "set agent failed due to:%v", true)
		return err
	}
	if grpcCode.Code(setResp.GetInfo().GetCode()) != grpcCode.Code_OK ||
		(grpcCode.Code(setResp.GetInfo().GetCode()) == grpcCode.Code_INVALID_ARGUMENT && setResp.GetInfo().GetMessage() != "Gather is not working") {
		errorCh.Report(err, exceptCode.SetAgentFailed, "set agent failed due to:%s", true, setResp.GetInfo().GetMessage())
		return err
	}

	_ = ctrl.Start(context.Background())
	if err = ctrl.StartGather(); err != nil {
		errorCh.Report(err, exceptCode.StartAgentFailed, "online agent failed", true)
		return err
	}
	return nil
}
func (a *agentApis) SetAgent(id string, info *agentpb.AgentInfo) error {
	panic("implement me")
}
func (a *agentApis) GetAgentInfo(id string) (*agentpb.AgentInfo, error) {
	panic("implement me")
}
func (a *agentApis) GetAgentModel(id string) (*model.Agent, error) {
	panic("implement me")
}
func (a *agentApis) AddAgent(req *model.CreateAgentRequest) error {
	return a.CreateAgent(req)
}
func (a *agentApis) CreateAgent(req *model.CreateAgentRequest) error {
	txn := a.db.Begin()
	errorCh := errCh.NewErrorChan()
	defer errCh.HandleErrorCh(errorCh).ErrorWillDo(func(error) {
		a.db.Rollback(txn)
	}).SuccessWillDo(func() {
		a.db.Commit(txn)
	}).Do()

	var instance model.AgentInstance
	instance.CreateConf, _ = sonic.MarshalString(req.CreateAgentConf)

	// 处理添加不在本机agent的情况
	if req.CreateAgentConf.AgentContainerInfo == nil && req.CreateAgentConf.DriverContainerInfo == nil &&
		req.AgentInfo.Address != "" {
		instance.AgentId = req.AgentInfo.AgentId
		instance.IP = req.AgentInfo.Address
		instance.Online = true
	} else {
		var driverIP, agentIP string
		if req.CreateAgentConf.DriverContainerInfo != nil {
			resp, err := docker.Create(context.Background(), req.CreateAgentConf.DriverContainerInfo)
			if err != nil {
				errorCh.Report(err, exceptCode.AddAgentFailed, "add edge driver failed", true)
				return err
			}
			driverIP, err = docker.GetContainerIPAddress(context.Background(), resp.ID)
			if err != nil {
				errorCh.Report(err, exceptCode.AddAgentFailed, "failed to get driver container ip", true)
				return err
			}
		}
		if req.CreateAgentConf.AgentContainerInfo != nil {
			req.CreateAgentConf.AgentContainerInfo.Env["DRIVER_IP"] = driverIP
			resp, err := docker.Create(context.Background(), req.CreateAgentConf.AgentContainerInfo)
			if err != nil {
				errorCh.Report(err, exceptCode.AddAgentFailed, "add edge agent failed", true)
				return err
			}
			agentIP, err = docker.GetContainerIPAddress(context.Background(), resp.ID)
			if err != nil {
				errorCh.Report(err, exceptCode.AddAgentFailed, "failed to get agent container ip", true)
				return err
			}
			instance.ContainerID = resp.ID
			instance.Local = true
			instance.AgentId = req.AgentInfo.AgentId
			instance.IP = fmt.Sprintf("%s:%d", agentIP, req.CreateAgentConf.AgentContainerInfo.ServicePort)
		}
	}

	err := a.db.AddAgentInstance(txn, instance)
	if err != nil {
		errorCh.Report(err, exceptCode.AddAgentFailed, "add agent instance failed", true)
		return err
	}
	err = a.pool.JoinAgent(context.Background(), newAgentControl(&req.AgentInfo))
	if err != nil {
		errorCh.Report(err, exceptCode.AddAgentFailed, "add agent instance failed", true)
		return err
	}
	return nil
}

func (a *agentApis) isLocalAgentDevice(id string) bool {
	return a.db.GetAgentInstance(nil, id).Local
}
