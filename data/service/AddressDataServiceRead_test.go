package service_test

import (
	"testing"

	"github.com/microbusinesses/AddressService/data/service"
	. "github.com/microbusinesses/Micro-Businesses-Core/system"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Read method input parameters", func() {
	var (
		addressDataService service.AddressDataService
		tenantId           UUID
		applicationId      UUID
		addressId          UUID
	)

	BeforeEach(func() {
		addressDataService = service.AddressDataService{}
		tenantId, _ = RandomUUID()
		applicationId, _ = RandomUUID()
		addressId, _ = RandomUUID()
	})

	Context("when UUID generator service not provided", func() {
		It("should panic", func() {
			addressDataService.UUIDGeneratorService = nil

			Ω(func() { addressDataService.Read(tenantId, applicationId, addressId) }).Should(Panic())
		})
	})

	Context("when empty tenant unique identifier provided", func() {
		It("should panic", func() {
			Ω(func() { addressDataService.Read(EmptyUUID, applicationId, addressId) }).Should(Panic())
		})
	})

	Context("when empty tenant unique identifier provided", func() {
		It("should panic", func() {
			Ω(func() { addressDataService.Read(tenantId, EmptyUUID, addressId) }).Should(Panic())
		})
	})

	Context("when empty address unique identifier provided", func() {
		It("should panic", func() {
			Ω(func() { addressDataService.Read(tenantId, applicationId, EmptyUUID) }).Should(Panic())
		})
	})
})

func TestRead(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Read method input parameters")
}
