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
func (addressDataService AddressDataService) Create(tenantId, applicationId system.UUID, address shared.Address) (system.UUID, error) {
	diagnostics.IsNotNil(addressDataService.UUIDGeneratorService, "addressDataService.UUIDGeneratorService", "UUIDGeneratorService must be provided.")
	diagnostics.IsNotNil(addressDataService.ClusterConfig, "addressDataService.ClusterConfig", "ClusterConfig must be provided.")
	diagnostics.IsNotNilOrEmpty(tenantId, "tenantId", "tenantId must be provided.")
	diagnostics.IsNotNilOrEmpty(applicationId, "applicationId", "applicationId must be provided.")

	addressDetailsCount := len(address.AddressDetails)

	if addressDetailsCount == 0 {
		panic("Address does not contain any address details.")
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

	if err = addNewAddress(tenantId, applicationId, address, addressId, session); err != nil {
		return system.EmptyUUID, err
	}

	return addressId, nil
}

// Update updates an existing address.
// tenantId: Mandatory. The unique identifier of the tenant owning the address.
// applicationId: Mandatory. The unique identifier of the tenant's application will be owning the address.
// addressId: Mandatory. The unique identifier of the existing address.
// address: Mandatory. The reeference to the updated address information.
// Returns error if something goes wrong.
func (addressDataService AddressDataService) Update(tenantId, applicationId, addressId system.UUID, address shared.Address) error {
	diagnostics.IsNotNil(addressDataService.ClusterConfig, "addressDataService.ClusterConfig", "ClusterConfig must be provided.")
	diagnostics.IsNotNilOrEmpty(tenantId, "tenantId", "tenantId must be provided.")
	diagnostics.IsNotNilOrEmpty(applicationId, "applicationId", "applicationId must be provided.")
	diagnostics.IsNotNilOrEmpty(addressId, "addressId", "addressId must be provided.")

	if len(address.AddressDetails) == 0 {
		panic("Address does not contain any address details.")
	}

	err := addressDataService.Delete(tenantId, applicationId, addressId)

	if err != nil {
		return err
	}

	session, err := addressDataService.ClusterConfig.CreateSession()

	if err != nil {
		return err
	}

	defer session.Close()

	return addNewAddress(tenantId, applicationId, address, addressId, session)
}

// Read retrieves an existing address information and returns the detail of it.
// tenantId: Mandatory. The unique identifier of the tenant owning the address.
// applicationId: Mandatory. The unique identifier of the tenant's application will be owning the address.
// addressId: Mandatory. The unique identifier of the existing address.
// Returns either the address information or error if something goes wrong.
func (addressDataService AddressDataService) ReadAll(tenantId, applicationId, addressId system.UUID) (shared.Address, error) {
	diagnostics.IsNotNil(addressDataService.ClusterConfig, "addressDataService.ClusterConfig", "ClusterConfig must be provided.")
	diagnostics.IsNotNilOrEmpty(tenantId, "tenantId", "tenantId must be provided.")
	diagnostics.IsNotNilOrEmpty(applicationId, "applicationId", "applicationId must be provided.")
	diagnostics.IsNotNilOrEmpty(addressId, "addressId", "addressId must be provided.")

	session, err := addressDataService.ClusterConfig.CreateSession()

	if err != nil {
		return shared.Address{}, err
	}

	defer session.Close()

	iter := session.Query(
		"SELECT address_key, address_value"+
			" FROM address"+
			" WHERE"+
			" tenant_id = ?"+
			" AND application_id = ?"+
			" AND address_id = ?",
		tenantId.String(),
		applicationId.String(),
		addressId.String()).Iter()

	var key string
	var value string

	address := shared.Address{AddressDetails: make(map[string]string)}

	for iter.Scan(&key, &value) {
		address.AddressDetails[key] = value
	}

	return address, nil
}

// Delete deletes an existing address information.
// tenantId: Mandatory. The unique identifier of the tenant owning the address.
// applicationId: Mandatory. The unique identifier of the tenant's application will be owning the address.
// addressId: Mandatory. The unique identifier of the existing address to remove.
// Returns error if something goes wrong.
func (addressDataService AddressDataService) Delete(tenantId, applicationId, addressId system.UUID) error {
	diagnostics.IsNotNil(addressDataService.ClusterConfig, "addressDataService.ClusterConfig", "ClusterConfig must be provided.")
	diagnostics.IsNotNilOrEmpty(tenantId, "tenantId", "tenantId must be provided.")
	diagnostics.IsNotNilOrEmpty(applicationId, "applicationId", "applicationId must be provided.")
	diagnostics.IsNotNilOrEmpty(addressId, "addressId", "addressId must be provided.")

	address, err := addressDataService.ReadAll(tenantId, applicationId, addressId)

	if err != nil {
		return err
	}

	session, err := addressDataService.ClusterConfig.CreateSession()

	if err != nil {
		return err
	}

	defer session.Close()

	return removeExistingAddress(tenantId, applicationId, address, addressId, session)
}

// mapSystemUUIDToGocqlUUID maps the system type UUID to gocql UUID type
func mapSystemUUIDToGocqlUUID(uuid system.UUID) gocql.UUID {
	mappedUUID, _ := gocql.UUIDFromBytes(uuid.Bytes())

	return mappedUUID
}

// addNewAddress adds new address to address table
func addNewAddress(
	tenantId, applicationId system.UUID,
	address shared.Address,
	addressId system.UUID,
	session *gocql.Session) error {
	addressDetailsCount := len(address.AddressDetails)

	errorChannel := make(chan error, addressDetailsCount*2)

	mappedTenantId := mapSystemUUIDToGocqlUUID(tenantId)
	mappedApplicationId := mapSystemUUIDToGocqlUUID(applicationId)
	mappedAddressId := mapSystemUUIDToGocqlUUID(addressId)

	var waitGroup sync.WaitGroup

	for key, value := range address.AddressDetails {
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
		return errors.New(errorMessage)
	}

	return nil

}

// addNewAddress adds new address to address table
func removeExistingAddress(
	tenantId, applicationId system.UUID,
	address shared.Address,
	addressId system.UUID,
	session *gocql.Session) error {
	addressDetailsCount := len(address.AddressDetails)

	errorChannel := make(chan error, addressDetailsCount+1)

	mappedTenantId := mapSystemUUIDToGocqlUUID(tenantId)
	mappedApplicationId := mapSystemUUIDToGocqlUUID(applicationId)
	mappedAddressId := mapSystemUUIDToGocqlUUID(addressId)

	var waitGroup sync.WaitGroup

	waitGroup.Add(1)

	go removeFromAddressTable(
		session,
		errorChannel,
		&waitGroup,
		mappedTenantId,
		mappedApplicationId,
		mappedAddressId)

	for key, _ := range address.AddressDetails {
		waitGroup.Add(1)

		go removeFromIndexByAddressKeyTable(
			session,
			errorChannel,
			&waitGroup,
			mappedTenantId,
			mappedApplicationId,
			mappedAddressId,
			key)
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
		return errors.New(errorMessage)
	}

	return nil
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

// removeFromAddressTable removes an existing address from address table using provided address unique identifier.
func removeFromAddressTable(
	session *gocql.Session,
	errorChannel chan<- error,
	waitGroup *sync.WaitGroup,
	tenantId, applicationId, addressId gocql.UUID) {

	defer waitGroup.Done()

	if err := session.Query(
		"DELETE FROM address"+
			" WHERE"+
			" tenant_id = ?"+
			" AND application_id = ?"+
			" AND address_id = ?",
		tenantId,
		applicationId,
		addressId).
		Exec(); err != nil {
		errorChannel <- err
	} else {
		errorChannel <- nil
	}
}

// removeFromIndexByAddressKeyTable removes an address key from index table.
func removeFromIndexByAddressKeyTable(
	session *gocql.Session,
	errorChannel chan<- error,
	waitGroup *sync.WaitGroup,
	tenantId, applicationId, addressId gocql.UUID,
	key string) {

	defer waitGroup.Done()

	if err := session.Query(
		"DELETE FROM address_indexed_by_address_key"+
			" WHERE"+
			" tenant_id = ? "+
			" AND application_id = ?"+
			" AND address_id = ?"+
			" AND address_key = ?",
		tenantId,
		applicationId,
		addressId,
		key).
		Exec(); err != nil {
		errorChannel <- err
	} else {
		errorChannel <- nil
	}
}
