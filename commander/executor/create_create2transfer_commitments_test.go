package executor

import (
	"context"
	"math/big"
	"testing"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type Create2TransferCommitmentsTestSuite struct {
	*require.Assertions
	suite.Suite
	storage                *st.TestStorage
	client                 *eth.TestClient
	cfg                    *config.RollupConfig
	transactionExecutor    *TransactionExecutor
	maxTxBytesInCommitment int
}

func (s *Create2TransferCommitmentsTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *Create2TransferCommitmentsTestSuite) SetupTest() {
	var err error
	s.storage, err = st.NewTestStorageWithBadger()
	s.NoError(err)
	s.client, err = eth.NewTestClient()
	s.NoError(err)
	s.cfg = &config.RollupConfig{
		MinTxsPerCommitment:    1,
		MaxTxsPerCommitment:    4,
		FeeReceiverPubKeyID:    2,
		MaxCommitmentsPerBatch: 1,
	}
	s.maxTxBytesInCommitment = encoder.Create2TransferLength * int(s.cfg.MaxTxsPerCommitment)

	err = populateAccounts(s.storage.Storage, genesisBalances)
	s.NoError(err)

	s.transactionExecutor = NewTestTransactionExecutor(s.storage.Storage, s.client.Client, s.cfg, context.Background())
}

func (s *Create2TransferCommitmentsTestSuite) TearDownTest() {
	s.client.Close()
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *Create2TransferCommitmentsTestSuite) TestCreateCreate2TransferCommitments_WithMinTxsPerCommitment() {
	transfers := generateValidCreate2Transfers(1)
	s.addCreate2Transfers(transfers)

	preRoot, err := s.transactionExecutor.storage.StateTree.Root()
	s.NoError(err)

	expectedTxsLength := encoder.Create2TransferLength * len(transfers)
	commitments, err := s.transactionExecutor.CreateCreate2TransferCommitments(testDomain)
	s.NoError(err)
	s.Len(commitments, 1)
	s.Len(commitments[0].Transactions, expectedTxsLength)

	postRoot, err := s.transactionExecutor.storage.StateTree.Root()
	s.NoError(err)
	s.NotEqual(preRoot, postRoot)
	s.Equal(commitments[0].PostStateRoot, *postRoot)
}

func (s *Create2TransferCommitmentsTestSuite) TestCreateCreate2TransferCommitments_WithMoreThanMinTxsPerCommitment() {
	transfers := generateValidCreate2Transfers(3)
	s.addCreate2Transfers(transfers)

	preRoot, err := s.transactionExecutor.storage.StateTree.Root()
	s.NoError(err)

	expectedTxsLength := encoder.Create2TransferLength * len(transfers)
	commitments, err := s.transactionExecutor.CreateCreate2TransferCommitments(testDomain)
	s.NoError(err)
	s.Len(commitments, 1)
	s.Len(commitments[0].Transactions, expectedTxsLength)

	postRoot, err := s.transactionExecutor.storage.StateTree.Root()
	s.NoError(err)
	s.NotEqual(preRoot, postRoot)
	s.Equal(commitments[0].PostStateRoot, *postRoot)
}

func (s *Create2TransferCommitmentsTestSuite) TestCreateCreate2TransferCommitments_QueriesForMorePendingTransfersUntilSatisfied() {
	addAccountWithHighNonce(s.Assertions, s.storage.Storage, 124)

	transfers := generateValidCreate2Transfers(6)
	s.invalidateCreate2Transfers(transfers[3:6])

	highNonceTransfer := testutils.MakeCreate2Transfer(124, nil, 10, 1, &models.PublicKey{5, 4, 3, 2, 1})
	transfers = append(transfers, highNonceTransfer)

	s.addCreate2Transfers(transfers)

	preRoot, err := s.transactionExecutor.storage.StateTree.Root()
	s.NoError(err)

	commitments, err := s.transactionExecutor.CreateCreate2TransferCommitments(testDomain)
	s.NoError(err)
	s.Len(commitments, 1)
	s.Len(commitments[0].Transactions, s.maxTxBytesInCommitment)

	postRoot, err := s.transactionExecutor.storage.StateTree.Root()
	s.NoError(err)
	s.NotEqual(preRoot, postRoot)
	s.Equal(commitments[0].PostStateRoot, *postRoot)
}

func (s *Create2TransferCommitmentsTestSuite) TestCreateCreate2TransferCommitments_ForMultipleCommitmentsInBatch() {
	s.cfg = &config.RollupConfig{
		MinTxsPerCommitment:    1,
		MaxTxsPerCommitment:    4,
		FeeReceiverPubKeyID:    2,
		MaxCommitmentsPerBatch: 3,
	}

	s.transactionExecutor = NewTestTransactionExecutor(s.storage.Storage, s.client.Client, s.cfg, context.Background())

	addAccountWithHighNonce(s.Assertions, s.storage.Storage, 124)

	transfers := generateValidCreate2Transfers(9)
	s.invalidateCreate2Transfers(transfers[7:9])

	highNonceTransfers := []models.Create2Transfer{
		testutils.MakeCreate2Transfer(124, nil, 10, 1, &models.PublicKey{5, 4, 3, 2, 1}),
		testutils.MakeCreate2Transfer(124, nil, 11, 1, &models.PublicKey{5, 4, 3, 2, 1}),
	}

	transfers = append(transfers, highNonceTransfers...)
	s.addCreate2Transfers(transfers)

	preRoot, err := s.transactionExecutor.storage.StateTree.Root()
	s.NoError(err)

	commitments, err := s.transactionExecutor.CreateCreate2TransferCommitments(testDomain)
	s.NoError(err)
	s.Len(commitments, 3)
	s.Len(commitments[0].Transactions, s.maxTxBytesInCommitment)
	s.Len(commitments[1].Transactions, s.maxTxBytesInCommitment)
	s.Len(commitments[2].Transactions, encoder.Create2TransferLength)

	postRoot, err := s.transactionExecutor.storage.StateTree.Root()
	s.NoError(err)
	s.NotEqual(preRoot, postRoot)
	s.Equal(commitments[2].PostStateRoot, *postRoot)
}

func (s *Create2TransferCommitmentsTestSuite) invalidateCreate2Transfers(transfers []models.Create2Transfer) {
	for i := range transfers {
		tx := &transfers[i]
		tx.Amount = *genesisBalances[tx.FromStateID].MulN(10)
	}
}

func (s *Create2TransferCommitmentsTestSuite) TestCreateCreate2TransferCommitments_ReturnsErrorWhenThereAreNotEnoughPendingTransfers() {
	preRoot, err := s.transactionExecutor.storage.StateTree.Root()
	s.NoError(err)

	commitments, err := s.transactionExecutor.CreateCreate2TransferCommitments(testDomain)
	s.Nil(commitments)
	s.Equal(ErrNotEnoughC2Transfers, err)

	postRoot, err := s.transactionExecutor.storage.StateTree.Root()
	s.NoError(err)

	s.Equal(preRoot, postRoot)
}

func (s *Create2TransferCommitmentsTestSuite) TestCreateCreate2TransferCommitments_ReturnsErrorWhenThereAreNotEnoughValidTransfers() {
	s.cfg = &config.RollupConfig{
		MinTxsPerCommitment:    32,
		MaxTxsPerCommitment:    32,
		FeeReceiverPubKeyID:    2,
		MaxCommitmentsPerBatch: 1,
	}

	s.transactionExecutor = NewTestTransactionExecutor(s.storage.Storage, s.client.Client, s.cfg, context.Background())

	transfers := generateValidCreate2Transfers(2)
	s.addCreate2Transfers(transfers)

	pendingTransfers, err := s.storage.GetPendingCreate2Transfers(32)
	s.NoError(err)
	s.Len(pendingTransfers, 2)

	preRoot, err := s.transactionExecutor.storage.StateTree.Root()
	s.NoError(err)

	commitments, err := s.transactionExecutor.CreateCreate2TransferCommitments(testDomain)
	s.Nil(commitments)
	s.Equal(ErrNotEnoughC2Transfers, err)

	postRoot, err := s.transactionExecutor.storage.StateTree.Root()
	s.NoError(err)

	s.Equal(preRoot, postRoot)
}

func (s *Create2TransferCommitmentsTestSuite) TestCreateCreate2TransferCommitments_StoresCorrectCommitment() {
	transfersCount := uint32(4)
	s.preparePendingCreate2Transfers(transfersCount)

	preRoot, err := s.transactionExecutor.storage.StateTree.Root()
	s.NoError(err)

	expectedTxsLength := encoder.Create2TransferLength * int(transfersCount)
	commitments, err := s.transactionExecutor.CreateCreate2TransferCommitments(testDomain)
	s.NoError(err)
	s.Len(commitments, 1)
	s.Len(commitments[0].Transactions, expectedTxsLength)
	s.Equal(commitments[0].FeeReceiver, uint32(2))
	s.Nil(commitments[0].IncludedInBatch)

	postRoot, err := s.transactionExecutor.storage.StateTree.Root()
	s.NoError(err)
	s.NotEqual(preRoot, postRoot)
	s.Equal(commitments[0].PostStateRoot, *postRoot)
}

func (s *Create2TransferCommitmentsTestSuite) TestCreateCreate2TransferCommitments_CreatesMaximallyAsManyCommitmentsAsSpecifiedInConfig() {
	s.preparePendingCreate2Transfers(5)

	commitments, err := s.transactionExecutor.CreateCreate2TransferCommitments(testDomain)
	s.NoError(err)
	s.Len(commitments, 1)
}

func (s *Create2TransferCommitmentsTestSuite) TestCreateCreate2TransferCommitments_UpdateTransfers() {
	s.preparePendingCreate2Transfers(2)

	pendingTransfers, err := s.storage.GetPendingCreate2Transfers(2)
	s.NoError(err)
	s.Len(pendingTransfers, 2)

	commitments, err := s.transactionExecutor.CreateCreate2TransferCommitments(testDomain)
	s.NoError(err)
	s.Len(commitments, 1)

	for i := range pendingTransfers {
		tx, err := s.storage.GetCreate2Transfer(pendingTransfers[i].Hash)
		s.NoError(err)
		s.Equal(int32(1), *tx.IncludedInCommitment)
		s.Equal(uint32(i+3), *tx.ToStateID)
	}
}

func (s *Create2TransferCommitmentsTestSuite) TestCreateCreate2TransferCommitments_RegistersAccounts() {
	transfers := generateValidCreate2Transfers(1)
	s.addCreate2Transfers(transfers)

	expectedTxsLength := encoder.Create2TransferLength * len(transfers)
	commitments, err := s.transactionExecutor.CreateCreate2TransferCommitments(testDomain)
	s.NoError(err)
	s.Len(commitments, 1)
	s.Len(commitments[0].Transactions, expectedTxsLength)

	s.client.Commit()
	accounts := s.getRegisteredAccounts(0)
	s.Len(accounts, 16)
	s.Equal(transfers[0].ToPublicKey, accounts[0].PublicKey)
}

func (s *Create2TransferCommitmentsTestSuite) TestRemoveCreate2Transfer() {
	transfer1 := models.Create2Transfer{
		TransactionBase: models.TransactionBase{
			Hash: utils.RandomHash(),
		},
	}
	transfer2 := models.Create2Transfer{
		TransactionBase: models.TransactionBase{
			Hash: utils.RandomHash(),
		},
	}
	transfer3 := models.Create2Transfer{
		TransactionBase: models.TransactionBase{
			Hash: utils.RandomHash(),
		},
	}

	transfers := []models.Create2Transfer{transfer1, transfer2, transfer3}
	toRemove := []models.Create2Transfer{transfer2}

	s.Equal([]models.Create2Transfer{transfer1, transfer3}, removeC2Ts(transfers, toRemove))
}

func (s *Create2TransferCommitmentsTestSuite) getRegisteredAccounts(startBlockNumber uint64) []models.AccountLeaf {
	it, err := s.client.AccountRegistry.FilterBatchPubkeyRegistered(&bind.FilterOpts{Start: startBlockNumber})
	s.NoError(err)

	registeredAccounts := make([]models.AccountLeaf, 0)
	for it.Next() {
		tx, _, err := s.client.ChainConnection.GetBackend().TransactionByHash(context.Background(), it.Event.Raw.TxHash)
		s.NoError(err)

		unpack, err := s.client.AccountRegistryABI.Methods["registerBatch"].Inputs.Unpack(tx.Data()[4:])
		s.NoError(err)

		pubKeyIDs := eth.ExtractPubKeyIDsFromBatchAccountEvent(it.Event)
		pubKeys := unpack[0].([16][4]*big.Int)
		for i := range pubKeys {
			registeredAccounts = append(registeredAccounts, models.AccountLeaf{
				PubKeyID:  pubKeyIDs[i],
				PublicKey: models.MakePublicKeyFromInts(pubKeys[i]),
			})
		}
	}
	return registeredAccounts
}

func TestCreate2TransferCommitmentsTestSuite(t *testing.T) {
	suite.Run(t, new(Create2TransferCommitmentsTestSuite))
}

func (s *Create2TransferCommitmentsTestSuite) addCreate2Transfers(transfers []models.Create2Transfer) {
	for i := range transfers {
		_, err := s.storage.AddCreate2Transfer(&transfers[i])
		s.NoError(err)
	}
}

func (s *Create2TransferCommitmentsTestSuite) preparePendingCreate2Transfers(transfersAmount uint32) {
	transfers := generateValidCreate2Transfers(transfersAmount)
	s.addCreate2Transfers(transfers)
}
