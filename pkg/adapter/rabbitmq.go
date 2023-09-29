package adapter

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

type MessageBroker interface {
	Publish(ctx context.Context, exchange, routing string, payload []byte) error
	Consume(queue string, cb func([]byte)) error
}

type QueueConfig struct {
	QueueName    string
	ExchangeName string
	ExchangeType string
}

type rabbitmq struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func NewRabbitMQ(url string, config QueueConfig) (*rabbitmq, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	q := &rabbitmq{
		conn,
		ch,
	}

	q.createAndBindQueue(config)

	return q, nil
}

func (q rabbitmq) createAndBindQueue(queueConfig QueueConfig) error {
	err := q.declareExchange(queueConfig.ExchangeName, queueConfig.ExchangeType)
	if err != nil {
		return err
	}
	err = q.declareQueue(queueConfig.QueueName)
	if err != nil {
		return err
	}
	routingKey := ""
	if queueConfig.ExchangeType == "direct" {
		routingKey = queueConfig.QueueName
	}
	err = q.bindQueue(queueConfig.QueueName, routingKey, queueConfig.ExchangeName)
	if err != nil {
		return err
	}
	return nil
}

func (q rabbitmq) declareExchange(exchangeName string, exchangeType string) error {
	err := q.ch.ExchangeDeclare(
		exchangeName, // name
		exchangeType, // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	return err
}

func (q rabbitmq) declareQueue(queueName string) error {
	_, err := q.ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	return err
}

func (q rabbitmq) bindQueue(queueName string, routingKey string, exchangeName string) error {
	err := q.ch.QueueBind(
		queueName,    // queue name
		routingKey,   // routing key
		exchangeName, // exchange
		false,
		nil)
	return err
}

func (r *rabbitmq) CleanUp() {
	r.ch.Close()
	r.conn.Close()
}

func (r rabbitmq) Consume(queue string, cb func([]byte)) error {
	msgs, err := r.ch.Consume(
		queue, // queue
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return err
	}

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			cb(d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever

	return nil
}

func (r rabbitmq) Publish(ctx context.Context, exchange, routing string, payload []byte) error {
	err := r.ch.PublishWithContext(ctx,
		exchange,
		routing,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        payload,
		})

	return err
}
