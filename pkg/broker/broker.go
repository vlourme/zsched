package broker

// BrokerQueue represents a queue in the message broker
type BrokerQueue interface {
	// New creates a new broker queue
	New(url, queueName string) (BrokerQueue, error)

	// Publish publishes a message to the message broker
	Publish(body []byte) error

	// Consume consumes a message from the message broker
	Consume(autoAck bool, concurrency int, handler func(body []byte) error) error

	// Close closes the connection to the message broker
	Close() error
}
