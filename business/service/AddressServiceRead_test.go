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

var _ = Describe("Read method input parameters and dependency test", func() {
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

	DescribeTable("Input Parameters",
		func(tenantId, applicationId, addressId system.UUID) {
			Ω(func() { addressService.Read(tenantId, applicationId, addressId) }).Should(Panic())
		},
		Entry("should panic when empty tenant unique identifier provided", system.EmptyUUID, applicationId, addressId),
		Entry("should panic when empty application unique identifier provided", tenantId, system.EmptyUUID, addressId),
		Entry("should panic when empty address unique identifier provided", tenantId, applicationId, system.EmptyUUID))
})

var _ = Describe("Read method behaviour", func() {
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

	It("should call address data service Read function", func() {
		mockAddressDataService.EXPECT().Read(tenantId, applicationId, addressId)

		addressService.Read(tenantId, applicationId, addressId)
	})

	Context("when address data service succeeds to read the requested address", func() {
		It("should return no error", func() {
			addressKeysValues := make(map[string]string)

			for idx := 0; idx < rand.Intn(10)+1; idx++ {
				key, _ := system.RandomUUID()
				value, _ := system.RandomUUID()

				addressKeysValues[key.String()] = value.String()
			}

			expectedAddress := domain.Address{AddressKeysValues: addressKeysValues}
			mockAddressDataService.
				EXPECT().
				Read(tenantId, applicationId, addressId).
				Return(dataShared.Address{AddressKeysValues: expectedAddress.AddressKeysValues}, nil)

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
	RunSpecs(t, "Read method input parameters and dependency test")
	RunSpecs(t, "Read method behaviour")
}
