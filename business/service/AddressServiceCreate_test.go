package service_test

import (
	"errors"
	"math/rand"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/microbusinesses/AddressService/business/domain"
	"github.com/microbusinesses/AddressService/business/service"
	"github.com/microbusinesses/AddressService/data/contract"
	"github.com/microbusinesses/Micro-Businesses-Core/system"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Create method input parameters and dependency test", func() {
	var (
		mockCtrl               *gomock.Controller
		addressService         *service.AddressService
		mockAddressDataService *MockAddressDataService
		tenantID               system.UUID
		applicationID          system.UUID
		validAddress           domain.Address
		emptyAddress           domain.Address
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockAddressDataService = NewMockAddressDataService(mockCtrl)

		addressService = &service.AddressService{AddressDataService: mockAddressDataService}

		tenantID, _ = system.RandomUUID()
		applicationID, _ = system.RandomUUID()
		validAddress = domain.Address{AddressDetails: map[string]string{"City": "Christchurch"}}
		emptyAddress = domain.Address{}
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("when address data service not provided", func() {
		It("should panic", func() {
			addressService.AddressDataService = nil

			Ω(func() { addressService.Create(tenantID, applicationID, validAddress) }).Should(Panic())
		})
	})

	DescribeTable("Input Parameters",
		func(tenantID, applicationID system.UUID, address domain.Address) {
			Ω(func() { addressService.Create(tenantID, applicationID, address) }).Should(Panic())
		},
		Entry("should panic when empty tenant unique identifier provided", system.EmptyUUID, applicationID, validAddress),
		Entry("should panic when empty application unique identifier provided", tenantID, system.EmptyUUID, validAddress),
		Entry("should panic when address without address key provided", tenantID, applicationID, emptyAddress))
})

var _ = Describe("Create method behaviour", func() {
	var (
		mockCtrl               *gomock.Controller
		addressService         *service.AddressService
		mockAddressDataService *MockAddressDataService
		tenantID               system.UUID
		applicationID          system.UUID
		validAddress           domain.Address
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockAddressDataService = NewMockAddressDataService(mockCtrl)

		addressService = &service.AddressService{AddressDataService: mockAddressDataService}

		tenantID, _ = system.RandomUUID()
		applicationID, _ = system.RandomUUID()
		validAddress = domain.Address{AddressDetails: map[string]string{"City": "Christchurch"}}
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	It("should call address data service Create function", func() {
		mappedAddress := contract.Address{AddressDetails: validAddress.AddressDetails}

		mockAddressDataService.EXPECT().Create(tenantID, applicationID, mappedAddress)

		addressService.Create(tenantID, applicationID, validAddress)
	})

	Context("when address data service succeeds to create the new address", func() {
		It("should return the returned address unique identifier by address data service and no error", func() {
			addressDetails := make(map[string]string)

			for idx := 0; idx < rand.Intn(10)+1; idx++ {
				key, _ := system.RandomUUID()
				value, _ := system.RandomUUID()

				addressDetails[key.String()] = value.String()
			}

			mappedAddress := contract.Address{AddressDetails: addressDetails}

			expectedAddressID, _ := system.RandomUUID()
			mockAddressDataService.
				EXPECT().
				Create(tenantID, applicationID, mappedAddress).
				Return(expectedAddressID, nil)

			dnewAddressIDd, err := addressService.Create(tenantID, applicationID, domain.Address{AddressDetails: addressDetails})

			Expect(expectedAddressID).To(Equal(dnewAddressIDd))
			Expect(err).To(BeNil())
		})
	})

	Context("when address data service fails to create the new address", func() {
		It("should return address unique identifier as empty UUID and the returned error by address data service", func() {
			mappedAddress := contract.Address{AddressDetails: validAddress.AddressDetails}

			expectedErrorID, _ := system.RandomUUID()
			expectedError := errors.New(expectedErrorID.String())
			mockAddressDataService.
				EXPECT().
				Create(tenantID, applicationID, mappedAddress).
				Return(system.EmptyUUID, expectedError)

			dnewAddressIDd, err := addressService.Create(tenantID, applicationID, validAddress)

			Expect(dnewAddressIDd).To(Equal(system.EmptyUUID))
			Expect(err).To(Equal(expectedError))
		})
	})
})

func TestCreate(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Create method input parameters and dependency test")
	RunSpecs(t, "Create method behaviour")
}
