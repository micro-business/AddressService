package main

import (
	"fmt"

	"github.com/gocql/gocql"
	"github.com/microbusinesses/AddressService/data/service"
	"github.com/microbusinesses/AddressService/data/shared"
)

func main() {
	uuidGeneratorService := service.UUIDGeneratorServiceImpl{}

	cluster := gocql.NewCluster("cassandra-test")
	cluster.Keyspace = "address"
	cluster.Consistency = gocql.Quorum

	addressDataService := service.AddressDataService{UUIDGeneratorService: &uuidGeneratorService, ClusterConfig: cluster}

	tenantId, _ := uuidGeneratorService.GenerateRandomUUID()
	applicationId, _ := uuidGeneratorService.GenerateRandomUUID()

	address := shared.Address{AddressKeysValues: map[string]string{"Address Line 1": "32 Hope st", "Address Line 2": "Shirley", "Town/City": "Christchurch", "Postcode": "8013", "Country": "New Zealand"}}

	if addressId, err := addressDataService.Create(tenantId, applicationId, address); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(addressId.String())
	}
}
