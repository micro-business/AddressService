// Defines all request message used in address service
package message

// CreateAddressRequest defines the message that is used to create new address
type CreateAddressRequest struct {
	AddressKeysValues map[string]string `json:AddressKeysValues`
}
