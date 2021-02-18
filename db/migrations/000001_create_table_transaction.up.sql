CREATE TABLE "transaction" (
    tx_hash    bytea,
    from_index NUMERIC(78),
    to_index   NUMERIC(78),
    amount    NUMERIC(78),
    fee       NUMERIC(78),
    nonce     NUMERIC(78),
    "signature" bytea
);
