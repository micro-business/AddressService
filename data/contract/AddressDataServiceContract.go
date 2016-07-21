// Defines the address data service contract.
package contract

import (
	"github.com/microbusinesses/AddressService/data/shared"
	"github.com/microbusinesses/Micro-Businesses-Core/system"
)

// The address data service contract, it can add new address and update/retrieve/remove an existing address.
type AddressDataService interface {
	// Create creates a new address.
	// tenantId: Mandatory. The unique identifier of the tenant owning the address.
	// applicationId: Mandatory. The unique identifier of the tenant's application will be owning the address.
	// address: Mandatory. The reference to the new address information.
	// Returns either the unique identifier of the new address or error if something goes wrong.
	Create(tenantId, applicationId system.UUID, address shared.Address) (system.UUID, error)

	// Update updates an existing address.
	// tenantId: Mandatory. The unique identifier of the tenant owning the address.
	// applicationId: Mandatory. The unique identifier of the tenant's application will be owning the address.
	// addressId: Mandatory. The unique identifier of the existing address.
	// address: Mandatory. The reeference to the updated address information.
	// Returns error if something goes wrong.
	Update(tenantId, applicationId, addressId system.UUID, address shared.Address) error

	// Read retrieves an existing address information and returns only the detail which the keys provided by the detailsKeys.
	// tenantId: Mandatory. The unique identifier of the tenant owning the address.
	// applicationId: Mandatory. The unique identifier of the tenant's application will be owning the address.
	// addressId: Mandatory. The unique identifier of the existing address.
	// detailsKeys: Mandatory. The interested address details keys to return.
	// Returns either the address information or error if something goes wrong.
	Read(tenantId, applicationId, addressId system.UUID, detailsKeys map[string]string) (shared.Address, error)

	// ReadAll retrieves an existing address information and returns all the detail of it.
	// tenantId: Mandatory. The unique identifier of the tenant owning the address.
	// applicationId: Mandatory. The unique identifier of the tenant's application will be owning the address.
	// addressId: Mandatory. The unique identifier of the existing address.
	// Returns either the address information or error if something goes wrong.
	ReadAll(tenantId, applicationId, addressId system.UUID) (shared.Address, error)

	// Delete deletes an existing address information.
	// tenantId: Mandatory. The unique identifier of the tenant owning the address.
	// applicationId: Mandatory. The unique identifier of the tenant's application will be owning the address.
	// addressId: Mandatory. The unique identifier of the existing address to remove.
	// Returns error if something goes wrong.
	Delete(tenantId, applicationId, addressId system.UUID) error
}
