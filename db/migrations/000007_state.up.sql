-- the root hash is stored as [0], all merkle paths are prepended with 0
CREATE TABLE state_node (
    state_id    BIT VARYING(33) PRIMARY KEY,
    data_hash   BYTEA NOT NULL
);

-- this table is append only
CREATE TABLE state_leaf (
    data_hash     BYTEA PRIMARY KEY,
    pubkey_id     BIGINT      NOT NULL,
    token_index   NUMERIC(78) NOT NULL,
    balance       NUMERIC(78) NOT NULL,
    nonce         NUMERIC(78) NOT NULL
);

-- this table is append only
CREATE TABLE state_update (
    id           BIGSERIAL PRIMARY KEY,
    state_id     BIT(33) NOT NULL,
    current_hash BYTEA   NOT NULL,
    current_root BYTEA   NOT NULL,
    prev_hash    BYTEA   NOT NULL,
    prev_root    BYTEA   NOT NULL
);
