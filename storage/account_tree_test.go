package storage

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type AccountTreeTestSuite struct {
	*require.Assertions
	suite.Suite
	storage *TestStorage
	tree    *AccountTree
	leaf    *models.AccountLeaf
}

func (s *AccountTreeTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *AccountTreeTestSuite) SetupTest() {
	var err error
	s.storage, err = NewTestStorageWithoutPostgres()
	s.NoError(err)
	s.tree = NewAccountTree(s.storage.Storage)

	s.leaf = &models.AccountLeaf{
		PubKeyID: 0,
		PublicKey: models.MakePublicKeyFromInts([4]*big.Int{
			big.NewInt(91237583),
			big.NewInt(43129487),
			big.NewInt(54351448),
			big.NewInt(12347495),
		}),
	}
}

func (s *AccountTreeTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *AccountTreeTestSuite) TestSetSingle_StoresAccountLeafRecord() {
	err := s.tree.SetSingle(s.leaf)
	s.NoError(err)

	actualLeaf, err := s.storage.GetAccountLeaf(s.leaf.PubKeyID)
	s.NoError(err)
	s.Equal(s.leaf, actualLeaf)
}

func (s *AccountTreeTestSuite) TestSetSingle_RootIsDifferentAfterSet() {
	leaf0 := &models.AccountLeaf{
		PubKeyID:  0,
		PublicKey: s.randomPublicKey(),
	}

	leaf1 := &models.AccountLeaf{
		PubKeyID:  1,
		PublicKey: s.randomPublicKey(),
	}

	err := s.tree.SetSingle(leaf0)
	s.NoError(err)

	accountRootAfter0, err := s.tree.Root()
	s.NoError(err)

	err = s.tree.SetSingle(leaf1)
	s.NoError(err)

	accountRootAfter1, err := s.tree.Root()
	s.NoError(err)

	s.NotEqual(accountRootAfter0, accountRootAfter1)
}

func (s *AccountTreeTestSuite) TestSetSingle_StoresLeafMerkleTreeNodeRecord() {
	err := s.tree.SetSingle(s.leaf)
	s.NoError(err)

	expectedNode := &models.MerkleTreeNode{
		MerklePath: models.MerklePath{
			Path:  0,
			Depth: AccountTreeDepth,
		},
		DataHash: crypto.Keccak256Hash(s.leaf.PublicKey.Bytes()),
	}

	node, err := s.tree.merkleTree.Get(expectedNode.MerklePath)
	s.NoError(err)
	s.Equal(expectedNode, node)
}

func (s *AccountTreeTestSuite) TestSetSingle_UpdatesRootMerkleTreeNodeRecord() {
	err := s.tree.SetSingle(s.leaf)
	s.NoError(err)

	root, err := s.tree.Root()
	s.NoError(err)
	s.Equal(common.HexToHash("0x6e082faf2fd8ce5accb1e08a15061f2c443ea5e9cb42d493050275d644bb51b9"), *root)
}

func (s *AccountTreeTestSuite) TestSetSingle_CalculatesCorrectRootForLeafOfId1() {
	s.leaf.PubKeyID = 1
	err := s.tree.SetSingle(s.leaf)
	s.NoError(err)

	root, err := s.tree.Root()
	s.NoError(err)
	s.Equal(common.HexToHash("0xcd6164584f02a9c4c9f88c2613d7ff2b709e0951369f9bd28528712e3fa96daa"), *root)
}

func (s *AccountTreeTestSuite) TestSetSingle_CalculatesCorrectRootForTwoLeaves() {
	err := s.tree.SetSingle(s.leaf)
	s.NoError(err)

	leaf1 := &models.AccountLeaf{
		PubKeyID: 1,
		PublicKey: models.MakePublicKeyFromInts([4]*big.Int{
			big.NewInt(83457234),
			big.NewInt(74928472),
			big.NewInt(11920493),
			big.NewInt(40048372),
		}),
	}
	err = s.tree.SetSingle(leaf1)
	s.NoError(err)

	root, err := s.tree.Root()
	s.NoError(err)
	s.Equal(common.HexToHash("0x3a7a7ff21991ccfcbf8a4580862def7c498253ad398e967f270ff421db1d4833"), *root)
}

