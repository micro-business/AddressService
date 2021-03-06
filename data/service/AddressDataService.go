package service

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/gocql/gocql"
	"github.com/micro-business/AddressService/data/contract"
	"github.com/micro-business/Micro-Business-Core/common/diagnostics"
	"github.com/micro-business/Micro-Business-Core/system"
)

// AddressDataService provides access to add new address and update/retrieve/remove an existing address.
type AddressDataService struct {
	UUIDGeneratorService system.UUIDGeneratorService
	ClusterConfig        *gocql.ClusterConfig
}

// Create creates a new address.
// tenantID: Mandatory. The unique identifier of the tenant owning the address.
// applicationID: Mandatory. The unique identifier of the tenant's application will be owning the address.
// address: Mandatory. The reference to the new address information.
// Returns either the unique identifier of the new address or error if something goes wrong.
func (addressDataService AddressDataService) Create(tenantID, applicationID system.UUID, address contract.Address) (system.UUID, error) {
	diagnostics.IsNotNil(addressDataService.UUIDGeneratorService, "addressDataService.UUIDGeneratorService", "UUIDGeneratorService must be provided.")
	diagnostics.IsNotNil(addressDataService.ClusterConfig, "addressDataService.ClusterConfig", "ClusterConfig must be provided.")

	addressID, err := addressDataService.UUIDGeneratorService.GenerateRandomUUID()

	if err != nil {
		return system.EmptyUUID, err
	}

	session, err := addressDataService.ClusterConfig.CreateSession()

	if err != nil {
		return system.EmptyUUID, err
	}

	defer session.Close()

	if err = addNewAddress(tenantID, applicationID, address, addressID, session); err != nil {
		return system.EmptyUUID, err
	}

	return addressID, nil
}

// Update updates an existing address.
// tenantID: Mandatory. The unique identifier of the tenant owning the address.
// applicationID: Mandatory. The unique identifier of the tenant's application will be owning the address.
// addressID: Mandatory. The unique identifier of the existing address.
// address: Mandatory. The reeference to the updated address information.
// Returns error if something goes wrong.
func (addressDataService AddressDataService) Update(tenantID, applicationID, addressID system.UUID, address contract.Address) error {
	diagnostics.IsNotNil(addressDataService.ClusterConfig, "addressDataService.ClusterConfig", "ClusterConfig must be provided.")

	session, err := addressDataService.ClusterConfig.CreateSession()

	if err != nil {
		return err
	}

	defer session.Close()

	if !doesAddressExist(tenantID, applicationID, addressID, session) {
		return fmt.Errorf("Address not found. Address ID: %s", addressID.String())
	}

	if err := deleteExistingAddress(tenantID, applicationID, addressID, session); err != nil {
		return err
	}

	return addNewAddress(tenantID, applicationID, address, addressID, session)
}

// Read retrieves an existing address information and returns only the detail which the keys provided by the keys.
// tenantID: Mandatory. The unique identifier of the tenant owning the address.
// applicationID: Mandatory. The unique identifier of the tenant's application will be owning the address.
// addressID: Mandatory. The unique identifier of the existing address.
// keys: Mandatory. The interested address details keys to return.
// Returns either the address information or error if something goes wrong.
func (addressDataService AddressDataService) Read(tenantID, applicationID, addressID system.UUID, keys []string) (contract.Address, error) {
	diagnostics.IsNotNil(addressDataService.ClusterConfig, "addressDataService.ClusterConfig", "ClusterConfig must be provided.")

	session, err := addressDataService.ClusterConfig.CreateSession()

	if err != nil {
		return contract.Address{}, err
	}

	defer session.Close()

	iter := session.Query(
		"SELECT address_key, address_value"+
			" FROM address"+
			" WHERE"+
			" tenant_id = ?"+
			" AND application_id = ?"+
			" AND address_id = ?"+
			" AND address_key IN "+
			" ('"+strings.Join(keys, "','")+"')",
		tenantID.String(),
		applicationID.String(),
		addressID.String()).Iter()

	defer iter.Close()

	var key string
	var value string

	address := contract.Address{AddressDetails: make(map[string]string)}

	for iter.Scan(&key, &value) {
		address.AddressDetails[key] = value
	}

	if len(address.AddressDetails) == 0 {
		return contract.Address{}, fmt.Errorf("Address not found. Address ID: %s", addressID.String())
	}

	return address, nil
}

