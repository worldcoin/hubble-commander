CREATE TABLE batch (
    batch_id           NUMERIC(78) PRIMARY KEY,
    type               SMALLINT     NOT NULL,
    transaction_hash   BYTEA UNIQUE NOT NULL,
    batch_hash         BYTEA,
    finalisation_block BIGINT,
    account_tree_root  BYTEA,
    prev_state_root    BYTEA,
    submission_time    TIMESTAMP
);
