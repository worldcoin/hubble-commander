-- the root hash is stored with empty merkle_path
CREATE TABLE "state_node" (
    merkle_path  BIT VARYING(32) primary key,
    data_hash    BYTEA
);

-- this table is append only
CREATE TABLE "state_leaf" (
    data_hash      BYTEA primary key,
    account_index  NUMERIC(78),
    token_index    NUMERIC(78),
    balance        NUMERIC(78),
    nonce          NUMERIC(78)
);

CREATE TABLE "state_update" (
    id            BIGSERIAL primary key,
    merkle_path   BIT(32),
    current_hash  BYTEA,
    current_root  BYTEA,
    prev_hash     BYTEA,
    prev_root     BYTEA
);
