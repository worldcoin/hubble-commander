CREATE TABLE transfer (
    tx_hash     BYTEA PRIMARY KEY REFERENCES transaction_base,
    to_state_id BIGINT NOT NULL
);
