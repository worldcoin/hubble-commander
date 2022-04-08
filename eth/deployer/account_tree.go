package deployer

import (
	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/storage"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type Tree struct {
	Depth uint8

	Smt *storage.StoredMerkleTree
	LeafIndexLeft uint32
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
		LeafIndexLeft: 0,
		Subtrees: subtrees,
	}
}

// lithp-TODO: check that this is the correct encoding method
func (t *Tree) RegisterAccount(pubkey *models.PublicKey) {
	bytes := pubkey.Bytes()  // abi.encodePacked(pubkey)
	leaf := crypto.Keccak256Hash(bytes)
	t.Insert(leaf)
}

func (t *Tree) Insert(hash common.Hash) {
	path := &models.MerklePath { Path: t.LeafIndexLeft, Depth: t.Depth}
	t.Smt.SetNode(path, hash)
	t.LeafIndexLeft += 1

	// This duplicates logic from AccountTree._updateSingle
	// - When _updateSingle computes the new root it needs to know all the
	//   subtrees which have a full left tree and are waiting for their right
	//   half.
	// - The least-significant 0 of our path is the left subtree which we have just
	//   completed, so we record it!
	// - If there is no least-significant 0 then we have added 2**Depth
	//   accounts, which will not happen in any of the situations this scripts
	//   expects to be called

	for {
		// search for the deepest part of our path which is a left-branch

		if path.Path & 1 == 0 {
			break
		}

		path, _ = path.Parent()

		if path.Depth == 0 {
			panic("we completely filled the tree")
		}
	}

	// Invert, because these two variables use different numbering schemes:
	// - subtrees[0] is a leaf node
	// - path.Depth = 0 is a root node
	newSubtreeDepth := t.Depth - path.Depth  // subtrees[0] is a leaf node

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
