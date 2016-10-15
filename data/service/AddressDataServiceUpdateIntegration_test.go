// +build integration

package service_test

import (
	"testing"

	"github.com/gocql/gocql"
	"github.com/golang/mock/gomock"
	"github.com/microbusinesses/AddressService/data/contract"
	"github.com/microbusinesses/AddressService/data/service"
	"github.com/microbusinesses/Micro-Businesses-Core/system"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Update method behaviour", func() {
	var (
		mockCtrl                 *gomock.Controller
		addressDataService       *service.AddressDataService
		mockUUIDGeneratorService *MockUUIDGeneratorService
		tenantID                 system.UUID
		applicationID            system.UUID
		addressID                system.UUID
		clusterConfig            *gocql.ClusterConfig
		keyspace                 string
	)

	BeforeEach(func() {
		keyspace = createRandomKeyspace()

		createAddressKeyspaceAndAllRequiredTables(keyspace)

		clusterConfig = getClusterConfig()
		clusterConfig.Keyspace = keyspace

		mockCtrl = gomock.NewController(GinkgoT())
		mockUUIDGeneratorService = NewMockUUIDGeneratorService(mockCtrl)

		addressDataService = &service.AddressDataService{UUIDGeneratorService: mockUUIDGeneratorService, ClusterConfig: clusterConfig}

		tenantID, _ = system.RandomUUID()
		applicationID, _ = system.RandomUUID()
		addressID, _ = system.RandomUUID()
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
				Return(addressID, nil)

			addressDetailsToAdd := createRandomAddressDetails()

			returnedAddressID, err := addressDataService.Create(
				tenantID,
				applicationID,
				contract.Address{AddressDetails: addressDetailsToAdd})

			Expect(err).To(BeNil())

			err = addressDataService.Update(
				tenantID,
				applicationID,
				returnedAddressID,
				contract.Address{AddressDetails: createRandomAddressDetails()})

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

			for iter.Scan(&key, &value) {
				_, ok := addressDetailsToAdd[key]

				Expect(ok).To(BeFalse())
			}
		})

		It("should remove all old records from address_indexed_by_address_key table", func() {
			mockUUIDGeneratorService.
				EXPECT().
				GenerateRandomUUID().
				Return(addressID, nil)

			addressDetailsToAdd := createRandomAddressDetails()

			returnedAddressID, err := addressDataService.Create(
				tenantID,
				applicationID,
				contract.Address{AddressDetails: addressDetailsToAdd})

			Expect(err).To(BeNil())

			err = addressDataService.Update(
				tenantID,
				applicationID,
				returnedAddressID,
				contract.Address{AddressDetails: createRandomAddressDetails()})

			config := getClusterConfig()
			config.Keyspace = keyspace

			session, err := config.CreateSession()

			Expect(err).To(BeNil())

			defer session.Close()

			for key := range addressDetailsToAdd {
				var addressValue string

				iter := session.Query(
					"SELECT address_id, address_value"+
						" FROM address_indexed_by_address_key"+
						" WHERE"+
						" tenant_id = ?"+
						" AND application_id = ?"+
						" AND address_key = ?"+
						" AND address_id = ?",
					tenantID.String(),
					applicationID.String(),
					key,
					addressID.String()).Iter()

				defer iter.Close()

				Expect(iter.Scan(&addressValue)).To(BeFalse())
			}
		})

		It("should update the records in address table", func() {
			mockUUIDGeneratorService.
				EXPECT().
				GenerateRandomUUID().
				Return(addressID, nil)

			returnedAddressID, err := addressDataService.Create(
				tenantID,
				applicationID,
				contract.Address{AddressDetails: createRandomAddressDetails()})

			Expect(err).To(BeNil())

			expectedAddressDetails := createRandomAddressDetails()

			err = addressDataService.Update(
				tenantID,
				applicationID,
				returnedAddressID,
				contract.Address{AddressDetails: expectedAddressDetails})

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

		It("should update the records in address_indexed_by_address_key table", func() {
			mockUUIDGeneratorService.
				EXPECT().
				GenerateRandomUUID().
				Return(addressID, nil)

			returnedAddressID, err := addressDataService.Create(
				tenantID,
				applicationID,
				contract.Address{AddressDetails: createRandomAddressDetails()})

			Expect(err).To(BeNil())

			expectedAddressDetails := createRandomAddressDetails()

			err = addressDataService.Update(
				tenantID,
				applicationID,
				returnedAddressID,
				contract.Address{AddressDetails: expectedAddressDetails})

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

func TestUpdateBehaviour(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Update method behaviour")
}
