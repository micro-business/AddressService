package service_test

import (
	"testing"

	"github.com/gocql/gocql"
	"github.com/microbusinesses/AddressService/data/contract"
	"github.com/microbusinesses/AddressService/data/service"
	"github.com/microbusinesses/Micro-Businesses-Core/system"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Update method input parameters and dependency test", func() {
	var (
		addressDataService *service.AddressDataService
		tenantID           system.UUID
		applicationID      system.UUID
		addressID          system.UUID
		validAddress       contract.Address
		emptyAddress       contract.Address
	)

	BeforeEach(func() {
		addressDataService = &service.AddressDataService{ClusterConfig: &gocql.ClusterConfig{}}
		tenantID, _ = system.RandomUUID()
		applicationID, _ = system.RandomUUID()
		addressID, _ = system.RandomUUID()
		validAddress = contract.Address{AddressDetails: map[string]string{"City": "Christchurch"}}
		emptyAddress = contract.Address{}
	})

	Context("when cluster configuration not provided", func() {
		It("should panic", func() {
			addressDataService.ClusterConfig = nil

			Ω(func() { addressDataService.Update(tenantID, applicationID, addressID, validAddress) }).Should(Panic())
		})
	})

	DescribeTable("Input Parameters",
		func(tenantID, applicationID, addressID system.UUID, address contract.Address) {
			Ω(func() { addressDataService.Update(tenantID, applicationID, addressID, address) }).Should(Panic())
		},
		Entry("should panic when empty tenant unique identifier provided", system.EmptyUUID, applicationID, addressID, validAddress),
		Entry("should panic when empty application unique identifier provided", tenantID, system.EmptyUUID, addressID, validAddress),
		Entry("should panic when empty address unique identifier provided", tenantID, applicationID, system.EmptyUUID, validAddress),
		Entry("should panic when address without address key provided", tenantID, applicationID, addressID, emptyAddress))
})

func TestUpdate(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Update method input parameters and dependency test")
}
