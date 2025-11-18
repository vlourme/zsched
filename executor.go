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
	taskLogger  *taskLogger[T]
	broker      broker.BrokerQueue
	logger      logger.Logger
	hooks       []Hook
	userContext T
}

// Publish publishes one or many executions to the broker
func (e *executor[T]) Publish(task *Task[T], state *State) error {
	state.ID = uuid.New()
	body, err := state.Serialize()
	if err != nil {
		return err
	}

	if err := e.taskLogger.LogTasks(task, state); err != nil {
		log.Printf("failed to log execution: %v", err)
	}

	if err := e.runBeforeExecuteHooks(task, state); err != nil {
		log.Printf("failed to run before execute hooks: %v", err)
	}

	if err := e.broker.Publish(body); err != nil {
		log.Printf("failed to publish task: %v", err)
		return err
	}

	return nil
}

// Consume listens for events from the broker and executes the task
func (e *executor[T]) Consume(task *Task[T]) error {
	if task.collectorAction != nil {
		go task.collectorAction(task.collector, e.userContext)
	}

	return e.broker.Consume(
		task.MaxRetries == 0, // prevent re-shipping on broker restart
		task.Concurrency,
		func(body []byte) error {
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

			if err := e.taskLogger.LogTasks(task, s); err != nil {
				log.Printf("failed to log execution: %v", err)
			}

			if err := e.runBeforeExecuteHooks(task, s); err != nil {
				log.Printf("failed to run before execute hooks: %v", err)
			}

			ctx := newContext(task, *s, e.logger, e.userContext)
			err = task.Action(ctx)
			if err != nil {
				ctx.WithField("error", err.Error()).WithField("task_name", task.Name()).Error("task execution failed")
				s.LastError = err.Error()
				s.Status = StatusFailed

				if task.MaxRetries == -1 || s.Iterations < task.MaxRetries {
					if err := e.taskLogger.LogTasks(task, s); err != nil {
						log.Printf("failed to log execution: %v", err)
					}
					if err := e.runAfterExecuteHooks(task, s); err != nil {
						log.Printf("failed to run after execute hooks: %v", err)
					}
					if err := e.Publish(task, s); err != nil {
						log.Printf("failed to re-publish task: %v", err)
					}
					return err
				}
			} else {
				s.Status = StatusSuccess
			}

			if err := e.taskLogger.LogTasks(task, s); err != nil {
				log.Printf("failed to log execution: %v", err)
			}
			if err := e.runAfterExecuteHooks(task, s); err != nil {
				log.Printf("failed to run after execute hooks: %v", err)
			}

			return nil
		},
	)
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
