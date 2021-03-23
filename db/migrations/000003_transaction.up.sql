CREATE TABLE "transaction" (
    tx_hash                BYTEA PRIMARY KEY,
    from_index             BIGINT NOT NULL,
    to_index               BIGINT NOT NULL,
    amount                 NUMERIC(78) NOT NULL,
    fee                    NUMERIC(78) NOT NULL,
    nonce                  NUMERIC(78) NOT NULL,
    signature              BYTEA       NOT NULL,
    included_in_commitment INTEGER REFERENCES commitment,
    error_message          TEXT
);
