package api

import (
	"github.com/gin-gonic/gin"
	"github.com/vlourme/scheduler/pkg/storage"
	"github.com/vlourme/scheduler/pkg/task"
)

type RegisterFn = func(router *gin.Engine)

var groups = []RegisterFn{
	RegisterTasks,
}

func NewAPI(tasks map[string]*task.Task, storage storage.Storage) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.Use(func(ctx *gin.Context) {
		ctx.Set("tasks", tasks)
		ctx.Set("storage", storage)
	})

	for _, group := range groups {
		group(router)
	}

	return router
}
