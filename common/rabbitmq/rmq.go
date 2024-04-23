package rabbitmq

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

// rabbitMQClient struct holds necessary components to interact with RMQ client
type rabbitMQClient struct {
	conn     *amqp.Connection
	channel  *amqp.Channel
	exchange string
	queue    string
}

// NewRabbitMQClient client will return a new RMQ client based on the given configuration
func NewRabbitMQClient(config *Config) (RabbitMQ, error) {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s/", config.Username, config.Password, config.URL))
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	err = ch.ExchangeDeclare(
		config.Exchange, // exchange name
		"direct",        // exchange type
		true,            // durable
		false,           // auto-deleted
		false,           // internal
		false,           // no-wait
		nil,             // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, err
	}

	_, err = ch.QueueDeclare(
		config.Queue, // queue name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, err
	}

	err = ch.QueueBind(
		config.Queue,    // queue name
		config.Queue,    // routing key
		config.Exchange, // exchange name
		false,           // no-wait
		nil,             // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, err
	}

	return &rabbitMQClient{
		conn:     conn,
		channel:  ch,
		exchange: config.Exchange,
		queue:    config.Queue,
	}, nil
}

// Publish implements `Publish` from RMQ's interface
func (c *rabbitMQClient) Publish(ctx context.Context, message []byte) error {
	err := c.channel.PublishWithContext(ctx,
		c.exchange, // exchange name
		c.queue,    // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        message,
		},
	)
	if err != nil {
		return err
	}
	return nil
}

// Consume implements `Consume` from RMQ's interface
func (c *rabbitMQClient) Consume(ctx context.Context, handler func(message []byte) error) error {
	msgs, err := c.channel.Consume(
		c.queue, // queue name
		"",      // consumer
		false,   // auto-ack
		false,   // exclusive
		false,   // no-local
		false,   // no-wait
		nil,     // args
	)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg := <-msgs:
			err := handler(msg.Body)
			if err != nil {
				_ = msg.Nack(false, true)
			} else {
				_ = msg.Ack(false)
			}
		}
	}
}

// Close implements `Close` from RMQ's interface
func (c *rabbitMQClient) Close() error {
	err := c.channel.Close()
	if err != nil {
		return err
	}
	err = c.conn.Close()
	if err != nil {
		return err
	}
	return nil
}
