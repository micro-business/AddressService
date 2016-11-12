// Package contract defines the address service contract.
package contract

import (
	"github.com/micro-business/AddressService/business/domain"
	"github.com/micro-business/Micro-Business-Core/system"
)

// AddressService contract, it can add new address and update/retrieve/remove an existing address.
type AddressService interface {
	// Create creates a new address.
	// tenantID: Mandatory. The unique identifier of the tenant owning the address.
	// applicationID: Mandatory. The unique identifier of the tenant's application will be owning the address.
	// address: Mandatory. The reference to the new address information.
	// Returns either the unique identifier of the new address or error if something goes wrong.
	Create(tenantID, applicationID system.UUID, address domain.Address) (system.UUID, error)

	// Update updates an existing address.
	// tenantID: Mandatory. The unique identifier of the tenant owning the address.
	// applicationID: Mandatory. The unique identifier of the tenant's application will be owning the address.
	// addressID: Mandatory. The unique identifier of the existing address.
	// address: Mandatory. The reeference to the updated address information.
	// Returns error if something goes wrong.
	Update(tenantID, applicationID, addressID system.UUID, address domain.Address) error

	// Read retrieves an existing address information and returns only the detail which the keys provided by the keys.
	// tenantID: Mandatory. The unique identifier of the tenant owning the address.
	// applicationID: Mandatory. The unique identifier of the tenant's application will be owning the address.
	// addressID: Mandatory. The unique identifier of the existing address.
	// keys: Mandatory. The interested address details keys to return.
	// Returns either the address information or error if something goes wrong.
	Read(tenantID, applicationID, addressID system.UUID, keys []string) (domain.Address, error)

	// ReadAll retrieves an existing address information and returns all the detail of it.
	// tenantID: Mandatory. The unique identifier of the tenant owning the address.
	// applicationID: Mandatory. The unique identifier of the tenant's application will be owning the address.
	// addressID: Mandatory. The unique identifier of the existing address.
	// Returns either the address information or error if something goes wrong.
	ReadAll(tenantID, applicationID, addressID system.UUID) (domain.Address, error)

	// Delete deletes an existing address information.
	// tenantID: Mandatory. The unique identifier of the tenant owning the address.
	// applicationID: Mandatory. The unique identifier of the tenant's application will be owning the address.
	// addressID: Mandatory. The unique identifier of the existing address to remove.
	// Returns error if something goes wrong.
	Delete(tenantID, applicationID, addressID system.UUID) error
}
