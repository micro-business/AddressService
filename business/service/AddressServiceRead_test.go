package service_test

import (
	"errors"
	"math/rand"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/microbusinesses/AddressService/business/domain"
	"github.com/microbusinesses/AddressService/business/service"
	dataServiceMocks "github.com/microbusinesses/AddressService/data/contract/mocks"
	dataShared "github.com/microbusinesses/AddressService/data/shared"
	"github.com/microbusinesses/Micro-Businesses-Core/system"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Read method input parameters", func() {
	var (
		mockCtrl               *gomock.Controller
		addressService         *service.AddressService
		mockAddressDataService *dataServiceMocks.MockAddressDataService
		tenantId               system.UUID
		applicationId          system.UUID
		addressId              system.UUID
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockAddressDataService = dataServiceMocks.NewMockAddressDataService(mockCtrl)

		addressService = &service.AddressService{AddressDataService: mockAddressDataService}

		tenantId, _ = system.RandomUUID()
		applicationId, _ = system.RandomUUID()
		addressId, _ = system.RandomUUID()
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("when address data service not provided", func() {
		It("should panic", func() {
			addressService.AddressDataService = nil

			Ω(func() { addressService.Read(tenantId, applicationId, addressId) }).Should(Panic())
		})
	})

	Context("when empty tenant unique identifier provided", func() {
		It("should panic", func() {
			Ω(func() { addressService.Read(system.EmptyUUID, applicationId, addressId) }).Should(Panic())
		})
	})

	Context("when empty tenant unique identifier provided", func() {
		It("should panic", func() {
			Ω(func() { addressService.Read(tenantId, system.EmptyUUID, addressId) }).Should(Panic())
		})
	})

	Context("when empty address unique identifier provided", func() {
		It("should panic", func() {
			Ω(func() { addressService.Read(tenantId, applicationId, system.EmptyUUID) }).Should(Panic())
		})
	})

	It("should call address data service Read function", func() {
		mockAddressDataService.EXPECT().Read(tenantId, applicationId, addressId)

		addressService.Read(tenantId, applicationId, addressId)
	})

	Context("when address data service succeeds to read the requested address", func() {
		It("should return no error", func() {
			addressParts := make(map[string]string)

			for idx := 0; idx < rand.Intn(10); idx++ {
				addressPartKey, _ := system.RandomUUID()
				addressPartValue, _ := system.RandomUUID()

				addressParts[addressPartKey.String()] = addressPartValue.String()
			}

			expectedAddress := domain.Address{AddressParts: addressParts}
			mockAddressDataService.
				EXPECT().
				Read(tenantId, applicationId, addressId).
				Return(dataShared.Address{AddressParts: expectedAddress.AddressParts}, nil)

			address, err := addressService.Read(tenantId, applicationId, addressId)

			Expect(address).To(Equal(expectedAddress))
			Expect(err).To(BeNil())
		})
	})

	Context("when address data service fails to read the requested address", func() {
		It("should return the error returned by address data service", func() {
			expectedErrorId, _ := system.RandomUUID()
			expectedError := errors.New(expectedErrorId.String())
			mockAddressDataService.
				EXPECT().
				Read(tenantId, applicationId, addressId).
				Return(dataShared.Address{}, expectedError)

			expectedAddress, err := addressService.Read(tenantId, applicationId, addressId)

			Expect(expectedAddress).To(Equal(domain.Address{}))
			Expect(err).To(Equal(expectedError))
		})
	})
})

func TestRead(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Read method input parameters")
}
