package entity

import (
	"time"
)

type Task struct {
	Id        *int32    `json:"id"`
	Files     []File    `json:"files"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type File struct {
	Id     *int32     `json:"id"`
	Url    string     `json:"url"`
	Status FileStatus `json:"status"`
}

type FileStatus string

const (
	StatusPending    FileStatus = "pending"
	StatusProcessing FileStatus = "processing"
	StatusCompleted  FileStatus = "completed"
	StatusFailed     FileStatus = "failed"
)

func NewTask(urls []string) *Task {
	files := make([]File, 0, len(urls))

	for _, url := range urls {
		files = append(files, File{Url: url, Status: StatusPending})
	}

	now := time.Now()

	return &Task{
		Id:        nil,
		Files:     files,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
