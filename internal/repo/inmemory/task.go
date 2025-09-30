package inmemory

import (
	"maps"
	"slices"
	"sync"
	"sync/atomic"
	"time"

	"github.com/charbuffer/download-manager/internal/entity"
	"github.com/charbuffer/download-manager/internal/repo"
)

var _ repo.TaskRepo = (*TaskRepo)(nil)

type TaskRepo struct {
	data   map[int32]*entity.Task
	mu     *sync.RWMutex
	length int32
}

func NewTaskRepo() *TaskRepo {
	return &TaskRepo{
		data:   map[int32]*entity.Task{},
		mu:     &sync.RWMutex{},
		length: 0,
	}
}

func (r *TaskRepo) AddTask(task entity.Task) entity.Task {
	r.mu.Lock()
	defer r.mu.Unlock()

	newLength := atomic.AddInt32(&r.length, 1)
	r.length = newLength

	task.Id = &newLength
	for i := range task.Files {
		id := int32(i)
		task.Files[i].Id = &id
	}
	r.data[newLength] = &task

	return task
}

func (r *TaskRepo) GetAllTasks() []*entity.Task {
	r.mu.Lock()
	defer r.mu.Unlock()

	tasks := maps.Values(r.data)

	return slices.Collect(tasks)
}

func (r *TaskRepo) GetTask(id int32) *entity.Task {
	r.mu.Lock()
	defer r.mu.Unlock()

	task, ok := r.data[id]

	if !ok {
		return nil
	}

	return task
}

func (r *TaskRepo) UpdateFileStatus(taskId int32, fileId int32, status entity.FileStatus) *entity.Task {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, ok := r.data[taskId]
	if !ok {
		return nil
	}

	r.data[taskId].Files[fileId].Status = status
	r.data[taskId].CreatedAt = time.Now()

	task := r.data[taskId]

	return task
}
