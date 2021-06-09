ALTER TABLE batch DROP account_tree_root;
ALTER TABLE commitment ADD COLUMN account_tree_root BYTEA;
