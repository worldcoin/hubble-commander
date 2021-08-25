ALTER TABLE transaction_base
    DROP CONSTRAINT transaction_base_included_in_commitment_fkey;
DROP TABLE commitment;

ALTER TABLE transaction_base DROP COLUMN included_in_commitment;
ALTER TABLE transaction_base ADD COLUMN batch_id NUMERIC(78);
ALTER TABLE transaction_base ADD COLUMN index_in_batch BIGINT;
