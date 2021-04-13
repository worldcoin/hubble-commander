CREATE TABLE "transaction_base" (
    tx_hash                BYTEA PRIMARY KEY,
    tx_type                SMALLINT NOT NULL,
    from_id                BIGINT NOT NULL,
    amount                 BIGINT NOT NULL,
    fee                    BIGINT NOT NULL,
    nonce                  BIGINT NOT NULL
);
