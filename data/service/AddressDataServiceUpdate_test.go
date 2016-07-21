package service_test

import (
	"testing"

	"github.com/gocql/gocql"
	"github.com/microbusinesses/AddressService/data/service"
	"github.com/microbusinesses/AddressService/data/shared"
	"github.com/microbusinesses/Micro-Businesses-Core/system"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Update method input parameters and dependency test", func() {
	var (
		addressDataService *service.AddressDataService
		tenantId           system.UUID
		applicationId      system.UUID
		addressId          system.UUID
		validAddress       shared.Address
		emptyAddress       shared.Address
	)

	BeforeEach(func() {
		addressDataService = &service.AddressDataService{ClusterConfig: &gocql.ClusterConfig{}}
		tenantId, _ = system.RandomUUID()
		applicationId, _ = system.RandomUUID()
		addressId, _ = system.RandomUUID()
		validAddress = shared.Address{AddressDetails: map[string]string{"City": "Christchurch"}}
		emptyAddress = shared.Address{}
	})

	Context("when cluster configuration not provided", func() {
		It("should panic", func() {
			addressDataService.ClusterConfig = nil

			Ω(func() { addressDataService.Update(tenantId, applicationId, addressId, validAddress) }).Should(Panic())
		})
	})

	DescribeTable("Input Parameters",
		func(tenantId, applicationId, addressId system.UUID, address shared.Address) {
			Ω(func() { addressDataService.Update(tenantId, applicationId, addressId, address) }).Should(Panic())
		},
		Entry("should panic when empty tenant unique identifier provided", system.EmptyUUID, applicationId, addressId, validAddress),
		Entry("should panic when empty application unique identifier provided", tenantId, system.EmptyUUID, addressId, validAddress),
		Entry("should panic when empty address unique identifier provided", tenantId, applicationId, system.EmptyUUID, validAddress),
		Entry("should panic when address without address key provided", tenantId, applicationId, addressId, emptyAddress))
})

func TestUpdate(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Update method input parameters and dependency test")
}
