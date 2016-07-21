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

		address := domain.Address{AddressDetails: req.AddressKeysValues}

		addressId, err := service.Create(tenantId, applicationId, address)

		if err != nil {
			return message.CreateAddressResponse{system.EmptyUUID, err.Error()}, err
		} else {
			return message.CreateAddressResponse{addressId, ""}, nil
		}
	}
}

func createUpdateAddressEndpoint(service contract.AddressService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(message.UpdateAddressRequest)
		tenantId, _ := system.ParseUUID("02365c33-43d5-4bf8-b220-25563443960b")
		applicationId, _ := system.ParseUUID("02365c33-43d5-4bf8-b220-25563443960c")

		address := domain.Address{AddressDetails: req.AddressKeysValues}

		err := service.Update(tenantId, applicationId, req.AddressId, address)

		if err != nil {
			return message.UpdateAddressResponse{err.Error()}, err
		} else {
			return message.UpdateAddressResponse{""}, nil
		}
	}
}

func createReadAddressEndpoint(service contract.AddressService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(message.ReadAddressRequest)
		tenantId, _ := system.ParseUUID("02365c33-43d5-4bf8-b220-25563443960b")
		applicationId, _ := system.ParseUUID("02365c33-43d5-4bf8-b220-25563443960c")

		address, err := service.ReadAll(tenantId, applicationId, req.AddressId)

		if err != nil {
			return message.ReadAddressResponse{make(map[string]string), err.Error()}, err
		} else {
			return message.ReadAddressResponse{address.AddressDetails, ""}, nil
		}
	}
}

func createDeleteAddressEndpoint(service contract.AddressService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(message.DeleteAddressRequest)
		tenantId, _ := system.ParseUUID("02365c33-43d5-4bf8-b220-25563443960b")
		applicationId, _ := system.ParseUUID("02365c33-43d5-4bf8-b220-25563443960c")

		err := service.Delete(tenantId, applicationId, req.AddressId)

		if err != nil {
			return message.DeleteAddressResponse{err.Error()}, err
		} else {
			return message.DeleteAddressResponse{""}, nil
		}
	}
}
