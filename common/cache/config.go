package cache

// Config holds aerospike's configuration
type Config struct {
	// Aerospike connection URL
	URL string

	// Aerospike logging level
	LogLevel string

	// Port on which Aerospike is running
	Port int

	// Namespace
	Namespace string

	// Set
	Set string

	// Bin
	Bin string
}
