package service

import (
	"errors"
	"strings"
	"sync"

	"github.com/gocql/gocql"
	"github.com/microbusinesses/AddressService/data/contract"
	"github.com/microbusinesses/Micro-Businesses-Core/common/diagnostics"
	"github.com/microbusinesses/Micro-Businesses-Core/system"
)

// AddressDataService provides access to add new address and update/retrieve/remove an existing address.
type AddressDataService struct {
	UUIDGeneratorService UUIDGeneratorService
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
	diagnostics.IsNotNilOrEmpty(tenantID, "tenantID", "tenantID must be provided.")
	diagnostics.IsNotNilOrEmpty(applicationID, "applicationID", "applicationID must be provided.")

	addressDetailsCount := len(address.AddressDetails)

	if addressDetailsCount == 0 {
		panic("Address does not contain any address details.")
	}

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
	diagnostics.IsNotNilOrEmpty(tenantID, "tenantID", "tenantID must be provided.")
	diagnostics.IsNotNilOrEmpty(applicationID, "applicationID", "applicationID must be provided.")
	diagnostics.IsNotNilOrEmpty(addressID, "addressID", "addressID must be provided.")

	if len(address.AddressDetails) == 0 {
		panic("Address does not contain any address details.")
	}

	err := addressDataService.Delete(tenantID, applicationID, addressID)

	if err != nil {
		return err
	}

	session, err := addressDataService.ClusterConfig.CreateSession()

	if err != nil {
		return err
	}

	defer session.Close()

	return addNewAddress(tenantID, applicationID, address, addressID, session)
}

// Read retrieves an existing address information and returns only the detail which the keys provided by the detailsKeys.
// tenantID: Mandatory. The unique identifier of the tenant owning the address.
// applicationID: Mandatory. The unique identifier of the tenant's application will be owning the address.
// addressID: Mandatory. The unique identifier of the existing address.
// detailsKeys: Mandatory. The interested address details keys to return.
// Returns either the address information or error if something goes wrong.
func (addressDataService AddressDataService) Read(tenantID, applicationID, addressID system.UUID, detailsKeys []string) (contract.Address, error) {
	diagnostics.IsNotNil(addressDataService.ClusterConfig, "addressDataService.ClusterConfig", "ClusterConfig must be provided.")
	diagnostics.IsNotNilOrEmpty(tenantID, "tenantID", "tenantID must be provided.")
	diagnostics.IsNotNilOrEmpty(applicationID, "applicationID", "applicationID must be provided.")
	diagnostics.IsNotNilOrEmpty(addressID, "addressID", "addressID must be provided.")

	detailsKeysCount := len(detailsKeys)

	if detailsKeysCount == 0 {
		panic("No address details key provided.")
	}

	session, err := addressDataService.ClusterConfig.CreateSession()

	if err != nil {
		return contract.Address{}, err
	}

	defer session.Close()

	keys := []string{}

	for _, key := range detailsKeys {
		keys = append(keys, key)
	}
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

	var key string
	var value string

	address := contract.Address{AddressDetails: make(map[string]string)}

	for iter.Scan(&key, &value) {
		address.AddressDetails[key] = value
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
	diagnostics.IsNotNilOrEmpty(tenantID, "tenantID", "tenantID must be provided.")
	diagnostics.IsNotNilOrEmpty(applicationID, "applicationID", "applicationID must be provided.")
	diagnostics.IsNotNilOrEmpty(addressID, "addressID", "addressID must be provided.")

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
			" AND address_id = ?",
		tenantID.String(),
		applicationID.String(),
		addressID.String()).Iter()

	var key string
	var value string

	address := contract.Address{AddressDetails: make(map[string]string)}

	for iter.Scan(&key, &value) {
		address.AddressDetails[key] = value
	}

	return address, nil
}

// Delete deletes an existing address information.
// tenantID: Mandatory. The unique identifier of the tenant owning the address.
// applicationID: Mandatory. The unique identifier of the tenant's application will be owning the address.
// addressID: Mandatory. The unique identifier of the existing address to remove.
// Returns error if something goes wrong.
func (addressDataService AddressDataService) Delete(tenantID, applicationID, addressID system.UUID) error {
	diagnostics.IsNotNil(addressDataService.ClusterConfig, "addressDataService.ClusterConfig", "ClusterConfig must be provided.")
	diagnostics.IsNotNilOrEmpty(tenantID, "tenantID", "tenantID must be provided.")
	diagnostics.IsNotNilOrEmpty(applicationID, "applicationID", "applicationID must be provided.")
	diagnostics.IsNotNilOrEmpty(addressID, "addressID", "addressID must be provided.")

	address, err := addressDataService.ReadAll(tenantID, applicationID, addressID)

	if err != nil {
		return err
	}

	session, err := addressDataService.ClusterConfig.CreateSession()

	if err != nil {
		return err
	}

	defer session.Close()

	return removeExistingAddress(tenantID, applicationID, address, addressID, session)
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

// addNewAddress adds new address to address table
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
