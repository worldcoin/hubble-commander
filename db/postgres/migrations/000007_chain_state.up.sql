CREATE TABLE chain_state (
    chain_id         NUMERIC(78) PRIMARY KEY,
    account_registry BYTEA   NOT NULL,
    rollup           BYTEA   NOT NULL,
    genesis_accounts JSON    NOT NULL,
    synced_block     INTEGER NOT NULL,
    deployment_block INTEGER NOT NULL
);
