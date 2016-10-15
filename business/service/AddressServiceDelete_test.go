package service_test

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/microbusinesses/AddressService/business/service"
	"github.com/microbusinesses/Micro-Businesses-Core/system"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Delete method input parameters and dependency test", func() {
	var (
		mockCtrl               *gomock.Controller
		addressService         *service.AddressService
		mockAddressDataService *MockAddressDataService
		tenantID               system.UUID
		applicationID          system.UUID
		addressID              system.UUID
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockAddressDataService = NewMockAddressDataService(mockCtrl)

		addressService = &service.AddressService{AddressDataService: mockAddressDataService}

		tenantID, _ = system.RandomUUID()
		applicationID, _ = system.RandomUUID()
		addressID, _ = system.RandomUUID()
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("when address data service not provided", func() {
		It("should panic", func() {
			addressService.AddressDataService = nil

			Ω(func() { addressService.Delete(tenantID, applicationID, addressID) }).Should(Panic())
		})
	})

	DescribeTable("Input Parameters",
		func(tenantID, applicationID, addressID system.UUID) {
			Ω(func() { addressService.Delete(tenantID, applicationID, addressID) }).Should(Panic())
		},
		Entry("should panic when empty tenant unique identifier provided", system.EmptyUUID, applicationID, addressID),
		Entry("should panic when empty application unique identifier provided", tenantID, system.EmptyUUID, addressID),
		Entry("should panic when empty address unique identifier provided", tenantID, applicationID, system.EmptyUUID))
})

var _ = Describe("Delete method behaviour", func() {
	var (
		mockCtrl               *gomock.Controller
		addressService         *service.AddressService
		mockAddressDataService *MockAddressDataService
		tenantID               system.UUID
		applicationID          system.UUID
		addressID              system.UUID
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockAddressDataService = NewMockAddressDataService(mockCtrl)

		addressService = &service.AddressService{AddressDataService: mockAddressDataService}

		tenantID, _ = system.RandomUUID()
		applicationID, _ = system.RandomUUID()
		addressID, _ = system.RandomUUID()
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	It("should call address data service Delete function", func() {
		mockAddressDataService.EXPECT().Delete(tenantID, applicationID, addressID)

		addressService.Delete(tenantID, applicationID, addressID)
	})

	Context("when address data service succeeds to delete the requested address", func() {
		It("should return no error", func() {
			mockAddressDataService.
				EXPECT().
				Delete(tenantID, applicationID, addressID).
				Return(nil)

			err := addressService.Delete(tenantID, applicationID, addressID)

			Expect(err).To(BeNil())
		})
	})

	Context("when address data service fails to delete the requested address", func() {
		It("should return the error returned by address data service", func() {
			expectedErrorID, _ := system.RandomUUID()
			expectedError := errors.New(expectedErrorID.String())
			mockAddressDataService.
				EXPECT().
				Delete(tenantID, applicationID, addressID).
				Return(expectedError)

			err := addressService.Delete(tenantID, applicationID, addressID)

			Expect(err).To(Equal(expectedError))
		})
	})
})

func TestDelete(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Delete method input parameters and dependency test")
	RunSpecs(t, "Delete method behaviour")
}
