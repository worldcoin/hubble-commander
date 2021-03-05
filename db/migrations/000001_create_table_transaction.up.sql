CREATE TABLE "transaction" (
    tx_hash      BYTEA primary key,
    from_index   NUMERIC(78),
    to_index     NUMERIC(78),
    amount       NUMERIC(78),
    fee          NUMERIC(78),
    nonce        NUMERIC(78),
    "signature"  BYTEA,
    included_in_commitment BYTEA
);
