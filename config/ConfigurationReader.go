package config

type ConfigurationReader interface {
	// GetListeningPort returns the port the application should start listening on.
	// Returns either the listening port or error if something goes wrong.
	GetListeningPort() (int, error)

	GetCassandraHosts() ([]string, error)

	GetCassandraKeyspace() (string, error)

	// GetCassandraProtocolVersion returns the cassandra procotol version.
	// Returns either the cassandra protcol version or error if something goes wrong.
	GetCassandraProtocolVersion() (int, error)
}
