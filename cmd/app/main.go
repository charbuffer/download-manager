package main

import (
	"github.com/charbuffer/download-manager/internal/app"
	"github.com/gin-gonic/gin"
)

func main() {
	config := app.NewConfig(4000, 10)
	router := gin.Default()

	app := app.NewApp(router, config)
	app.Run()
}
