package deployer

import (
	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/storage"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// This file duplicates logic found in AccountTree.sol and BLSAccountRegistry.sol
// It allows us to simulate the result of adding many leaves to AccountTree.sol so
// that during deploys we can initialize our AccountTree.sol with a pre-computed state.

type Tree struct {
	Depth uint8

	Smt *storage.StoredMerkleTree
	AccountCount uint32

	// This is not a witness, this is an accumulator which allows AccountTree.sol to
	// correctly update the root node. It is possible to compute the new root using
	// constant space because AccountTree.sol only inserts leaves into the leftmost
	// open slot.
	Subtrees []common.Hash
}

func NewTree(depth uint8) *Tree {
	db, err := db.NewInMemoryDatabase()
	if err != nil {
		panic("err creating db")
	}

	database := &storage.Database{
		Badger: db,
	}

	smt := storage.NewStoredMerkleTree("tree", database, depth)

	var subtrees []common.Hash
	path := &models.MerklePath { Path: 0, Depth: depth}

	for {
		node, err := smt.Get(*path)
		if err != nil {
			panic("failed to lookup node")
		}

		subtrees = append(subtrees, node.DataHash)

		path, err = path.Parent()
		if err != nil {
			panic(err)
		}

		if path.Depth == 1 {
				break
		}
	}

	return &Tree{
		Depth: depth,
		Smt: smt,
		AccountCount: 0,
		Subtrees: subtrees,
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
	path := &models.MerklePath { Path: t.AccountCount, Depth: t.Depth}
	t.Smt.SetNode(path, hash)
	t.AccountCount += 1

	// This duplicates AccountTree.sol:_updateSingle
	// - When _updateSingle computes the new root it needs to know all the
	//   subtrees which have a full left tree and are waiting for their right
	//   half.
	// - The least-significant 0 of our path is the left subtree which we have just
	//   completed, so we record it!

	for {
		// search for the deepest part of our path which is a left-branch

		if path.Path & 1 == 0 {
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

	subtreeRoot, err := t.Smt.Get(*path)
	if err != nil {
		panic("failed to lookup node")
	}

	t.Subtrees[newSubtreeDepth] = subtreeRoot.DataHash
}

func (t *Tree) Get(path models.MerklePath) common.Hash {
	res, err := t.Smt.Get(path)

	if err != nil {
		panic("failed to lookup node")
	}

	return res.DataHash
}

func (t *Tree) LeftRoot() common.Hash {
	return t.Get(models.MerklePath { Path: 0, Depth: 1 })
}
