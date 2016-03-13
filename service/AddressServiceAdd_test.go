package service_test

import (
	"testing"

	. "github.com/microbusinesses/AddressService/domain"
	. "github.com/microbusinesses/AddressService/service"
	. "github.com/microbusinesses/Micro-Businesses-Core/system"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Add method input parameters", func() {
	var (
		addressService AddressService
		tenantId       UUID
		validAddress   Address
		emptyAddress   Address
	)

	BeforeEach(func() {
		addressService = AddressService{}
		tenantId, _ = RandomUUID()
		validAddress = Address{"FirstName": "Morteza"}
		emptyAddress = make(Address)
	})

	Context("when empty tenant unique identifier provided", func() {
		It("should panic", func() {
			Ω(func() { addressService.Add(EmptyUUID, validAddress) }).Should(Panic())
		})
	})

	Context("when nil address provided", func() {
		It("should panic", func() {
			Ω(func() { addressService.Add(tenantId, nil) }).Should(Panic())
		})
	})

	Context("when empty address provided", func() {
		It("should panic", func() {
			Ω(func() { addressService.Add(tenantId, emptyAddress) }).Should(Panic())
		})
	})
})

func TestAdd(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Add method input parameters")
}
