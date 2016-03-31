package endpoint

import (
	"github.com/go-kit/kit/endpoint"
	"github.com/microbusinesses/AddressService/business/contract"
	"github.com/microbusinesses/AddressService/business/domain"
	"github.com/microbusinesses/AddressService/endpoint/message"
	"github.com/microbusinesses/Micro-Businesses-Core/system"
	"golang.org/x/net/context"
)

func createCreateAddressEndpoint(service contract.AddressService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(message.CreateAddressRequest)
		tenantId, _ := system.ParseUUID("02365c33-43d5-4bf8-b220-25563443960b")
		applicationId, _ := system.ParseUUID("02365c33-43d5-4bf8-b220-25563443960c")

		address := domain.Address{AddressKeysValues: req.AddressKeysValues}

		addressId, err := service.Create(tenantId, applicationId, address)

		if err != nil {
			return message.CreateAddressResponse{system.EmptyUUID, err.Error()}, err
		} else {
			return message.CreateAddressResponse{addressId, ""}, nil
		}
	}
}
