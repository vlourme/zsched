package ctx

import (
	"log"

	"github.com/sirupsen/logrus"
	"github.com/vlourme/zsched/pkg/logger"
	"github.com/vlourme/zsched/pkg/state"
)

type C struct {
	logger.Logger
	state.State

	task        Task
	userContext any
}

func New(log logger.Logger, task Task, state state.State, userContext any) *C {
	return &C{
		task:        task,
		userContext: userContext,
		State:       state,
		Logger: log.
			WithFields(logrus.Fields{
				"state_id": state.ID,
				"task_id":  state.TaskID,
			}),
	}
}

// Execute starts a new task of the same job, with different parameters
func (c *C) Execute(params map[string]any) error {
	return c.task.Execute(params, &c.State)
}

// Push pushes a value to the collector
func (c *C) Push(value any) {
	if c.task.Collector() == nil {
		log.Printf("task %s has no collector", c.task.Name())
		return
	}

	c.task.Collector().Push(value)
}

// UserContext returns the user context of the task
func (c *C) UserContext() any {
	return c.userContext
}
