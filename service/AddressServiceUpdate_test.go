package service_test

import (
	"testing"

	. "github.com/microbusinesses/AddressService/domain"
	. "github.com/microbusinesses/AddressService/service"
	. "github.com/microbusinesses/Micro-Businesses-Core/system"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Update method input parameters", func() {
	var (
		addressService AddressService
		tenantId       UUID
		addressId      UUID
		validAddress   Address
		emptyAddress   Address
	)

	BeforeEach(func() {
		addressService = AddressService{}
		tenantId, _ = RandomUUID()
		addressId, _ = RandomUUID()
		validAddress = Address{AddressParts: map[string]string{"FirstName": "Morteza"}}
		emptyAddress = Address{}
	})

	Context("when empty tenant unique identifier provided", func() {
		It("should panic", func() {
			Ω(func() { addressService.Update(EmptyUUID, addressId, validAddress) }).Should(Panic())
		})
	})

	Context("when empty address unique identifier provided", func() {
		It("should panic", func() {
			Ω(func() { addressService.Update(tenantId, EmptyUUID, validAddress) }).Should(Panic())
		})
	})

	Context("when address without address parts provided", func() {
		It("should panic", func() {
			Ω(func() { addressService.Update(tenantId, addressId, emptyAddress) }).Should(Panic())
		})
	})
})

func TestUpdate(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Update method input parameters")
}
