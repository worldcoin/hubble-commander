ALTER TABLE transaction_base
    DROP CONSTRAINT transaction_base_included_in_commitment_fkey;
DROP TABLE commitment;
