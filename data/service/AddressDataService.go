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
func (addressDataService *AddressDataService) Create(tenantId, applicationId system.UUID, address shared.Address) (system.UUID, error) {
	diagnostics.IsNotNil(addressDataService.UUIDGeneratorService, "addressDataService.UUIDGeneratorService", "UUIDGeneratorService must be provided.")
	diagnostics.IsNotNil(addressDataService.ClusterConfig, "addressDataService.ClusterConfig", "ClusterConfig must be provided.")
	diagnostics.IsNotNilOrEmpty(tenantId, "tenantId", "tenantId must be provided.")
	diagnostics.IsNotNilOrEmpty(applicationId, "applicationId", "applicationId must be provided.")

	addressKeysValuesCount := len(address.AddressKeysValues)

	if addressKeysValuesCount == 0 {
		panic("Address does not contain any address key.")
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

	errorChannel := make(chan error, addressKeysValuesCount*2)

	mappedTenantId := mapSystemUUIDToGocqlUUID(tenantId)
	mappedApplicationId := mapSystemUUIDToGocqlUUID(applicationId)
	mappedAddressId := mapSystemUUIDToGocqlUUID(addressId)

	var waitGroup sync.WaitGroup

	for key, value := range address.AddressKeysValues {
		waitGroup.Add(1)

		go addToAddressTable(
			session,
			errorChannel,
			&waitGroup,
			mappedTenantId,
			mappedApplicationId,
			mappedAddressId,
			key,
			value)

		waitGroup.Add(1)

		go addToAddressIndexByAddressKeyTable(
			session,
			errorChannel,
			&waitGroup,
			mappedTenantId,
			mappedApplicationId,
			mappedAddressId,
			key,
			value)
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
func (addressDataService *AddressDataService) Update(tenantId, applicationId, addressId system.UUID, address shared.Address) error {
	diagnostics.IsNotNilOrEmpty(applicationId, "applicationId", "applicationId must be provided.")
	diagnostics.IsNotNilOrEmpty(addressId, "addressId", "addressId must be provided.")

	if len(address.AddressKeysValues) == 0 {
		panic("Address does not contain any address key.")
	}

	panic("Not Implemented")
}

// Read retrieves an existing address information and returns the detail of it.
// tenantId: Mandatory. The unique identifier of the tenant owning the address.
// applicationId: Mandatory. The unique identifier of the tenant's application will be owning the address.
// addressId: Mandatory. The unique identifier of the existing address.
// Returns either the address information or error if something goes wrong.
func (addressDataService *AddressDataService) Read(tenantId, applicationId, addressId system.UUID) (shared.Address, error) {
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
func (addressDataService *AddressDataService) Delete(tenantId, applicationId, addressId system.UUID) error {
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

// addToAddressTable adds new address key/value to address table using provided address unique identifier.
func addToAddressTable(
	session *gocql.Session,
	errorChannel chan<- error,
	waitGroup *sync.WaitGroup,
	tenantId, applicationId, addressId gocql.UUID,
	key, value string) {

	defer waitGroup.Done()

	if err := session.Query(
		"INSERT INTO address"+
			" (tenant_id, application_id, address_id, address_key, address_value)"+
			" VALUES(?, ?, ?, ?, ?)",
		tenantId,
		applicationId,
		addressId,
		key,
		value).
		Exec(); err != nil {
		errorChannel <- err
	} else {
		errorChannel <- nil
	}
}

// addToAddressIndexByAddressKeyTable adds address key/value to index table, so running query on address key will be faster.
func addToAddressIndexByAddressKeyTable(
	session *gocql.Session,
	errorChannel chan<- error,
	waitGroup *sync.WaitGroup,
	tenantId, applicationId, addressId gocql.UUID,
	key, value string) {

	defer waitGroup.Done()

	if err := session.Query(
		"INSERT INTO address_indexed_by_address_key"+
			" (tenant_id, application_id, address_id, address_key, address_value)"+
			" VALUES(?, ?, ?, ?, ?)",
		tenantId,
		applicationId,
		addressId,
		key,
		value).
		Exec(); err != nil {
		errorChannel <- err
	} else {
		errorChannel <- nil
	}
}
