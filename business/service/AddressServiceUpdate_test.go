package service_test

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/microbusinesses/AddressService/business/domain"
	"github.com/microbusinesses/AddressService/business/service"
	"github.com/microbusinesses/AddressService/data/contract"
	"github.com/microbusinesses/Micro-Businesses-Core/system"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Update method input parameters and dependency test", func() {
	var (
		mockCtrl                   *gomock.Controller
		addressService             *service.AddressService
		mockAddressDataService     *MockAddressDataService
		tenantID                   system.UUID
		applicationID              system.UUID
		addressID                  system.UUID
		validAddress               domain.Address
		emptyAddress               domain.Address
		addressWithEmptyKey        domain.Address
		addressWithWhitespaceKey   domain.Address
		addressWithEmptyValue      domain.Address
		addressWithWhitespaceValue domain.Address
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockAddressDataService = NewMockAddressDataService(mockCtrl)

		addressService = &service.AddressService{AddressDataService: mockAddressDataService}

		tenantID, _ = system.RandomUUID()
		applicationID, _ = system.RandomUUID()
		addressID, _ = system.RandomUUID()
		validAddress = domain.Address{AddressDetails: map[string]string{"City": "Christchurch"}}
		emptyAddress = domain.Address{}
		addressWithEmptyKey = domain.Address{AddressDetails: map[string]string{"": "Christchurch"}}
		addressWithWhitespaceKey = domain.Address{AddressDetails: map[string]string{"    ": "Christchurch"}}
		addressWithEmptyValue = domain.Address{AddressDetails: map[string]string{"City": ""}}
		addressWithWhitespaceValue = domain.Address{AddressDetails: map[string]string{"City": "    "}}
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("when address data service not provided", func() {
		It("should panic", func() {
			addressService.AddressDataService = nil

			Ω(func() { addressService.Update(tenantID, applicationID, addressID, validAddress) }).Should(Panic())
		})
	})

	Describe("Input Parameters", func() {
		It("should panic when empty tenant unique identifier provided", func() {
			Ω(func() { addressService.Update(system.EmptyUUID, applicationID, addressID, validAddress) }).Should(Panic())
		})

		It("should panic when empty application unique identifier provided", func() {
			Ω(func() { addressService.Update(tenantID, system.EmptyUUID, addressID, validAddress) }).Should(Panic())
		})

		It("should panic when empty address unique identifier provided", func() {
			Ω(func() { addressService.Update(tenantID, applicationID, system.EmptyUUID, validAddress) }).Should(Panic())
		})

		It("should panic when address without address key provided", func() {
			Ω(func() { addressService.Update(tenantID, applicationID, addressID, emptyAddress) }).Should(Panic())
		})

		It("should panic when address with empty key provided", func() {
			Ω(func() { addressService.Update(tenantID, applicationID, addressID, addressWithEmptyKey) }).Should(Panic())
		})

		It("should panic when address with key contains whitespace only provided", func() {
			Ω(func() { addressService.Update(tenantID, applicationID, addressID, addressWithWhitespaceKey) }).Should(Panic())
		})

		It("should panic when address with empty value provided", func() {
			Ω(func() { addressService.Update(tenantID, applicationID, addressID, addressWithEmptyValue) }).Should(Panic())
		})

		It("should panic when address with value contains whitespace only provided", func() {
			Ω(func() { addressService.Update(tenantID, applicationID, addressID, addressWithWhitespaceValue) }).Should(Panic())
		})
	})
})

var _ = Describe("Update method behaviour", func() {
	var (
		mockCtrl               *gomock.Controller
		addressService         *service.AddressService
		mockAddressDataService *MockAddressDataService
		tenantID               system.UUID
		applicationID          system.UUID
		addressID              system.UUID
		validAddress           domain.Address
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockAddressDataService = NewMockAddressDataService(mockCtrl)

		addressService = &service.AddressService{AddressDataService: mockAddressDataService}

		tenantID, _ = system.RandomUUID()
		applicationID, _ = system.RandomUUID()
		addressID, _ = system.RandomUUID()
		validAddress = domain.Address{AddressDetails: map[string]string{"City": "Christchurch"}}
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	It("should call address data service Update function", func() {
		mappedAddress := contract.Address{AddressDetails: validAddress.AddressDetails}

		mockAddressDataService.EXPECT().Update(tenantID, applicationID, addressID, mappedAddress)

		addressService.Update(tenantID, applicationID, addressID, validAddress)
	})

	Context("when address data service succeeds to update the requested address", func() {
		It("should return no error", func() {
			mappedAddress := contract.Address{AddressDetails: validAddress.AddressDetails}

			mockAddressDataService.
				EXPECT().
				Update(tenantID, applicationID, addressID, mappedAddress).
				Return(nil)

			err := addressService.Update(tenantID, applicationID, addressID, validAddress)

			Expect(err).To(BeNil())
		})
	})

	Context("when address data service fails to update the requested address", func() {
		It("should return the error returned by address data service", func() {
			mappedAddress := contract.Address{AddressDetails: validAddress.AddressDetails}

			expectedErrorID, _ := system.RandomUUID()
			expectedError := errors.New(expectedErrorID.String())
			mockAddressDataService.
				EXPECT().
				Update(tenantID, applicationID, addressID, mappedAddress).
				Return(expectedError)

			err := addressService.Update(tenantID, applicationID, addressID, validAddress)

			Expect(err).To(Equal(expectedError))
		})
	})
})

func TestUpdate(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Update method input parameters and dependency test")
	RunSpecs(t, "Update method behaviour")
}
