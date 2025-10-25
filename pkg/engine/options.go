package engine

import (
	"log"

	"github.com/vlourme/scheduler/pkg/broker"
	"github.com/vlourme/scheduler/pkg/hooks"
	"github.com/vlourme/scheduler/pkg/logger"
	"github.com/vlourme/scheduler/pkg/storage"
)

// WithRabbitMQBroker creates a new RabbitMQ broker and sets it as the engine's broker
func WithRabbitMQBroker(url string) func(*Engine) {
	return func(e *Engine) {
		broker, err := broker.NewRabbitMQBroker(url)
		if err != nil {
			log.Fatalf("failed to create RabbitMQ broker: %v", err)
		}

		e.Broker = broker
	}
}

// WithLogger sets the logger for the engine
func WithLogger(logger logger.Logger) func(*Engine) {
	return func(e *Engine) {
		e.logger = logger
	}
}

// WithQuestDBStorage sets the QuestDB storage for the engine
func WithQuestDBStorage(dsn string) func(*Engine) {
	return func(e *Engine) {
		storage, err := storage.NewQuestDBStorage(dsn)
		if err != nil {
			log.Fatalf("failed to create QuestDB storage: %v", err)
		}
		e.Storage = storage
	}
}

// WithAPI sets the API for the engine
func WithAPI(address string) func(*Engine) {
	return func(e *Engine) {
		e.apiAddress = address
	}
}

func WithHooks(hooks ...hooks.Hook) func(*Engine) {
	return func(e *Engine) {
		e.hooks = hooks
	}
}

// WithUserContext sets the user context for the engine
func WithUserContext[T any](userContext T) func(*Engine) {
	return func(e *Engine) {
		e.userContext = userContext
	}
}
