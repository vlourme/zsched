package broker

import (
	"github.com/wagslane/go-rabbitmq"
)

type RabbitMQBroker struct {
	connection *rabbitmq.Conn
	publisher  *rabbitmq.Publisher
	consumers  []*rabbitmq.Consumer
}

func NewRabbitMQBroker(url string) (*RabbitMQBroker, error) {
	conn, err := rabbitmq.NewConn(url)
	if err != nil {
		return nil, err
	}

	publisher, err := rabbitmq.NewPublisher(conn)
	if err != nil {
		return nil, err
	}

	return &RabbitMQBroker{
		connection: conn,
		publisher:  publisher,
	}, nil
}

func (b *RabbitMQBroker) Publish(body []byte, routingKey ...string) error {
	return b.publisher.Publish(
		body,
		routingKey,
		rabbitmq.WithPublishOptionsContentType("application/json"),
	)
}

func (b *RabbitMQBroker) Consume(queue string, concurrency int, handler func(body []byte) error) error {
	consumer, err := rabbitmq.NewConsumer(
		b.connection,
		queue,
		rabbitmq.WithConsumerOptionsConcurrency(concurrency),
		rabbitmq.WithConsumerOptionsQOSPrefetch(concurrency),
	)
	if err != nil {
		return err
	}

	b.consumers = append(b.consumers, consumer)

	return consumer.Run(func(d rabbitmq.Delivery) (action rabbitmq.Action) {
		if err := handler(d.Body); err != nil {
			return rabbitmq.NackDiscard
		}

		return rabbitmq.Ack
	})
}

func (b *RabbitMQBroker) Close() error {
	b.publisher.Close()
	for _, consumer := range b.consumers {
		consumer.Close()
	}
	return b.connection.Close()
}
