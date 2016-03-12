package service_test

import (
	"testing"

	"github.com/microbusinesses/AddressService/service"
	"github.com/microbusinesses/Micro-Businesses-Core/system"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Add method input parameters", func() {
	var (
		addressService service.AddressService
		tenantId       system.UUID
	)

	BeforeEach(func() {
		addressService = service.AddressService{}
		tenantId, _ = system.RandomUUID()
	})

	Context("when empty tenant unique identifier provided", func() {
		It("should panic", func() {
			defer func() {
				if r := recover(); r == nil {
					Fail("Should have paniced")
				}
			}()

			addressService.Add(system.EmptyUUID, nil)
		})
	})
})

func TestAdd(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Add method input parameters")
}
