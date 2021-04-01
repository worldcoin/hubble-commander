CREATE TABLE batch (
    batch_hash         BYTEA PRIMARY KEY,
    batch_id           NUMERIC(78) NOT NULL,
    finalisation_block BIGINT NOT NULL
);
