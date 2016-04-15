package config

type ConfigurationReader interface {
	GetCassandraHosts() ([]string, error)

	GetCassandraKeyspace() (string, error)

	GetCassandraProtocolVersion() (int, error)
}
