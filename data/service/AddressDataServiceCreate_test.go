package service_test

import (
	"testing"

	"github.com/microbusinesses/AddressService/data/service"
	"github.com/microbusinesses/AddressService/data/shared"
	"github.com/microbusinesses/Micro-Businesses-Core/system"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Create method input parameters", func() {
	var (
		addressDataService *service.AddressDataService
		tenantId           system.UUID
		applicationId      system.UUID
		validAddress       shared.Address
		emptyAddress       shared.Address
	)

	BeforeEach(func() {
		addressDataService = &service.AddressDataService{}
		tenantId, _ = system.RandomUUID()
		applicationId, _ = system.RandomUUID()
		validAddress = shared.Address{AddressParts: map[string]string{"FirstName": "Morteza"}}
		emptyAddress = shared.Address{}
	})

	Context("when UUID generator service not provided", func() {
		It("should panic", func() {
			addressDataService.UUIDGeneratorService = nil

			Ω(func() { addressDataService.Create(tenantId, applicationId, validAddress) }).Should(Panic())
		})
	})

	Context("when empty tenant unique identifier provided", func() {
		It("should panic", func() {
			Ω(func() { addressDataService.Create(system.EmptyUUID, applicationId, validAddress) }).Should(Panic())
		})
	})

	Context("when empty application unique identifier provided", func() {
		It("should panic", func() {
			Ω(func() { addressDataService.Create(tenantId, system.EmptyUUID, validAddress) }).Should(Panic())
		})
	})

	Context("when address without address parts provided", func() {
		It("should panic", func() {
			Ω(func() { addressDataService.Create(tenantId, applicationId, emptyAddress) }).Should(Panic())
		})
	})
})

func TestCreate(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Create method input parameters")
}
