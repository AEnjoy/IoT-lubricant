package async

import (
	"sync/atomic"

	"github.com/aenjoy/iot-lubricant/pkg/cache"
	"github.com/aenjoy/iot-lubricant/pkg/logger"
	taskCode "github.com/aenjoy/iot-lubricant/pkg/types/task"
	corepb "github.com/aenjoy/iot-lubricant/protobuf/core"
	"github.com/panjf2000/ants/v2"
	"google.golang.org/genproto/googleapis/rpc/status"
)

var _ Task = (*task)(nil)

type task struct {
	result   *cache.MemoryCache[*corepb.QueryTaskResultResponse] // id -> QueryTaskResultResponse
	queue    chan *corepb.TaskDetail
	notify   chan string
	actor    func(*corepb.TaskDetail, *cache.MemoryCache[*corepb.QueryTaskResultResponse])
	pool     *ants.Pool
	finished int32
}

func (r *task) init() {
	r.result = cache.NewMemoryCache[*corepb.QueryTaskResultResponse]()
	r.queue = make(chan *corepb.TaskDetail, 100)
	r.notify = make(chan string, 5)
	var err error
	r.pool, err = ants.NewPool(100, ants.WithPreAlloc(true))
	if err != nil {
		panic(err)
	}
	go r.run()
}
func (r *task) run() {
	for detail := range r.queue {
		err := r.pool.Submit(func() {
			r.actor(detail, r.result)
			if detail.MessageId == "notice" {
				r.notify <- detail.TaskId
			}
			atomic.AddInt32(&r.finished, 1)
		})
		if err != nil {
			logger.Errorf("submit task failed: %v", err)
			result := cache.NewStoreResult[*corepb.QueryTaskResultResponse](cache.NeverExpired, &corepb.QueryTaskResultResponse{
				TaskId: detail.TaskId,
				Result: &corepb.QueryTaskResultResponse_NotFound{
					NotFound: &status.Status{
						Message: "Task execute failed",
					},
				},
			})
			r.result.Set(detail.TaskId, "", result)
		}
	}
}

func (r *task) AddTask(task *corepb.TaskDetail, notice bool) {
	result := cache.NewStoreResult[*corepb.QueryTaskResultResponse](cache.NeverExpired,
		&corepb.QueryTaskResultResponse{
			TaskId: task.TaskId,
			Result: &corepb.QueryTaskResultResponse_Pending{
				Pending: &status.Status{
					Code:    taskCode.TaskStatusPending,
					Message: "Task is pending due to the actuator does not run this task",
				},
			},
		})
	if notice {
		task.MessageId = "notice"
	}
	logger.Debugf("Add task result cache: id: %s", task.TaskId)
	r.result.Set(task.TaskId, task.TaskId, result)
	r.queue <- task
}
func (r *task) Query(id string) *corepb.QueryTaskResultResponse {
	if v, ok := r.result.GetCache(id); ok {
		return v
	}
	return &corepb.QueryTaskResultResponse{
		TaskId: id,
		Result: &corepb.QueryTaskResultResponse_NotFound{
			NotFound: &status.Status{
				Message: "Task not found",
			},
		},
	}
}
func (r *task) SetActor(f func(*corepb.TaskDetail, *cache.MemoryCache[*corepb.QueryTaskResultResponse])) {
	r.actor = f
}
func (r *task) Release() {
	r.pool.Release()
}
func (r *task) GetNotifyCh() <-chan string {
	return r.notify
}
func (r *task) RemoveResult(id string) {
	r.result.Delete(id)
}
