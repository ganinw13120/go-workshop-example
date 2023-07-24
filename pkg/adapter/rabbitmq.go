package adapter

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

type IRabbitMQAdapter interface {
	Publish(ctx context.Context, exchange, routing string, payload []byte) error
	Consume(queue string, cb func([]byte)) error
}

type rabbitmqAdapter struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func NewRabbitMQAdapter() *rabbitmqAdapter {
	conn, err := amqp.Dial("amqp://root:root@localhost:5672/")
	if err != nil {
		panic(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}

	return &rabbitmqAdapter{
		conn,
		ch,
	}
}

func (r *rabbitmqAdapter) CleanUp() {
	r.ch.Close()
	r.conn.Close()
}

func (r rabbitmqAdapter) Consume(queue string, cb func([]byte)) error {
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

func (r rabbitmqAdapter) Publish(ctx context.Context, exchange, routing string, payload []byte) error {
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
