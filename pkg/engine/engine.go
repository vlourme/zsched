package engine

import (
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/robfig/cron/v3"
	"github.com/vlourme/zsched/pkg/api"
	"github.com/vlourme/zsched/pkg/broker"
	"github.com/vlourme/zsched/pkg/executor"
	"github.com/vlourme/zsched/pkg/hooks"
	"github.com/vlourme/zsched/pkg/logger"
	"github.com/vlourme/zsched/pkg/storage"
	"github.com/vlourme/zsched/pkg/task"
)

// Scheduler is the scheduler for the tasks
type Engine struct {
	// Broker is the Broker for the tasks
	Broker broker.Broker

	// Storage is the Storage for the engine
	Storage storage.Storage

	// hooks are the hooks for the engine
	hooks []hooks.Hook

	// logger is the logger for the engine
	logger logger.Logger

	// apiAddress is the address for the API
	apiAddress string

	// tasks are the registered tasks
	tasks map[string]*task.Task

	// executor is the executor for the tasks
	executor *executor.Executor

	// cron is the cron for the engine
	cron *cron.Cron

	// wg is the wait group for the engine
	wg *sync.WaitGroup

	// userContext is the user context for the engine
	userContext any
}

// NewScheduler creates a new scheduler
func NewEngine(options ...func(*Engine)) (*Engine, error) {
	engine := &Engine{
		tasks: make(map[string]*task.Task),
		cron:  cron.New(cron.WithSeconds()),
		wg:    &sync.WaitGroup{},
	}

	for _, option := range options {
		option(engine)
	}

	if engine.Broker == nil {
		return nil, errors.New("broker is required")
	}

	if engine.Storage == nil {
		return nil, errors.New("storage is required")
	}

	if engine.logger == nil {
		engine.logger = logger.NewLogger(engine.Storage)
	}

	engine.hooks = append(
		engine.hooks,
		hooks.NewTaskLoggerHook(engine.Storage),
	)

	return engine, nil
}

// RegisterHooks registers new hooks to the engine
func (e *Engine) RegisterHooks(hooks ...hooks.Hook) {
	e.hooks = append(e.hooks, hooks...)
}

// Register registers new tasks to the scheduler
func (e *Engine) Register(task ...*task.Task) {
	for _, t := range task {
		if _, ok := e.tasks[t.Name()]; ok {
			log.Printf("task %s already registered", t.Name())
		}

		e.tasks[t.Name()] = t
	}
}

// Start starts the engine, this function is blocking until the engine is stopped
func (e *Engine) Start() {
	e.logger.Info("Starting engine...")
	e.executor = executor.NewExecutor(e.Broker, e.logger, e.hooks, e.userContext)
	for _, task := range e.tasks {
		task.Executor = e.executor

		for _, schedule := range task.Schedules {
			_, err := e.cron.AddFunc(schedule.Schedule, func() {
				task.Execute(schedule.Parameters)
			})
			if err != nil {
				e.logger.
					WithField("name", task.Name()).
					WithError(err).
					Error("failed to add schedule for task")
			}
		}

		e.wg.Add(1)
		go func() {
			defer e.wg.Done()
			e.executor.Consume(task)
			fmt.Println("Task consumed")
		}()
	}
	e.cron.Start()

	if e.apiAddress != "" {
		e.logger.WithField("address", e.apiAddress).Info("Starting API...")
		go api.NewAPI(e.tasks, e.Storage).Run(e.apiAddress)
	}

	e.wg.Wait()
}

// Close closes the engine
func (e *Engine) Close() error {
	err := e.Broker.Close()
	if err != nil {
		return err
	}

	err = e.Storage.Close()
	if err != nil {
		return err
	}

	e.cron.Stop()

	return nil
}
