package database

// Config struct holds necessary fields for
// a new client to connect to the DB
type Config struct {
	Driver             string
	ConnectionString   string
	Port               string
	Username           string
	Password           string
	MaxOpenConnections int
	MaxIdleConnections int
	SSLEnabled         bool
	SchemaName         string
	DBName             string
}
