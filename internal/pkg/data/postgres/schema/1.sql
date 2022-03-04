/* Dropping tables is only for dev */
DROP TABLE IF EXISTS schema_version;
DROP TABLE IF EXISTS unregistered_users;
DROP TABLE IF EXISTS registration_keys;
/* ******************************* */

CREATE TABLE schema_version
(
    version_number BIGSERIAL NOT NULL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO schema_version (created_at) VALUES (CURRENT_TIMESTAMP);

CREATE TABLE registration_keys
(
    id BIGSERIAL NOT NULL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    max_claims SMALLINT NOT NULL DEFAULT 1,
    key_value TEXT NOT NULL UNIQUE,
    instant_registration BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE unregistered_users
(
    id BIGSERIAL NOT NULL PRIMARY KEY,
    registration_key_id BIGINT REFERENCES registration_keys (id) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    email TEXT NOT NULL
);
