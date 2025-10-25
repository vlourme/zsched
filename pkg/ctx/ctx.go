package ctx

import (
	"github.com/sirupsen/logrus"
	"github.com/vlourme/scheduler/pkg/logger"
	"github.com/vlourme/scheduler/pkg/state"
)

type Task interface {
	// Execute the same task with parameters, pass a state to log as children task
	Execute(params map[string]any, state ...*state.State) error

	// Task name in the queue
	Name() string
}

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

// UserContext returns the user context of the task
func (c *C) UserContext() any {
	return c.userContext
}
