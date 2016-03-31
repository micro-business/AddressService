package service

import (
	"github.com/microbusinesses/AddressService/business/domain"
	dataContract "github.com/microbusinesses/AddressService/data/contract"
	dataShared "github.com/microbusinesses/AddressService/data/shared"
	"github.com/microbusinesses/Micro-Businesses-Core/common/diagnostics"
	"github.com/microbusinesses/Micro-Businesses-Core/system"
)

// The address service provides access to add new address and update/retrieve/remove an existing address.
type AddressService struct {
	AddressDataService dataContract.AddressDataService
}

// Create creates a new address.
// tenantId: Mandatory. The unique identifier of the tenant owning the address.
// applicationId: Mandatory. The unique identifier of the tenant's application will be owning the address.
// address: Mandatory. The reference to the new address information.
// Returns either the unique identifier of the new address or error if something goes wrong.
func (addressService AddressService) Create(tenantId, applicationId system.UUID, address domain.Address) (system.UUID, error) {
	diagnostics.IsNotNil(addressService.AddressDataService, "addressService.AddressDataService", "AddressDataService must be provided.")
	diagnostics.IsNotNilOrEmpty(tenantId, "tenantId", "tenantId must be provided.")
	diagnostics.IsNotNilOrEmpty(applicationId, "applicationId", "applicationId must be provided.")

	if len(address.AddressKeysValues) == 0 {
		panic("Address does not contain any address key.")
	}

	return addressService.AddressDataService.Create(tenantId, applicationId, mapToDataAddress(address))
}

// Update updates an existing address.
// tenantId: Mandatory. The unique identifier of the tenant owning the address.
// applicationId: Mandatory. The unique identifier of the tenant's application will be owning the address.
// addressId: Mandatory. The unique identifier of the existing address.
// address: Mandatory. The reeference to the updated address information.
// Returns error if something goes wrong.
func (addressService AddressService) Update(tenantId, applicationId, addressId system.UUID, address domain.Address) error {
	diagnostics.IsNotNil(addressService.AddressDataService, "addressService.AddressDataService", "AddressDataService must be provided.")
	diagnostics.IsNotNilOrEmpty(tenantId, "tenantId", "tenantId must be provided.")
	diagnostics.IsNotNilOrEmpty(applicationId, "applicationId", "applicationId must be provided.")
	diagnostics.IsNotNilOrEmpty(addressId, "addressId", "addressId must be provided.")

	if len(address.AddressKeysValues) == 0 {
		panic("Address does not contain any address key.")
	}

	return addressService.AddressDataService.Update(tenantId, applicationId, addressId, mapToDataAddress(address))
}

// Read retrieves an existing address information and returns the detail of it.
// tenantId: Mandatory. The unique identifier of the tenant owning the address.
// applicationId: Mandatory. The unique identifier of the tenant's application will be owning the address.
// addressId: Mandatory. The unique identifier of the existing address.
// Returns either the address information or error if something goes wrong.
func (addressService AddressService) Read(tenantId, applicationId, addressId system.UUID) (domain.Address, error) {
	diagnostics.IsNotNil(addressService.AddressDataService, "addressService.AddressDataService", "AddressDataService must be provided.")
	diagnostics.IsNotNilOrEmpty(tenantId, "tenantId", "tenantId must be provided.")
	diagnostics.IsNotNilOrEmpty(applicationId, "applicationId", "applicationId must be provided.")
	diagnostics.IsNotNilOrEmpty(addressId, "addressId", "addressId must be provided.")

	if address, err := addressService.AddressDataService.Read(tenantId, applicationId, addressId); err != nil {
		return domain.Address{}, err
	} else {
		return mapFromDataAddress(address), nil
	}
}

// Delete deletes an existing address information.
// tenantId: Mandatory. The unique identifier of the tenant owning the address.
// applicationId: Mandatory. The unique identifier of the tenant's application will be owning the address.
// addressId: Mandatory. The unique identifier of the existing address to remove.
// Returns error if something goes wrong.
func (addressService AddressService) Delete(tenantId, applicationId, addressId system.UUID) error {
	diagnostics.IsNotNil(addressService.AddressDataService, "addressService.AddressDataService", "AddressDataService must be provided.")
	diagnostics.IsNotNilOrEmpty(tenantId, "tenantId", "tenantId must be provided.")
	diagnostics.IsNotNilOrEmpty(applicationId, "applicationId", "applicationId must be provided.")
	diagnostics.IsNotNilOrEmpty(addressId, "addressId", "addressId  must be provided.")

	return addressService.AddressDataService.Delete(tenantId, applicationId, addressId)
}

// mapToDataAddress Maps the domain address object to the Address object used in data layer.
// address: Mandatory. The address domain object
// Returns the converted address object used in data layer
func mapToDataAddress(address domain.Address) dataShared.Address {
	return dataShared.Address{AddressKeysValues: address.AddressKeysValues}
}

// mapFromDataAddress Maps the address object used in data layer to the Address domain object.
// address: Mandatory. The address object used in data layer
// Returns the converted address domain object
func mapFromDataAddress(address dataShared.Address) domain.Address {
	return domain.Address{AddressKeysValues: address.AddressKeysValues}
}
