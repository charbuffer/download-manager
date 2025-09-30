package repo

import "github.com/charbuffer/download-manager/internal/entity"

type TaskRepo interface {
	AddTask(task entity.Task) entity.Task
	GetAllTasks() []*entity.Task
	GetTask(id int32) *entity.Task
	UpdateFileStatus(taskId int32, fileId int32, status entity.FileStatus) *entity.Task
}
