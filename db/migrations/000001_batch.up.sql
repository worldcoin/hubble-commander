CREATE TABLE batch (
    batch_hash         BYTEA PRIMARY KEY,
    batch_id           NUMERIC(78),
    finalisation_block NUMERIC(78)
);
