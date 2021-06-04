package commander

import (
	"context"
	"testing"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type SyncTestSuite struct {
	*require.Assertions
	suite.Suite
	teardown            func() error
	storage             *st.Storage
	tree                *st.StateTree
	client              *eth.TestClient
	cfg                 *config.RollupConfig
	transactionExecutor *transactionExecutor
}

func (s *SyncTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *SyncTestSuite) SetupTest() {
	var err error
	s.client, err = eth.NewTestClient()
	s.NoError(err)

	s.cfg = &config.RollupConfig{
		MinCommitmentsPerBatch: 1,
		MaxCommitmentsPerBatch: 32,
		TxsPerCommitment:       1,
	}

	s.setupDB()
}

func (s *SyncTestSuite) setupDB() {
	testStorage, err := st.NewTestStorageWithBadger()
	s.NoError(err)
	s.storage = testStorage.Storage
	s.teardown = testStorage.Teardown
	s.tree = st.NewStateTree(s.storage)
	s.transactionExecutor = newTestTransactionExecutor(s.storage, s.client.Client, s.cfg)

	s.seedDB()
}

func (s *SyncTestSuite) seedDB() {
	err := s.storage.AddAccountIfNotExists(&models.Account{
		PubKeyID:  0,
		PublicKey: models.PublicKey{1, 2, 3},
	})
	s.NoError(err)

	err = s.storage.AddAccountIfNotExists(&models.Account{
		PubKeyID:  1,
		PublicKey: models.PublicKey{2, 3, 4},
	})
	s.NoError(err)

	err = s.tree.Set(0, &models.UserState{
		PubKeyID:   0,
		TokenIndex: models.MakeUint256(0),
		Balance:    models.MakeUint256(1000),
		Nonce:      models.MakeUint256(0),
	})
	s.NoError(err)

	err = s.tree.Set(1, &models.UserState{
		PubKeyID:   1,
		TokenIndex: models.MakeUint256(0),
		Balance:    models.MakeUint256(0),
		Nonce:      models.MakeUint256(0),
	})
	s.NoError(err)
}

func (s *SyncTestSuite) TearDownTest() {
	s.client.Close()
	err := s.teardown()
	s.NoError(err)
}

func (s *SyncTestSuite) TestSyncBatches_Transfer() {
	tx := models.Transfer{
		TransactionBase: models.TransactionBase{
			TxType:      txtype.Transfer,
			FromStateID: 0,
			Amount:      models.MakeUint256(400),
			Fee:         models.MakeUint256(0),
			Nonce:       models.MakeUint256(0),
			Signature:   *mockSignature(s.T()),
		},
		ToStateID: 1,
	}
	err := s.storage.AddTransfer(&tx)
	s.NoError(err)

	commitments, err := s.transactionExecutor.createTransferCommitments([]models.Transfer{tx}, testDomain)
	s.NoError(err)
	s.Len(commitments, 1)

	err = s.transactionExecutor.submitBatch(txtype.Transfer, commitments)
	s.NoError(err)

	s.client.Commit()

	// Recreate database
	err = s.teardown()
	s.NoError(err)
	s.setupDB()

	latestBlockNumber, err := s.client.GetLatestBlockNumber()
	s.NoError(err)
	err = s.transactionExecutor.SyncBatches(0, *latestBlockNumber)
	s.NoError(err)

	// Begin db transaction
	transactionExecutor, err := newTransactionExecutorWithCtx(context.Background(), s.storage, s.client.Client, s.cfg)
	s.NoError(err)

	tx2 := models.Transfer{
		TransactionBase: models.TransactionBase{
			TxType:      txtype.Transfer,
			FromStateID: 1,
			Amount:      models.MakeUint256(100),
			Fee:         models.MakeUint256(0),
			Nonce:       models.MakeUint256(0),
			Signature:   *mockSignature(s.T()),
		},
		ToStateID: 0,
	}
	err = transactionExecutor.storage.AddTransfer(&tx2)
	s.NoError(err)

	commitments, err = transactionExecutor.createTransferCommitments([]models.Transfer{tx2}, testDomain)
	s.NoError(err)
	s.Len(commitments, 1)

	err = transactionExecutor.submitBatch(txtype.Transfer, commitments)
	s.NoError(err)

	s.client.Commit()

	batches, err := transactionExecutor.storage.GetBatchesInRange(nil, nil)
	s.NoError(err)
	s.Len(batches, 2)

	transactionExecutor.Rollback(nil)

	batches, err = s.storage.GetBatchesInRange(nil, nil)
	s.NoError(err)
	s.Len(batches, 1)

	err = s.transactionExecutor.SyncBatches(0, *latestBlockNumber+2)
	s.NoError(err)

	state0, err := s.storage.GetStateLeaf(0)
	s.NoError(err)
	s.Equal(models.MakeUint256(700), state0.Balance)

	state1, err := s.storage.GetStateLeaf(1)
	s.NoError(err)
	s.Equal(models.MakeUint256(300), state1.Balance)

	batches, err = s.storage.GetBatchesInRange(nil, nil)
	s.NoError(err)
	s.Len(batches, 2)
}

func (s *SyncTestSuite) TestSyncBatches_Transfer2() {
	txs := []models.Transfer{
		{
			TransactionBase: models.TransactionBase{
				TxType:      txtype.Transfer,
				FromStateID: 0,
				Amount:      models.MakeUint256(400),
				Fee:         models.MakeUint256(0),
				Nonce:       models.MakeUint256(0),
				Signature:   *mockSignature(s.T()),
			},
			ToStateID: 1,
		}, {
			TransactionBase: models.TransactionBase{
				TxType:      txtype.Transfer,
				FromStateID: 1,
				Amount:      models.MakeUint256(100),
				Fee:         models.MakeUint256(0),
				Nonce:       models.MakeUint256(0),
				Signature:   *mockSignature(s.T()),
			},
			ToStateID: 0,
		},
	}
	for i := range txs {
		transferHash, err := encoder.HashTransfer(&txs[i])
		s.NoError(err)
		txs[i].Hash = *transferHash
		err = s.storage.AddTransfer(&txs[i])
		s.NoError(err)
	}

	commitments := make([]models.Commitment, 2)
	for i := range commitments {
		commitments1, err := s.transactionExecutor.createTransferCommitments([]models.Transfer{txs[i]}, testDomain)
		s.NoError(err)
		s.Len(commitments1, 1)
		commitments[i] = commitments1[0]
		err = s.transactionExecutor.submitBatch(txtype.Transfer, commitments1)
		s.NoError(err)
	}

	// Recreate database
	err := s.teardown()
	s.NoError(err)
	s.setupDB()

	latestBlockNumber, err := s.client.GetLatestBlockNumber()
	s.NoError(err)
	err = s.transactionExecutor.SyncBatches(0, *latestBlockNumber)
	s.NoError(err)

	batches, err := s.storage.GetBatchesInRange(nil, nil)
	s.NoError(err)
	s.Len(batches, 2)

	for i := range commitments {
		commitment, err := s.storage.GetCommitment(commitments[i].ID)
		s.NoError(err)
		commitments[i].IncludedInBatch = &batches[i].ID
		commitments[i].AccountTreeRoot = commitment.AccountTreeRoot
		s.Equal(commitments[i], *commitment)

		transfer, err := s.storage.GetTransfer(txs[i].Hash)
		s.NoError(err)
		txs[i].IncludedInCommitment = &commitments[i].ID
		s.Equal(txs[i], *transfer)
	}
}

func (s *SyncTestSuite) TestSyncBatches_PendingBatch() {
	tx := models.Transfer{
		TransactionBase: models.TransactionBase{
			TxType:      txtype.Transfer,
			FromStateID: 0,
			Amount:      models.MakeUint256(400),
			Fee:         models.MakeUint256(0),
			Nonce:       models.MakeUint256(0),
			Signature:   *mockSignature(s.T()),
		},
		ToStateID: 1,
	}
	err := s.storage.AddTransfer(&tx)
	s.NoError(err)

	commitments, err := s.transactionExecutor.createTransferCommitments([]models.Transfer{tx}, testDomain)
	s.NoError(err)
	s.Len(commitments, 1)

	err = s.transactionExecutor.submitBatch(txtype.Transfer, commitments)
	s.NoError(err)

	s.client.Commit()

	// Recreate database
	err = s.teardown()
	s.NoError(err)
	s.setupDB()
}

func (s *SyncTestSuite) TestSyncBatches_Create2Transfer() {
	// register sender account on chain
	registrations, unsubscribe, err := s.client.WatchRegistrations(&bind.WatchOpts{})
	s.NoError(err)
	defer unsubscribe()
	senderPubKeyID, err := s.client.RegisterAccount(&models.PublicKey{1, 2, 3}, registrations)
	s.NoError(err)
	s.Equal(uint32(0), *senderPubKeyID)

	tx := models.Create2Transfer{
		TransactionBase: models.TransactionBase{
			TxType:      txtype.Create2Transfer,
			FromStateID: *senderPubKeyID,
			Amount:      models.MakeUint256(400),
			Fee:         models.MakeUint256(0),
			Nonce:       models.MakeUint256(0),
			Signature:   *mockSignature(s.T()),
		},
		ToStateID:   ref.Uint32(5),
		ToPublicKey: models.PublicKey{2, 3, 4},
	}
	err = s.storage.AddCreate2Transfer(&tx)
	s.NoError(err)

	commitments, err := s.transactionExecutor.createCreate2TransferCommitments([]models.Create2Transfer{tx}, testDomain)
	s.NoError(err)
	s.Len(commitments, 1)

	err = s.transactionExecutor.submitBatch(txtype.Create2Transfer, commitments)
	s.NoError(err)

	s.client.Commit()

	// Recreate database
	err = s.teardown()
	s.NoError(err)
	s.setupDB()

	latestBlockNumber, err := s.client.GetLatestBlockNumber()
	s.NoError(err)
	err = s.transactionExecutor.SyncBatches(0, *latestBlockNumber)
	s.NoError(err)

	state0, err := s.storage.GetStateLeaf(0)
	s.NoError(err)
	s.Equal(models.MakeUint256(600), state0.Balance)

	state1, err := s.storage.GetStateLeaf(5)
	s.NoError(err)
	s.Equal(models.MakeUint256(400), state1.Balance)
	s.Equal(uint32(1), state1.PubKeyID)

	batches, err := s.storage.GetBatchesInRange(nil, nil)
	s.NoError(err)
	s.Len(batches, 1)
}

func mockSignature(t *testing.T) *models.Signature {
	wallet, err := bls.NewRandomWallet(*testDomain)
	require.NoError(t, err)
	signature, err := wallet.Sign(utils.RandomBytes(4))
	require.NoError(t, err)
	return signature.ModelsSignature()
}

func TestSyncTestSuite(t *testing.T) {
	suite.Run(t, new(SyncTestSuite))
}
