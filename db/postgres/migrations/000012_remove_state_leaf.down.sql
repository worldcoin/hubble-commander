CREATE TABLE state_leaf (
    data_hash   BYTEA PRIMARY KEY,
    pub_key_id  BIGINT REFERENCES account,
    token_index NUMERIC(78) NOT NULL,
    balance     NUMERIC(78) NOT NULL,
    nonce       NUMERIC(78) NOT NULL
);
