package zsched

import "github.com/vlourme/zsched/pkg/storage"

// AnyTask is the interface that represents any task, no matter
// the inner type of the user context
type AnyTask interface {
	Name() string
}

// Hook is the interface for the hooks
type Hook interface {
	// Initialize is called when the engine is initialized
	Initialize(storage storage.Storage) error

	// BeforeExecute is called before the task is executed
	BeforeExecute(task AnyTask, state *State) error

	// AfterExecute is called after the task is executed
	AfterExecute(task AnyTask, state *State) error
}
