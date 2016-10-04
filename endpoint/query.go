package endpoint

import (
	"github.com/graphql-go/graphql"
	"github.com/microbusinesses/Micro-Businesses-Core/common/query"
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
				// Resolve: func(resolveParams graphql.ResolveParams) (interface{}, error) {
				// 	return nil, nil
				// },
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

						address := address{Id: id, Details: make([]query.StringKeyValuePair, 2)}

						address.Details[0].Key = "Key 1"
						address.Details[0].Value = "Value 1"

						address.Details[1].Key = "Key 2"
						address.Details[1].Value = "Value 2"

						return address, nil
					}

					return nil, nil
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
