CREATE TABLE "transfer" (
    tx_hash                BYTEA PRIMARY KEY,
    tx_type                SMALLINT NOT NULL,
    to_state_id            BIGINT NOT NULL,
    signature              BYTEA NOT NULL,
    included_in_commitment INTEGER REFERENCES commitment,
    error_message          TEXT
);
