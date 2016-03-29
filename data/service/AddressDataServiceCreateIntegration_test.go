// +build integration

package service_test

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gocql/gocql"
	"github.com/golang/mock/gomock"
	"github.com/microbusinesses/AddressService/data/service"
	dataServiceMocks "github.com/microbusinesses/AddressService/data/service/mocks"
	"github.com/microbusinesses/AddressService/data/shared"
	"github.com/microbusinesses/Micro-Businesses-Core/system"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const databasePreparationMaxTimeout = time.Minute

var _ = Describe("Create method behaviour", func() {
	var (
		mockCtrl                 *gomock.Controller
		addressDataService       *service.AddressDataService
		mockUUIDGeneratorService *dataServiceMocks.MockUUIDGeneratorService
		tenantId                 system.UUID
		applicationId            system.UUID
		validAddress             shared.Address
		clusterConfig            *gocql.ClusterConfig
		keyspace                 string
	)

	BeforeSuite(func() {
		config := getClusterConfig()
		config.Timeout = databasePreparationMaxTimeout

		keyspaceRandomValue, _ := system.RandomUUID()
		keyspace = strings.ToLower("a" + strings.Replace(keyspaceRandomValue.String(), "-", "", -1))

		session, err := config.CreateSession()

		Expect(err).To(BeNil())

		defer session.Close()

		err = session.Query(
			"CREATE KEYSPACE " +
				keyspace +
				" with replication = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };").
			Exec()

		Expect(err).To(BeNil())

		err = session.Query(
			"CREATE TABLE " +
				keyspace +
				".address(tenant_id UUID, application_id UUID, address_id UUID, address_key text, address_value text, PRIMARY KEY(tenant_id, application_id, address_id, address_key));").
			Exec()

		Expect(err).To(BeNil())

		err = session.Query(
			"CREATE TABLE " +
				keyspace +
				".address_indexed_by_address_key(tenant_id UUID, application_id UUID, address_id UUID, address_key text, address_value text, PRIMARY KEY(tenant_id, application_id, address_key, address_id));").
			Exec()

		Expect(err).To(BeNil())

		clusterConfig = getClusterConfig()
		clusterConfig.Keyspace = keyspace
	})

	AfterSuite(func() {
		config := getClusterConfig()
		config.Timeout = databasePreparationMaxTimeout
		session, err := config.CreateSession()

		Expect(err).To(BeNil())

		defer session.Close()

		err = session.Query("DROP KEYSPACE " + keyspace + " ;").Exec()

		Expect(err).To(BeNil())
	})

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockUUIDGeneratorService = dataServiceMocks.NewMockUUIDGeneratorService(mockCtrl)

		addressDataService = &service.AddressDataService{UUIDGeneratorService: mockUUIDGeneratorService, ClusterConfig: clusterConfig}

		tenantId, _ = system.RandomUUID()
		applicationId, _ = system.RandomUUID()
		validAddress = shared.Address{AddressKeysValues: map[string]string{"City": "Christchurch"}}
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("when UUID generator service succeeds to create the new UUID", func() {
		It("should return the new UUID as address uniuqe identifier and no error", func() {
			expectedAddressId, _ := system.RandomUUID()
			mockUUIDGeneratorService.
				EXPECT().
				GenerateRandomUUID().
				Return(expectedAddressId, nil)

			newAddressId, err := addressDataService.Create(tenantId, applicationId, validAddress)

			Expect(newAddressId).To(Equal(newAddressId))
			Expect(err).To(BeNil())
		})
	})

	Context("when UUID generator service fails to create the new UUID", func() {
		It("should return address unique identifier as empty UUID and the returned error by address data service", func() {
			expectedErrorId, _ := system.RandomUUID()
			expectedError := errors.New(expectedErrorId.String())
			mockUUIDGeneratorService.
				EXPECT().
				GenerateRandomUUID().
				Return(system.EmptyUUID, expectedError)

			newAddressId, err := addressDataService.Create(tenantId, applicationId, validAddress)

			Expect(newAddressId).To(Equal(system.EmptyUUID))
			Expect(err).To(Equal(expectedError))
		})
	})
})

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

func TestCreateIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Create method behaviour")
}
