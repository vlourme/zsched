package api

import (
	"maps"
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/vlourme/scheduler/pkg/task"
)

func RegisterTasks(router *gin.Engine) {
	g := router.Group("/tasks")
	g.GET("/", GetTasks)
	g.GET("/:name", GetTask)
	g.POST("/:name", PostTask)
}

// GetTasks returns all tasks
func GetTasks(c *gin.Context) {
	tasks := c.MustGet("tasks").(map[string]*task.Task)
	c.JSON(http.StatusOK, slices.Collect(maps.Values(tasks)))
}

// GetTask returns a task by name
func GetTask(c *gin.Context) {
	tasks := c.MustGet("tasks").(map[string]*task.Task)

	var t *task.Task
	t, ok := tasks[c.Param("name")]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(http.StatusOK, t)
}

// PostTask dispatches a task
func PostTask(c *gin.Context) {
	tasks := c.MustGet("tasks").(map[string]*task.Task)
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
