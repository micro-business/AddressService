// +build integration

package service_test

import (
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gocql/gocql"
	"github.com/microbusinesses/Micro-Businesses-Core/system"
	. "github.com/onsi/gomega"
)

const databasePreparationMaxTimeout = time.Minute

func getClusterConfig() *gocql.ClusterConfig {
	cassandraIPAddress := os.Getenv("CASSANDRA_IP_ADDRESS")

	if len(cassandraIPAddress) == 0 {
		cassandraIPAddress = "127.0.0.1"
	}

	config := gocql.NewCluster(cassandraIPAddress)

	cassandraProtocolVersion := os.Getenv("CASSANDRA_PROTOCOL_VERSION")

	if len(cassandraProtocolVersion) != 0 {
		if protocolVersion, err := strconv.Atoi(cassandraProtocolVersion); err == nil {
			config.ProtoVersion = protocolVersion
		}
	}

	config.Consistency = gocql.Quorum

	return config
}

func createRandomKeyspace() string {
	keyspaceRandomValue, _ := system.RandomUUID()

	return strings.ToLower("a" + strings.Replace(keyspaceRandomValue.String(), "-", "", -1))
}

func createAddressKeyspaceAndAllRequiredTables(keyspace string) {
	config := getClusterConfig()
	config.Timeout = databasePreparationMaxTimeout
	session, err := config.CreateSession()

	Expect(err).To(BeNil())

	defer session.Close()

	Expect(session.Query(
		"CREATE KEYSPACE " +
			keyspace +
			" with replication = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };").
		Exec()).To(BeNil())

	Expect(session.Query(
		"CREATE TABLE " +
			keyspace +
			".address(tenant_id UUID, application_id UUID, address_id UUID, address_key text, address_value text," +
			" PRIMARY KEY(tenant_id, application_id, address_id, address_key));").
		Exec()).To(BeNil())

	Expect(session.Query(
		"CREATE TABLE " +
			keyspace +
			".address_indexed_by_address_key(tenant_id UUID, application_id UUID, address_id UUID, address_key text, address_value text," +
			" PRIMARY KEY(tenant_id, application_id, address_key, address_id));").
		Exec()).To(BeNil())
}

func dropKeyspace(keyspace string) {
	config := getClusterConfig()
	config.Timeout = databasePreparationMaxTimeout
	session, err := config.CreateSession()

	Expect(err).To(BeNil())

	defer session.Close()

	err = session.Query("DROP KEYSPACE " + keyspace + " ;").Exec()

	Expect(err).To(BeNil())
}

func mapSystemUUIDToGocqlUUID(uuid system.UUID) gocql.UUID {
	mappedUUID, _ := gocql.UUIDFromBytes(uuid.Bytes())

	return mappedUUID
}

func mapGocqlUUIDToSystemUUID(uuid gocql.UUID) system.UUID {
	mappedUUID, _ := system.UUIDFromBytes(uuid.Bytes())

	return mappedUUID
}

func createRandomAddressKeyValues() map[string]string {
	keyValues := make(map[string]string)

	for idx := 0; idx < rand.Intn(10)+1; idx++ {
		key, _ := system.RandomUUID()
		value, _ := system.RandomUUID()

		keyValues[key.String()] = value.String()
	}

	return keyValues
}
