CREATE TABLE "transfer" (
    tx_hash                BYTEA PRIMARY KEY,
    to_state_id            BIGINT NOT NULL
);
