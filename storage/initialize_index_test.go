package storage

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/stored"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	bh "github.com/timshannon/badgerhold/v4"
)

type InitializeIndexTestSuite struct {
	*require.Assertions
	suite.Suite
	storage *TestStorage
}

func (s *InitializeIndexTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *InitializeIndexTestSuite) SetupTest() {
	var err error
	s.storage, err = NewTestStorage()
	s.NoError(err)
}

func (s *InitializeIndexTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *InitializeIndexTestSuite) TestStoredBatchedTx_CommitmentID_IndexWorks() {
	id1 := models.CommitmentID{BatchID: models.MakeUint256(1), IndexInBatch: 0}
	id2 := models.CommitmentID{BatchID: models.MakeUint256(2), IndexInBatch: 0}
	s.addStoredBatchedTx(&id1)
	s.addStoredBatchedTx(&id2)
	s.addStoredBatchedTx(&id1)

	indexValues := s.getCommitmentIDIndexValues()
	s.Len(indexValues, 2)
	s.Len(indexValues[id1], 2)
	s.Len(indexValues[id2], 1)
}

func (s *InitializeIndexTestSuite) TestAccountLeaf_PublicKey_IndexWorks() {
	pk1 := models.PublicKey{1, 2, 3}
	pk2 := models.PublicKey{4, 5, 6}
	_, err := s.storage.AccountTree.unsafeSet(&models.AccountLeaf{
		PubKeyID:  0,
		PublicKey: pk1,
	})
	s.NoError(err)
	_, err = s.storage.AccountTree.unsafeSet(&models.AccountLeaf{
		PubKeyID:  1,
		PublicKey: pk1,
	})
	s.NoError(err)
	_, err = s.storage.AccountTree.unsafeSet(&models.AccountLeaf{
		PubKeyID:  2,
		PublicKey: pk2,
	})
	s.NoError(err)

	indexValues := s.getPublicKeyIndexValues(models.AccountLeafName)
	s.Len(indexValues, 3)
	s.Len(indexValues[models.ZeroPublicKey], 0) // value set due to index initialization, see NewAccountTree
	s.Len(indexValues[pk1], 2)
	s.Len(indexValues[pk2], 1)
}

func (s *InitializeIndexTestSuite) TestAccountLeaf_PublicKey_AccountsWithZeroPublicKeysAreNotIndexed() {
	_, err := s.storage.AccountTree.unsafeSet(&models.AccountLeaf{
		PubKeyID:  0,
		PublicKey: models.ZeroPublicKey,
	})
	s.NoError(err)

	indexValues := s.getPublicKeyIndexValues(models.AccountLeafName)
	s.Len(indexValues, 1)
	s.Len(indexValues[models.ZeroPublicKey], 0) // value set due to index initialization, see NewAccountTree
}

// This test checks an edge case that we introduced by not indexing AccountLeaf structs with zero public keys.
func (s *InitializeIndexTestSuite) TestAccountLeaf_PublicKey_FindUsingIndexWorksWhenThereAreOnlyAccountsWithZeroPublicKey() {
	_, err := s.storage.AccountTree.unsafeSet(&models.AccountLeaf{
		PubKeyID:  1,
		PublicKey: models.ZeroPublicKey, // zero public key values are not indexed
	})
	s.NoError(err)

	accounts := make([]models.AccountLeaf, 0, 1)
	err = s.storage.database.Badger.Find(
		&accounts,
		bh.Where("PublicKey").Ge(models.ZeroPublicKey).Index("PublicKey"),
	)
	s.NoError(err)
	s.Len(accounts, 0)
}

func (s *InitializeIndexTestSuite) getCommitmentIDIndexValues() map[models.CommitmentID]bh.KeyList {
	indexValues := make(map[models.CommitmentID]bh.KeyList)

	s.iterateIndex(stored.BatchedTxName, "CommitmentID", func(encodedKey []byte, keyList bh.KeyList) {
		var commitmentID models.CommitmentID
		err := db.Decode(encodedKey, &commitmentID)
		s.NoError(err)

		indexValues[commitmentID] = keyList
	})

	return indexValues
}

func (s *InitializeIndexTestSuite) getPublicKeyIndexValues(typeName []byte) map[models.PublicKey]bh.KeyList {
	indexValues := make(map[models.PublicKey]bh.KeyList)

	s.iterateIndex(typeName, "PublicKey", func(encodedKey []byte, keyList bh.KeyList) {
		var publicKey models.PublicKey
		err := db.Decode(encodedKey, &publicKey)
		s.NoError(err)

		indexValues[publicKey] = keyList
	})

	return indexValues
}

func (s *InitializeIndexTestSuite) iterateIndex(
	typeName []byte,
	indexName string,
	handleIndex func(encodedKey []byte, keyList bh.KeyList),
) {
	testutils.IterateIndex(s.Assertions, s.storage.database.Badger, typeName, indexName, handleIndex)
}

func (s *InitializeIndexTestSuite) addStoredBatchedTx(commitmentID *models.CommitmentID) {
	tx := testutils.MakeTransfer(11, 0, 1, 10)
	tx.CommitmentID = commitmentID

	batchedTx := stored.NewBatchedTx(&tx)
	err := s.storage.database.Badger.Insert(batchedTx.Hash, *batchedTx)
	s.NoError(err)
}

func TestInitializeIndexTestSuite(t *testing.T) {
	suite.Run(t, new(InitializeIndexTestSuite))
}
