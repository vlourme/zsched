package task

import (
	"regexp"

	"github.com/vlourme/zsched/pkg/ctx"
	"github.com/vlourme/zsched/pkg/state"
)

// nameRegex is the regex to validate the task name
var nameRegex = regexp.MustCompile(`[^a-zA-Z0-9-._]+`)

// TaskAction is the function that performs the task
type TaskAction func(ctx *ctx.C) error

// TaskExecutor is the interface for the task executor
type taskExecutor interface {
	// Publish publishes a message to the message broker
	Publish(task *Task, s *state.State) error
}

type TaskSchedule struct {
	// Schedule is a cron expression with seconds precision (e.g. "0 0 * * * *").
	Schedule string `json:"schedule"`

	// Parameters is the parameters for the task
	Parameters map[string]any `json:"parameters"`
}

type Task struct {
	// Name of the task, should be unique without any spaces or special characters
	TaskName string `json:"name"`

	// Action to be performed by the task
	Action TaskAction `json:"-"`

	// Concurrency of the task
	// The number of concurrent instances of the task that can run at the same time
	Concurrency int `json:"concurrency"`

	// MaxRetries is the maximum number of retries of the task
	// If the task fails, it will be retried up to this number of times
	MaxRetries int `json:"max_retries"`

	// executor is the executor for the task
	Executor taskExecutor `json:"-"`

	// Schedules is the schedules for the task
	Schedules []TaskSchedule `json:"schedules"`
}

// NewTask creates a new task
func NewTask(name string, action TaskAction, options ...func(*Task)) *Task {
	t := &Task{
		TaskName:    name,
		Action:      action,
		Concurrency: 1,
		MaxRetries:  3,
		Schedules:   make([]TaskSchedule, 0),
	}

	for _, option := range options {
		option(t)
	}

	return t
}

// Execute executes the task
func (t *Task) Execute(params map[string]any, parentId ...*state.State) error {
	s := state.NewState(params)

	if len(parentId) > 0 {
		s.ParentID = parentId[0].TaskID
	}

	return t.Executor.Publish(t, s)
}

// formatName formats the name of the task to a valid RabbitMQ queue name
// TODO: Move to a separate package
func (t *Task) Name() string {
	return nameRegex.ReplaceAllString(t.TaskName, "")
}

// WithConcurrency sets the concurrency for the task
func WithConcurrency(concurrency int) func(*Task) {
	return func(t *Task) {
		t.Concurrency = concurrency
	}
}

// WithMaxRetries sets the max retries for the task, default is 3. -1 means infinite retries.
func WithMaxRetries(maxRetries int) func(*Task) {
	return func(t *Task) {
		t.MaxRetries = maxRetries
	}
}

// WithSchedule adds a schedule to the task.
// The schedule is a cron expression with seconds precision (e.g. "0 0 * * * *").
func WithSchedule(schedule string, parameters map[string]any) func(*Task) {
	return func(t *Task) {
		t.Schedules = append(t.Schedules, TaskSchedule{
			Schedule:   schedule,
			Parameters: parameters,
		})
	}
}
