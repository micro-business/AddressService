package service_test

import (
	"testing"

	"github.com/microbusinesses/AddressService/business/domain"
	"github.com/microbusinesses/AddressService/business/service"
	. "github.com/microbusinesses/Micro-Businesses-Core/system"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Create method input parameters", func() {
	var (
		addressService service.AddressService
		tenantId       UUID
		applicationId  UUID
		validAddress   domain.Address
		emptyAddress   domain.Address
	)

	BeforeEach(func() {
		addressService = service.AddressService{}
		tenantId, _ = RandomUUID()
		applicationId, _ = RandomUUID()
		validAddress = domain.Address{AddressParts: map[string]string{"FirstName": "Morteza"}}
		emptyAddress = domain.Address{}
	})

	Context("when address data service not provided", func() {
		It("should panic", func() {
			addressService.AddressDataService = nil

			Ω(func() { addressService.Create(tenantId, applicationId, validAddress) }).Should(Panic())
		})
	})

	Context("when empty tenant unique identifier provided", func() {
		It("should panic", func() {
			Ω(func() { addressService.Create(EmptyUUID, applicationId, validAddress) }).Should(Panic())
		})
	})

	Context("when empty application unique identifier provided", func() {
		It("should panic", func() {
			Ω(func() { addressService.Create(tenantId, EmptyUUID, validAddress) }).Should(Panic())
		})
	})

	Context("when address without address parts provided", func() {
		It("should panic", func() {
			Ω(func() { addressService.Create(tenantId, applicationId, emptyAddress) }).Should(Panic())
		})
	})
})

func TestCreate(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Create method input parameters")
}
