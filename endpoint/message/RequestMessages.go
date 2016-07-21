// Defines all request message used in address service
package message

import (
	"github.com/microbusinesses/Micro-Businesses-Core/common/query"
	"github.com/microbusinesses/Micro-Businesses-Core/system"
)

// ApiRequest defines the message that will be used to talk to the API interface
type ApiRequest struct {
	RequestQuery query.RequestQuery `json:RequestQuery`
}

// CreateAddressRequest defines the message that is used to create a new address
type CreateAddressRequest struct {
	AddressDetails map[string]string `json:AddressDetails`
}

// UpdateAddressRequest defines the message that is used to update an existing address
type UpdateAddressRequest struct {
	AddressId      system.UUID       `json:AddressId`
	AddressDetails map[string]string `json:AddressDetails`
}

// ReadAllAddressRequest defines the message that is used to read an existing address
type ReadAllAddressRequest struct {
	AddressId system.UUID `json:AddressId`
}

// DeleteAddressRequest defines the message that is used to delete an existing address
type DeleteAddressRequest struct {
	AddressId system.UUID `json:AddressId`
}
