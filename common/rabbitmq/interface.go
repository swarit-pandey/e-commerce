package rabbitmq

import "context"

// RabbitMQ interface defines methods to interact with an instance
type RabbitMQ interface {
	// Publish will flush the messages
	Publish(ctx context.Context, message []byte) error

	// Consume will consume the messages
	Consume(ctx context.Context, handler func(message []byte) error) error

	// Close will close an active connection
	Close() error
}