func (s *AccountTreeTestSuite) TestSetSingle_ThrowsOnSettingAlreadySetLeaf() {
	err := s.tree.SetSingle(s.leaf)
	s.NoError(err)

	s.leaf.PublicKey = s.randomPublicKey()
	err = s.tree.SetSingle(s.leaf)
	s.ErrorIs(err, ErrPubKeyIDAlreadyExists)
}

func (s *AccountTreeTestSuite) TestSetSingle_InvalidPubKeyID() {
	account := &models.AccountLeaf{
		PubKeyID:  rightSubtreeMaxValue,
		PublicKey: models.PublicKey{1, 2, 3},
	}

	errMsg := fmt.Sprintf("invalid pubKeyID value: %d", account.PubKeyID)
	err := s.tree.SetSingle(account)
	s.Error(err)
	s.Equal(errMsg, err.Error())
}

func (s *AccountTreeTestSuite) TestUnsafeSet_ReturnsWitness() {
	witness, err := s.tree.unsafeSet(s.leaf)
	s.NoError(err)
	s.Len(witness, AccountTreeDepth)

	node, err := s.tree.getMerkleTreeNodeByPath(&models.MerklePath{Depth: AccountTreeDepth, Path: 1})
	s.NoError(err)
	s.Equal(node.DataHash, witness[0])

	node, err = s.tree.getMerkleTreeNodeByPath(&models.MerklePath{Depth: 1, Path: 1})
	s.NoError(err)
	s.Equal(node.DataHash, witness[31])
}

func (s *AccountTreeTestSuite) TestSetBatch_AddsAccountLeaves() {
	leaves := make([]models.AccountLeaf, 16)
	for i := range leaves {
		leaves[i] = models.AccountLeaf{
			PubKeyID:  uint32(i + accountBatchOffset),
			PublicKey: models.PublicKey{1, 2, byte(i)},
		}
	}

	err := s.tree.SetBatch(leaves)
	s.NoError(err)

	for i := range leaves {
		accountLeaf, err := s.tree.Leaf(leaves[i].PubKeyID)
		s.NoError(err)
		s.Equal(leaves[i], *accountLeaf)
	}
}

func (s *AccountTreeTestSuite) TestSetBatch_ChangesStateRoot() {
	leaves := make([]models.AccountLeaf, 16)
	for i := range leaves {
		leaves[i] = models.AccountLeaf{
			PubKeyID:  uint32(i + accountBatchOffset),
			PublicKey: models.PublicKey{1, 2, byte(i)},
		}
	}

	rootBeforeSet, err := s.tree.Root()
	s.NoError(err)

	err = s.tree.SetBatch(leaves)
	s.NoError(err)

	rootAfterSet, err := s.tree.Root()
	s.NoError(err)

	s.NotEqual(rootBeforeSet, rootAfterSet)
}

func (s *AccountTreeTestSuite) TestSetBatch_InvalidLeavesLength() {
	leaves := make([]models.AccountLeaf, 3)
	for i := range leaves {
		leaves[i] = models.AccountLeaf{
			PubKeyID:  uint32(i + accountBatchOffset),
			PublicKey: models.PublicKey{1, 2, byte(i)},
		}
	}

	err := s.tree.SetBatch(leaves)
	s.ErrorIs(err, ErrInvalidAccountsLength)
}

func (s *AccountTreeTestSuite) TestSetBatch_InvalidPubKeyIDValue() {
	leaves := make([]models.AccountLeaf, 16)
	for i := range leaves {
		leaves[i] = models.AccountLeaf{
			PubKeyID:  uint32(i + accountBatchOffset),
			PublicKey: models.PublicKey{1, 2, byte(i)},
		}
	}

	leaves[7].PubKeyID = 12

	errMsg := fmt.Sprintf("invalid pubKeyID value: %d", leaves[7].PubKeyID)
	err := s.tree.SetBatch(leaves)
	s.Error(err)
	s.Equal(errMsg, err.Error())

	_, err = s.tree.Leaf(leaves[0].PubKeyID)
	s.Equal(NewNotFoundError("account leaf"), err)
}

func (s *AccountTreeTestSuite) randomPublicKey() models.PublicKey {
	publicKey := models.PublicKey{}
	randomBytes := make([]byte, models.PublicKeyLength)
	_, err := rand.Read(randomBytes)
	s.NoError(err)
	err = publicKey.SetBytes(randomBytes)
	s.NoError(err)

	return publicKey
}

func TestAccountTreeTestSuite(t *testing.T) {
	suite.Run(t, new(AccountTreeTestSuite))
}
