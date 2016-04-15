package config

import (
	"strconv"
	"strings"

	"github.com/hashicorp/consul/api"
)

type ConsulConfigurationReader struct {
	ConsulAddress string
	ConsulScheme  string
}

func (consul ConsulConfigurationReader) GetCassandraHosts() ([]string, error) {
	kv, err := getKV(consul)

	if err != nil {
		return nil, err
	}

	keyPair, _, err := kv.Get("services/address-service/data/cassandra/hosts", nil)

	if err != nil {
		return nil, err
	}

	return strings.Split(string(keyPair.Value), ","), nil
}

func (consul ConsulConfigurationReader) GetCassandraKeyspace() (string, error) {
	kv, err := getKV(consul)

	if err != nil {
		return "", err
	}

	keyPair, _, err := kv.Get("services/address-service/data/cassandra/keyspace", nil)

	if err != nil {
		return "", err
	}

	keyspace := ""

	if keyPair == nil {
		keyspace = "address"
	} else {
		keyspace = string(keyPair.Value)

		if len(keyspace) == 0 {
			keyspace = "address"
		}
	}

	return keyspace, nil
}

func (consul ConsulConfigurationReader) GetCassandraProtocolVersion() (int, error) {
	kv, err := getKV(consul)

	if err != nil {
		return 0, err
	}

	keyPair, _, err := kv.Get("services/address-service/data/cassandra/protocol-version", nil)

	if err != nil {
		return 0, err
	}

	protocolVersion := 4

	if keyPair != nil {
		protocolVersionInString := string(keyPair.Value)

		if len(protocolVersionInString) != 0 {
			protocolVersion, err = strconv.Atoi(protocolVersionInString)

			if err != nil {
				return 0, err
			}

		}
	}

	return protocolVersion, nil
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
