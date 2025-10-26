package executor

import (
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/vlourme/zsched/pkg/broker"
	"github.com/vlourme/zsched/pkg/ctx"
	"github.com/vlourme/zsched/pkg/hooks"
	"github.com/vlourme/zsched/pkg/logger"
	"github.com/vlourme/zsched/pkg/state"
	"github.com/vlourme/zsched/pkg/task"
)

// Executor is the executor for the tasks
type Executor struct {
	// broker is the broker for the tasks
	broker broker.Broker

	// logger is the logger for the executor
	logger logger.Logger

	// hooks are the hooks for the executor
	hooks []hooks.Hook

	// userContext is the user context for the executor
	userContext any
}

// NewExecutor creates a new executor
func NewExecutor(broker broker.Broker, logger logger.Logger, hooks []hooks.Hook, userContext any) *Executor {
	return &Executor{
		broker:      broker,
		logger:      logger,
		hooks:       hooks,
		userContext: userContext,
	}
}

// Publish publishes the task to the broker
func (e *Executor) Publish(task *task.Task, s *state.State) error {
	body, err := state.Serialize(s)
	if err != nil {
		return err
	}

	return e.broker.Publish(body, task.Name())
}

// Consume listens for events from the broker and executes the task
func (e *Executor) Consume(task *task.Task) error {
	if task.CollectorAction != nil {
		go task.CollectorAction(task.Collector())
	}

	return e.broker.Consume(task.Name(), task.Concurrency, func(body []byte) error {
		s, err := state.Deserialize(body)
		if err != nil {
			return err
		}

		s.ID = uuid.New()
		s.StartedAt = time.Now()
		s.Iterations++

		defer func() {
			if r := recover(); r != nil {
				log.Printf("recovered from a panic: [%s]%s: %v", task.Name(), s.ID, r)
			}
		}()

		if err := e.runBeforeExecuteHooks(task, s); err != nil {
			log.Printf("failed to run before execute hooks: %v", err)
		}

		s.Status = state.StatusRunning
		ctx := ctx.New(e.logger, task, *s, e.userContext)
		if err := task.Action(ctx); err != nil && (task.MaxRetries == -1 || s.Iterations < task.MaxRetries) {
			e.logger.WithError(err).WithField("task_name", task.Name()).Print("task errored")
			s.LastError = err.Error()
			s.Status = state.StatusFailed
			if err := e.runAfterExecuteHooks(task, s); err != nil {
				log.Printf("failed to run after execute hooks: %v", err)
			}
			return e.Publish(task, s)
		}

		s.Status = state.StatusSuccess
		if err := e.runAfterExecuteHooks(task, s); err != nil {
			log.Printf("failed to run after execute hooks: %v", err)
		}

		return nil
	})
}

// runBeforeExecuteHooks runs the before execute hooks
func (e *Executor) runBeforeExecuteHooks(task *task.Task, s *state.State) error {
	for _, hook := range e.hooks {
		if err := hook.BeforeExecute(task, s); err != nil {
			return err
		}
	}
	return nil
}

// runAfterExecuteHooks runs the after execute hooks
func (e *Executor) runAfterExecuteHooks(task *task.Task, s *state.State) error {
	for _, hook := range e.hooks {
		if err := hook.AfterExecute(task, s); err != nil {
			return err
		}
	}
	return nil
}
