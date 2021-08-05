CREATE TABLE account (
    pub_key_id BIGINT PRIMARY KEY,
    public_key BYTEA NOT NULL
);
CREATE INDEX account_public_key_idx ON account (public_key);
