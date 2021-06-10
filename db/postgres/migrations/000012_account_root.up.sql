ALTER TABLE commitment DROP account_tree_root;
ALTER TABLE batch ADD COLUMN account_tree_root BYTEA;
