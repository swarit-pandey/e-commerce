package rabbitmq

// Config defines configuration for
// a client to connect to RabbitMQ messaging system
type Config struct {
	URL      string
	Username string
	Password string
	Exchange string
	Queue    string
}
