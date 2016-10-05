// Defines all reply messages used in address service
package message

import "github.com/microbusinesses/Micro-Businesses-Core/system"

// CreateAddressResponse defines the message that contains the result of creating a new address
type CreateAddressResponse struct {
	AddressId system.UUID `json:AddressId`
	Error     string      `json:"error,omitempty"`
}

// UpdateAddressResponse defines the message that contains the result of updating an existing address
type UpdateAddressResponse struct {
	Error string `json:"error:omitempty"`
}

// DeleteAddressResponse defines the message that contains the result of deleting an existing address
type DeleteAddressResponse struct {
	Error string `json:"error,omitempty"`
}
