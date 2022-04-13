package deployment

import (
	"fmt"

	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/storage"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// This file duplicates logic found in AccountTree.sol and BLSAccountRegistry.sol
// It allows us to simulate the result of adding many leaves to AccountTree.sol so
// that during deploys we can initialize our AccountTree.sol with a pre-computed state.

func get(smt *storage.StoredMerkleTree, path models.MerklePath) common.Hash {
	// throw away the error because this file only ever looks for paths which exist
	// also: we use an in-memory database which vastly reduces the number of errors
	//   lookups might experience

	res, err := smt.Get(path)

	if err != nil {
		panic(fmt.Errorf("failed to lookup node: %w", err))
	}

	return res.DataHash
}

type Tree struct {
	Depth uint8

	Smt          *storage.StoredMerkleTree
	AccountCount uint32

	// This is not a witness, this is an accumulator which allows AccountTree.sol to
	// correctly update the root node. It is possible to compute the new root using
	// constant space because AccountTree.sol only inserts leaves into the leftmost
	// open slot.
	Subtrees []common.Hash
}

func NewTree(depth uint8) *Tree {
	if depth < 2 {
		panic("tree must have a depth of at least 2")
	}

	db, err := db.NewInMemoryDatabase()
	if err != nil {
		panic(fmt.Errorf("err creating db: %w", err))
	}

	database := &storage.Database{
		Badger: db,
	}

	smt := storage.NewStoredMerkleTree("tree", database, depth)

	// For an empty tree all leaves have the same witness. This witness is the
	// initial value of Subtrees.
	var subtrees []common.Hash
	path := &models.MerklePath{Path: 0, Depth: depth}
	for {
		node := get(smt, *path)
		subtrees = append(subtrees, node)

		path, err = path.Parent()
		if err != nil {
			// unreachable, only happens when path was at depth 0
			panic(err)
		}

		if path.Depth == 1 {
			break
		}
	}

	return &Tree{
		Depth:        depth,
		Smt:          smt,
		AccountCount: 0,
		Subtrees:     subtrees,
	}
}

func (t *Tree) RegisterAccount(pubkey *models.PublicKey) {
	// duplicates BLSAccountRegistry.sol:register(pubkey)
	//   leaf = keccak256(abi.encodePacked(pubkey))

	bytes := pubkey.Bytes()
	leaf := crypto.Keccak256Hash(bytes)
	t.Insert(leaf)
}

func (t *Tree) Insert(hash common.Hash) {
	path := &models.MerklePath{Path: t.AccountCount, Depth: t.Depth}
	t.Smt.SetNode(path, hash)

	// This duplicates AccountTree.sol:_updateSingle
	// - Accounts are inserted into the tree from left to right, `AccountCount` is the
	//   location of our new leaf. When interpreted as a bitstring it becomes a path:
	//   0 means left and 1 means right
	// - When _updateSingle computes the new root it starts at the leaf which was
	//   just added and works up.
	// - As it works up, if our new account is on the left subtree for a given depth
	//   then we know the right subtree is empty and can easily lookup its hash.
	// - If our new account is in the right subtree for a given depth then a previous
	//   call to _updateSingle has already finished creating the left subtree for the
	//   depth. If we have remembered it then we can compute the new parent hash.
	// - So, to compute the new root we only need to remember the hash of each full
	//   left subtree.
	// - Unless we are inserting the final leaf, inserting a leaf always completes a
	//   single left subtree. The least-significant 0 of our path is the left subtree
	//   we are completing. If we save it to t.Subtrees then we are done!

	for {
		// search for the deepest part of our path which is a left-branch

		if path.Path&1 == 0 {
			break
		}

		path, _ = path.Parent()

		if path.Depth == 0 {
			// If there is no least-significant 0 then we have `Insert`ed
			// 2**Depth accounts.
			panic("we completely filled the tree")
		}
	}

	// Invert, because these two variables use different numbering schemes:
	// - subtrees[0] is a leaf node
	// - path.Depth = 0 is a root node
	newSubtreeDepth := t.Depth - path.Depth

	subtreeRoot := get(t.Smt, *path)
	t.Subtrees[newSubtreeDepth] = subtreeRoot
	t.AccountCount += 1
}

func (t *Tree) LeftRoot() common.Hash {
	return get(t.Smt, models.MerklePath{Path: 0, Depth: 1})
}
