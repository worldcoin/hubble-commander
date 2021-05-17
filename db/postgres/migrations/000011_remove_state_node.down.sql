CREATE TABLE state_node (
    merkle_path BIT VARYING(33) PRIMARY KEY,
    data_hash   BYTEA NOT NULL
);
