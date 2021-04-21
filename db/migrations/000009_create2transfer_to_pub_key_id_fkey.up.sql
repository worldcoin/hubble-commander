ALTER TABLE create2transfer
    ADD CONSTRAINT create2transfer_to_pub_key_id_fkey
        FOREIGN KEY (to_pub_key_id)
            REFERENCES account (pub_key_id);
