# Data schema

# Performance analysis

**state_leaf** with data_hash as key (append only table holding all historical states)

**state_update** holding data_hash of the previous state_leaf

VS

**state_leaf** with state_id as key (table holding only current states)

**state_update** holding the whole previous state_leaf object

### data_hash key

**O(1) / O(log(n))**
AddStateLeaf
GetStateLeafByHash
GetStateLeafByPath
GetUserStateByID (duplicate of GetStateLeafByPath)

AddStateUpdate
DeleteStateUpdate
GetStateUpdateByRootHash

**O(sth)**
GetUserStatesByPublicKey
- GetAccounts = O(n) - could be O(log(n))
- get all state_nodes = O(n) - n is num of nodes in merkle tree
- A LOT
GetUnusedPubKeyID -> should take token index as well
- GetAccounts = O(n) - could be O(log(n))
- get all state_leaves = O(n) - n is num of leaves in merkle tree
- A LOT
GetTransfersByPublicKey
- GetAccounts = O(n) - could be O(log(n))
- get all historical state_leaves with given pubkeyid using index - O(log(n))
- for each possibly historical state_leaf get from state_node using index - O(nlog(n))

**Easily cached**
GetUserStateByPubKeyIDAndTokenIndex

### state_id key

**O(1) / O(log(n))**
AddStateLeaf
GetStateLeafByPath
GetUserStatesByPublicKey
- GetAccounts = O(n) - could be O(log(n))
- index on pubkeyID
- get by StateID
GetUserStateByID (duplicate of GetStateLeafByPath)
GetUnusedPubKeyID
- GetAccounts = O(n) - could be O(log(n))
- for each PubKeyID use {token_index, pub_key_id} => {state_id} index
GetTransfersByPublicKey
- GetAccounts = O(n) - could be O(log(n))
- get stateIDs using index on pub_key_id (probably using Badger directly)

AddStateUpdate
DeleteStateUpdate
GetStateUpdateByRootHash

**Easily cached**
GetUserStateByPubKeyIDAndTokenIndex

**Not necessary**
GetStateLeafByHash

```sql
CREATE TABLE state_leaf (
    state_id    BYTEA PRIMARY KEY,
    pub_key_id  BIGINT,
    token_index NUMERIC(78) NOT NULL,
    balance     NUMERIC(78) NOT NULL,
    nonce       NUMERIC(78) NOT NULL
);
```
```ignore
GetByStateID: (stateID) -> state_leaf
GetStatesByPubkeyId: (pubkeyId) -> state_leaf[]
# GetStatesByPubkeyIdAndToken: (pubkeyId, token) -> state_leaf[]
Set: (state_leaf)

Mappings

state_leaf:<state_id> -> { state_id, pub_key_id, token_index, balance, nonce }

# state_leaf:<state_id>:pub_key_id -> <pub_key_id>
# state_leaf:<state_id>:token_index -> <token_index>
# state_leaf:<state_id>:balance -> <balance>
# state_leaf:<state_id>:nonce -> <nonce>

index:state_leaf:pubkeyId:<pubkeyId> -> [<state_id_0>, <state_id_1>]
...

state_node:<merkle_path> -> {merkle_path, hash}

State update:
id, batch_id, prev_root, current_root, state_id, prev_leaf_obj

State node:

StateTree.Set(index, node):

  leaf_hash = hash(node)

  state_leafs.insert(leaf_hash, node)

  for path in path_to_root(index): 

    state_nodes.set(path, hash)

  state_update.insert(...)

StateTree.Set(index, node):
  prev_node = get("state_leaf:<leaf_index>")

  set("state_leaf:<leaf_index>", node)
  append("index:state_leaf:pubkeyId:<node.pubkeyId>", index)

  leaf_hash = hash(node)

  for path in path_to_root(index): 

    set("state_node:<path>", hash)

  set("state_update:<seq_id>", { pre_root_hash, post_root_hash, index, prev_node })
```
