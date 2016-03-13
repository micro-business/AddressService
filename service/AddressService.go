// The default implementation of the address service.
package service

import (
	. "github.com/microbusinesses/AddressService/domain"
	. "github.com/microbusinesses/Micro-Businesses-Core/common/diagnostics"
	. "github.com/microbusinesses/Micro-Businesses-Core/system"
)

// The address service provides access to add new address and update/retrieve/remove an existing address.
type AddressService struct{}

// Adds a new address.
// tenantId: Mandatory. The unique identifier of the tenant owning the address.
// address: Mandatory. The reference to the new address information.
// Returns either the unique identifier of the new address or error if something goes wrong.
func (AddressService) Add(tenantId UUID, address Address) (UUID, error) {
	IsNotNilOrEmpty(tenantId, "tenantId", "tenantId must provided.")
	IsNotNilOrEmpty(address, "address", "address must provided.")

	// TODO: 20160313: Morteza: Should not do the following checking, it must be done in IsNotNillOrEmpty function!!!
	if len(address) == 0 {
		panic("Address does not contain any address part.")
	}

	return EmptyUUID, nil
}

// Updates an existing address.
// tenantId: Mandatory. The unique identifier of the tenant owning the address.
// id: Mandatory. The unique identifier of the existing address.
// address: Mandatory. The reeference to the updated address information.
// Returns error if something goes wrong.
func (AddressService) Update(tenantId UUID, id UUID, address Address) error {
	panic("Not Implemented")
}

// Retrieve an existing address information and returns the detail of it.
// tenantId: Mandatory. The unique identifier of the tenant owning the address.
// id: Mandatory. The unique identifier of the existing address.
// Returns either the address information or error if something goes wrong.
func (AddressService) Get(tenantId UUID, id UUID) (Address, error) {
	panic("Not Implemented")
}

// Removes an existing address information.
// tenantId: Mandatory. The unique identifier of the tenant owning the address.
// id: Mandatory. The unique identifier of the existing address to remove.
// Returns error if something goes wrong.
func (AddressService) Remove(tenantId UUID, id UUID) error {
	panic("Not Implemented")
}
