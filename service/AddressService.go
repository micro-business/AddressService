package service

import (
	. "github.com/microbusinesses/AddressService/domain"
	. "github.com/microbusinesses/Micro-Businesses-Core/common/diagnostics"
	. "github.com/microbusinesses/Micro-Businesses-Core/system"
)

// The address service provides access to add new address and update/retrieve/remove an existing address.
type AddressService struct {
	UUIDGeneratorService UUIDGeneratorService
}

// Creates a new address.
// tenantId: Mandatory. The unique identifier of the tenant owning the address.
// applicationId: Mandatory. The unique identifier of the tenant's application will be owning the address.
// address: Mandatory. The reference to the new address information.
// Returns either the unique identifier of the new address or error if something goes wrong.
func (AddressService) Create(tenantId UUID, applicationId UUID, address Address) (UUID, error) {
	IsNotNilOrEmpty(tenantId, "tenantId", "tenantId must be provided.")
	IsNotNilOrEmpty(applicationId, "applicationId", "applicationId must be provided.")

	if len(address.AddressParts) == 0 {
		panic("Address does not contain any address part.")
	}

	if addressId, err := RandomUUID(); err != nil {
		return EmptyUUID, err
	} else {
		return addressId, nil
	}

}

// Updates an existing address.
// tenantId: Mandatory. The unique identifier of the tenant owning the address.
// applicationId: Mandatory. The unique identifier of the tenant's application will be owning the address.
// addressId: Mandatory. The unique identifier of the existing address.
// address: Mandatory. The reeference to the updated address information.
// Returns error if something goes wrong.
func (AddressService) Update(tenantId UUID, applicationId UUID, addressId UUID, address Address) error {
	IsNotNilOrEmpty(tenantId, "tenantId", "tenantId must be provided.")
	IsNotNilOrEmpty(applicationId, "applicationId", "applicationId must be provided.")
	IsNotNilOrEmpty(addressId, "addressId", "addressId must be provided.")

	if len(address.AddressParts) == 0 {
		panic("Address does not contain any address part.")
	}

	return nil
}

// Retrieves an existing address information and returns the detail of it.
// tenantId: Mandatory. The unique identifier of the tenant owning the address.
// applicationId: Mandatory. The unique identifier of the tenant's application will be owning the address.
// addressId: Mandatory. The unique identifier of the existing address.
// Returns either the address information or error if something goes wrong.
func (AddressService) Read(tenantId UUID, applicationId UUID, addressId UUID) (Address, error) {
	IsNotNilOrEmpty(tenantId, "tenantId", "tenantId must be provided.")
	IsNotNilOrEmpty(applicationId, "applicationId", "applicationId must be provided.")
	IsNotNilOrEmpty(addressId, "addressId", "addressId must be provided.")

	panic("Not Implemented")

	return Address{}, nil
}

// Deletes an existing address information.
// tenantId: Mandatory. The unique identifier of the tenant owning the address.
// applicationId: Mandatory. The unique identifier of the tenant's application will be owning the address.
// addressId: Mandatory. The unique identifier of the existing address to remove.
// Returns error if something goes wrong.
func (AddressService) Delete(tenantId UUID, applicationId UUID, addressId UUID) error {
	IsNotNilOrEmpty(tenantId, "tenantId", "tenantId must be provided.")
	IsNotNilOrEmpty(applicationId, "applicationId", "applicationId must be provided.")
	IsNotNilOrEmpty(addressId, "addressId", "addressId  must be provided.")

	return nil
}
