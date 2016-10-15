package service_test

import (
	"testing"

	"github.com/gocql/gocql"
	"github.com/microbusinesses/AddressService/data/service"
	"github.com/microbusinesses/Micro-Businesses-Core/system"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Delete method input parameters and dependency test", func() {
	var (
		addressDataService *service.AddressDataService
		tenantID           system.UUID
		applicationID      system.UUID
		addressID          system.UUID
	)

	BeforeEach(func() {
		addressDataService = &service.AddressDataService{ClusterConfig: &gocql.ClusterConfig{}}
		tenantID, _ = system.RandomUUID()
		applicationID, _ = system.RandomUUID()
		addressID, _ = system.RandomUUID()
	})

	Context("when cluster configuration not provided", func() {
		It("should panic", func() {
			addressDataService.ClusterConfig = nil

			Ω(func() { addressDataService.Delete(tenantID, applicationID, addressID) }).Should(Panic())
		})
	})

	DescribeTable("Input Parameters",
		func(tenantID, applicationID, addressID system.UUID) {
			Ω(func() { addressDataService.Delete(tenantID, applicationID, addressID) }).Should(Panic())
		},
		Entry("should panic when empty tenant unique identifier provided", system.EmptyUUID, applicationID, addressID),
		Entry("should panic when empty application unique identifier provided", tenantID, system.EmptyUUID, addressID),
		Entry("should panic when empty address unique identifier provided", tenantID, applicationID, system.EmptyUUID))
})

func TestDelete(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Delete method input parameters and dependency test")
}
