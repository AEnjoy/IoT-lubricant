package services

import (
	"context"

	"github.com/aenjoy/iot-lubricant/pkg/model/request"
	"github.com/aenjoy/iot-lubricant/services/corepkg/datastore"
)

var _ IProjectService = (*ProjectService)(nil)

type ProjectService struct {
	*datastore.DataStore
}

func (p *ProjectService) AddDataStoreEngine(ctx context.Context, projectid, dsn, dataBaseType, description string) error {
	//TODO implement me
	panic("implement me")
}

func (p *ProjectService) GetProjectDataStoreEngineStatus(ctx context.Context, projectid string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (p *ProjectService) UpdateEngineInfo(ctx context.Context, projectid, dsn, dataBaseType, description string) error {
	//TODO implement me
	panic("implement me")
}

func (p *ProjectService) BindProject(ctx context.Context, projectid string, agents []string) error {
	//TODO implement me
	panic("implement me")
}

func (p *ProjectService) AddWasher(ctx context.Context, req *request.AddWasherRequest) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (p *ProjectService) BindWasher(ctx context.Context, projectid string, washerID int, agentIDs []string) error {
	//TODO implement me
	panic("implement me")
}