// ReadAll retrieves an existing address information and returns all the detail of it.
// tenantID: Mandatory. The unique identifier of the tenant owning the address.
// applicationID: Mandatory. The unique identifier of the tenant's application will be owning the address.
// addressID: Mandatory. The unique identifier of the existing address.
// Returns either the address information or error if something goes wrong.
func (addressDataService AddressDataService) ReadAll(tenantID, applicationID, addressID system.UUID) (contract.Address, error) {
	diagnostics.IsNotNil(addressDataService.ClusterConfig, "addressDataService.ClusterConfig", "ClusterConfig must be provided.")

	session, err := addressDataService.ClusterConfig.CreateSession()

	if err != nil {
		return contract.Address{}, err
	}

	defer session.Close()

	return readAllAddressDetails(tenantID, applicationID, addressID, session)
}

// Delete deletes an existing address information.
// tenantID: Mandatory. The unique identifier of the tenant owning the address.
// applicationID: Mandatory. The unique identifier of the tenant's application will be owning the address.
// addressID: Mandatory. The unique identifier of the existing address to remove.
// Returns error if something goes wrong.
func (addressDataService AddressDataService) Delete(tenantID, applicationID, addressID system.UUID) error {
	diagnostics.IsNotNil(addressDataService.ClusterConfig, "addressDataService.ClusterConfig", "ClusterConfig must be provided.")

	session, err := addressDataService.ClusterConfig.CreateSession()

	if err != nil {
		return err
	}

	defer session.Close()

	if !doesAddressExist(tenantID, applicationID, addressID, session) {
		return fmt.Errorf("Address not found. Address ID: %s", addressID.String())
	}

	return deleteExistingAddress(tenantID, applicationID, addressID, session)
}

// mapSystemUUIDToGocqlUUID maps the system type UUID to gocql UUID type
func mapSystemUUIDToGocqlUUID(uuid system.UUID) gocql.UUID {
	mappedUUID, _ := gocql.UUIDFromBytes(uuid.Bytes())

	return mappedUUID
}

