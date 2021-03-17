-- the root hash is stored as [0], all merkle paths are prepended with 0
CREATE TABLE "state_node" (
    merkle_path  BIT VARYING(33) PRIMARY KEY,
    data_hash    BYTEA
);

-- this table is append only
CREATE TABLE "state_leaf" (
    data_hash      BYTEA PRIMARY KEY,
    account_index  BIGINT,
    token_index    NUMERIC(78),
    balance        NUMERIC(78),
    nonce          NUMERIC(78)
);

-- this table is append only
CREATE TABLE "state_update" (
    id            BIGSERIAL PRIMARY KEY,
    merkle_path   BIT(33),
    current_hash  BYTEA,
    current_root  BYTEA,
    prev_hash     BYTEA,
    prev_root     BYTEA
);
