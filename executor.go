package zsched

import (
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/vlourme/zsched/pkg/broker"
	"github.com/vlourme/zsched/pkg/logger"
)

// Executor is the executor for the tasks
type executor[T any] struct {
	broker      broker.Broker
	logger      logger.Logger
	hooks       []Hook
	userContext T
}

// Publish publishes the task to the broker
func (e *executor[T]) Publish(task *Task[T], s *State) error {
	s.ID = uuid.New()
	body, err := s.Serialize()
	if err != nil {
		return err
	}

	if err := e.runBeforeExecuteHooks(task, s); err != nil {
		log.Printf("failed to run before execute hooks: %v", err)
	}

	return e.broker.Publish(body, task.Name())
}

// Consume listens for events from the broker and executes the task
func (e *executor[T]) Consume(task *Task[T]) error {
	if task.collectorAction != nil {
		go task.collectorAction(task.collector, e.userContext)
	}

	return e.broker.Consume(task.Name(), task.Concurrency, func(body []byte) error {
		s, err := deserializeState(body)
		if err != nil {
			return err
		}

		s.Status = StatusRunning
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

		ctx := newContext(task, *s, e.logger, e.userContext)
		if err := task.Action(ctx); err != nil && (task.MaxRetries == -1 || s.Iterations < task.MaxRetries) {
			e.logger.WithError(err).WithField("task_name", task.Name()).Print("task errored")
			s.LastError = err.Error()
			s.Status = StatusFailed
			if err := e.runAfterExecuteHooks(task, s); err != nil {
				log.Printf("failed to run after execute hooks: %v", err)
			}
			return e.Publish(task, s)
		}

		s.Status = StatusSuccess
		if err := e.runAfterExecuteHooks(task, s); err != nil {
			log.Printf("failed to run after execute hooks: %v", err)
		}

		return nil
	})
}

// runBeforeExecuteHooks runs the before execute hooks
func (e *executor[T]) runBeforeExecuteHooks(task *Task[T], s *State) error {
	for _, hook := range e.hooks {
		if err := hook.BeforeExecute(task, s); err != nil {
			return err
		}
	}
	return nil
}

// runAfterExecuteHooks runs the after execute hooks
func (e *executor[T]) runAfterExecuteHooks(task *Task[T], s *State) error {
	for _, hook := range e.hooks {
		if err := hook.AfterExecute(task, s); err != nil {
			return err
		}
	}
	return nil
}
