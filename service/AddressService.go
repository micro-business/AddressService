// The default implementation of the address service.
package service

import (
	"github.com/microbusinesses/AddressService/domain"
	"github.com/microbusinesses/Micro-Businesses-Core/system"
)

// The address service provides access to add new address and update/retrieve/remove an existing address.
type AddressService struct{}

// Adds a new address.
// tenantId: Mandatory. The unique identifier of the tenant owning the address.
// address: Mandatory. The reference to the new address information.
// Returns either the unique identifier of the new address or error if something goes wrong.
func (AddressService) Add(tenantId system.UUID, address domain.Address) (system.UUID, error) {
	panic("Not Implemented")
}

// Updates an existing address.
// tenantId: Mandatory. The unique identifier of the tenant owning the address.
// id: Mandatory. The unique identifier of the existing address.
// address: Mandatory. The reeference to the updated address information.
// Returns error if something goes wrong.
func (AddressService) Update(tenantId system.UUID, id system.UUID, address domain.Address) error {
	panic("Not Implemented")
}

// Retrieve an existing address information and returns the detail of it.
// tenantId: Mandatory. The unique identifier of the tenant owning the address.
// id: Mandatory. The unique identifier of the existing address.
// Returns either the address information or error if something goes wrong.
func (AddressService) Get(tenantId system.UUID, id system.UUID) (domain.Address, error) {
	panic("Not Implemented")
}

// Removes an existing address information.
// tenantId: Mandatory. The unique identifier of the tenant owning the address.
// id: Mandatory. The unique identifier of the existing address to remove.
// Returns error if something goes wrong.
func (AddressService) Remove(tenantId system.UUID, id system.UUID) error {
	panic("Not Implemented")
}
