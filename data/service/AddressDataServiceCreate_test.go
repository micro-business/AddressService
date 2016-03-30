package service_test

import (
	"testing"

	"github.com/gocql/gocql"
	"github.com/golang/mock/gomock"
	"github.com/microbusinesses/AddressService/data/service"
	dataServiceMocks "github.com/microbusinesses/AddressService/data/service/mocks"
	"github.com/microbusinesses/AddressService/data/shared"
	"github.com/microbusinesses/Micro-Businesses-Core/system"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Create method input parameters and dependency test", func() {
	var (
		mockCtrl                 *gomock.Controller
		addressDataService       *service.AddressDataService
		mockUUIDGeneratorService *dataServiceMocks.MockUUIDGeneratorService
		tenantId                 system.UUID
		applicationId            system.UUID
		validAddress             shared.Address
		emptyAddress             shared.Address
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockUUIDGeneratorService = dataServiceMocks.NewMockUUIDGeneratorService(mockCtrl)

		addressDataService = &service.AddressDataService{UUIDGeneratorService: mockUUIDGeneratorService, ClusterConfig: &gocql.ClusterConfig{}}

		tenantId, _ = system.RandomUUID()
		applicationId, _ = system.RandomUUID()
		validAddress = shared.Address{AddressKeysValues: map[string]string{"City": "Christchurch"}}
		emptyAddress = shared.Address{}
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("when UUID generator service not provided", func() {
		It("should panic", func() {
			addressDataService.UUIDGeneratorService = nil

			Ω(func() { addressDataService.Create(tenantId, applicationId, validAddress) }).Should(Panic())
		})
	})

	Context("when cluster configuration not provided", func() {
		It("should panic", func() {
			addressDataService.ClusterConfig = nil

			Ω(func() { addressDataService.Create(tenantId, applicationId, validAddress) }).Should(Panic())
		})
	})

	DescribeTable("Input Parameters",
		func(tenantId, applicationId system.UUID, address shared.Address) {
			Ω(func() { addressDataService.Create(tenantId, applicationId, address) }).Should(Panic())
		},
		Entry("should panic when empty tenant unique identifier provided", system.EmptyUUID, applicationId, validAddress),
		Entry("should panic when empty application unique identifier provided", tenantId, system.EmptyUUID, validAddress),
		Entry("should panic when address without address key provided", tenantId, applicationId, emptyAddress))
})

func TestCreate(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Create method input parameters and dependency test")
}
