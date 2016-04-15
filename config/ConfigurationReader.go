package config

type ConfigurationReader interface {
	GetListeningPort() (int, error)

	GetCassandraHosts() ([]string, error)

	GetCassandraKeyspace() (string, error)

	GetCassandraProtocolVersion() (int, error)
}
