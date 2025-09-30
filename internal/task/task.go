package task

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/charbuffer/download-manager/internal/entity"
	"github.com/charbuffer/download-manager/internal/repo"
	"github.com/charbuffer/download-manager/internal/repo/inmemory"
	"github.com/charbuffer/download-manager/internal/task/request"
	"github.com/charbuffer/download-manager/internal/task/response"
	"github.com/charbuffer/download-manager/internal/worker"
	"github.com/charbuffer/download-manager/pkg/utils"
	"github.com/gin-gonic/gin"
)

type TaskService struct {
	repo   repo.TaskRepo
	worker worker.TaskWorkerPool
}

func NewTaskHandler(repo repo.TaskRepo, workersCount int) *TaskService {
	return &TaskService{
		repo:   inmemory.NewTaskRepo(),
		worker: *worker.NewTaskWorkerPool(workersCount),
	}
}

func (s *TaskService) AddTask(ctx *gin.Context) {
	var req request.AddTask

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
	}

	urls := utils.RemoveURLDuplicates(req.Urls)

	task := s.repo.AddTask(*entity.NewTask(urls))

	for _, file := range task.Files {
		log.Printf("saving file %s...\n", file.Url)

		go s.repo.UpdateFileStatus(*task.Id, *file.Id, entity.StatusProcessing)
		s.worker.Submit(worker.NewJob(*task.Id, *file.Id, file.Url))
	}

	go func() {
		for r := range s.worker.Results() {
			s.repo.UpdateFileStatus(r.TaskId, r.FileId, r.FileStatus)
		}
	}()

	ctx.JSON(http.StatusCreated, response.AddTask{Task: task})
}

func (s *TaskService) GetAllTasks(ctx *gin.Context) {
	fmt.Println("GetAll")

	tasks := s.repo.GetAllTasks()

	ctx.JSON(http.StatusOK, response.GetAllTasks{Tasks: tasks})
}

func (s *TaskService) GetTask(ctx *gin.Context) {
	stringId := ctx.Param("id")

	id, err := strconv.Atoi(stringId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
		return
	}

	task := s.repo.GetTask(int32(id))
	if task == nil {
		ctx.JSON(http.StatusNotFound, gin.H{})
		return
	}

	ctx.JSON(http.StatusOK, response.GetTask{Task: *task})
}

func (s *TaskService) Shutdown() {
	s.worker.Shutdown()
}
