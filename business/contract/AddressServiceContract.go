// Defines the address service contract.
package contract

import (
	"github.com/microbusinesses/AddressService/business/domain"
	"github.com/microbusinesses/Micro-Businesses-Core/common/query"
	"github.com/microbusinesses/Micro-Businesses-Core/system"
)

// The address service contract, it can add new address and update/retrieve/remove an existing address.
type AddressService interface {

	// ProcessQuery processes the provided query through API interface.
	// tenantId: Mandatory. The unique identifier of the tenant owning the address.
	// applicationId: Mandatory. The unique identifier of the tenant's application will be owning the address.
	// requestQuery: Mandatory. The reference to the request query.
	ProcessQuery(tenantId, applicationId system.UUID, requestQuery query.RequestQuery) (query.ResponseQuery, error)

	// Create creates a new address.
	// tenantId: Mandatory. The unique identifier of the tenant owning the address.
	// applicationId: Mandatory. The unique identifier of the tenant's application will be owning the address.
	// address: Mandatory. The reference to the new address information.
	// Returns either the unique identifier of the new address or error if something goes wrong.
	Create(tenantId, applicationId system.UUID, address domain.Address) (system.UUID, error)

	// Update updates an existing address.
	// tenantId: Mandatory. The unique identifier of the tenant owning the address.
	// applicationId: Mandatory. The unique identifier of the tenant's application will be owning the address.
	// addressId: Mandatory. The unique identifier of the existing address.
	// address: Mandatory. The reeference to the updated address information.
	// Returns error if something goes wrong.
	Update(tenantId, applicationId, addressId system.UUID, address domain.Address) error

	// Read retrieves an existing address information and returns the detail of it.
	// tenantId: Mandatory. The unique identifier of the tenant owning the address.
	// applicationId: Mandatory. The unique identifier of the tenant's application will be owning the address.
	// addressId: Mandatory. The unique identifier of the existing address.
	// Returns either the address information or error if something goes wrong.
	ReadAll(tenantId, applicationId, addressId system.UUID) (domain.Address, error)

	// Delete deletes an existing address information.
	// tenantId: Mandatory. The unique identifier of the tenant owning the address.
	// applicationId: Mandatory. The unique identifier of the tenant's application will be owning the address.
	// addressId: Mandatory. The unique identifier of the existing address to remove.
	// Returns error if something goes wrong.
	Delete(tenantId, applicationId, addressId system.UUID) error
}
