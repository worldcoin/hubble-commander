CREATE TABLE "create2transfer" (
    tx_hash                BYTEA PRIMARY KEY REFERENCES transaction_base,
    to_state_id            BIGINT,
    to_pubkey_id           BIGINT NOT NULL
);

