// +build integration

package service_test

import (
	"fmt"
	"testing"

	"github.com/gocql/gocql"
	"github.com/golang/mock/gomock"
	"github.com/micro-business/AddressService/data/contract"
	"github.com/micro-business/AddressService/data/service"
	"github.com/micro-business/Micro-Business-Core/system"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ReadAll method behaviour", func() {
	var (
		mockCtrl                 *gomock.Controller
		addressDataService       *service.AddressDataService
		mockUUIDGeneratorService *MockUUIDGeneratorService
		tenantID                 system.UUID
		applicationID            system.UUID
		addressID                system.UUID
		clusterConfig            *gocql.ClusterConfig
	)

	BeforeEach(func() {
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
	})

	Context("when reading existing address", func() {
		It("should return error if address does not exist", func() {
			address, err := addressDataService.ReadAll(tenantID, applicationID, addressID)

			Expect(err).To(Equal(fmt.Errorf("Address not found. Address ID: %s", addressID.String())))
			Expect(address).To(Equal(contract.Address{}))
		})

		It("should return the existing address keys and values", func() {
			mockUUIDGeneratorService.
				EXPECT().
				GenerateRandomUUID().
				Return(addressID, nil)

			expectedAddressDetails := createRandomAddressDetails()

			expectedAddress := contract.Address{AddressDetails: expectedAddressDetails}
			returnedAddressID, err := addressDataService.Create(
				tenantID,
				applicationID,
				expectedAddress)

			Expect(err).To(BeNil())

			returnedAddress, err := addressDataService.ReadAll(
				tenantID,
				applicationID,
				returnedAddressID)

			Expect(err).To(BeNil())
			Expect(expectedAddress).To(Equal(returnedAddress))
		})
	})
})

func TestReadAllBehaviour(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ReadAll method behaviour")
}
