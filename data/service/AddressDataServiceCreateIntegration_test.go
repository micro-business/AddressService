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
		tenantID                 system.UUID
		applicationID            system.UUID
		addressID                system.UUID
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

		tenantID, _ = system.RandomUUID()
		applicationID, _ = system.RandomUUID()
		addressID, _ = system.RandomUUID()
		validAddress = shared.Address{AddressDetails: map[string]string{"City": "Christchurch"}}
	})

	AfterEach(func() {
		mockCtrl.Finish()
		dropKeyspace(keyspace)
	})

	Context("when UUID generator service succeeds to create the new UUID", func() {
		It("should return the new UUID as address uniuqe identifier and no error", func() {
			expectedAddressID, _ := system.RandomUUID()
			mockUUIDGeneratorService.
				EXPECT().
				GenerateRandomUUID().
				Return(expectedAddressID, nil)

			newAddressID, err := addressDataService.Create(tenantID, applicationID, validAddress)

			Expect(expectedAddressID).To(Equal(newAddressID))
			Expect(err).To(BeNil())
		})
	})

	Context("when UUID generator service fails to create the new UUID", func() {
		It("should return address unique identifier as empty UUID and the returned error by address data service", func() {
			expectedErrorID, _ := system.RandomUUID()
			expectedError := errors.New(expectedErrorID.String())
			mockUUIDGeneratorService.
				EXPECT().
				GenerateRandomUUID().
				Return(system.EmptyUUID, expectedError)

			newAddressID, err := addressDataService.Create(tenantID, applicationID, validAddress)

			Expect(newAddressID).To(Equal(system.EmptyUUID))
			Expect(err).To(Equal(expectedError))
		})
	})

	Context("when creating new address", func() {
		It("should insert the records into address table", func() {
			mockUUIDGeneratorService.
				EXPECT().
				GenerateRandomUUID().
				Return(addressID, nil)

			expectedAddressDetails := createRandomAddressDetails()

			returnedAddressID, err := addressDataService.Create(
				tenantID,
				applicationID,
				shared.Address{AddressDetails: expectedAddressDetails})

			Expect(addressID).To(Equal(returnedAddressID))
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
				tenantID.String(),
				applicationID.String(),
				addressID.String()).Iter()

			defer iter.Close()

			var key string
			var value string

			addressDetails := make(map[string]string)

			for iter.Scan(&key, &value) {
				addressDetails[key] = value
			}

			Expect(expectedAddressDetails).To(Equal(addressDetails))
		})

		It("should insert the records into address_indexed_by_address_key table", func() {
			mockUUIDGeneratorService.
				EXPECT().
				GenerateRandomUUID().
				Return(addressID, nil)

			expectedAddressDetails := createRandomAddressDetails()

			addressDataService.Create(tenantID, applicationID, shared.Address{AddressDetails: expectedAddressDetails})

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
					tenantID.String(),
					applicationID.String(),
					key).Scan(&id, &addressValue)

				Expect(err).To(BeNil())

				Expect(addressID).To(Equal(mapGocqlUUIDToSystemUUID(id)))
				Expect(value).To(Equal(addressValue))
			}
		})
	})
})

func TestCreateBehaviour(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Create method behaviour")
}
