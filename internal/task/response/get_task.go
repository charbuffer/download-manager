package response

import "github.com/charbuffer/download-manager/internal/entity"

type GetTask struct {
	Task entity.Task `json:"task"`
}
