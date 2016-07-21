// +build integration

package service_test

import (
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

var _ = Describe("Update method behaviour", func() {
	var (
		mockCtrl                 *gomock.Controller
		addressDataService       *service.AddressDataService
		mockUUIDGeneratorService *dataServiceMocks.MockUUIDGeneratorService
		tenantId                 system.UUID
		applicationId            system.UUID
		addressId                system.UUID
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
	})

	AfterEach(func() {
		mockCtrl.Finish()
		dropKeyspace(keyspace)
	})

	Context("when updating existing address", func() {
		It("should remove all old records from address table", func() {
			mockUUIDGeneratorService.
				EXPECT().
				GenerateRandomUUID().
				Return(addressId, nil)

			addressKeysValuesToAdd := createRandomAddressDetails()

			returnedAddressId, err := addressDataService.Create(
				tenantId,
				applicationId,
				shared.Address{AddressDetails: addressKeysValuesToAdd})

			Expect(err).To(BeNil())

			err = addressDataService.Update(
				tenantId,
				applicationId,
				returnedAddressId,
				shared.Address{AddressDetails: createRandomAddressDetails()})

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

			defer iter.Close()

			var key string
			var value string

			for iter.Scan(&key, &value) {
				_, ok := addressKeysValuesToAdd[key]

				Expect(ok).To(BeFalse())
			}
		})

		It("should remove all old records from address_indexed_by_address_key table", func() {
			mockUUIDGeneratorService.
				EXPECT().
				GenerateRandomUUID().
				Return(addressId, nil)

			addressKeysValuesToAdd := createRandomAddressDetails()

			returnedAddressId, err := addressDataService.Create(
				tenantId,
				applicationId,
				shared.Address{AddressDetails: addressKeysValuesToAdd})

			Expect(err).To(BeNil())

			err = addressDataService.Update(
				tenantId,
				applicationId,
				returnedAddressId,
				shared.Address{AddressDetails: createRandomAddressDetails()})

			config := getClusterConfig()
			config.Keyspace = keyspace

			session, err := config.CreateSession()

			Expect(err).To(BeNil())

			defer session.Close()

			for key, _ := range addressKeysValuesToAdd {
				var addressValue string

				iter := session.Query(
					"SELECT address_id, address_value"+
						" FROM address_indexed_by_address_key"+
						" WHERE"+
						" tenant_id = ?"+
						" AND application_id = ?"+
						" AND address_key = ?"+
						" AND address_id = ?",
					tenantId.String(),
					applicationId.String(),
					key,
					addressId.String()).Iter()

				defer iter.Close()

				Expect(iter.Scan(&addressValue)).To(BeFalse())
			}
		})

		It("should update the records in address table", func() {
			mockUUIDGeneratorService.
				EXPECT().
				GenerateRandomUUID().
				Return(addressId, nil)

			returnedAddressId, err := addressDataService.Create(
				tenantId,
				applicationId,
				shared.Address{AddressDetails: createRandomAddressDetails()})

			Expect(err).To(BeNil())

			expectedAddressDetails := createRandomAddressDetails()

			err = addressDataService.Update(
				tenantId,
				applicationId,
				returnedAddressId,
				shared.Address{AddressDetails: expectedAddressDetails})

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

			defer iter.Close()

			var key string
			var value string

			addressKeysValues := make(map[string]string)

			for iter.Scan(&key, &value) {
				addressKeysValues[key] = value
			}

			Expect(expectedAddressDetails).To(Equal(addressKeysValues))
		})

		It("should update the records in address_indexed_by_address_key table", func() {
			mockUUIDGeneratorService.
				EXPECT().
				GenerateRandomUUID().
				Return(addressId, nil)

			returnedAddressId, err := addressDataService.Create(
				tenantId,
				applicationId,
				shared.Address{AddressDetails: createRandomAddressDetails()})

			Expect(err).To(BeNil())

			expectedAddressDetails := createRandomAddressDetails()

			err = addressDataService.Update(
				tenantId,
				applicationId,
				returnedAddressId,
				shared.Address{AddressDetails: expectedAddressDetails})

			Expect(err).To(BeNil())

			config := getClusterConfig()
			config.Keyspace = keyspace

			session, err := config.CreateSession()

			Expect(err).To(BeNil())

			defer session.Close()

			for key, value := range expectedAddressDetails {
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

func TestUpdateBehaviour(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Update method behaviour")
}
