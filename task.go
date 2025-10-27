package zsched

import (
	"regexp"
)

// nameRegex is the regex to validate the task name
var nameRegex = regexp.MustCompile(`[^a-zA-Z0-9-._]+`)

// TaskAction is the function that performs the task
type taskAction[T any] func(ctx *Context[T]) error

type taskSchedule struct {
	// Schedule is a cron expression with seconds precision (e.g. "0 0 * * * *").
	Schedule string `json:"schedule"`

	// Parameters is the parameters for the task
	Parameters map[string]any `json:"parameters"`
}

type taskConfig struct {
	// collectorAction is the action to be performed by the collector
	collectorAction TaskCollectorAction `json:"-"`

	// collector is the collector for the task
	collector *Collector `json:"-"`

	// Concurrency is the number of concurrent tasks to run
	Concurrency int `json:"concurrency"`

	// MaxRetries is the maximum number of retries for the task
	MaxRetries int `json:"max_retries"`

	// Schedules is the schedules for the task
	Schedules []taskSchedule `json:"schedules"`
}

type Task[T any] struct {
	taskConfig

	// Name of the task, should be unique without any spaces or special characters
	TaskName string `json:"name"`

	// Action to be performed by the task
	Action taskAction[T] `json:"-"`

	// executor is the executor for the task
	executor *executor[T] `json:"-"`
}

// NewTask creates a new task
func NewTask[T any](name string, action taskAction[T], opts ...func(*taskConfig)) *Task[T] {
	t := &Task[T]{
		TaskName: name,
		Action:   action,
		taskConfig: taskConfig{
			collectorAction: nil,
			collector:       nil,
			Concurrency:     1,
			MaxRetries:      3,
			Schedules:       make([]taskSchedule, 0),
		},
	}

	for _, opt := range opts {
		opt(&t.taskConfig)
	}

	return t
}

// Execute executes the task
func (t *Task[T]) Execute(params ...map[string]any) error {
	var p map[string]any
	if len(params) > 0 {
		p = params[0]
	}

	return t.ExecuteWithState(p)
}

// ExecuteWithState executes the task with a parent state
func (t *Task[T]) ExecuteWithState(params map[string]any, parentId ...*State) error {
	s := newState(params)

	if len(parentId) > 0 {
		s.ParentID = parentId[0].TaskID
	}

	return t.executor.Publish(t, s)
}

// formatName formats the name of the task to a valid RabbitMQ queue name
// TODO: Move to a separate package
func (t *Task[T]) Name() string {
	return nameRegex.ReplaceAllString(t.TaskName, "")
}

// Collector returns the collector for the task
func (t *Task[T]) Collector() *Collector {
	return t.collector
}

// WithConcurrency sets the concurrency for the task
func WithConcurrency(concurrency int) func(*taskConfig) {
	return func(t *taskConfig) {
		t.Concurrency = concurrency
	}
}

// WithMaxRetries sets the max retries for the task, default is 3. -1 means infinite retries.
func WithMaxRetries(maxRetries int) func(*taskConfig) {
	return func(t *taskConfig) {
		t.MaxRetries = maxRetries
	}
}

// WithCollector sets the collector for the task with optional buffer size
func WithCollector(collector TaskCollectorAction, bufferSize ...int) func(*taskConfig) {
	return func(t *taskConfig) {
		t.collector = newCollector(bufferSize...)
		t.collectorAction = collector
	}
}

// WithSchedule adds a schedule to the task.
// The schedule is a cron expression with seconds precision (e.g. "0 0 * * * *").
func WithSchedule(schedule string, parameters map[string]any) func(*taskConfig) {
	return func(t *taskConfig) {
		t.Schedules = append(t.Schedules, taskSchedule{
			Schedule:   schedule,
			Parameters: parameters,
		})
	}
}