// addNewAddress adds new address to address table
func addNewAddress(
	tenantID, applicationID system.UUID,
	address contract.Address,
	addressID system.UUID,
	session *gocql.Session) error {
	addressDetailsCount := len(address.AddressDetails)

	errorChannel := make(chan error, addressDetailsCount*2)

	mappedTenantID := mapSystemUUIDToGocqlUUID(tenantID)
	mappedApplicationID := mapSystemUUIDToGocqlUUID(applicationID)
	mappedAddressID := mapSystemUUIDToGocqlUUID(addressID)

	var waitGroup sync.WaitGroup

	for key, value := range address.AddressDetails {
		waitGroup.Add(1)

		go addToAddressTable(
			session,
			errorChannel,
			&waitGroup,
			mappedTenantID,
			mappedApplicationID,
			mappedAddressID,
			key,
			value)

		waitGroup.Add(1)

		go addToAddressIndexByAddressKeyTable(
			session,
			errorChannel,
			&waitGroup,
			mappedTenantID,
			mappedApplicationID,
			mappedAddressID,
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

// removeExistingAddress adds new address to address table
func removeExistingAddress(
	tenantID, applicationID system.UUID,
	address contract.Address,
	addressID system.UUID,
	session *gocql.Session) error {
	addressDetailsCount := len(address.AddressDetails)

	errorChannel := make(chan error, addressDetailsCount+1)

	mappedTenantID := mapSystemUUIDToGocqlUUID(tenantID)
	mappedApplicationID := mapSystemUUIDToGocqlUUID(applicationID)
	mappedAddressID := mapSystemUUIDToGocqlUUID(addressID)

	var waitGroup sync.WaitGroup

	waitGroup.Add(1)

	go removeFromAddressTable(
		session,
		errorChannel,
		&waitGroup,
		mappedTenantID,
		mappedApplicationID,
		mappedAddressID)

	for key := range address.AddressDetails {
		waitGroup.Add(1)

		go removeFromIndexByAddressKeyTable(
			session,
			errorChannel,
			&waitGroup,
			mappedTenantID,
			mappedApplicationID,
			mappedAddressID,
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
	tenantID, applicationID, addressID gocql.UUID,
	key, value string) {

	defer waitGroup.Done()

	if err := session.Query(
		"INSERT INTO address"+
			" (tenant_id, application_id, address_id, address_key, address_value)"+
			" VALUES(?, ?, ?, ?, ?)",
		tenantID,
		applicationID,
		addressID,
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
	tenantID, applicationID, addressID gocql.UUID,
	key, value string) {

	defer waitGroup.Done()

	if err := session.Query(
		"INSERT INTO address_indexed_by_address_key"+
			" (tenant_id, application_id, address_id, address_key, address_value)"+
			" VALUES(?, ?, ?, ?, ?)",
		tenantID,
		applicationID,
		addressID,
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
	tenantID, applicationID, addressID gocql.UUID) {

	defer waitGroup.Done()

	if err := session.Query(
		"DELETE FROM address"+
			" WHERE"+
			" tenant_id = ?"+
			" AND application_id = ?"+
			" AND address_id = ?",
		tenantID,
		applicationID,
		addressID).
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
	tenantID, applicationID, addressID gocql.UUID,
	key string) {

	defer waitGroup.Done()

	if err := session.Query(
		"DELETE FROM address_indexed_by_address_key"+
			" WHERE"+
			" tenant_id = ? "+
			" AND application_id = ?"+
			" AND address_id = ?"+
			" AND address_key = ?",
		tenantID,
		applicationID,
		addressID,
		key).
		Exec(); err != nil {
		errorChannel <- err
	} else {
		errorChannel <- nil
	}
}

// doesAddressExist checks whether the provided addressID exists in database
func doesAddressExist(tenantID, applicationID, addressID system.UUID, session *gocql.Session) bool {
	iter := session.Query(
		"SELECT address_key"+
			" FROM address"+
			" WHERE"+
			" tenant_id = ?"+
			" AND application_id = ?"+
			" AND address_id = ?"+
			" LIMIT 1",
		tenantID.String(),
		applicationID.String(),
		addressID.String()).Iter()

	defer iter.Close()

	var addressKey string

	return iter.Scan(&addressKey)
}

func deleteExistingAddress(tenantID, applicationID, addressID system.UUID, session *gocql.Session) error {
	address, err := readAllAddressDetails(tenantID, applicationID, addressID, session)

	if err != nil {
		return err
	}

	return removeExistingAddress(tenantID, applicationID, address, addressID, session)
}

func readAllAddressDetails(tenantID, applicationID, addressID system.UUID, session *gocql.Session) (contract.Address, error) {
	iter := session.Query(
		"SELECT address_key, address_value"+
			" FROM address"+
			" WHERE"+
			" tenant_id = ?"+
			" AND application_id = ?"+
			" AND address_id = ?",
		tenantID.String(),
		applicationID.String(),
		addressID.String()).Iter()

	defer iter.Close()

	var key string
	var value string

	address := contract.Address{AddressDetails: make(map[string]string)}

	for iter.Scan(&key, &value) {
		address.AddressDetails[key] = value
	}

	if len(address.AddressDetails) == 0 {
		return contract.Address{}, fmt.Errorf("Address not found. Address ID: %s", addressID.String())
	}

	return address, nil
}
