CREATE TABLE state_update (
    id           BIGSERIAL PRIMARY KEY,
    state_id     BIT(33) NOT NULL,
    current_hash BYTEA   NOT NULL,
    current_root BYTEA   NOT NULL,
    prev_hash    BYTEA   NOT NULL,
    prev_root    BYTEA   NOT NULL
);
