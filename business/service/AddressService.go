package service

import (
	"github.com/microbusinesses/AddressService/business/domain"
	"github.com/microbusinesses/AddressService/data/contract"
	"github.com/microbusinesses/Micro-Businesses-Core/common/diagnostics"
	"github.com/microbusinesses/Micro-Businesses-Core/system"
)

// AddressService provides access to add new address and update/retrieve/remove an existing address.
type AddressService struct {
	AddressDataService contract.AddressDataService
}

// Create creates a new address.
// tenantID: Mandatory. The unique identifier of the tenant owning the address.
// applicationID: Mandatory. The unique identifier of the tenant's application will be owning the address.
// address: Mandatory. The reference to the new address information.
// Returns either the unique identifier of the new address or error if something goes wrong.
func (addressService AddressService) Create(tenantID, applicationID system.UUID, address domain.Address) (system.UUID, error) {
	diagnostics.IsNotNil(addressService.AddressDataService, "addressService.AddressDataService", "AddressDataService must be provided.")
	diagnostics.IsNotNilOrEmpty(tenantID, "tenantID", "tenantID must be provided.")
	diagnostics.IsNotNilOrEmpty(applicationID, "applicationID", "applicationID must be provided.")

	if len(address.AddressDetails) == 0 {
		panic("Address does not contain any address key.")
	}

	return addressService.AddressDataService.Create(tenantID, applicationID, mapToDataAddress(address))
}

// Update updates an existing address.
// tenantID: Mandatory. The unique identifier of the tenant owning the address.
// applicationID: Mandatory. The unique identifier of the tenant's application will be owning the address.
// addressID: Mandatory. The unique identifier of the existing address.
// address: Mandatory. The reeference to the updated address information.
// Returns error if something goes wrong.
func (addressService AddressService) Update(tenantID, applicationID, addressID system.UUID, address domain.Address) error {
	diagnostics.IsNotNil(addressService.AddressDataService, "addressService.AddressDataService", "AddressDataService must be provided.")
	diagnostics.IsNotNilOrEmpty(tenantID, "tenantID", "tenantID must be provided.")
	diagnostics.IsNotNilOrEmpty(applicationID, "applicationID", "applicationID must be provided.")
	diagnostics.IsNotNilOrEmpty(addressID, "addressID", "addressID must be provided.")

	if len(address.AddressDetails) == 0 {
		panic("Address does not contain any address key.")
	}

	return addressService.AddressDataService.Update(tenantID, applicationID, addressID, mapToDataAddress(address))
}

// Read retrieves an existing address information and returns only the detail which the keys provided by the keys.
// tenantID: Mandatory. The unique identifier of the tenant owning the address.
// applicationID: Mandatory. The unique identifier of the tenant's application will be owning the address.
// addressID: Mandatory. The unique identifier of the existing address.
// keys: Mandatory. The interested address details keys to return.
// Returns either the address information or error if something goes wrong.
func (addressService AddressService) Read(tenantID, applicationID, addressID system.UUID, keys []string) (domain.Address, error) {
	diagnostics.IsNotNil(addressService.AddressDataService, "addressService.AddressDataService", "AddressDataService must be provided.")
	diagnostics.IsNotNilOrEmpty(tenantID, "tenantID", "tenantID must be provided.")
	diagnostics.IsNotNilOrEmpty(applicationID, "applicationID", "applicationID must be provided.")
	diagnostics.IsNotNilOrEmpty(addressID, "addressID", "addressID must be provided.")

	//TODO: 20160722: Add details keys validation here

	address, err := addressService.AddressDataService.Read(tenantID, applicationID, addressID, keys)

	if err != nil {
		return domain.Address{}, err
	}

	return mapFromDataAddress(address), nil
}

// ReadAll retrieves an existing address information and returns all the detail of it.
// tenantID: Mandatory. The unique identifier of the tenant owning the address.
// applicationID: Mandatory. The unique identifier of the tenant's application will be owning the address.
// addressID: Mandatory. The unique identifier of the existing address.
// Returns either the address information or error if something goes wrong.
func (addressService AddressService) ReadAll(tenantID, applicationID, addressID system.UUID) (domain.Address, error) {
	diagnostics.IsNotNil(addressService.AddressDataService, "addressService.AddressDataService", "AddressDataService must be provided.")
	diagnostics.IsNotNilOrEmpty(tenantID, "tenantID", "tenantID must be provided.")
	diagnostics.IsNotNilOrEmpty(applicationID, "applicationID", "applicationID must be provided.")
	diagnostics.IsNotNilOrEmpty(addressID, "addressID", "addressID must be provided.")

	address, err := addressService.AddressDataService.ReadAll(tenantID, applicationID, addressID)

	if err != nil {
		return domain.Address{}, err
	}

	return mapFromDataAddress(address), nil
}

// Delete deletes an existing address information.
// tenantID: Mandatory. The unique identifier of the tenant owning the address.
// applicationID: Mandatory. The unique identifier of the tenant's application will be owning the address.
// addressID: Mandatory. The unique identifier of the existing address to remove.
// Returns error if something goes wrong.
func (addressService AddressService) Delete(tenantID, applicationID, addressID system.UUID) error {
	diagnostics.IsNotNil(addressService.AddressDataService, "addressService.AddressDataService", "AddressDataService must be provided.")
	diagnostics.IsNotNilOrEmpty(tenantID, "tenantID", "tenantID must be provided.")
	diagnostics.IsNotNilOrEmpty(applicationID, "applicationID", "applicationID must be provided.")
	diagnostics.IsNotNilOrEmpty(addressID, "addressID", "addressID  must be provided.")

	return addressService.AddressDataService.Delete(tenantID, applicationID, addressID)
}

// mapToDataAddress Maps the domain address object to the Address object used in data layer.
// address: Mandatory. The address domain object
// Returns the converted address object used in data layer
func mapToDataAddress(address domain.Address) contract.Address {
	return contract.Address{AddressDetails: address.AddressDetails}
}

// mapFromDataAddress Maps the address object used in data layer to the Address domain object.
// address: Mandatory. The address object used in data layer
// Returns the converted address domain object
func mapFromDataAddress(address contract.Address) domain.Address {
	return domain.Address{AddressDetails: address.AddressDetails}
}
