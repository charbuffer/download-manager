package request

type AddTask struct {
	Urls []string `json:"urls" binding:"required"`
}
