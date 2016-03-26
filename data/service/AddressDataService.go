package service

import (
	"errors"
	"sync"

	"github.com/gocql/gocql"
	"github.com/microbusinesses/AddressService/data/shared"
	"github.com/microbusinesses/Micro-Businesses-Core/common/diagnostics"
	"github.com/microbusinesses/Micro-Businesses-Core/system"
)

// The address service provides access to add new address and update/retrieve/remove an existing address.
type AddressDataService struct {
	UUIDGeneratorService UUIDGeneratorService
	ClusterConfig        *gocql.ClusterConfig
}

// Create creates a new address.
// tenantId: Mandatory. The unique identifier of the tenant owning the address.
// applicationId: Mandatory. The unique identifier of the tenant's application will be owning the address.
// address: Mandatory. The reference to the new address information.
// Returns either the unique identifier of the new address or error if something goes wrong.
func (addressDataService *AddressDataService) Create(tenantId system.UUID, applicationId system.UUID, address shared.Address) (system.UUID, error) {
	diagnostics.IsNotNil(addressDataService.UUIDGeneratorService, "addressDataService.UUIDGeneratorService", "UUIDGeneratorService must be provided.")
	diagnostics.IsNotNil(addressDataService.ClusterConfig, "addressDataService.ClusterConfig", "ClusterConfig must be provided.")
	diagnostics.IsNotNilOrEmpty(tenantId, "tenantId", "tenantId must be provided.")
	diagnostics.IsNotNilOrEmpty(applicationId, "applicationId", "applicationId must be provided.")

	addressPartsCount := len(address.AddressParts)

	if addressPartsCount == 0 {
		panic("Address does not contain any address part.")
	}

	addressId, err := addressDataService.UUIDGeneratorService.GenerateRandomUUID()

	if err != nil {
		return system.EmptyUUID, err
	}

	session, err := addressDataService.ClusterConfig.CreateSession()

	if err != nil {
		return system.EmptyUUID, err
	}

	defer session.Close()

	errorChannel := make(chan error, addressPartsCount)

	var waitGroup sync.WaitGroup

	for addressPart, addressValue := range address.AddressParts {
		waitGroup.Add(1)

		go func(tenantId, applicationId, addressId gocql.UUID, addressPart, addressValue string) {
			defer waitGroup.Done()

			if err := session.Query(
				"INSERT INTO address (tenant_id, application_id, address_id, address_part, address_value) VALUES(?, ?, ?, ?, ?)",
				tenantId,
				applicationId,
				addressId,
				addressPart,
				addressValue).
				Exec(); err != nil {
				errorChannel <- err
			} else {
				errorChannel <- nil
			}
		}(
			mapSystemUUIDToGocqlUUID(tenantId),
			mapSystemUUIDToGocqlUUID(applicationId),
			mapSystemUUIDToGocqlUUID(addressId),
			addressPart,
			addressValue)
	}

	go func() {
		waitGroup.Wait()
		close(errorChannel)
	}()

	errorMessage := ""
	errorFound := false

	for err := range errorChannel {
		if err != nil {
			errorMessage += err.Error()
			errorMessage += "\n"
			errorFound = true
		}
	}

	if errorFound {
		return system.EmptyUUID, errors.New(errorMessage)
	}

	return addressId, nil

}

// Update updates an existing address.
// tenantId: Mandatory. The unique identifier of the tenant owning the address.
// applicationId: Mandatory. The unique identifier of the tenant's application will be owning the address.
// addressId: Mandatory. The unique identifier of the existing address.
// address: Mandatory. The reeference to the updated address information.
// Returns error if something goes wrong.
func (addressDataService *AddressDataService) Update(tenantId system.UUID, applicationId system.UUID, addressId system.UUID, address shared.Address) error {
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
	diagnostics.IsNotNilOrEmpty(tenantId, "tenantId", "tenantId must be provided.")
	diagnostics.IsNotNilOrEmpty(applicationId, "applicationId", "applicationId must be provided.")
	diagnostics.IsNotNilOrEmpty(addressId, "addressId", "addressId  must be provided.")

	panic("Not Implemented")
}

// mapSystemUUIDToGocqlUUID maps the system type UUID to gocql UUID type
func mapSystemUUIDToGocqlUUID(uuid system.UUID) gocql.UUID {
	mappedUUID, _ := gocql.UUIDFromBytes(uuid.Bytes())

	return mappedUUID
}
