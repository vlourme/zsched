package ctx

import "github.com/vlourme/zsched/pkg/state"

type Task interface {
	// Execute the same task with parameters, pass a state to log as children task
	Execute(params map[string]any, state ...*state.State) error

	// Collector returns the collector for the task
	Collector() Collector

	// Task name in the queue
	Name() string
}

type Collector interface {
	// Push a value to the collector
	Push(value any)

	// Pull pulls a value from the collector
	Pull() any

	// Channel returns the channel of the collector
	Channel() <-chan any
}
