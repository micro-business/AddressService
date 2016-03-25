package service

import (
	"github.com/microbusinesses/AddressService/data/shared"
	"github.com/microbusinesses/Micro-Businesses-Core/common/diagnostics"
	"github.com/microbusinesses/Micro-Businesses-Core/system"
)

// The address service provides access to add new address and update/retrieve/remove an existing address.
type AddressDataService struct {
	UUIDGeneratorService UUIDGeneratorService
}

// Create creates a new address.
// tenantId: Mandatory. The unique identifier of the tenant owning the address.
// applicationId: Mandatory. The unique identifier of the tenant's application will be owning the address.
// address: Mandatory. The reference to the new address information.
// Returns either the unique identifier of the new address or error if something goes wrong.
func (addressDataService *AddressDataService) Create(tenantId system.UUID, applicationId system.UUID, address shared.Address) (system.UUID, error) {
	diagnostics.IsNotNil(addressDataService.UUIDGeneratorService, "addressDataService.UUIDGeneratorService", "UUIDGeneratorService must be provided.")
	diagnostics.IsNotNilOrEmpty(tenantId, "tenantId", "tenantId must be provided.")
	diagnostics.IsNotNilOrEmpty(applicationId, "applicationId", "applicationId must be provided.")

	if len(address.AddressParts) == 0 {
		panic("Address does not contain any address part.")
	}

	panic("Not Implemented")
}

// Update updates an existing address.
// tenantId: Mandatory. The unique identifier of the tenant owning the address.
// applicationId: Mandatory. The unique identifier of the tenant's application will be owning the address.
// addressId: Mandatory. The unique identifier of the existing address.
// address: Mandatory. The reeference to the updated address information.
// Returns error if something goes wrong.
func (addressDataService *AddressDataService) Update(tenantId system.UUID, applicationId system.UUID, addressId system.UUID, address shared.Address) error {
	diagnostics.IsNotNil(addressDataService.UUIDGeneratorService, "addressDataService.UUIDGeneratorService", "UUIDGeneratorService must be provided.")
	diagnostics.IsNotNilOrEmpty(applicationId, "applicationId", "applicationId must be provided.")
	diagnostics.IsNotNilOrEmpty(addressId, "addressId", "addressId must be provided.")

	if len(address.AddressParts) == 0 {
		panic("Address does not contain any address part.")
	}

	panic("Not Implemented")
}

// Read retrieves an existing address information and returns the detail of it.
// tenantId: Mandatory. The unique identifier of the tenant owning the address.
// applicationId: Mandatory. The unique identifier of the tenant's application will be owning the address.
// addressId: Mandatory. The unique identifier of the existing address.
// Returns either the address information or error if something goes wrong.
func (addressDataService *AddressDataService) Read(tenantId system.UUID, applicationId system.UUID, addressId system.UUID) (shared.Address, error) {
	diagnostics.IsNotNil(addressDataService.UUIDGeneratorService, "addressDataService.UUIDGeneratorService", "UUIDGeneratorService must be provided.")
	diagnostics.IsNotNilOrEmpty(tenantId, "tenantId", "tenantId must be provided.")
	diagnostics.IsNotNilOrEmpty(applicationId, "applicationId", "applicationId must be provided.")
	diagnostics.IsNotNilOrEmpty(addressId, "addressId", "addressId must be provided.")

	panic("Not Implemented")
}

// Delete deletes an existing address information.
// tenantId: Mandatory. The unique identifier of the tenant owning the address.
// applicationId: Mandatory. The unique identifier of the tenant's application will be owning the address.
// addressId: Mandatory. The unique identifier of the existing address to remove.
// Returns error if something goes wrong.
func (addressDataService *AddressDataService) Delete(tenantId system.UUID, applicationId system.UUID, addressId system.UUID) error {
	diagnostics.IsNotNil(addressDataService.UUIDGeneratorService, "addressDataService.UUIDGeneratorService", "UUIDGeneratorService must be provided.")
	diagnostics.IsNotNilOrEmpty(tenantId, "tenantId", "tenantId must be provided.")
	diagnostics.IsNotNilOrEmpty(applicationId, "applicationId", "applicationId must be provided.")
	diagnostics.IsNotNilOrEmpty(addressId, "addressId", "addressId  must be provided.")

	panic("Not Implemented")
}
