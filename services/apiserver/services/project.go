package services

import (
	"context"
	"sync"

	"github.com/aenjoy/iot-lubricant/pkg/model"
	"github.com/aenjoy/iot-lubricant/pkg/model/request"
	exceptionCode "github.com/aenjoy/iot-lubricant/pkg/types/exception/code"
	"github.com/aenjoy/iot-lubricant/services/corepkg/datastore"
	"github.com/aenjoy/iot-lubricant/services/datastoreAssistant/api"
	logg "github.com/aenjoy/iot-lubricant/services/logg/api"
)

var _ IProjectService = (*ProjectService)(nil)

type ProjectService struct {
	*datastore.DataStore
}

func (p *ProjectService) GetProject(ctx context.Context, projectid string) (model.Project, error) {
	return p.DataStore.GetProject(ctx, projectid)
}

func (p *ProjectService) ListProject(ctx context.Context, userID string) ([]model.Project, error) {
	return p.DataStore.ListProject(ctx, userID)
}

func (p *ProjectService) AddProject(ctx context.Context, userid, projectid, projectname, description string) (string, error) {
	txn, _, commit := p.txnHelper()
	defer commit()
	return projectid, p.DataStore.AddProject(ctx, txn, userid, projectid, projectname, description)
}

func (p *ProjectService) RemoveProject(ctx context.Context, projectid string, removeAgent, removeGateway *bool) error {
	txn, errCh, commit := p.txnHelper()
	defer commit()

	var gatewayIds []string
	if removeAgent != nil && *removeAgent {
		ids, err := p.DataStore.ICoreDb.GetAgentsByProjectID(ctx, txn, projectid)
		if err != nil {
			errCh.Report(err, exceptionCode.ListAgentFailed, "failed to list agents at RemoveProject", false)
			return err
		}
		for _, id := range ids {
			gatewayIds = append(gatewayIds, id.GatewayId)
			err = p.DataStore.DeleteAgent(ctx, txn, id.AgentId)
			if err != nil {
				errCh.Report(err, exceptionCode.RemoveAgentFailed, "failed to remove agent at RemoveProject", false)
				return err
			}
		}
	}

	if removeGateway != nil && *removeGateway {
		for _, id := range gatewayIds {
			err := p.DataStore.DeleteGateway(ctx, txn, id)
			if err != nil {
				errCh.Report(err, exceptionCode.RemoveGatewayFailed, "failed to remove gateway at RemoveProject", false)
				return err
			}
		}
	}

	return p.DataStore.RemoveProject(ctx, txn, projectid)
}

func (p *ProjectService) AddDataStoreEngine(ctx context.Context, projectid, dsn, dataBaseType, description string) error {
	txn, _, commit := p.txnHelper()
	defer commit()

	err := p.DataStore.AddDataStoreEngine(ctx, txn, projectid, dsn, dataBaseType, description)
	if err != nil {
		return err
	}
	return nil
}

func (p *ProjectService) GetProjectDataStoreEngineStatus(ctx context.Context, projectid, userId string) (string, error) {
	engine, err := p.DataStore.GetEngineByProjectID(ctx, projectid)
	if err != nil {
		logg.L.Errorf("failed to get project data store engine status for project %s", projectid)
		return "Failed", err
	}

	return api.DsnTest(engine.DataBaseType, engine.DSN, userId), nil
}

func (p *ProjectService) UpdateEngineInfo(ctx context.Context, projectid, dsn, dataBaseType, description string) error {
	txn, _, commit := p.txnHelper()
	defer commit()

	return p.DataStore.UpdateEngineInfo(ctx, txn, projectid, dsn, dataBaseType, description)
}

func (p *ProjectService) BindProject(ctx context.Context, projectid string, agents []string) error {
	txn, _, commit := p.txnHelper()
	defer commit()

	return p.DataStore.BindProject(ctx, txn, projectid, agents)
}

var addWasherMutex sync.Mutex // todo:其实应该用分布式锁
func (p *ProjectService) AddWasher(ctx context.Context, req *request.AddWasherRequest) (int, error) {
	txn, errCh, commit := p.txnHelper()
	defer commit()

	addWasherMutex.Lock()
	number, err := p.DataStore.GetProjectAgentNumber(ctx, txn, req.ProjectID)
	addWasherMutex.Unlock()

	if err != nil {
		errCh.Report(err, exceptionCode.GetProjectAgentNumberFailed, "failed to get project agent number", false)
		return 0, err
	}
	if len(req.ToAgents) == 0 {
		req.ToAgents = append(req.ToAgents, "") // 代表稍后手动添加
	}
	for i := 0; i < len(req.ToAgents); i++ {
		err = p.DataStore.AddWasher(ctx, txn, &model.Clean{
			WasherID:    number + 1 + i,
			AgentID:     req.ToAgents[i],
			Description: req.Description,
			ProjectID:   req.ProjectID,
			Table:       req.Table,
			Interpreter: req.Interpreter,
			Script:      req.Script,
			Command:     req.Command,
		})
		if err != nil {
			errCh.Report(err, exceptionCode.AddWasherFailed, "failed to add washer to agents:%s", false, req.ToAgents[i])
			return 0, err
		}
	}
	return number + 1, nil
}

func (p *ProjectService) BindWasher(ctx context.Context, projectid string, washerID int, agentIDs []string) error {
	txn, _, commit := p.txnHelper()
	defer commit()

	return p.DataStore.BindWasher(ctx, txn, projectid, washerID, agentIDs)
}
