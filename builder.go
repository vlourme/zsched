package zsched

import (
	"errors"
	"sync"

	"github.com/robfig/cron/v3"
	"github.com/vlourme/zsched/pkg/broker"
	"github.com/vlourme/zsched/pkg/logger"
	"github.com/vlourme/zsched/pkg/storage"
)

// builder is the builder for the engine
type builder[T any] struct {
	engine *Engine[T]
	err    error
}

// NewBuilder creates a new builder for the engine
func NewBuilder[T any](userContext T) *builder[T] {
	return &builder[T]{
		engine: &Engine[T]{
			userContext: userContext,
			tasks:       make(map[string]*Task[T]),
			wg:          &sync.WaitGroup{},
			cron:        cron.New(cron.WithSeconds()),
			hooks:       make([]Hook, 0),
		},
	}
}

// WithBroker sets the broker for the engine
func (b *builder[T]) WithBroker(broker broker.BrokerQueue) *builder[T] {
	b.engine.broker = broker
	return b
}

// WithRabbitMQBroker sets the RabbitMQ broker for the engine
func (b *builder[T]) WithRabbitMQBroker(url string) *builder[T] {
	b.engine.brokerUrl = url
	return b.WithBroker(&broker.RabbitMQBroker{})
}

// WithStorage sets the storage for the engine
func (b *builder[T]) WithStorage(storage storage.Storage) *builder[T] {
	b.engine.storage = storage
	return b
}

// WithTimescaleDBStorage sets the TimescaleDB storage for the engine
func (b *builder[T]) WithTimescaleDBStorage(url string) *builder[T] {
	storage, err := storage.NewTimescaleDBStorage(url)
	if err != nil {
		b.err = err
		return b
	}
	return b.WithStorage(storage)
}

// WithHooks sets the hooks for the engine
func (b *builder[T]) WithHooks(hooks ...Hook) *builder[T] {
	b.engine.hooks = hooks
	return b
}

// WithLogger sets the logger for the engine
func (b *builder[T]) WithLogger(logger logger.Logger) *builder[T] {
	b.engine.logger = logger
	return b
}

func (b *builder[T]) WithAPI(address string) *builder[T] {
	b.engine.apiAddress = address
	return b
}

// Build builds the engine
func (b *builder[T]) Build() (*Engine[T], error) {
	if b.err != nil {
		return nil, errors.New("failed to build engine: " + b.err.Error())
	}

	if b.engine.broker == nil {
		return nil, errors.New("broker is required")
	}

	if b.engine.storage == nil {
		return nil, errors.New("storage is required")
	}

	if b.engine.logger == nil {
		b.engine.logger = logger.NewLogger(b.engine.storage)
	}

	for _, hook := range b.engine.hooks {
		if err := hook.Initialize(b.engine.storage); err != nil {
			return nil, errors.New("failed to initialize hook: " + err.Error())
		}
	}

	return b.engine, nil
}
