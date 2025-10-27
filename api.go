package zsched

import (
	"maps"
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/vlourme/zsched/pkg/storage"
)

func newRouter[T any](tasks map[string]*Task[T], storage storage.Storage) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.Use(func(ctx *gin.Context) {
		ctx.Set("tasks", tasks)
		ctx.Set("storage", storage)
	})

	router.GET("/tasks", GetTasks[T])
	router.GET("/tasks/:name", GetTask[T])
	router.POST("/tasks/:name", PostTask[T])

	return router
}

// GetTasks returns all tasks
func GetTasks[T any](c *gin.Context) {
	tasks := c.MustGet("tasks").(map[string]*Task[T])
	c.JSON(http.StatusOK, slices.Collect(maps.Values(tasks)))
}

// GetTask returns a task by name
func GetTask[T any](c *gin.Context) {
	tasks := c.MustGet("tasks").(map[string]*Task[T])

	var t *Task[T]
	t, ok := tasks[c.Param("name")]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(http.StatusOK, t)
}

// PostTask dispatches a task
func PostTask[T any](c *gin.Context) {
	tasks := c.MustGet("tasks").(map[string]*Task[T])
	t, ok := tasks[c.Param("name")]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	body := make(map[string]any)
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := t.Execute(body); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task dispatched successfully"})
}
