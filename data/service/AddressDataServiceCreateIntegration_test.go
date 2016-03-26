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
		validAddress             shared.Address
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockUUIDGeneratorService = dataServiceMocks.NewMockUUIDGeneratorService(mockCtrl)

		addressDataService = &service.AddressDataService{UUIDGeneratorService: mockUUIDGeneratorService, ClusterConfig: &gocql.ClusterConfig{}}

		tenantId, _ = system.RandomUUID()
		applicationId, _ = system.RandomUUID()
		validAddress = shared.Address{AddressParts: map[string]string{"FirstName": "Morteza"}}
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	It("should call UUID generator service to create the address unique identifier", func() {
		mockUUIDGeneratorService.EXPECT().GenerateRandomUUID()

		addressDataService.Create(tenantId, applicationId, validAddress)
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

func TestCreateIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Create method behaviour")
}
