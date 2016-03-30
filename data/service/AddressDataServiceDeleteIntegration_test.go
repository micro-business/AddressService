// +build integration

package service_test

import (
	"math/rand"
	"testing"

	"github.com/gocql/gocql"
	"github.com/golang/mock/gomock"
	"github.com/microbusinesses/AddressService/data/service"
	dataServiceMocks "github.com/microbusinesses/AddressService/data/service/mocks"
	"github.com/microbusinesses/AddressService/data/shared"
	"github.com/microbusinesses/Micro-Businesses-Core/system"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Delete method behaviour", func() {
	var (
		mockCtrl                 *gomock.Controller
		addressDataService       *service.AddressDataService
		mockUUIDGeneratorService *dataServiceMocks.MockUUIDGeneratorService
		tenantId                 system.UUID
		applicationId            system.UUID
		addressId                system.UUID
		validAddress             shared.Address
		clusterConfig            *gocql.ClusterConfig
		keyspace                 string
	)

	BeforeEach(func() {
		keyspace = createRandomKeyspace()

		createAddressKeyspaceAndAllRequiredTablesForDeleteTest(keyspace)

		clusterConfig = getClusterConfig()
		clusterConfig.Keyspace = keyspace

		mockCtrl = gomock.NewController(GinkgoT())
		mockUUIDGeneratorService = dataServiceMocks.NewMockUUIDGeneratorService(mockCtrl)

		addressDataService = &service.AddressDataService{UUIDGeneratorService: mockUUIDGeneratorService, ClusterConfig: clusterConfig}

		tenantId, _ = system.RandomUUID()
		applicationId, _ = system.RandomUUID()
		addressId, _ = system.RandomUUID()
		validAddress = shared.Address{AddressKeysValues: map[string]string{"City": "Christchurch"}}
	})

	AfterEach(func() {
		mockCtrl.Finish()
		dropKeyspace(keyspace)
	})

	Context("when deleting existing address", func() {
		It("should remove the record from address table", func() {
			mockUUIDGeneratorService.
				EXPECT().
				GenerateRandomUUID().
				Return(addressId, nil)

			expectedAddressKeysValues := make(map[string]string)

			for idx := 0; idx < rand.Intn(10)+1; idx++ {
				key, _ := system.RandomUUID()
				value, _ := system.RandomUUID()

				expectedAddressKeysValues[key.String()] = value.String()
			}

			returnedAddressId, err := addressDataService.Create(
				tenantId,
				applicationId,
				shared.Address{AddressKeysValues: expectedAddressKeysValues})

			Expect(err).To(BeNil())

			err = addressDataService.Delete(
				tenantId,
				applicationId,
				returnedAddressId)

			Expect(err).To(BeNil())

			config := getClusterConfig()
			config.Keyspace = keyspace

			session, err := config.CreateSession()

			Expect(err).To(BeNil())

			defer session.Close()

			var key string
			var value string

			Expect(session.Query(
				"SELECT address_key, address_value"+
					" FROM address"+
					" WHERE"+
					" tenant_id = ?"+
					" AND application_id = ?"+
					" AND address_id = ?",
				tenantId.String(),
				applicationId.String(),
				returnedAddressId.String()).Iter().Scan(&key, &value)).To(BeFalse())
		})

		It("should remove all the index records from address_indexed_by_address_key table", func() {
			mockUUIDGeneratorService.
				EXPECT().
				GenerateRandomUUID().
				Return(addressId, nil)

			expectedAddressKeysValues := make(map[string]string)

			for idx := 0; idx < rand.Intn(10)+1; idx++ {
				key, _ := system.RandomUUID()
				value, _ := system.RandomUUID()

				expectedAddressKeysValues[key.String()] = value.String()
			}

			returnedAddressId, err := addressDataService.Create(
				tenantId,
				applicationId,
				shared.Address{AddressKeysValues: expectedAddressKeysValues})

			Expect(err).To(BeNil())

			err = addressDataService.Delete(
				tenantId,
				applicationId,
				returnedAddressId)

			Expect(err).To(BeNil())

			config := getClusterConfig()
			config.Keyspace = keyspace

			session, err := config.CreateSession()

			Expect(err).To(BeNil())

			defer session.Close()

			for key, _ := range expectedAddressKeysValues {
				var id gocql.UUID
				var addressValue string

				Expect(session.Query(
					"SELECT address_id, address_value"+
						" FROM address_indexed_by_address_key"+
						" WHERE"+
						" tenant_id = ?"+
						" AND application_id = ?"+
						" AND address_key = ?",
					tenantId.String(),
					applicationId.String(),
					key).Iter().Scan(&id, &addressValue)).To(BeFalse())

			}
		})
	})
})

func createAddressKeyspaceAndAllRequiredTablesForDeleteTest(keyspace string) {
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

func TestDeleteBehaviour(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Delete method behaviour")
}
