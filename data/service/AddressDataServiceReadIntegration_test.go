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

var _ = Describe("Read method behaviour", func() {
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

	It("should return the existing address keys and values", func() {
		mockUUIDGeneratorService.
			EXPECT().
			GenerateRandomUUID().
			Return(addressId, nil)

		expectedAddressKeysValues := createRandomAddressKeyValues()

		expectedAddress := shared.Address{AddressKeysValues: expectedAddressKeysValues}
		returnedAddressId, err := addressDataService.Create(
			tenantId,
			applicationId,
			expectedAddress)

		Expect(err).To(BeNil())

		returnedAddress, err := addressDataService.Read(
			tenantId,
			applicationId,
			returnedAddressId)

		Expect(err).To(BeNil())
		Expect(expectedAddress).To(Equal(returnedAddress))
	})
})

func TestReadBehaviour(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Read method behaviour")
}
