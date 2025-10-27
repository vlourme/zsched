package zsched

import (
	"errors"
	"log"
	"net/http"
	"sync"

	"github.com/robfig/cron/v3"
	"github.com/vlourme/zsched/pkg/broker"
	"github.com/vlourme/zsched/pkg/logger"
	"github.com/vlourme/zsched/pkg/storage"
)

// Engine is the main engine for Zsched scheduler
type Engine[T any] struct {
	broker      broker.Broker
	storage     storage.Storage
	hooks       []Hook
	logger      logger.Logger
	tasks       map[string]*Task[T]
	wg          *sync.WaitGroup
	cron        *cron.Cron
	userContext T
	executor    *executor[T]
	apiAddress  string
}

// Register registers new tasks to the scheduler
func (e *Engine[T]) Register(task ...*Task[T]) {
	for _, t := range task {
		if _, ok := e.tasks[t.Name()]; ok {
			log.Printf("task %s already registered", t.Name())
		}

		e.tasks[t.Name()] = t
	}
}

// Start starts the engine, this function is blocking until the engine is stopped
func (e *Engine[T]) Start() error {
	e.logger.Info("Starting engine...")

	e.executor = &executor[T]{
		broker:      e.broker,
		logger:      e.logger,
		hooks:       e.hooks,
		userContext: e.userContext,
	}

	for _, task := range e.tasks {
		task.executor = e.executor

		for _, schedule := range task.Schedules {
			_, err := e.cron.AddFunc(schedule.Schedule, func() {
				if err := task.Execute(schedule.Parameters); err != nil {
					e.logger.WithError(err).WithField("task_name", task.Name()).Error("failed to execute task")
				}
			})
			if err != nil {
				return errors.Join(errors.New("failed to add schedule for task"), err)
			}
		}

		e.wg.Go(func() {
			e.executor.Consume(task)
		})
	}

	if e.apiAddress != "" {
		router := newRouter(e.tasks, e.storage)
		e.logger.WithField("listen_addr", e.apiAddress).Info("Starting API server...")
		go http.ListenAndServe(e.apiAddress, router)
	}

	e.cron.Start()
	e.wg.Wait()

	return nil
}

// Close closes the engine
func (e *Engine[T]) Close() error {
	err := e.broker.Close()
	if err != nil {
		return err
	}

	err = e.storage.Close()
	if err != nil {
		return err
	}

	e.cron.Stop()

	return nil
}
