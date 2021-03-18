CREATE TABLE commitment (
    leaf_hash          BYTEA PRIMARY KEY,
    post_state_root    BYTEA       NOT NULL,
    body_hash          BYTEA       NOT NULL,
    account_tree_root  BYTEA       NOT NULL,
    combined_signature BYTEA       NOT NULL,
    fee_receiver       NUMERIC(78) NOT NULL, -- state index of tree receiver
    transactions       BYTEA       NOT NULL, -- 32 transactions encoded in compact format
    included_in_batch  BYTEA REFERENCES batch
);
