package service_test

import (
	"testing"

	"github.com/gocql/gocql"
	"github.com/microbusinesses/AddressService/data/service"
	"github.com/microbusinesses/Micro-Businesses-Core/system"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Read method input parameters and dependency test", func() {
	var (
		addressDataService *service.AddressDataService
		tenantID           system.UUID
		applicationID      system.UUID
		addressID          system.UUID
		validKeys          []string
	)

	BeforeEach(func() {
		addressDataService = &service.AddressDataService{ClusterConfig: &gocql.ClusterConfig{}}

		addressDataService = &service.AddressDataService{}
		tenantID, _ = system.RandomUUID()
		applicationID, _ = system.RandomUUID()
		addressID, _ = system.RandomUUID()
		validKeys = make([]string, 1)
		validKeys[0] = "Line1"
	})

	Context("when cluster configuration not provided", func() {
		It("should panic", func() {
			addressDataService.ClusterConfig = nil

			Ω(func() { addressDataService.Read(tenantID, applicationID, addressID, validKeys) }).Should(Panic())
		})
	})
})

func TestRead(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Read method input parameters and dependency test")
}
