package config

// ConfigurationReader defines the interface that provides access to all configurations parameters required by the service.
type ConfigurationReader interface {
	// GetListeningPort returns the port the application should start listening on.
	GetListeningPort() (int, error)

	// GetCassandraHosts returns the list of Cassandra host addresses.
	GetCassandraHosts() ([]string, error)

	// GetCassandraKeyspace returns the name of Cassandra key space that the service data is stored under.
	GetCassandraKeyspace() (string, error)

	// GetCassandraProtocolVersion returns the cassandra procotol version.
	GetCassandraProtocolVersion() (int, error)
}
