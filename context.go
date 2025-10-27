package zsched

import (
	"github.com/sirupsen/logrus"
	"github.com/vlourme/zsched/pkg/logger"
)

// Context is a temporary object into the task execution context
// It allow logging, outputting values and accessing the user context
type Context[T any] struct {
	logger.Logger
	State

	task        Task[T]
	userContext T
}

func newContext[T any](task *Task[T], state State, logger logger.Logger, userContext T) *Context[T] {
	return &Context[T]{
		Logger: logger.WithFields(logrus.Fields{
			"scope":    task.Name(),
			"state_id": state.ID,
			"task_id":  state.TaskID,
		}),
		userContext: userContext,
		task:        *task,
		State:       state,
	}
}

// Execute starts the same task with given parameters
func (c *Context[T]) Execute(params ...map[string]any) error {
	var p map[string]any
	if len(params) > 0 {
		p = params[0]
	}

	return c.task.ExecuteWithState(p, &c.State)
}

// Push pushes a value to the collector
func (c *Context[T]) Push(value any) {
	c.task.Collector().Push(value)
}

// UserContext returns the user context of the task
func (c *Context[T]) UserContext() T {
	return c.userContext
}
