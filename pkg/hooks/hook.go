package hooks

import (
	"github.com/vlourme/scheduler/pkg/state"
	"github.com/vlourme/scheduler/pkg/task"
)

// Hook is the interface for the hooks
type Hook interface {
	// BeforeExecute is called before the task is executed
	BeforeExecute(task *task.Task, state *state.State) error

	// AfterExecute is called after the task is executed
	AfterExecute(task *task.Task, state *state.State) error
}
