package response

import "github.com/charbuffer/download-manager/internal/entity"

type AddTask struct {
	Task entity.Task `json:"task"`
}
