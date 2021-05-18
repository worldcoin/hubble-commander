CREATE INDEX account_public_key_idx ON account (public_key);
CREATE INDEX batch_batch_id_idx ON batch (batch_id);
CREATE INDEX commitment_included_in_batch_idx ON commitment (included_in_batch);
CREATE INDEX transaction_base_from_state_id_idx ON transaction_base (from_state_id);
CREATE INDEX transaction_base_included_in_commitment_idx ON transaction_base (included_in_commitment);
