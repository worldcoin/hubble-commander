ALTER TABLE commitment
    DROP CONSTRAINT commitment_included_in_batch_fkey;
DROP TABLE batch;
