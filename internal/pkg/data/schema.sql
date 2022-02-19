CREATE TABLE IF NOT EXISTS registration_keys
(
    id BIGSERIAL NOT NULL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    claimed_at TIMESTAMP,
    key_value TEXT NOT NULL UNIQUE,
);

CREATE TABLE IF NOT EXISTS unregistered_users
(
    id BIGSERIAL NOT NULL PRIMARY KEY,
    registration_key_id BIGINT REFERENCES registration_keys (id) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    first_name BYTEA NOT NULL,
    last_name BYTEA NOT NULL,
    email BYTEA NOT NULL,
    password BYTEA NOT NULL,
);
