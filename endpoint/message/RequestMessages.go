// Defines all request message used in address service
package message

import "github.com/microbusinesses/Micro-Businesses-Core/system"

// CreateAddressRequest defines the message that is used to create a new address
type CreateAddressRequest struct {
	AddressKeysValues map[string]string `json:AddressKeysValues`
}

// UpdateAddressRequest defines the message that is used to update an existing address
type UpdateAddressRequest struct {
	AddressId         system.UUID       `json:AddressId`
	AddressKeysValues map[string]string `json:AddressKeysValues`
}

// UpdateAddressRequest defines the message that is used to read an existing address
type ReadAddressRequest struct {
	AddressId system.UUID `json:AddressId`
}

// DeleteAddressRequest defines the message that is used to delete an existing address
type DeleteAddressRequest struct {
	AddressId system.UUID `json:AddressId`
}
