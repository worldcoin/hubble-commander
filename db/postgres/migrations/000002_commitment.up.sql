CREATE TABLE commitment (
    commitment_id      SERIAL PRIMARY KEY,
    type               SMALLINT NOT NULL,
    transactions       BYTEA    NOT NULL, -- 32 transactions encoded in compact format
    fee_receiver       BIGINT   NOT NULL, -- state id of fee receiver
    combined_signature BYTEA    NOT NULL,
    post_state_root    BYTEA    NOT NULL,
    included_in_batch  NUMERIC(78) REFERENCES batch
);
