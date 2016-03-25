package service_test

import (
	"testing"

	"github.com/microbusinesses/AddressService/data/service"
	"github.com/microbusinesses/AddressService/data/shared"
	. "github.com/microbusinesses/Micro-Businesses-Core/system"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Update method input parameters", func() {
	var (
		addressDataService service.AddressDataService
		tenantId           UUID
		applicationId      UUID
		addressId          UUID
		validAddress       shared.Address
		emptyAddress       shared.Address
	)

	BeforeEach(func() {
		addressDataService = service.AddressDataService{}
		tenantId, _ = RandomUUID()
		applicationId, _ = RandomUUID()
		addressId, _ = RandomUUID()
		validAddress = shared.Address{AddressParts: map[string]string{"FirstName": "Morteza"}}
		emptyAddress = shared.Address{}
	})

	Context("when UUID generator service not provided", func() {
		It("should panic", func() {
			addressDataService.UUIDGeneratorService = nil

			Ω(func() { addressDataService.Update(tenantId, applicationId, addressId, validAddress) }).Should(Panic())
		})
	})

	Context("when empty tenant unique identifier provided", func() {
		It("should panic", func() {
			Ω(func() { addressDataService.Update(EmptyUUID, applicationId, addressId, validAddress) }).Should(Panic())
		})
	})

	Context("when empty application unique identifier provided", func() {
		It("should panic", func() {
			Ω(func() { addressDataService.Update(tenantId, EmptyUUID, addressId, validAddress) }).Should(Panic())
		})
	})

	Context("when empty address unique identifier provided", func() {
		It("should panic", func() {
			Ω(func() { addressDataService.Update(tenantId, applicationId, EmptyUUID, validAddress) }).Should(Panic())
		})
	})

	Context("when address without address parts provided", func() {
		It("should panic", func() {
			Ω(func() { addressDataService.Update(tenantId, applicationId, addressId, emptyAddress) }).Should(Panic())
		})
	})
})

func TestUpdate(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Update method input parameters")
}
