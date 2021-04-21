CREATE TABLE transaction_base (
    tx_hash                BYTEA PRIMARY KEY,
    tx_type                SMALLINT    NOT NULL,
    from_state_id          BIGINT      NOT NULL,
    amount                 NUMERIC(78) NOT NULL,
    fee                    NUMERIC(78) NOT NULL,
    nonce                  NUMERIC(78) NOT NULL,
    signature              BYTEA       NOT NULL,
    included_in_commitment INTEGER REFERENCES commitment,
    error_message          TEXT
);
