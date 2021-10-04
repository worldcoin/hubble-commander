package storage

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"

	"github.com/Worldcoin/hubble-commander/config"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type AccountTreeTestSuite struct {
	*require.Assertions
	suite.Suite
	storage   *TestStorage
	leaf      *models.AccountLeaf
	treeDepth uint8
}

func (s *AccountTreeTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
	s.treeDepth = 32
}

func (s *AccountTreeTestSuite) SetupTest() {
	var err error
	s.storage, err = NewTestStorageWithoutPostgres()
	s.NoError(err)

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

func (s *AccountTreeTestSuite) TestLeaf_NonExistentLeaf() {
	_, err := s.storage.AccountTree.Leaf(0)
	s.Equal(NewNotFoundError("account leaf"), err)
}

func (s *AccountTreeTestSuite) TestLeaves_NonExistentLeaves() {
	_, err := s.storage.AccountTree.Leaves(&models.PublicKey{1, 2, 3})
	s.Equal(NewNotFoundError("account leaves"), err)
}

func (s *AccountTreeTestSuite) TestLeaves_ReturnsAllLeaves() {
	pubKey := models.PublicKey{1, 2, 3}
	accounts := []models.AccountLeaf{{
		PubKeyID:  0,
		PublicKey: pubKey,
	}, {
		PubKeyID:  1,
		PublicKey: pubKey,
	}}

	err := s.storage.AccountTree.SetSingle(&accounts[0])
	s.NoError(err)
	err = s.storage.AccountTree.SetSingle(&accounts[1])
	s.NoError(err)

	res, err := s.storage.AccountTree.Leaves(&pubKey)
	s.NoError(err)

	s.Equal(accounts, res)
}

func (s *AccountTreeTestSuite) TestSetSingle_StoresAccountLeafRecord() {
	err := s.storage.AccountTree.SetSingle(s.leaf)
	s.NoError(err)

	actualLeaf, err := s.storage.AccountTree.Leaf(s.leaf.PubKeyID)
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

	err := s.storage.AccountTree.SetSingle(leaf0)
	s.NoError(err)

	accountRootAfter0, err := s.storage.AccountTree.Root()
	s.NoError(err)

	err = s.storage.AccountTree.SetSingle(leaf1)
	s.NoError(err)

	accountRootAfter1, err := s.storage.AccountTree.Root()
	s.NoError(err)

	s.NotEqual(accountRootAfter0, accountRootAfter1)
}

func (s *AccountTreeTestSuite) TestSetSingle_StoresLeafMerkleTreeNodeRecord() {
	err := s.storage.AccountTree.SetSingle(s.leaf)
	s.NoError(err)

	expectedNode := &models.MerkleTreeNode{
		MerklePath: models.MerklePath{
			Path:  0,
			Depth: s.treeDepth,
		},
		DataHash: crypto.Keccak256Hash(s.leaf.PublicKey.Bytes()),
	}

	node, err := s.storage.AccountTree.merkleTree.Get(expectedNode.MerklePath)
	s.NoError(err)
	s.Equal(expectedNode, node)
}

func (s *AccountTreeTestSuite) TestSetSingle_UpdatesRootMerkleTreeNodeRecord() {
	err := s.storage.AccountTree.SetSingle(s.leaf)
	s.NoError(err)

	root, err := s.storage.AccountTree.Root()
	s.NoError(err)
	s.Equal(common.HexToHash("0x6e082faf2fd8ce5accb1e08a15061f2c443ea5e9cb42d493050275d644bb51b9"), *root)
}

func (s *AccountTreeTestSuite) TestSetSingle_CalculatesCorrectRootForLeafOfId1() {
	s.leaf.PubKeyID = 1
	err := s.storage.AccountTree.SetSingle(s.leaf)
	s.NoError(err)

	root, err := s.storage.AccountTree.Root()
	s.NoError(err)
	s.Equal(common.HexToHash("0xcd6164584f02a9c4c9f88c2613d7ff2b709e0951369f9bd28528712e3fa96daa"), *root)
}

func (s *AccountTreeTestSuite) TestSetSingle_CalculatesCorrectRootForTwoLeaves() {
	err := s.storage.AccountTree.SetSingle(s.leaf)
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
	err = s.storage.AccountTree.SetSingle(leaf1)
	s.NoError(err)

	root, err := s.storage.AccountTree.Root()
	s.NoError(err)
	s.Equal(common.HexToHash("0x3a7a7ff21991ccfcbf8a4580862def7c498253ad398e967f270ff421db1d4833"), *root)
}

func (s *AccountTreeTestSuite) TestSetSingle_ReturnsErrorOnSettingAlreadySetLeaf() {
	err := s.storage.AccountTree.SetSingle(s.leaf)
	s.NoError(err)

	s.leaf.PublicKey = s.randomPublicKey()
	err = s.storage.AccountTree.SetSingle(s.leaf)

	var accountAlreadyExistsError *AccountAlreadyExistsError
	s.ErrorAs(err, &accountAlreadyExistsError)
	s.Equal(s.leaf, accountAlreadyExistsError.Account)
}

func (s *AccountTreeTestSuite) TestSetSingle_InvalidPubKeyID() {
	account := &models.AccountLeaf{
		PubKeyID:  rightSubtreeMaxValue,
		PublicKey: models.PublicKey{1, 2, 3},
	}

	err := s.storage.AccountTree.SetSingle(account)

	var invalidPubKeyIDError *InvalidPubKeyIDError
	s.ErrorAs(err, &invalidPubKeyIDError)
	s.Equal(account.PubKeyID, invalidPubKeyIDError.value)
}

func (s *AccountTreeTestSuite) TestUnsafeSet_ReturnsWitness() {
	witness, err := s.storage.AccountTree.unsafeSet(s.leaf)
	s.NoError(err)
	s.Len(witness, int(s.treeDepth))

	node, err := s.storage.AccountTree.merkleTree.Get(models.MerklePath{Depth: s.treeDepth, Path: 1})
	s.NoError(err)
	s.Equal(node.DataHash, witness[0])

	node, err = s.storage.AccountTree.merkleTree.Get(models.MerklePath{Depth: 1, Path: 1})
	s.NoError(err)
	s.Equal(node.DataHash, witness[31])
}

func (s *AccountTreeTestSuite) TestSetBatch_AddsAccountLeaves() {
	leaves := make([]models.AccountLeaf, AccountBatchSize)
	for i := range leaves {
		leaves[i] = models.AccountLeaf{
			PubKeyID:  uint32(i + AccountBatchOffset),
			PublicKey: models.PublicKey{1, 2, byte(i)},
		}
	}

	err := s.storage.AccountTree.SetBatch(leaves)
	s.NoError(err)

	for i := range leaves {
		accountLeaf, err := s.storage.AccountTree.Leaf(leaves[i].PubKeyID)
		s.NoError(err)
		s.Equal(leaves[i], *accountLeaf)
	}
}

func (s *AccountTreeTestSuite) TestSetBatch_ChangesStateRoot() {
	leaves := make([]models.AccountLeaf, AccountBatchSize)
	for i := range leaves {
		leaves[i] = models.AccountLeaf{
			PubKeyID:  uint32(i + AccountBatchOffset),
			PublicKey: models.PublicKey{1, 2, byte(i)},
		}
	}

	rootBeforeSet, err := s.storage.AccountTree.Root()
	s.NoError(err)

	err = s.storage.AccountTree.SetBatch(leaves)
	s.NoError(err)

	rootAfterSet, err := s.storage.AccountTree.Root()
	s.NoError(err)

	s.NotEqual(rootBeforeSet, rootAfterSet)
}

func (s *AccountTreeTestSuite) TestSetBatch_InvalidLeavesLength() {
	leaves := make([]models.AccountLeaf, 3)
	for i := range leaves {
		leaves[i] = models.AccountLeaf{
			PubKeyID:  uint32(i + AccountBatchOffset),
			PublicKey: models.PublicKey{1, 2, byte(i)},
		}
	}

	err := s.storage.AccountTree.SetBatch(leaves)
	s.ErrorIs(err, ErrInvalidAccountsLength)
}

func (s *AccountTreeTestSuite) TestSetBatch_ReturnsErrorOnSettingAlreadySetLeaf() {
	leaves := make([]models.AccountLeaf, AccountBatchSize)
	for i := range leaves {
		leaves[i] = models.AccountLeaf{
			PubKeyID:  uint32(i + AccountBatchOffset),
			PublicKey: models.PublicKey{1, 2, byte(i)},
		}
	}
	err := s.storage.AccountTree.SetBatch(leaves)
	s.NoError(err)

	err = s.storage.AccountTree.SetBatch(leaves)

	var accountBatchExistsError *AccountBatchAlreadyExistsError
	s.ErrorAs(err, &accountBatchExistsError)
	s.Equal(leaves, accountBatchExistsError.Accounts)
}

func (s *AccountTreeTestSuite) TestSetBatch_InvalidPubKeyIDValue() {
	leaves := make([]models.AccountLeaf, AccountBatchSize)
	for i := range leaves {
		leaves[i] = models.AccountLeaf{
			PubKeyID:  uint32(i + AccountBatchOffset),
			PublicKey: models.PublicKey{1, 2, byte(i)},
		}
	}

	leaves[7].PubKeyID = 12

	err := s.storage.AccountTree.SetBatch(leaves)

	var invalidPubKeyIDError *InvalidPubKeyIDError
	s.ErrorAs(err, &invalidPubKeyIDError)
	s.Equal(leaves[7].PubKeyID, invalidPubKeyIDError.value)

	_, err = s.storage.AccountTree.Leaf(leaves[0].PubKeyID)
	s.Equal(NewNotFoundError("account leaf"), err)
}

func (s *AccountTreeTestSuite) TestNextBatchAccountPubKeyID() {
	leaves := make([]models.AccountLeaf, AccountBatchSize)
	for i := range leaves {
		leaves[i] = models.AccountLeaf{
			PubKeyID:  uint32(i + AccountBatchOffset),
			PublicKey: models.PublicKey{1, 2, byte(i)},
		}
	}

	err := s.storage.AccountTree.SetBatch(leaves)
	s.NoError(err)

	pubKeyID, err := s.storage.AccountTree.NextBatchAccountPubKeyID()
	s.NoError(err)
	s.EqualValues(AccountBatchOffset+16, *pubKeyID)
}

func (s *AccountTreeTestSuite) TestNextBatchAccountPubKeyID_OnlySingleAccountsRegistered() {
	err := s.storage.AccountTree.SetSingle(s.leaf)
	s.NoError(err)

	pubKeyID, err := s.storage.AccountTree.NextBatchAccountPubKeyID()
	s.NoError(err)
	s.EqualValues(AccountBatchOffset, *pubKeyID)
}

func (s *AccountTreeTestSuite) TestNextBatchAccountPubKeyID_NoAccounts() {
	pubKeyID, err := s.storage.AccountTree.NextBatchAccountPubKeyID()
	s.NoError(err)
	s.EqualValues(AccountBatchOffset, *pubKeyID)
}

func (s *AccountTreeTestSuite) randomPublicKey() models.PublicKey {
	publicKey := models.PublicKey{}
	err := publicKey.SetBytes(utils.RandomBytes(models.PublicKeyLength))
	s.NoError(err)
	return publicKey
}

func TestAccountTreeTestSuite(t *testing.T) {
	suite.Run(t, new(AccountTreeTestSuite))
}

func TestReadAllRegisterAccounts(t *testing.T) {
	cfg := config.GetConfig()
	//cfg.Badger.Path = "badgerPath"
	database, err := NewDatabase(cfg)
	require.NoError(t, err)

	accountTree := NewAccountTree(database)

	accounts := make([]models.AccountLeaf, 0, 100)
	for i := AccountBatchOffset; ; i++ {
		leaf, err := accountTree.Leaf(uint32(i))
		if IsNotFoundError(err) {
			break
		}
		require.NoError(t, err)
		accounts = append(accounts, *leaf)
	}

	// jsonAccounts, err := json.MarshalIndent(accounts, "", "\t")
	// require.NoError(t, err)
	// fmt.Println(string(jsonAccounts))

	for _, account := range accounts {
		fmt.Printf("%d = 0x%s\n", account.PubKeyID, hex.EncodeToString(account.PublicKey.Bytes()))
	}

	// counter := 0
	// for i := range accounts {
	// 	counter++
	// 	if counter == 1 {
	// 		fmt.Printf("%d -", accounts[i].PubKeyID)
	// 	}
	// 	if counter == 16 {
	// 		fmt.Printf(" %d\n", accounts[i].PubKeyID)
	// 		counter = 0
	// 	}
	// }

	fmt.Printf("\n\n\n")
}
