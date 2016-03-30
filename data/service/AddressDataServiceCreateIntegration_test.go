// +build integration

package service_test

import (
	"errors"
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

var _ = Describe("Create method behaviour", func() {
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

		createAddressKeyspaceAndAllRequiredTables(keyspace)

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

	Context("when creating new address", func() {
		It("should insert the record into address table", func() {
			mockUUIDGeneratorService.
				EXPECT().
				GenerateRandomUUID().
				Return(addressId, nil)

			expectedAddressKeysValues := createRandomAddressKeyValues()

			returnedAddressId, err := addressDataService.Create(
				tenantId,
				applicationId,
				shared.Address{AddressKeysValues: expectedAddressKeysValues})

			Expect(addressId).To(Equal(returnedAddressId))
			Expect(err).To(BeNil())

			config := getClusterConfig()
			config.Keyspace = keyspace

			session, err := config.CreateSession()

			Expect(err).To(BeNil())

			defer session.Close()

			iter := session.Query(
				"SELECT address_key, address_value"+
					" FROM address"+
					" WHERE"+
					" tenant_id = ?"+
					" AND application_id = ?"+
					" AND address_id = ?",
				tenantId.String(),
				applicationId.String(),
				addressId.String()).Iter()

			var key string
			var value string

			addressKeysValues := make(map[string]string)

			for iter.Scan(&key, &value) {
				addressKeysValues[key] = value
			}

			err = iter.Close()

			Expect(err).To(BeNil())

			Expect(expectedAddressKeysValues).To(Equal(addressKeysValues))
		})

		It("should insert the record into address_indexed_by_address_key table", func() {
			mockUUIDGeneratorService.
				EXPECT().
				GenerateRandomUUID().
				Return(addressId, nil)

			expectedAddressKeysValues := createRandomAddressKeyValues()

			addressDataService.Create(tenantId, applicationId, shared.Address{AddressKeysValues: expectedAddressKeysValues})

			config := getClusterConfig()
			config.Keyspace = keyspace

			session, err := config.CreateSession()

			Expect(err).To(BeNil())

			defer session.Close()

			for key, value := range expectedAddressKeysValues {
				var id gocql.UUID
				var addressValue string

				err = session.Query(
					"SELECT address_id, address_value"+
						" FROM address_indexed_by_address_key"+
						" WHERE"+
						" tenant_id = ?"+
						" AND application_id = ?"+
						" AND address_key = ?",
					tenantId.String(),
					applicationId.String(),
					key).Scan(&id, &addressValue)

				Expect(err).To(BeNil())

				Expect(addressId).To(Equal(mapGocqlUUIDToSystemUUID(id)))
				Expect(value).To(Equal(addressValue))
			}
		})
	})
})

func TestCreateBehaviour(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Create method behaviour")
}
