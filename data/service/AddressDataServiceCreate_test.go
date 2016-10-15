package service_test

import (
	"testing"

	"github.com/gocql/gocql"
	"github.com/golang/mock/gomock"
	"github.com/microbusinesses/AddressService/data/contract"
	"github.com/microbusinesses/AddressService/data/service"
	"github.com/microbusinesses/Micro-Businesses-Core/system"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Create method input parameters and dependency test", func() {
	var (
		mockCtrl                 *gomock.Controller
		addressDataService       *service.AddressDataService
		mockUUIDGeneratorService *MockUUIDGeneratorService
		tenantID                 system.UUID
		applicationID            system.UUID
		validAddress             contract.Address
		emptyAddress             contract.Address
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockUUIDGeneratorService = NewMockUUIDGeneratorService(mockCtrl)

		addressDataService = &service.AddressDataService{UUIDGeneratorService: mockUUIDGeneratorService, ClusterConfig: &gocql.ClusterConfig{}}

		tenantID, _ = system.RandomUUID()
		applicationID, _ = system.RandomUUID()
		validAddress = contract.Address{AddressDetails: map[string]string{"City": "Christchurch"}}
		emptyAddress = contract.Address{}
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("when UUID generator service not provided", func() {
		It("should panic", func() {
			addressDataService.UUIDGeneratorService = nil

			Ω(func() { addressDataService.Create(tenantID, applicationID, validAddress) }).Should(Panic())
		})
	})

	Context("when cluster configuration not provided", func() {
		It("should panic", func() {
			addressDataService.ClusterConfig = nil

			Ω(func() { addressDataService.Create(tenantID, applicationID, validAddress) }).Should(Panic())
		})
	})

	DescribeTable("Input Parameters",
		func(tenantID, applicationID system.UUID, address contract.Address) {
			Ω(func() { addressDataService.Create(tenantID, applicationID, address) }).Should(Panic())
		},
		Entry("should panic when empty tenant unique identifier provided", system.EmptyUUID, applicationID, validAddress),
		Entry("should panic when empty application unique identifier provided", tenantID, system.EmptyUUID, validAddress),
		Entry("should panic when address without address key provided", tenantID, applicationID, emptyAddress))
})

func TestCreate(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Create method input parameters and dependency test")
}
