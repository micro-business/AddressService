package main

import (
	"fmt"

	"github.com/gocql/gocql"
	"github.com/microbusinesses/AddressService/data/service"
	"github.com/microbusinesses/AddressService/data/shared"
	"github.com/microbusinesses/Micro-Businesses-Core/system"
)

func main() {
	uuidGeneratorService := service.UUIDGeneratorServiceImpl{}

	cluster := gocql.NewCluster("cassandra-test")
	cluster.Keyspace = "address"
	cluster.Consistency = gocql.Quorum

	addressDataService := service.AddressDataService{UUIDGeneratorService: &uuidGeneratorService, ClusterConfig: cluster}

	tenantId, _ := system.ParseUUID("e7c5d067-63f8-499f-b6de-79288903ff77")
	applicationId, _ := system.ParseUUID("e7c5d067-63f8-499f-b6de-79288903ff78")

	address := shared.Address{
		AddressKeysValues: map[string]string{
			"Address Line 1": "32 Hope st",
			"Address Line 2": "Shirley",
			"Town/City":      "Christchurch",
			"Postcode":       "8013",
			"Country":        "New Zealand"}}

	if addressId, err := addressDataService.Create(tenantId, applicationId, address); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(addressId.String())
	}
}
