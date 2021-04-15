CREATE TABLE "transaction_base" (
    tx_hash                BYTEA PRIMARY KEY,
    tx_type                SMALLINT NOT NULL,
    from_state_id          BIGINT NOT NULL,
    amount                 BIGINT NOT NULL,
    fee                    BIGINT NOT NULL,
    nonce                  BIGINT NOT NULL,
    signature              BYTEA NOT NULL,
    included_in_commitment INTEGER REFERENCES commitment,
    error_message          TEXT
);
