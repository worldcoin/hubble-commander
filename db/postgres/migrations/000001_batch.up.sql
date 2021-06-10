CREATE TABLE batch (
    batch_id           SERIAL PRIMARY KEY,
    type               SMALLINT    NOT NULL,
    transaction_hash   BYTEA       NOT NULL,
    batch_hash         BYTEA,
    batch_number       NUMERIC(78) NOT NULL,
    finalisation_block BIGINT,
    account_tree_root  BYTEA
);
