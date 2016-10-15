package config

import (
	"fmt"
	"strings"

	"github.com/microbusinesses/Micro-Businesses-Core/common/config"
)

// ConsulConfigurationReader implements ConfigurationReader using Consul to provide access to all configurations parameters required by the service.
type ConsulConfigurationReader struct {
	ConsulAddress                      string
	ConsulScheme                       string
	ListeningPortToOverride            int
	CassandraHostsToOverride           []string
	CassandraKeyspaceToOverride        string
	CassandraProtocolVersionToOverride int
}

const serviceListeningPortKey = "services/address-service/endpoint/listening-port"
const cassandraHostsKey = "services/address-service/data/cassandra/hosts"
const cassandraKeyspaceKey = "services/address-service/data/cassandra/keyspace"
const cassandraProtocolVersionKey = "services/address-service/data/cassandra/protocol-version"

func (consul ConsulConfigurationReader) GetListeningPort() (int, error) {
	if consul.ListeningPortToOverride != 0 {
		return consul.ListeningPortToOverride, nil

	}
	consulHelper := config.ConsulHelper{ConsulAddress: consul.ConsulAddress, ConsulScheme: consul.ConsulScheme}

	return consulHelper.GetInt(serviceListeningPortKey)
}

func (consul ConsulConfigurationReader) GetCassandraHosts() ([]string, error) {
	if len(consul.CassandraHostsToOverride) != 0 {
		return consul.CassandraHostsToOverride, nil
	}

	consulHelper := config.ConsulHelper{ConsulAddress: consul.ConsulAddress, ConsulScheme: consul.ConsulScheme}
	keyPair, err := consulHelper.GetKeyPair(cassandraHostsKey)

	if err != nil {
		return nil, err
	}

	if keyPair == nil {
		return nil, fmt.Errorf("Consul key %s does not exist.", cassandraHostsKey)

	}

	valueInString := string(keyPair.Value)

	if len(valueInString) == 0 {
		return nil, fmt.Errorf("Consul key %s is empty.", cassandraHostsKey)

	}

	return strings.Split(string(keyPair.Value), ","), nil
}

func (consul ConsulConfigurationReader) GetCassandraKeyspace() (string, error) {
	if len(consul.CassandraKeyspaceToOverride) != 0 {
		return consul.CassandraKeyspaceToOverride, nil
	}

	consulHelper := config.ConsulHelper{ConsulAddress: consul.ConsulAddress, ConsulScheme: consul.ConsulScheme}

	return consulHelper.GetString(cassandraKeyspaceKey)
}

func (consul ConsulConfigurationReader) GetCassandraProtocolVersion() (int, error) {
	if consul.CassandraProtocolVersionToOverride != 0 {
		return consul.CassandraProtocolVersionToOverride, nil
	}

	consulHelper := config.ConsulHelper{ConsulAddress: consul.ConsulAddress, ConsulScheme: consul.ConsulScheme}

	return consulHelper.GetInt(cassandraProtocolVersionKey)
}
