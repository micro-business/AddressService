package endpoint

import (
	"errors"
	"strings"

	"github.com/go-kit/kit/endpoint"
	"github.com/graphql-go/graphql"
	"github.com/microbusinesses/AddressService/business/contract"
	"github.com/microbusinesses/AddressService/business/domain"
	"github.com/microbusinesses/AddressService/endpoint/message"
	"github.com/microbusinesses/Micro-Businesses-Core/common/query"
	"github.com/microbusinesses/Micro-Businesses-Core/system"
	"golang.org/x/net/context"
)

type address struct {
	Id      string                     `json:"id"`
	Details []query.StringKeyValuePair `json:"details"`
}

var addressType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Address",
		Fields: graphql.Fields{
			"id": &graphql.Field{Type: graphql.ID},
			"details": &graphql.Field{
				Type: graphql.NewList(query.StringKeyValuePairType),
				Args: graphql.FieldConfigArgument{
					"keys": &graphql.ArgumentConfig{
						Type: graphql.NewList(graphql.String),
					},
				},
				Resolve: func(resolveParams graphql.ResolveParams) (interface{}, error) {
					executionContext := resolveParams.Context.Value("ExecutionContext").(executionContext)

					keys, keysArgProvided := resolveParams.Args["keys"].([]interface{})
					currentAddress, _ := resolveParams.Source.(address)
					addressId, _ := system.ParseUUID(currentAddress.Id)

					var address domain.Address
					var err error

					if keysArgProvided {

						detailsKeys := make([]string, 0, len(keys))

						for _, key := range keys {
							castedKey, _ := key.(string)
							if len(strings.TrimSpace(castedKey)) != 0 {
								detailsKeys = append(detailsKeys, castedKey)
							}
						}

						if address, err = executionContext.addressService.Read(
							executionContext.tenantId,
							executionContext.applicationId,
							addressId,
							detailsKeys); err != nil {
							return nil, err
						}

					} else {

						if address, err = executionContext.addressService.ReadAll(
							executionContext.tenantId,
							executionContext.applicationId,
							addressId); err != nil {
							return nil, err
						}
					}

					if len(address.AddressDetails) == 0 {
						return nil, errors.New("Provided AddressId not found!!!")
					}

					currentAddress.Details = make([]query.StringKeyValuePair, 0, len(address.AddressDetails))

					for key, value := range address.AddressDetails {
						currentAddress.Details = append(currentAddress.Details, query.StringKeyValuePair{Key: key, Value: value})
					}

					return currentAddress, nil
				},
			},
		},
	},
)

var rootQueryType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "RootQuery",
		Fields: graphql.Fields{
			"address": &graphql.Field{
				Type: addressType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(resolveParams graphql.ResolveParams) (interface{}, error) {
					id, idArgProvided := resolveParams.Args["id"].(string)

					if idArgProvided {
						if addressId, err := system.ParseUUID(id); err != nil {
							return nil, err
						} else {
							return address{Id: addressId.String()}, nil
						}
					}

					return nil, errors.New("Address Id must be provided!!!")
				}},
		},
	},
)

var rootMutationType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "RootMutation",
	},
)

var addressSchema, _ = graphql.NewSchema(graphql.SchemaConfig{Query: rootQueryType, Mutation: rootMutationType})

type executionContext struct {
	addressService contract.AddressService
	tenantId       system.UUID
	applicationId  system.UUID
}

func createApiEndpoint(addressService contract.AddressService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		tenantId, _ := system.ParseUUID("02365c33-43d5-4bf8-b220-25563443960b")
		applicationId, _ := system.ParseUUID("02365c33-43d5-4bf8-b220-25563443960c")

		result := executeQuery(request.(string), addressService, tenantId, applicationId)

		if result.HasErrors() {
			errorMessages := []string{}

			for _, err := range result.Errors {
				errorMessages = append(errorMessages, err.Error())
			}

			return nil, errors.New(strings.Join(errorMessages, "\n"))
		}

		return result, nil
	}
}

func executeQuery(query string, addressService contract.AddressService, tenantId system.UUID, applicationId system.UUID) *graphql.Result {

	return graphql.Do(
		graphql.Params{
			Schema:        addressSchema,
			RequestString: query,
			Context:       context.WithValue(context.Background(), "ExecutionContext", executionContext{addressService, tenantId, applicationId}),
		})
}

func createCreateAddressEndpoint(service contract.AddressService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(message.CreateAddressRequest)
		tenantId, _ := system.ParseUUID("02365c33-43d5-4bf8-b220-25563443960b")
		applicationId, _ := system.ParseUUID("02365c33-43d5-4bf8-b220-25563443960c")

		address := domain.Address{AddressDetails: req.AddressDetails}

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

		address := domain.Address{AddressDetails: req.AddressDetails}

		err := service.Update(tenantId, applicationId, req.AddressId, address)

		if err != nil {
			return message.UpdateAddressResponse{err.Error()}, err
		} else {
			return message.UpdateAddressResponse{""}, nil
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
