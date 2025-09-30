package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/charbuffer/download-manager/internal/repo/inmemory"
	"github.com/charbuffer/download-manager/internal/task"
	"github.com/gin-gonic/gin"
)

type App struct {
	router  *gin.Engine
	server  *http.Server
	handler task.TaskService
}

func NewApp(router *gin.Engine, config Config) *App {
	handler := task.NewTaskHandler(inmemory.NewTaskRepo(), config.workers)
	router.GET("/task/:id", handler.GetTask)
	router.GET("/task", handler.GetAllTasks)
	router.POST("/task", handler.AddTask)

	return &App{
		router: router,
		server: &http.Server{
			Addr:    fmt.Sprintf(":%d", config.port),
			Handler: router.Handler(),
		},
		handler: *handler,
	}
}

func (a *App) Run() {
	go func() {
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down a server ..")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := a.Shutdown(ctx); err != nil {
		log.Println("Server Shutdown: ", err)
	}

	log.Println("Server exiting")
}

func (a *App) Shutdown(ctx context.Context) error {
	a.handler.Shutdown()
	return a.server.Shutdown(ctx)
}
