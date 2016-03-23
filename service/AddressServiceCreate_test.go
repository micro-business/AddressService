package service_test

import (
	"testing"

	. "github.com/microbusinesses/AddressService/domain"
	. "github.com/microbusinesses/AddressService/service"
	. "github.com/microbusinesses/Micro-Businesses-Core/system"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Create method input parameters", func() {
	var (
		addressService AddressService
		tenantId       UUID
		applicationId  UUID
		validAddress   Address
		emptyAddress   Address
	)

	BeforeEach(func() {
		addressService = AddressService{}
		tenantId, _ = RandomUUID()
		applicationId, _ = RandomUUID()
		validAddress = Address{AddressParts: map[string]string{"FirstName": "Morteza"}}
		emptyAddress = Address{}
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
