package storage

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/Worldcoin/hubble-commander/models/stored"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/Worldcoin/hubble-commander/utils/consts"
	"github.com/dgraph-io/badger/v3"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type MempoolTestSuite struct {
	*require.Assertions
	suite.Suite
	storage *TestStorage
}

func (s *MempoolTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *MempoolTestSuite) SetupTest() {
	var err error
	s.storage, err = NewTestStorage()
	s.NoError(err)

	err = s.storage.AccountTree.SetSingle(&models.AccountLeaf{
		PubKeyID:  123,
		PublicKey: models.PublicKey{1, 2, 3},
	})
	s.NoError(err)

	// sender
	_, err = s.storage.StateTree.Set(
		1,
		&models.UserState{
			PubKeyID: 1,
			TokenID:  models.MakeUint256(1),
			Balance:  models.MakeUint256(100),
			Nonce:    models.MakeUint256(0),
		},
	)
	s.NoError(err)

	// receiver
	_, err = s.storage.StateTree.Set(
		2,
		&models.UserState{
			PubKeyID: 1,
			TokenID:  models.MakeUint256(1),
			Balance:  models.MakeUint256(0),
			Nonce:    models.MakeUint256(0),
		},
	)
	s.NoError(err)
}

func (s *MempoolTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

// TODO: this should not be a unit test, this test belongs on the api/ directory
//       the api test will test the exact same thing but make the code less brittle
func (s *MempoolTestSuite) TestMempool_UsesPendingBalance() {
	// stateID=2 starts with Balance=0, so the API does not accept any transfers
	// from it until it receives some money

	transfer := testutils.NewTransfer(2, 1, 0, 10)
	err := s.storage.AddMempoolTx(transfer)
	s.ErrorContains(err, "balance too low")

	transfer = testutils.NewTransfer(1, 2, 0, 10)
	err = s.storage.AddMempoolTx(transfer)
	s.NoError(err)

	// we have 10 but this txn also pays a fee of 10 so it still fails
	transfer = testutils.NewTransfer(2, 1, 0, 10)
	err = s.storage.AddMempoolTx(transfer)
	s.ErrorContains(err, "balance too low")

	transfer = testutils.NewTransfer(1, 2, 1, 10)
	err = s.storage.AddMempoolTx(transfer)
	s.NoError(err)

	// now that our pending balance is high enough it is accepted
	transfer = testutils.NewTransfer(2, 1, 0, 10)
	err = s.storage.AddMempoolTx(transfer)
	s.NoError(err)
}

func (s *MempoolTestSuite) TestMempool_LowestNoncePendingTxs() {
	txs := []*models.Transfer{
		testutils.NewTransfer(1, 2, 0, 10),
		testutils.NewTransfer(1, 2, 1, 10),
		testutils.NewTransfer(2, 1, 0, 10),
	}

	for _, tx := range txs {
		s.NotNil(tx)

		err := s.storage.AddMempoolTx(tx)
		s.NoError(err)
	}

	// it skips the second txn from 1->2

	pendingTxs, err := s.storage.lowestNoncePendingTxs()
	s.NoError(err)
	s.Len(pendingTxs, 2)

	s.Equal(uint32(1), pendingTxs[0].FromStateID)
	s.Equal(models.MakeUint256(0), pendingTxs[0].Nonce)

	s.Equal(uint32(2), pendingTxs[1].FromStateID)
	s.Equal(models.MakeUint256(0), pendingTxs[1].Nonce)
}

func (s *MempoolTestSuite) TestMempool_PeekEmptyMempool() {
	mempoolHeap, err := s.storage.NewMempoolHeap(txtype.Transfer)
	s.NoError(err)
	s.NotNil(mempoolHeap)

	// if there's nothing in the mempool you won't crash when you ask for the next Tx
	firstTx := mempoolHeap.PeekHighestFeeExecutableTx()
	s.Nil(firstTx)
}

func (s *MempoolTestSuite) randomPublicKey() *models.PublicKey {
	domain := bls.Domain{1, 2, 3, 4}
	wallet, err := bls.NewRandomWallet(domain)
	s.NoError(err)

	return wallet.PublicKey()
}

// skips all pending state maintenance
func (s *MempoolTestSuite) rawInsert(tx models.GenericTransaction) {
	pendingTx := stored.NewPendingTx(tx)
	txKey := pendingTxKey(pendingTx.FromStateID, pendingTx.Nonce.Uint64())
	err := s.storage.rawSet(txKey, pendingTx.Bytes())
	s.NoError(err)
}

func (s *MempoolTestSuite) TestMempool_MigrateState() {
	destKey := s.randomPublicKey()

	s.rawInsert(testutils.NewCreate2Transfer(1, nil, 0, 10, destKey))
	s.rawInsert(testutils.NewCreate2Transfer(1, nil, 1, 10, destKey))

	pendingState, err := s.storage.GetPendingC2TState(destKey)
	s.NoError(err)
	s.Nil(pendingState)

	migrated, err := s.storage.alreadyRanPubKeyMigration()
	s.NoError(err)
	s.False(migrated)

	err = s.storage.MigratePubKeyPendingState()
	s.NoError(err)

	migrated, err = s.storage.alreadyRanPubKeyMigration()
	s.NoError(err)
	s.True(migrated)

	pendingState, err = s.storage.GetPendingC2TState(destKey)
	s.NoError(err)
	s.Equal(
		models.UserState{
			PubKeyID: consts.PendingID,
			TokenID:  models.MakeUint256(0),
			Balance:  models.MakeUint256(20),
			Nonce:    models.MakeUint256(0),
		},
		*pendingState,
	)
}

func TestMempoolTestSuite(t *testing.T) {
	suite.Run(t, new(MempoolTestSuite))
}

type ConflictTestSuite struct {
	*require.Assertions
	suite.Suite
	rollupStorage *Storage
	apiStorage    *Storage
}

func (s *ConflictTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *ConflictTestSuite) SetupTest() {
	memoryDB, err := db.NewInMemoryDatabase()
	s.NoError(err)

	database := &Database{
		Badger: memoryDB,
	}

	rollupStorage, err := newStorageFromDatabase(database)
	s.NoError(err)
	s.NotNil(rollupStorage)
	s.rollupStorage = rollupStorage

	apiStorage, err := newStorageFromDatabase(database)
	s.NoError(err)
	s.NotNil(apiStorage)
	s.apiStorage = apiStorage

	_, err = apiStorage.StateTree.Set(
		1,
		&models.UserState{
			PubKeyID: 1,
			TokenID:  models.MakeUint256(1),
			Balance:  models.MakeUint256(100),
			Nonce:    models.MakeUint256(0),
		},
	)
	s.NoError(err)

	_, err = apiStorage.StateTree.Set(
		2,
		&models.UserState{
			PubKeyID: 1,
			TokenID:  models.MakeUint256(1),
			Balance:  models.MakeUint256(0),
			Nonce:    models.MakeUint256(0),
		},
	)
	s.NoError(err)

	// confirm this immediately shows up in rollupStorage, they share a db
	leaf, err := s.rollupStorage.StateTree.Leaf(1)
	s.NoError(err)
	s.Equal(models.MakeUint256(100), leaf.UserState.Balance)
}

func (s *ConflictTestSuite) TestConflict_ConflictsWithGet() {
	// nextTxForAccount goes to some lengths to not call txn.Get() in order to avoid
	// badger.ErrConflict, this tests confirms that is a real concern.

	// (I) open rollupTx and do nothing with it

	rollupTxController, txRollupStorage := s.rollupStorage.BeginTransaction(TxOptions{})
	s.NotNil(rollupTxController)
	s.NotNil(txRollupStorage)

	// (II) open apiTx and insert a transaction

	apiTxController, txAPIStorage := s.apiStorage.BeginTransaction(TxOptions{})
	s.NotNil(apiTxController)
	s.NotNil(txAPIStorage)

	transfer := testutils.NewTransfer(1, 2, 0, 10)
	err := txAPIStorage.unsafeAddMempoolTx(transfer) // the safe version Commit()s
	s.NoError(err)

	// (III) txn.Get() in rollupTx

	key := pendingTxKey(1, 0)
	_, err = txRollupStorage.rawLookup(key)
	s.ErrorIs(err, badger.ErrKeyNotFound) // the API has not yet Commit()

	// (IV) do a write so that badger knows to try to fail this tx in case of conflict
	key = pendingTxKey(10, 10)
	err = txRollupStorage.rawSet(key, []byte("hello"))
	s.NoError(err)

	// (V) commit apiTX

	err = apiTxController.Commit()
	s.NoError(err)

	// (VI) commit rollupTx and notice that it fails

	err = rollupTxController.Commit()
	s.ErrorIs(err, badger.ErrConflict)
}

// now, confirm that there is no conflict even if the rollup loop has "read" the key.
// in TestConflict_ConflictsWithGet we saw that a call to `rawLookup` in the rollup loop
// will cause a conflict to occur. Here we throw a transaction into the mempool which is not
// executable (the balance is too low) and check that the rollup builds a batch which correctly
// excludes the transaction without causing a conflict.
func (s *ConflictTestSuite) TestConflict_NoConflict() {
	transfer := testutils.NewTransfer(1, 2, 0, 10)
	err := s.apiStorage.AddMempoolTx(transfer)
	s.NoError(err)

	// stateID=2 starts with a balance of 0 so this transaction is not
	// executable
	transfer = testutils.NewTransfer(2, 1, 0, 10)
	{
		// go through some hoops to insert the transaction because the normal
		// code paths very much so do not want to accept transactions which the
		// user does not have enough funds to pay for
		pending := stored.NewPendingTx(transfer)
		key := pendingTxKey(2, 0)
		err = s.apiStorage.rawSet(key, pending.Bytes())
		s.NoError(err)

		err = s.apiStorage.UnsafeSetPendingState(2, models.MakeUint256(1), models.MakeUint256(0))
		s.NoError(err)
	}

	// (I)   we start the rollupTx

	rollupTxController, txRollupStorage := s.rollupStorage.BeginTransaction(TxOptions{})
	s.NotNil(rollupTxController)
	s.NotNil(txRollupStorage)

	// (II)  we start the apiTx

	apiTxController, txAPIStorage := s.apiStorage.BeginTransaction(TxOptions{})
	s.NotNil(apiTxController)
	s.NotNil(txAPIStorage)

	// (III) we read out some transactions in the rollup loop

	mempoolHeap, err := txRollupStorage.NewMempoolHeap(txtype.Transfer)
	s.NoError(err)
	s.NotNil(mempoolHeap)

	firstTx := mempoolHeap.PeekHighestFeeExecutableTx()
	s.NotNil(firstTx)
	s.Equal(uint32(1), firstTx.GetFromStateID())

	// this attempts to read out txns from both accounts
	//  this is the line which will cause a conflict if mempool contains a bug
	err = mempoolHeap.DropHighestFeeExecutableTx()
	s.NoError(err)
	err = mempoolHeap.Savepoint()
	s.NoError(err)

	// (IV)  we insert some new txns from the api

	transfer = testutils.NewTransfer(1, 2, 1, 40)
	err = txAPIStorage.AddMempoolTx(transfer)
	s.NoError(err)

	transfer = testutils.NewTransfer(2, 1, 1, 10)
	err = txAPIStorage.AddMempoolTx(transfer)
	s.NoError(err)

	err = apiTxController.Commit()
	s.NoError(err)

	// (VII) we cleanly commit the rollup txn

	err = rollupTxController.Commit()
	s.NoError(err)
}

// TODO: are there more likely scenarios where this might conflict?

func TestConflictTestSuite(t *testing.T) {
	suite.Run(t, new(ConflictTestSuite))
}
