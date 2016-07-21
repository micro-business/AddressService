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
		tenantId           system.UUID
		applicationId      system.UUID
		addressId          system.UUID
	)

	BeforeEach(func() {
		addressDataService = &service.AddressDataService{ClusterConfig: &gocql.ClusterConfig{}}

		addressDataService = &service.AddressDataService{}
		tenantId, _ = system.RandomUUID()
		applicationId, _ = system.RandomUUID()
		addressId, _ = system.RandomUUID()
	})

	Context("when cluster configuration not provided", func() {
		It("should panic", func() {
			addressDataService.ClusterConfig = nil

			Ω(func() { addressDataService.ReadAll(tenantId, applicationId, addressId) }).Should(Panic())
		})
	})

	DescribeTable("Input Parameters",
		func(tenantId, applicationId, addressId system.UUID) {
			Ω(func() { addressDataService.ReadAll(tenantId, applicationId, addressId) }).Should(Panic())
		},
		Entry("should panic when empty tenant unique identifier provided", system.EmptyUUID, applicationId, addressId),
		Entry("should panic when empty application unique identifier provided", tenantId, system.EmptyUUID, addressId),
		Entry("should panic when empty address unique identifier provided", tenantId, applicationId, system.EmptyUUID))
})

func TestReadAll(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ReadAll method input parameters and dependency test")
}
