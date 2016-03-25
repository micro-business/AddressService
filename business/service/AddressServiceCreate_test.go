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

var _ = Describe("Create method input parameters", func() {
	var (
		mockCtrl               *gomock.Controller
		addressService         *service.AddressService
		mockAddressDataService *dataServiceMocks.MockAddressDataService
		tenantId               system.UUID
		applicationId          system.UUID
		validAddress           domain.Address
		emptyAddress           domain.Address
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockAddressDataService = dataServiceMocks.NewMockAddressDataService(mockCtrl)

		addressService = &service.AddressService{AddressDataService: mockAddressDataService}

		tenantId, _ = system.RandomUUID()
		applicationId, _ = system.RandomUUID()
		validAddress = domain.Address{AddressParts: map[string]string{"FirstName": "Morteza"}}
		emptyAddress = domain.Address{}
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("when address data service not provided", func() {
		It("should panic", func() {
			addressService.AddressDataService = nil

			Ω(func() { addressService.Create(tenantId, applicationId, validAddress) }).Should(Panic())
		})
	})

	Context("when empty tenant unique identifier provided", func() {
		It("should panic", func() {
			Ω(func() { addressService.Create(system.EmptyUUID, applicationId, validAddress) }).Should(Panic())
		})
	})

	Context("when empty application unique identifier provided", func() {
		It("should panic", func() {
			Ω(func() { addressService.Create(tenantId, system.EmptyUUID, validAddress) }).Should(Panic())
		})
	})

	Context("when address without address parts provided", func() {
		It("should panic", func() {
			Ω(func() { addressService.Create(tenantId, applicationId, emptyAddress) }).Should(Panic())
		})
	})
})

var _ = Describe("Create method behaviour", func() {
	var (
		mockCtrl               *gomock.Controller
		addressService         *service.AddressService
		mockAddressDataService *dataServiceMocks.MockAddressDataService
		tenantId               system.UUID
		applicationId          system.UUID
		validAddress           domain.Address
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockAddressDataService = dataServiceMocks.NewMockAddressDataService(mockCtrl)

		addressService = &service.AddressService{AddressDataService: mockAddressDataService}

		tenantId, _ = system.RandomUUID()
		applicationId, _ = system.RandomUUID()
		validAddress = domain.Address{AddressParts: map[string]string{"FirstName": "Morteza"}}
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	It("should call address data service Create function", func() {
		mappedAddress := dataShared.Address{AddressParts: validAddress.AddressParts}

		mockAddressDataService.EXPECT().Create(tenantId, applicationId, mappedAddress)

		addressService.Create(tenantId, applicationId, validAddress)
	})

	Context("when address data service succeeds to create the new address", func() {
		It("should return the returned address unique identifier by address data service and no error", func() {
			addressParts := make(map[string]string)

			for idx := 0; idx < rand.Intn(10); idx++ {
				addressPartKey, _ := system.RandomUUID()
				addressPartValue, _ := system.RandomUUID()

				addressParts[addressPartKey.String()] = addressPartValue.String()
			}

			mappedAddress := dataShared.Address{AddressParts: addressParts}

			expectedAddressId, _ := system.RandomUUID()
			mockAddressDataService.
				EXPECT().
				Create(tenantId, applicationId, mappedAddress).
				Return(expectedAddressId, nil)

			newAddressId, err := addressService.Create(tenantId, applicationId, domain.Address{AddressParts: addressParts})

			Expect(expectedAddressId).To(Equal(newAddressId))
			Expect(err).To(BeNil())
		})
	})

	Context("when address data service fails to create the new address", func() {
		It("should return address unique identifier as empty UUID and the returned error by address data service", func() {
			mappedAddress := dataShared.Address{AddressParts: validAddress.AddressParts}

			expectedErrorId, _ := system.RandomUUID()
			expectedError := errors.New(expectedErrorId.String())
			mockAddressDataService.
				EXPECT().
				Create(tenantId, applicationId, mappedAddress).
				Return(system.EmptyUUID, expectedError)

			newAddressId, err := addressService.Create(tenantId, applicationId, validAddress)

			Expect(newAddressId).To(Equal(system.EmptyUUID))
			Expect(err).To(Equal(expectedError))
		})
	})
})

func TestCreate(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Create method input parameters")
	RunSpecs(t, "Create method behaviour")
}
