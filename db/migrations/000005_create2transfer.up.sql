CREATE TABLE "create2transfer" (
    tx_hash                BYTEA PRIMARY KEY,
    to_state_id            BIGINT NOT NULL,
    to_pubkey_id           BIGINT NOT NULL
);

