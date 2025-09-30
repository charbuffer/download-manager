package response

import "github.com/charbuffer/download-manager/internal/entity"

type GetAllTasks struct {
	Tasks []*entity.Task `json:"tasks"`
}
