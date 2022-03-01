CREATE TABLE IF NOT EXISTS registration_keys
(
    id BIGSERIAL NOT NULL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    claimed_at TIMESTAMP,
    key_value TEXT NOT NULL UNIQUE,
    instant_registration BOOLEAN NOT NULL DEFAULT FALSE,
);

CREATE TABLE IF NOT EXISTS unregistered_users
(
    id BIGSERIAL NOT NULL PRIMARY KEY,
    registration_key_id BIGINT REFERENCES registration_keys (id) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    approved_at TIMESTAMP,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    email TEXT NOT NULL,
);
