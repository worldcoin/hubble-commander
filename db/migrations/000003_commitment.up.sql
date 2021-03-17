CREATE TABLE commitment (
    leaf_hash          BYTEA PRIMARY KEY,
    post_state_root    BYTEA,
    body_hash          BYTEA,
    account_tree_root  BYTEA,
    combined_signature BYTEA,
    fee_receiver       NUMERIC(78), -- state index of tree receiver
    transactions       BYTEA,       -- 32 transactions encoded in compact format
    included_in_batch  BYTEA
);
