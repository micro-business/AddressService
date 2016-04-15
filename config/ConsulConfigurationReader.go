package config

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/consul/api"
)

type ConsulConfigurationReader struct {
	ConsulAddress                      string
	ConsulScheme                       string
	ListeningPortToOverride            int
	CassandraHostsToOverride           []string
	CassandraKeyspaceToOverride        string
	CassandraProtocolVersionToOverride int
}

func (consul ConsulConfigurationReader) GetListeningPort() (int, error) {
	if consul.ListeningPortToOverride != 0 {
		return consul.ListeningPortToOverride, nil

	} else {
		return getInt(consul, "services/address-service/endpoint/listening-port")
	}
}

func (consul ConsulConfigurationReader) GetCassandraHosts() ([]string, error) {
	if len(consul.CassandraHostsToOverride) != 0 {
		return consul.CassandraHostsToOverride, nil
	}

	keyPair, err := getKeyPair(consul, "services/address-service/data/cassandra/hosts")

	if err != nil {
		return nil, err
	}

	return strings.Split(string(keyPair.Value), ","), nil
}

func (consul ConsulConfigurationReader) GetCassandraKeyspace() (string, error) {
	if len(consul.CassandraKeyspaceToOverride) != 0 {
		return consul.CassandraKeyspaceToOverride, nil
	} else {
		return getString(consul, "services/address-service/data/cassandra/keyspace")
	}
}

func (consul ConsulConfigurationReader) GetCassandraProtocolVersion() (int, error) {
	if consul.CassandraProtocolVersionToOverride != 0 {
		return consul.CassandraProtocolVersionToOverride, nil
	} else {
		return getInt(consul, "services/address-service/data/cassandra/protocol-version")
	}
}

func getKV(consul ConsulConfigurationReader) (*api.KV, error) {
	config := api.DefaultConfig()

	if len(consul.ConsulAddress) != 0 && len(consul.ConsulScheme) != 0 {
		config.Address = consul.ConsulAddress
		config.Scheme = consul.ConsulScheme
	}

	if client, err := api.NewClient(config); err != nil {
		return nil, err
	} else {
		return client.KV(), nil
	}
}

func getKeyPair(consul ConsulConfigurationReader, configKeyPath string) (*api.KVPair, error) {
	kv, err := getKV(consul)

	if err != nil {
		return nil, err
	}

	if keyPair, _, err := kv.Get(configKeyPath, nil); err != nil {
		return nil, err
	} else {
		return keyPair, nil
	}
}

func getInt(consul ConsulConfigurationReader, keyPath string) (int, error) {
	keyPair, err := getKeyPair(consul, keyPath)

	if err != nil {
		return 0, err
	}

	if keyPair == nil {
		return 0, errors.New(fmt.Sprintf("Consul key %s does not exist.", keyPath))

	}

	valueInString := string(keyPair.Value)

	if len(valueInString) == 0 {
		return 0, errors.New(fmt.Sprintf("Consul key %s is empty.", keyPath))

	}

	if value, err := strconv.Atoi(valueInString); err != nil {
		return 0, err
	} else {
		if value == 0 {
			return 0, errors.New(fmt.Sprintf("Consul key %s is zero.", keyPath))
		}

		return value, nil
	}
}

func getString(consul ConsulConfigurationReader, keyPath string) (string, error) {
	keyPair, err := getKeyPair(consul, keyPath)

	if err != nil {
		return "", err
	}

	if keyPair == nil {
		return "", errors.New(fmt.Sprintf("Consul key %s does not exist.", keyPath))

	}

	return string(keyPair.Value), nil

}
