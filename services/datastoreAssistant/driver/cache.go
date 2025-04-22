package driver

import (
	"context"
	"fmt"
	"sync"

	"github.com/aenjoy/iot-lubricant/services/corepkg/repo"
)

var _repo repo.ICoreDb

var nodeIdMap sync.Map // (nodeName,userID)->nodeId

func SetRepo(repo repo.ICoreDb) {
	_repo = repo
}
func nodeNameGetNodeId(nodeName, userID string) (string, error) {
	//todo: change to redis
	if nodeId, ok := nodeIdMap.Load(fmt.Sprintf("%s,%s", nodeName, userID)); ok {
		return nodeId.(string), nil
	}
	nodeId, err := _repo.GetAgentIDByAgentNameAndUserID(context.Background(), nodeName, userID)
	if err != nil {
		return "", err
	}
	nodeIdMap.Store(nodeName, nodeId)
	return nodeId, nil
}
