package service_test

import (
	"errors"
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

var _ = Describe("Update method input parameters and dependency test", func() {
	var (
		mockCtrl               *gomock.Controller
		addressService         *service.AddressService
		mockAddressDataService *dataServiceMocks.MockAddressDataService
		tenantID               system.UUID
		applicationID          system.UUID
		addressID              system.UUID
		validAddress           domain.Address
		emptyAddress           domain.Address
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockAddressDataService = dataServiceMocks.NewMockAddressDataService(mockCtrl)

		addressService = &service.AddressService{AddressDataService: mockAddressDataService}

		tenantID, _ = system.RandomUUID()
		applicationID, _ = system.RandomUUID()
		addressID, _ = system.RandomUUID()
		validAddress = domain.Address{AddressDetails: map[string]string{"City": "Christchurch"}}
		emptyAddress = domain.Address{}
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

	DescribeTable("Input Parameters",
		func(tenantID, applicationID, addressID system.UUID, address domain.Address) {
			Ω(func() { addressService.Update(tenantID, applicationID, addressID, address) }).Should(Panic())
		},
		Entry("should panic when empty tenant unique identifier provided", system.EmptyUUID, applicationID, addressID, validAddress),
		Entry("should panic when empty application unique identifier provided", tenantID, system.EmptyUUID, addressID, validAddress),
		Entry("should panic when empty address unique identifier provided", tenantID, applicationID, system.EmptyUUID, validAddress),
		Entry("should panic when address without address key provided", tenantID, applicationID, addressID, emptyAddress))
})

var _ = Describe("Update method behaviour", func() {
	var (
		mockCtrl               *gomock.Controller
		addressService         *service.AddressService
		mockAddressDataService *dataServiceMocks.MockAddressDataService
		tenantID               system.UUID
		applicationID          system.UUID
		addressID              system.UUID
		validAddress           domain.Address
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockAddressDataService = dataServiceMocks.NewMockAddressDataService(mockCtrl)

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
		mappedAddress := dataShared.Address{AddressDetails: validAddress.AddressDetails}

		mockAddressDataService.EXPECT().Update(tenantID, applicationID, addressID, mappedAddress)

		addressService.Update(tenantID, applicationID, addressID, validAddress)
	})

	Context("when address data service succeeds to update the requested address", func() {
		It("should return no error", func() {
			mappedAddress := dataShared.Address{AddressDetails: validAddress.AddressDetails}

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
			mappedAddress := dataShared.Address{AddressDetails: validAddress.AddressDetails}

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
