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

var _ = Describe("ReadAll method input parameters and dependency test", func() {
	var (
		addressDataService *service.AddressDataService
		tenantID           system.UUID
		applicationID      system.UUID
		addressID          system.UUID
	)

	BeforeEach(func() {
		addressDataService = &service.AddressDataService{ClusterConfig: &gocql.ClusterConfig{}}

		addressDataService = &service.AddressDataService{}
		tenantID, _ = system.RandomUUID()
		applicationID, _ = system.RandomUUID()
		addressID, _ = system.RandomUUID()
	})

	Context("when cluster configuration not provided", func() {
		It("should panic", func() {
			addressDataService.ClusterConfig = nil

			Ω(func() { addressDataService.ReadAll(tenantID, applicationID, addressID) }).Should(Panic())
		})
	})

	DescribeTable("Input Parameters",
		func(tenantID, applicationID, addressID system.UUID) {
			Ω(func() { addressDataService.ReadAll(tenantID, applicationID, addressID) }).Should(Panic())
		},
		Entry("should panic when empty tenant unique identifier provided", system.EmptyUUID, applicationID, addressID),
		Entry("should panic when empty application unique identifier provided", tenantID, system.EmptyUUID, addressID),
		Entry("should panic when empty address unique identifier provided", tenantID, applicationID, system.EmptyUUID))
})

func TestReadAll(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ReadAll method input parameters and dependency test")
}
