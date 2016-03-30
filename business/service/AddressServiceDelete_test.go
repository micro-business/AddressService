package service_test

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/microbusinesses/AddressService/business/service"
	dataServiceMocks "github.com/microbusinesses/AddressService/data/contract/mocks"
	"github.com/microbusinesses/Micro-Businesses-Core/system"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Delete method input parameters and dependency test", func() {
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

			Ω(func() { addressService.Delete(tenantId, applicationId, addressId) }).Should(Panic())
		})
	})

	DescribeTable("Input Parameters",
		func(tenantId, applicationId, addressId system.UUID) {
			Ω(func() { addressService.Delete(tenantId, applicationId, addressId) }).Should(Panic())
		},
		Entry("should panic when empty tenant unique identifier provided", system.EmptyUUID, applicationId, addressId),
		Entry("should panic when empty application unique identifier provided", tenantId, system.EmptyUUID, addressId),
		Entry("should panic when empty address unique identifier provided", tenantId, applicationId, system.EmptyUUID))
})

var _ = Describe("Delete method behaviour", func() {
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

	It("should call address data service Delete function", func() {
		mockAddressDataService.EXPECT().Delete(tenantId, applicationId, addressId)

		addressService.Delete(tenantId, applicationId, addressId)
	})

	Context("when address data service succeeds to delete the requested address", func() {
		It("should return no error", func() {
			mockAddressDataService.
				EXPECT().
				Delete(tenantId, applicationId, addressId).
				Return(nil)

			err := addressService.Delete(tenantId, applicationId, addressId)

			Expect(err).To(BeNil())
		})
	})

	Context("when address data service fails to delete the requested address", func() {
		It("should return the error returned by address data service", func() {
			expectedErrorId, _ := system.RandomUUID()
			expectedError := errors.New(expectedErrorId.String())
			mockAddressDataService.
				EXPECT().
				Delete(tenantId, applicationId, addressId).
				Return(expectedError)

			err := addressService.Delete(tenantId, applicationId, addressId)

			Expect(err).To(Equal(expectedError))
		})
	})
})

func TestDelete(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Delete method input parameters and dependency test")
	RunSpecs(t, "Delete method behaviour")
}
