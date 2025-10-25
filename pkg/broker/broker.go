package broker

type Broker interface {
	// Publish publishes a message to the message broker
	Publish(body []byte, routingKey ...string) error

	// Consume consumes a message from the message broker
	Consume(queue string, concurrency int, handler func(body []byte) error) error

	// Close closes the connection to the message broker
	Close() error
}
