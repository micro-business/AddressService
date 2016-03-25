CREATE KEYSPACE address with replication = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };
CREATE TABLE address.Address(tenant_id UUID, application_id UUID, address_id UUID, address_part text, address_value text, PRIMARY KEY(tenant_id, application_id, address_id, address_part));
