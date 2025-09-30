package worker

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/charbuffer/download-manager/internal/entity"
)

type Job struct {
	TaskId int32
	FileId int32
	Url    string
}

type Result struct {
	Job
	entity.FileStatus
}

func NewJob(taskId, fileId int32, url string) Job {
	return Job{
		TaskId: taskId,
		FileId: fileId,
		Url:    url,
	}
}

type TaskWorkerPool struct {
	wg      *sync.WaitGroup
	jobs    chan Job
	results chan Result
	ctx     context.Context
	cancel  context.CancelFunc
}

func NewTaskWorkerPool(workerCount int) *TaskWorkerPool {
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	wp := &TaskWorkerPool{
		wg:      wg,
		jobs:    make(chan Job),
		results: make(chan Result),
		ctx:     ctx,
		cancel:  cancel,
	}

	for i := range workerCount {
		wp.wg.Add(1)
		go wp.worker(i)
	}

	return wp
}

func (wp *TaskWorkerPool) Submit(job Job) {
	wp.jobs <- job
}

func (wp *TaskWorkerPool) Results() <-chan Result {
	return wp.results
}

func (wp *TaskWorkerPool) Shutdown() {
	close(wp.jobs)
	wp.wg.Wait()
	close(wp.results)
	wp.cancel()
}

func (wp *TaskWorkerPool) worker(i int) {
	defer func() {
		log.Printf("Worker %d done\n", i)
		wp.wg.Done()
	}()

	log.Printf("Worker %d stated\n", i)
	for {
		select {
		case <-wp.ctx.Done():
			fmt.Println("done")
			return
		case job, ok := <-wp.jobs:
			if !ok {
				return
			}
			_, err := downloadFile(job.Url, "downloads")
			if err != nil {
				wp.results <- Result{job, entity.StatusFailed}
			} else {
				wp.results <- Result{job, entity.StatusCompleted}
			}

		}
	}
}

func downloadFile(fileUrl, folder string) (string, error) {
	resp, err := http.Get(fileUrl)
	if err != nil {
		log.Println("can't get file from url")
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("can't get file from url: bad status")
		return "", err
	}

	var filename string
	url, _ := url.Parse(fileUrl)

	if resp.Header.Get("Content-Disposition") == "attachment" {
		path := strings.Split(url.Path, "/")
		filename = path[len(path)-1]
	} else {
		filename = url.Host
	}

	filename = fmt.Sprintf("%s-%d", filename, time.Now().Unix())

	out, err := os.Create(fmt.Sprintf("%s/%s", folder, filename))
	if err != nil {
		return "", err

	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Println("can't write response to file")
		return "", err
	}

	log.Printf("Saved: %v\n", filename)

	return filename, nil
}
