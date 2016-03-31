// Defines all reply messages used in address service
package message

import "github.com/microbusinesses/Micro-Businesses-Core/system"

// CreateAddressResponse defines the message that is used to create new address
type CreateAddressResponse struct {
	AddressId system.UUID `json:AddressId`
	Error     string      `json:"error:omitempty"`
}
