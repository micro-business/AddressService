package service

import (
	. "github.com/microbusinesses/AddressService/business/domain"
	. "github.com/microbusinesses/AddressService/data/contract"
	shared "github.com/microbusinesses/AddressService/data/shared"
	. "github.com/microbusinesses/Micro-Businesses-Core/common/diagnostics"
	. "github.com/microbusinesses/Micro-Businesses-Core/system"
)

// The address service provides access to add new address and update/retrieve/remove an existing address.
type AddressService struct {
	AddressDataService AddressDataService
}

// Creates a new address.
// tenantId: Mandatory. The unique identifier of the tenant owning the address.
// applicationId: Mandatory. The unique identifier of the tenant's application will be owning the address.
// address: Mandatory. The reference to the new address information.
// Returns either the unique identifier of the new address or error if something goes wrong.
func (addressService AddressService) Create(tenantId UUID, applicationId UUID, address Address) (UUID, error) {
	IsNotNil(addressService.AddressDataService, "addressService.AddressDataService", "AddressDataService must be provided.")
	IsNotNilOrEmpty(tenantId, "tenantId", "tenantId must be provided.")
	IsNotNilOrEmpty(applicationId, "applicationId", "applicationId must be provided.")

	if len(address.AddressParts) == 0 {
		panic("Address does not contain any address part.")
	}

	return addressService.AddressDataService.Create(tenantId, applicationId, mapToDataAddress(address))
}

// Updates an existing address.
// tenantId: Mandatory. The unique identifier of the tenant owning the address.
// applicationId: Mandatory. The unique identifier of the tenant's application will be owning the address.
// addressId: Mandatory. The unique identifier of the existing address.
// address: Mandatory. The reeference to the updated address information.
// Returns error if something goes wrong.
func (addressService AddressService) Update(tenantId UUID, applicationId UUID, addressId UUID, address Address) error {
	IsNotNil(addressService.AddressDataService, "addressService.AddressDataService", "AddressDataService must be provided.")
	IsNotNilOrEmpty(tenantId, "tenantId", "tenantId must be provided.")
	IsNotNilOrEmpty(applicationId, "applicationId", "applicationId must be provided.")
	IsNotNilOrEmpty(addressId, "addressId", "addressId must be provided.")

	if len(address.AddressParts) == 0 {
		panic("Address does not contain any address part.")
	}

	return addressService.AddressDataService.Update(tenantId, applicationId, addressId, mapToDataAddress(address))
}

// Retrieves an existing address information and returns the detail of it.
// tenantId: Mandatory. The unique identifier of the tenant owning the address.
// applicationId: Mandatory. The unique identifier of the tenant's application will be owning the address.
// addressId: Mandatory. The unique identifier of the existing address.
// Returns either the address information or error if something goes wrong.
func (addressService AddressService) Read(tenantId UUID, applicationId UUID, addressId UUID) (Address, error) {
	IsNotNil(addressService.AddressDataService, "addressService.AddressDataService", "AddressDataService must be provided.")
	IsNotNilOrEmpty(tenantId, "tenantId", "tenantId must be provided.")
	IsNotNilOrEmpty(applicationId, "applicationId", "applicationId must be provided.")
	IsNotNilOrEmpty(addressId, "addressId", "addressId must be provided.")

	if address, err := addressService.AddressDataService.Read(tenantId, applicationId, addressId); err != nil {
		return Address{}, err
	} else {
		return mapFromDataAddress(address), nil
	}
}

// Deletes an existing address information.
// tenantId: Mandatory. The unique identifier of the tenant owning the address.
// applicationId: Mandatory. The unique identifier of the tenant's application will be owning the address.
// addressId: Mandatory. The unique identifier of the existing address to remove.
// Returns error if something goes wrong.
func (addressService AddressService) Delete(tenantId UUID, applicationId UUID, addressId UUID) error {
	IsNotNil(addressService.AddressDataService, "addressService.AddressDataService", "AddressDataService must be provided.")
	IsNotNilOrEmpty(tenantId, "tenantId", "tenantId must be provided.")
	IsNotNilOrEmpty(applicationId, "applicationId", "applicationId must be provided.")
	IsNotNilOrEmpty(addressId, "addressId", "addressId  must be provided.")

	return addressService.AddressDataService.Delete(tenantId, applicationId, addressId)
}

// Maps the domain address object to the Address object used in data layer.
// address: Mandatory. The address domain object
// Returns the converted address object used in data layer
func mapToDataAddress(address Address) shared.Address {
	return shared.Address{AddressParts: address.AddressParts}
}

// Maps the address object used in data layer to the Address domain object.
// address: Mandatory. The address object used in data layer
// Returns the converted address domain object
func mapFromDataAddress(address shared.Address) Address {
	return Address{AddressParts: address.AddressParts}
}
