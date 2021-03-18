CREATE TABLE "transaction" (
    tx_hash                BYTEA PRIMARY KEY,
    from_index             NUMERIC(78) NOT NULL,
    to_index               NUMERIC(78) NOT NULL,
    amount                 NUMERIC(78) NOT NULL,
    fee                    NUMERIC(78) NOT NULL,
    nonce                  NUMERIC(78) NOT NULL,
    signature              BYTEA       NOT NULL,
    included_in_commitment BYTEA REFERENCES commitment,
    error_message          TEXT
);
