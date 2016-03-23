package service_test

import (
	"testing"

	. "github.com/microbusinesses/AddressService/service"
	. "github.com/microbusinesses/Micro-Businesses-Core/system"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Delete method input parameters", func() {
	var (
		addressService AddressService
		tenantId       UUID
		applicationId  UUID
		addressId      UUID
	)

	BeforeEach(func() {
		addressService = AddressService{}
		tenantId, _ = RandomUUID()
		applicationId, _ = RandomUUID()
		addressId, _ = RandomUUID()
	})

	Context("when empty tenant unique identifier provided", func() {
		It("should panic", func() {
			Ω(func() { addressService.Delete(EmptyUUID, applicationId, addressId) }).Should(Panic())
		})
	})

	Context("when empty application unique identifier provided", func() {
		It("should panic", func() {
			Ω(func() { addressService.Delete(tenantId, EmptyUUID, addressId) }).Should(Panic())
		})
	})

	Context("when empty address unique identifier provided", func() {
		It("should panic", func() {
			Ω(func() { addressService.Delete(tenantId, applicationId, EmptyUUID) }).Should(Panic())
		})
	})
})

func TestDelete(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Delete method input parameters")
}
