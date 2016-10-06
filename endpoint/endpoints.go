package endpoint

import (
	"errors"
	"strings"

	"github.com/go-kit/kit/endpoint"
	"github.com/graphql-go/graphql"
	ast "github.com/graphql-go/graphql/language/ast"
	"github.com/microbusinesses/AddressService/business/contract"
	"github.com/microbusinesses/AddressService/business/domain"
	"github.com/microbusinesses/AddressService/endpoint/message"
	"github.com/microbusinesses/Micro-Businesses-Core/system"
	"golang.org/x/net/context"
)

type address struct {
	Id             string `json:"Id"`
	BuildingNumber string `json:"BuildingNumber"`
	StreetNumber   string `json:"StreetNumber"`
	Line1          string `json:"Line1"`
	Line2          string `json:"Line2"`
	Line3          string `json:"Line3"`
	Line4          string `json:"Line4"`
	Line5          string `json:"Line5"`
	Suburb         string `json:"Suburb"`
	City           string `json:"City"`
	State          string `json:"State"`
	Postcode       string `json:"Postcode"`
	Country        string `json:"Country"`
}

var addressType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Address",
		Fields: graphql.Fields{
			"id":             &graphql.Field{Type: graphql.ID},
			"BuildingNumber": &graphql.Field{Type: graphql.String},
			"StreetNumber":   &graphql.Field{Type: graphql.String},
			"Line1":          &graphql.Field{Type: graphql.String},
			"Line2":          &graphql.Field{Type: graphql.String},
			"Line3":          &graphql.Field{Type: graphql.String},
			"Line4":          &graphql.Field{Type: graphql.String},
			"Line5":          &graphql.Field{Type: graphql.String},
			"Suburb":         &graphql.Field{Type: graphql.String},
			"City":           &graphql.Field{Type: graphql.String},
			"State":          &graphql.Field{Type: graphql.String},
			"Postcode":       &graphql.Field{Type: graphql.String},
			"Country":        &graphql.Field{Type: graphql.String},
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
					executionContext := resolveParams.Context.Value("ExecutionContext").(executionContext)
					id, idArgProvided := resolveParams.Args["id"].(string)

					if idArgProvided {
						if addressId, err := system.ParseUUID(id); err != nil {
							return nil, err
						} else {
							keys := getSelectedFields([]string{"address"}, resolveParams)

							var returnedAddress domain.Address
							var err error

							if len(keys) == 0 {
								if returnedAddress, err = executionContext.addressService.ReadAll(
									executionContext.tenantId,
									executionContext.applicationId,
									addressId); err != nil {
									return nil, err
								}

							} else {
								if returnedAddress, err = executionContext.addressService.Read(
									executionContext.tenantId,
									executionContext.applicationId,
									addressId,
									keys); err != nil {
									return nil, err
								}

							}

							if len(returnedAddress.AddressDetails) == 0 {
								return nil, errors.New("Provided AddressId not found!!!")
							}

							address := address{
								Id:             addressId.String(),
								BuildingNumber: returnedAddress.AddressDetails["BuildingNumber"],
								StreetNumber:   returnedAddress.AddressDetails["StreetNumber"],
								Line1:          returnedAddress.AddressDetails["Line1"],
								Line2:          returnedAddress.AddressDetails["Line2"],
								Line3:          returnedAddress.AddressDetails["Line3"],
								Line4:          returnedAddress.AddressDetails["Line4"],
								Line5:          returnedAddress.AddressDetails["Line5"],
								Suburb:         returnedAddress.AddressDetails["Suburb"],
								City:           returnedAddress.AddressDetails["City"],
								State:          returnedAddress.AddressDetails["State"],
								Postcode:       returnedAddress.AddressDetails["Postcode"],
								Country:        returnedAddress.AddressDetails["Country"],
							}

							return address, nil
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

func getSelectedFields(selectionPath []string,
	resolveParams graphql.ResolveParams) []string {
	fields := resolveParams.Info.FieldASTs
	for _, propName := range selectionPath {
		found := false
		for _, field := range fields {
			if field.Name.Value == propName {
				selections := field.SelectionSet.Selections
				fields = make([]*ast.Field, 0)
				for _, selection := range selections {
					fields = append(fields, selection.(*ast.Field))
				}
				found = true
				break
			}
		}
		if !found {
			return []string{}
		}
	}
	var collect []string
	for _, field := range fields {
		collect = append(collect, field.Name.Value)
	}
	return collect
}
