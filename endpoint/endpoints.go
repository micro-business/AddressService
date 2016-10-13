package endpoint

import (
	"errors"
	"strings"

	"github.com/go-kit/kit/endpoint"
	"github.com/graphql-go/graphql"
	"github.com/microbusinesses/AddressService/business/contract"
	"github.com/microbusinesses/AddressService/business/domain"
	"github.com/microbusinesses/Micro-Businesses-Core/common/query"
	"github.com/microbusinesses/Micro-Businesses-Core/system"
	"golang.org/x/net/context"
)

const (
	buildingNumber = "BuildingNumber"
	streetNumber   = "StreetNumber"
	line1          = "Line1"
	line2          = "Line2"
	line3          = "Line3"
	line4          = "Line4"
	line5          = "Line5"
	suburb         = "Suburb"
	city           = "City"
	state          = "State"
	postcode       = "Postcode"
	country        = "Country"
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
			"id":           &graphql.Field{Type: graphql.ID},
			buildingNumber: &graphql.Field{Type: graphql.String},
			streetNumber:   &graphql.Field{Type: graphql.String},
			line1:          &graphql.Field{Type: graphql.String},
			line2:          &graphql.Field{Type: graphql.String},
			line3:          &graphql.Field{Type: graphql.String},
			line4:          &graphql.Field{Type: graphql.String},
			line5:          &graphql.Field{Type: graphql.String},
			suburb:         &graphql.Field{Type: graphql.String},
			city:           &graphql.Field{Type: graphql.String},
			state:          &graphql.Field{Type: graphql.String},
			postcode:       &graphql.Field{Type: graphql.String},
			country:        &graphql.Field{Type: graphql.String},
		},
	},
)

var rootQueryType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "RootQuery",
		Fields: graphql.Fields{
			"address": &graphql.Field{
				Type:        addressType,
				Description: "Returns an existing address",
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
							keys := query.GetSelectedFields([]string{"address"}, resolveParams)

							var returnedAddress domain.Address

							if returnedAddress, err = executionContext.addressService.Read(
								executionContext.tenantId,
								executionContext.applicationId,
								addressId,
								keys); err != nil {
								return nil, err
							}

							if len(returnedAddress.AddressDetails) == 0 {
								return nil, errors.New("Provided AddressId not found!!!")
							}

							address := address{
								Id:             addressId.String(),
								BuildingNumber: returnedAddress.AddressDetails[buildingNumber],
								StreetNumber:   returnedAddress.AddressDetails[streetNumber],
								Line1:          returnedAddress.AddressDetails[line1],
								Line2:          returnedAddress.AddressDetails[line2],
								Line3:          returnedAddress.AddressDetails[line3],
								Line4:          returnedAddress.AddressDetails[line4],
								Line5:          returnedAddress.AddressDetails[line5],
								Suburb:         returnedAddress.AddressDetails[suburb],
								City:           returnedAddress.AddressDetails[city],
								State:          returnedAddress.AddressDetails[state],
								Postcode:       returnedAddress.AddressDetails[postcode],
								Country:        returnedAddress.AddressDetails[country],
							}

							return address, nil
						}
					}

					return nil, errors.New("Address Id must be provided.")
				}},
		},
	},
)

var rootMutationType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "RootMutation",
		Fields: graphql.Fields{
			"create": &graphql.Field{
				Type:        graphql.ID,
				Description: "Creates new address",
				Args: graphql.FieldConfigArgument{
					buildingNumber: &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					streetNumber: &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					line1: &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					line2: &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					line3: &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					line4: &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					line5: &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					suburb: &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					city: &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					state: &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					postcode: &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					country: &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(resolveParams graphql.ResolveParams) (interface{}, error) {

					var address domain.Address
					var err error

					if address, err = resolveAddressDetails(resolveParams); err != nil {
						return nil, err
					}

					executionContext := resolveParams.Context.Value("ExecutionContext").(executionContext)

					if addressId, err := executionContext.addressService.Create(
						executionContext.tenantId,
						executionContext.applicationId,
						address); err != nil {
						return nil, err
					} else {
						return addressId.String(), nil
					}
				},
			},

			"update": &graphql.Field{
				Type:        graphql.ID,
				Description: "Update existing address",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.ID),
					},
					buildingNumber: &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					streetNumber: &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					line1: &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					line2: &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					line3: &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					line4: &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					line5: &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					suburb: &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					city: &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					state: &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					postcode: &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					country: &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(resolveParams graphql.ResolveParams) (interface{}, error) {
					id, _ := resolveParams.Args["id"].(string)

					var addressId system.UUID
					var err error

					if addressId, err = system.ParseUUID(id); err != nil {
						return nil, err
					}

					var address domain.Address

					if address, err = resolveAddressDetails(resolveParams); err != nil {
						return nil, err
					}

					executionContext := resolveParams.Context.Value("ExecutionContext").(executionContext)

					if err := executionContext.addressService.Update(
						executionContext.tenantId,
						executionContext.applicationId,
						addressId,
						address); err != nil {
						return nil, err
					} else {
						return addressId.String(), nil
					}
				},
			},

			"delete": &graphql.Field{
				Type:        graphql.ID,
				Description: "Delete existing address",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.ID),
					},
				},
				Resolve: func(resolveParams graphql.ResolveParams) (interface{}, error) {
					id, _ := resolveParams.Args["id"].(string)

					var addressId system.UUID
					var err error

					if addressId, err = system.ParseUUID(id); err != nil {
						return nil, err
					}

					executionContext := resolveParams.Context.Value("ExecutionContext").(executionContext)

					if err := executionContext.addressService.Delete(
						executionContext.tenantId,
						executionContext.applicationId,
						addressId); err != nil {
						return nil, err
					} else {
						return addressId.String(), nil
					}
				},
			},
		},
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

func resolveAddressDetails(resolveParams graphql.ResolveParams) (domain.Address, error) {

	addressDetails := []string{
		buildingNumber,
		streetNumber,
		line1,
		line2,
		line3,
		line4,
		line5,
		suburb,
		city,
		state,
		postcode,
		country}

	address := domain.Address{AddressDetails: make(map[string]string)}

	for _, key := range addressDetails {
		keyArg, KeyArgProvided := resolveParams.Args[key].(string)

		if KeyArgProvided {
			address.AddressDetails[key] = keyArg
		}

	}

	if len(address.AddressDetails) == 0 {
		return domain.Address{}, errors.New("At least one address part must be provided.")
	}

	return address, nil
}
