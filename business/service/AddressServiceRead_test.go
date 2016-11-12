package service_test

import (
	"errors"
	"math/rand"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/micro-business/AddressService/business/domain"
	"github.com/micro-business/AddressService/business/service"
	"github.com/micro-business/AddressService/data/contract"
	"github.com/micro-business/Micro-Business-Core/system"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Read method input parameters and dependency test", func() {
	var (
		mockCtrl                *gomock.Controller
		addressService          *service.AddressService
		mockAddressDataService  *MockAddressDataService
		tenantID                system.UUID
		applicationID           system.UUID
		addressID               system.UUID
		validKeys               []string
		emptyKeys               []string
		keysWithEmptyValue      []string
		keysWithWhitespaceValue []string
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockAddressDataService = NewMockAddressDataService(mockCtrl)

		addressService = &service.AddressService{AddressDataService: mockAddressDataService}

		tenantID, _ = system.RandomUUID()
		applicationID, _ = system.RandomUUID()
		addressID, _ = system.RandomUUID()

		validKeys = make([]string, 1)
		randomKey, _ := system.RandomUUID()
		validKeys[0] = randomKey.String()

		emptyKeys = make([]string, 0)

		keysWithEmptyValue = make([]string, 1)
		keysWithEmptyValue[0] = ""

		keysWithWhitespaceValue = make([]string, 1)
		keysWithWhitespaceValue[0] = "    "
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("when address data service not provided", func() {
		It("should panic", func() {
			addressService.AddressDataService = nil

			Ω(func() { addressService.Read(tenantID, applicationID, addressID, validKeys) }).Should(Panic())
		})
	})

	Describe("Input Parameters", func() {
		It("should panic when empty tenant unique identifier provided", func() {
			Ω(func() { addressService.Read(system.EmptyUUID, applicationID, addressID, validKeys) }).Should(Panic())
		})

		It("should panic when empty application unique identifier provided", func() {
			Ω(func() { addressService.Read(tenantID, system.EmptyUUID, addressID, validKeys) }).Should(Panic())
		})

		It("should panic when empty address unique identifier provided", func() {
			Ω(func() { addressService.Read(tenantID, applicationID, system.EmptyUUID, validKeys) }).Should(Panic())
		})

		It("should panic when empty keys provided", func() {
			Ω(func() { addressService.Read(tenantID, applicationID, addressID, emptyKeys) }).Should(Panic())
		})

		It("should panic when keys with empty value provided", func() {
			Ω(func() { addressService.Read(tenantID, applicationID, addressID, keysWithEmptyValue) }).Should(Panic())
		})

		It("should panic when keys with whitespace only value provided", func() {
			Ω(func() { addressService.Read(tenantID, applicationID, addressID, keysWithWhitespaceValue) }).Should(Panic())
		})

	})
})

var _ = Describe("Read method behaviour", func() {
	var (
		mockCtrl               *gomock.Controller
		addressService         *service.AddressService
		mockAddressDataService *MockAddressDataService
		tenantID               system.UUID
		applicationID          system.UUID
		addressID              system.UUID
		validKeys              []string
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockAddressDataService = NewMockAddressDataService(mockCtrl)

		addressService = &service.AddressService{AddressDataService: mockAddressDataService}

		tenantID, _ = system.RandomUUID()
		applicationID, _ = system.RandomUUID()
		addressID, _ = system.RandomUUID()

		validKeys = make([]string, 1)
		randomKey, _ := system.RandomUUID()
		validKeys[0] = randomKey.String()
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	It("should call address data service Read function", func() {
		mockAddressDataService.EXPECT().Read(tenantID, applicationID, addressID, validKeys)

		addressService.Read(tenantID, applicationID, addressID, validKeys)
	})

	Context("when address data service succeeds to read the requested address", func() {
		It("should return no error", func() {
			addressDetails := make(map[string]string)

			for idx := 0; idx < rand.Intn(10)+1; idx++ {
				key, _ := system.RandomUUID()
				value, _ := system.RandomUUID()

				addressDetails[key.String()] = value.String()
			}

			keys := make([]string, 0, len(addressDetails))

			for key := range addressDetails {
				keys = append(keys, key)
			}

			expectedAddress := domain.Address{AddressDetails: addressDetails}
			mockAddressDataService.
				EXPECT().
				Read(tenantID, applicationID, addressID, keys).
				Return(contract.Address{AddressDetails: expectedAddress.AddressDetails}, nil)

			address, err := addressService.Read(tenantID, applicationID, addressID, keys)

			Expect(address).To(Equal(expectedAddress))
			Expect(err).To(BeNil())
		})
	})

	Context("when address data service fails to read the requested address", func() {
		It("should return the error returned by address data service", func() {
			expectedErrorID, _ := system.RandomUUID()
			expectedError := errors.New(expectedErrorID.String())
			mockAddressDataService.
				EXPECT().
				Read(tenantID, applicationID, addressID, validKeys).
				Return(contract.Address{}, expectedError)

			expectedAddress, err := addressService.Read(tenantID, applicationID, addressID, validKeys)

			Expect(expectedAddress).To(Equal(domain.Address{}))
			Expect(err).To(Equal(expectedError))
		})
	})
})

func TestRead(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Read method input parameters and dependency test")
	RunSpecs(t, "Read method behaviour")
}
