package async

import (
	"sync"

	"github.com/aenjoy/iot-lubricant/pkg/cache"
	"github.com/aenjoy/iot-lubricant/protobuf/core"
)

type Task interface {
	AddTask(task *core.TaskDetail, notice bool)
	RemoveResult(id string)
	Query(id string) *core.QueryTaskResultResponse
	SetActor(f func(*core.TaskDetail, *cache.MemoryCache[*core.QueryTaskResultResponse]))
	GetNotifyCh() <-chan string
	Release()
}

var _task *task
var _taskL sync.Mutex

func NewAsyncTask() Task {
	_taskL.Lock()
	defer _taskL.Unlock()
	if _task == nil {
		_task = new(task)
		_task.init()
	}
	return _task
}
