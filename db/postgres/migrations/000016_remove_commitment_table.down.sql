CREATE TABLE commitment
(
    commitment_id      SERIAL PRIMARY KEY,
    type               SMALLINT NOT NULL,
    transactions       BYTEA    NOT NULL,
    fee_receiver       BIGINT   NOT NULL,
    combined_signature BYTEA    NOT NULL,
    post_state_root    BYTEA    NOT NULL,
    included_in_batch  NUMERIC(78)
);
CREATE INDEX commitment_included_in_batch_idx ON commitment (included_in_batch);
