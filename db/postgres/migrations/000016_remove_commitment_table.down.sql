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

ALTER TABLE transaction_base ADD COLUMN included_in_commitment INTEGER;
CREATE INDEX transaction_base_included_in_commitment_idx ON transaction_base (included_in_commitment);
ALTER TABLE transaction_base DROP COLUMN batch_id;
ALTER TABLE transaction_base DROP COLUMN index_in_batch;
