CREATE KEYSPACE address with replication = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };
CREATE TABLE address.address(tenant_id UUID, application_id UUID, address_id UUID, address_key text, address_value text, PRIMARY KEY(tenant_id, application_id, address_id, address_key));
CREATE TABLE address.address_indexed_by_address_key(tenant_id UUID, application_id UUID, address_id UUID, address_key text, address_value text, PRIMARY KEY(tenant_id, application_id, address_key, address_id));
