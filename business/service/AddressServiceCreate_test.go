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
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Create method input parameters and dependency test", func() {
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
		validAddress = domain.Address{AddressDetails: map[string]string{"City": "Christchurch"}}
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

	DescribeTable("Input Parameters",
		func(tenantId, applicationId system.UUID, address domain.Address) {
			Ω(func() { addressService.Create(tenantId, applicationId, address) }).Should(Panic())
		},
		Entry("should panic when empty tenant unique identifier provided", system.EmptyUUID, applicationId, validAddress),
		Entry("should panic when empty application unique identifier provided", tenantId, system.EmptyUUID, validAddress),
		Entry("should panic when address without address key provided", tenantId, applicationId, emptyAddress))
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
		validAddress = domain.Address{AddressDetails: map[string]string{"City": "Christchurch"}}
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	It("should call address data service Create function", func() {
		mappedAddress := dataShared.Address{AddressDetails: validAddress.AddressDetails}

		mockAddressDataService.EXPECT().Create(tenantId, applicationId, mappedAddress)

		addressService.Create(tenantId, applicationId, validAddress)
	})

	Context("when address data service succeeds to create the new address", func() {
		It("should return the returned address unique identifier by address data service and no error", func() {
			addressKeysValues := make(map[string]string)

			for idx := 0; idx < rand.Intn(10)+1; idx++ {
				key, _ := system.RandomUUID()
				value, _ := system.RandomUUID()

				addressKeysValues[key.String()] = value.String()
			}

			mappedAddress := dataShared.Address{AddressDetails: addressKeysValues}

			expectedAddressId, _ := system.RandomUUID()
			mockAddressDataService.
				EXPECT().
				Create(tenantId, applicationId, mappedAddress).
				Return(expectedAddressId, nil)

			newAddressId, err := addressService.Create(tenantId, applicationId, domain.Address{AddressDetails: addressKeysValues})

			Expect(expectedAddressId).To(Equal(newAddressId))
			Expect(err).To(BeNil())
		})
	})

	Context("when address data service fails to create the new address", func() {
		It("should return address unique identifier as empty UUID and the returned error by address data service", func() {
			mappedAddress := dataShared.Address{AddressDetails: validAddress.AddressDetails}

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
	RunSpecs(t, "Create method input parameters and dependency test")
	RunSpecs(t, "Create method behaviour")
}
