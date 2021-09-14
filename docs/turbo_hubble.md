# TurboHubble

[Data schema](https://www.notion.so/Data-schema-e16f81e2c7a84e9dbffe6391eae0f1aa)

- Storing state tree nodes using [Badger](https://github.com/dgraph-io/badger)
    - TurboGeth uses [Bolt](https://github.com/boltdb/bolt) (Buckets of B+ Trees)
    - Geth uses LevelDB (LSM Trees)
    - RocksDB is also LSM
    - [Badger](https://github.com/dgraph-io/badger) provides the [best of both worlds](https://medium.com/@giulio.rebuffo/turbo-geth-whats-different-the-database-5916e8ec834b)
    - **Badger comes with atomicity**
- Parallel commitment creation
    - We can apply all transactions that don't share the same sender or receiver within a commitment in parallel (we can also sort txns like this to get this property)
- Signature aggregation can run in parallel to next commitment creation
- SendTransaction handler can run in parallel:
    - Signature verification
    - Encoding of transactions
- Syncing mempool to database only in regular intervals
    - Keep in memory, no continuous reads or writes from/to disk

