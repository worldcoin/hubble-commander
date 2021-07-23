package storage

import (
	"crypto/rand"
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
	s.tree = NewAccountTree(s.storage.InternalStorage)

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

func (s *AccountTreeTestSuite) TestSet_StoresAccountLeafRecord() {
	_, err := s.tree.Set(s.leaf)
	s.NoError(err)

	actualLeaf, err := s.storage.GetAccountLeaf(s.leaf.PubKeyID)
	s.NoError(err)
	s.Equal(s.leaf, actualLeaf)
}

func (s *AccountTreeTestSuite) TestSet_RootIsDifferentAfterSet() {
	leaf0 := &models.AccountLeaf{
		PubKeyID:  0,
		PublicKey: s.randomPublicKey(),
	}

	leaf1 := &models.AccountLeaf{
		PubKeyID:  1,
		PublicKey: s.randomPublicKey(),
	}

	_, err := s.tree.Set(leaf0)
	s.NoError(err)

	accountRootAfter0, err := s.tree.Root()
	s.NoError(err)

	_, err = s.tree.Set(leaf1)
	s.NoError(err)

	accountRootAfter1, err := s.tree.Root()
	s.NoError(err)

	s.NotEqual(accountRootAfter0, accountRootAfter1)
}

func (s *AccountTreeTestSuite) TestSet_StoresLeafMerkleTreeNodeRecord() {
	_, err := s.tree.Set(s.leaf)
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

func (s *AccountTreeTestSuite) TestSet_UpdatesRootMerkleTreeNodeRecord() {
	_, err := s.tree.Set(s.leaf)
	s.NoError(err)

	root, err := s.tree.Root()
	s.NoError(err)
	s.Equal(common.HexToHash("0x6e082faf2fd8ce5accb1e08a15061f2c443ea5e9cb42d493050275d644bb51b9"), *root)
}

func (s *AccountTreeTestSuite) TestSet_CalculatesCorrectRootForLeafOfId1() {
	s.leaf.PubKeyID = 1
	_, err := s.tree.Set(s.leaf)
	s.NoError(err)

	root, err := s.tree.Root()
	s.NoError(err)
	s.Equal(common.HexToHash("0xcd6164584f02a9c4c9f88c2613d7ff2b709e0951369f9bd28528712e3fa96daa"), *root)
}

func (s *AccountTreeTestSuite) TestSet_CalculatesCorrectRootForTwoLeaves() {
	_, err := s.tree.Set(s.leaf)
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
	_, err = s.tree.Set(leaf1)
	s.NoError(err)

	root, err := s.tree.Root()
	s.NoError(err)
	s.Equal(common.HexToHash("0x3a7a7ff21991ccfcbf8a4580862def7c498253ad398e967f270ff421db1d4833"), *root)
}

func (s *AccountTreeTestSuite) TestSet_ThrowsOnSettingAlreadySetLeaf() {
	_, err := s.tree.Set(s.leaf)
	s.NoError(err)

	s.leaf.PublicKey = s.randomPublicKey()
	_, err = s.tree.Set(s.leaf)
	s.ErrorIs(err, ErrPubKeyIDAlreadyExists)
}

func (s *AccountTreeTestSuite) TestSet_ReturnsWitness() {
	witness, err := s.tree.Set(s.leaf)
	s.NoError(err)
	s.Len(witness, AccountTreeDepth)

	node, err := s.tree.getMerkleTreeNodeByPath(&models.MerklePath{Depth: AccountTreeDepth, Path: 1})
	s.NoError(err)
	s.Equal(node.DataHash, witness[0])

	node, err = s.tree.getMerkleTreeNodeByPath(&models.MerklePath{Depth: 1, Path: 1})
	s.NoError(err)
	s.Equal(node.DataHash, witness[31])
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
