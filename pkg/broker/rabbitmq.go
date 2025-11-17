package broker

import (
	"github.com/wagslane/go-rabbitmq"
)

type RabbitMQBroker struct {
	connection *rabbitmq.Conn
	publisher  *rabbitmq.Publisher
	consumer   *rabbitmq.Consumer
	queueName  string
}

func (b *RabbitMQBroker) New(url, queueName string) (BrokerQueue, error) {
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
		queueName:  queueName,
	}, nil
}

func (b *RabbitMQBroker) Publish(body []byte) error {
	return b.publisher.Publish(
		body,
		[]string{b.queueName},
		rabbitmq.WithPublishOptionsContentType("application/json"),
	)
}

func (b *RabbitMQBroker) Consume(autoAck bool, concurrency int, handler func(body []byte) error) error {
	consumer, err := rabbitmq.NewConsumer(
		b.connection,
		b.queueName,
		rabbitmq.WithConsumerOptionsConcurrency(concurrency),
		rabbitmq.WithConsumerOptionsQOSPrefetch(concurrency),
		rabbitmq.WithConsumerOptionsConsumerAutoAck(autoAck),
	)
	if err != nil {
		return err
	}
	b.consumer = consumer

	return consumer.Run(func(d rabbitmq.Delivery) (action rabbitmq.Action) {
		if err := handler(d.Body); err != nil {
			return rabbitmq.NackDiscard
		}

		return rabbitmq.Ack
	})
}

func (b *RabbitMQBroker) Close() error {
	b.publisher.Close()
	b.consumer.Close()
	return b.connection.Close()
}
