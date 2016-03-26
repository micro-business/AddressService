package service_test

import (
	"testing"

	"github.com/microbusinesses/AddressService/data/service"
	"github.com/microbusinesses/Micro-Businesses-Core/system"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Delete method input parameters", func() {
	var (
		addressDataService *service.AddressDataService
		tenantId           system.UUID
		applicationId      system.UUID
		addressId          system.UUID
	)

	BeforeEach(func() {
		addressDataService = &service.AddressDataService{}
		tenantId, _ = system.RandomUUID()
		applicationId, _ = system.RandomUUID()
		addressId, _ = system.RandomUUID()
	})

	Context("when empty tenant unique identifier provided", func() {
		It("should panic", func() {
			Ω(func() { addressDataService.Delete(system.EmptyUUID, applicationId, addressId) }).Should(Panic())
		})
	})

	Context("when empty application unique identifier provided", func() {
		It("should panic", func() {
			Ω(func() { addressDataService.Delete(tenantId, system.EmptyUUID, addressId) }).Should(Panic())
		})
	})

	Context("when empty address unique identifier provided", func() {
		It("should panic", func() {
			Ω(func() { addressDataService.Delete(tenantId, applicationId, system.EmptyUUID) }).Should(Panic())
		})
	})
})

func TestDelete(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Delete method input parameters")
}
