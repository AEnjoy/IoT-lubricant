package async

import (
	"github.com/AEnjoy/IoT-lubricant/internal/cache"
	"github.com/AEnjoy/IoT-lubricant/protobuf/core"
)

type Task interface {
	AddTask(task *core.TaskDetail)
	Query(id string) *core.QueryTaskResultResponse
	SetActor(f func(*core.TaskDetail, *cache.MemoryCache[*core.QueryTaskResultResponse]))
	Release()
}

func NewAsyncTask() Task {
	t := new(task)
	t.init()
	return t
}
