package executor

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/stretchr/testify/suite"
)

type Create2TransferCommitmentsTestSuite struct {
	TestSuiteWithRollupContext
	maxTxBytesInCommitment int
}

func (s *Create2TransferCommitmentsTestSuite) SetupTest() {
	s.TestSuiteWithRollupContext.SetupTestWithConfig(batchtype.Create2Transfer, config.RollupConfig{
		MinTxsPerCommitment:    1,
		MaxTxsPerCommitment:    4,
		FeeReceiverPubKeyID:    2,
		MaxCommitmentsPerBatch: 1,
	})
	s.maxTxBytesInCommitment = encoder.Create2TransferLength * int(s.cfg.MaxTxsPerCommitment)

	err := populateAccounts(s.storage.Storage, genesisBalances)
	s.NoError(err)
}

func (s *Create2TransferCommitmentsTestSuite) TestCreateCommitments_WithMinTxsPerCommitment() {
	transfers := testutils.GenerateValidCreate2Transfers(1)
	s.addCreate2Transfers(transfers)

	preRoot, err := s.executionCtx.storage.StateTree.Root()
	s.NoError(err)

	expectedTxsLength := encoder.Create2TransferLength * len(transfers)
	commitments, err := s.rollupCtx.CreateCommitments()
	s.NoError(err)
	s.Len(commitments, 1)
	s.Len(commitments[0].Transactions, expectedTxsLength)

	postRoot, err := s.executionCtx.storage.StateTree.Root()
	s.NoError(err)
	s.NotEqual(preRoot, postRoot)
	s.Equal(commitments[0].PostStateRoot, *postRoot)
}

func (s *Create2TransferCommitmentsTestSuite) TestCreateCommitments_WithMoreThanMinTxsPerCommitment() {
	transfers := testutils.GenerateValidCreate2Transfers(3)
	s.addCreate2Transfers(transfers)

	preRoot, err := s.executionCtx.storage.StateTree.Root()
	s.NoError(err)

	expectedTxsLength := encoder.Create2TransferLength * len(transfers)
	commitments, err := s.rollupCtx.CreateCommitments()
	s.NoError(err)
	s.Len(commitments, 1)
	s.Len(commitments[0].Transactions, expectedTxsLength)

	postRoot, err := s.executionCtx.storage.StateTree.Root()
	s.NoError(err)
	s.NotEqual(preRoot, postRoot)
	s.Equal(commitments[0].PostStateRoot, *postRoot)
}

func (s *Create2TransferCommitmentsTestSuite) TestCreateCommitments_QueriesForMorePendingTransfersUntilSatisfied() {
	addAccountWithHighNonce(s.Assertions, s.storage.Storage, 124)

	transfers := testutils.GenerateValidCreate2Transfers(6)
	s.invalidateCreate2Transfers(transfers[3:6])

	highNonceTransfer := testutils.MakeCreate2Transfer(124, nil, 10, 1, &models.PublicKey{5, 4, 3, 2, 1})
	transfers = append(transfers, highNonceTransfer)

	s.addCreate2Transfers(transfers)

	preRoot, err := s.executionCtx.storage.StateTree.Root()
	s.NoError(err)

	commitments, err := s.rollupCtx.CreateCommitments()
	s.NoError(err)
	s.Len(commitments, 1)
	s.Len(commitments[0].Transactions, s.maxTxBytesInCommitment)

	postRoot, err := s.executionCtx.storage.StateTree.Root()
	s.NoError(err)
	s.NotEqual(preRoot, postRoot)
	s.Equal(commitments[0].PostStateRoot, *postRoot)
}

func (s *Create2TransferCommitmentsTestSuite) TestCreateCommitments_ForMultipleCommitmentsInBatch() {
	s.cfg = &config.RollupConfig{
		MinTxsPerCommitment:    1,
		MaxTxsPerCommitment:    4,
		FeeReceiverPubKeyID:    2,
		MaxCommitmentsPerBatch: 3,
	}

	executionCtx := NewTestExecutionContext(s.storage.Storage, s.client.Client, s.cfg)
	s.rollupCtx = NewTestRollupContext(executionCtx, s.rollupCtx.BatchType)

	addAccountWithHighNonce(s.Assertions, s.storage.Storage, 124)

	transfers := testutils.GenerateValidCreate2Transfers(9)
	s.invalidateCreate2Transfers(transfers[7:9])

	highNonceTransfers := []models.Create2Transfer{
		testutils.MakeCreate2Transfer(124, nil, 10, 1, &models.PublicKey{5, 4, 3, 2, 1}),
		testutils.MakeCreate2Transfer(124, nil, 11, 1, &models.PublicKey{5, 4, 3, 2, 1}),
	}

	transfers = append(transfers, highNonceTransfers...)
	s.addCreate2Transfers(transfers)

	preRoot, err := s.executionCtx.storage.StateTree.Root()
	s.NoError(err)

	commitments, err := s.rollupCtx.CreateCommitments()
	s.NoError(err)
	s.Len(commitments, 3)
	s.Len(commitments[0].Transactions, s.maxTxBytesInCommitment)
	s.Len(commitments[1].Transactions, s.maxTxBytesInCommitment)
	s.Len(commitments[2].Transactions, encoder.Create2TransferLength)

	postRoot, err := s.executionCtx.storage.StateTree.Root()
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

func (s *Create2TransferCommitmentsTestSuite) TestCreateCommitments_ReturnsErrorWhenThereAreNotEnoughPendingTransfers() {
	preRoot, err := s.executionCtx.storage.StateTree.Root()
	s.NoError(err)

	commitments, err := s.rollupCtx.CreateCommitments()
	s.Nil(commitments)
	s.ErrorIs(err, ErrNotEnoughTxs)

	postRoot, err := s.executionCtx.storage.StateTree.Root()
	s.NoError(err)

	s.Equal(preRoot, postRoot)
}

func (s *Create2TransferCommitmentsTestSuite) TestCreateCommitments_ReturnsErrorWhenThereAreNotEnoughValidTransfers() {
	s.cfg = &config.RollupConfig{
		MinTxsPerCommitment:    32,
		MaxTxsPerCommitment:    32,
		FeeReceiverPubKeyID:    2,
		MaxCommitmentsPerBatch: 1,
	}

	executionCtx := NewTestExecutionContext(s.storage.Storage, s.client.Client, s.cfg)
	s.rollupCtx = NewTestRollupContext(executionCtx, s.rollupCtx.BatchType)

	transfers := testutils.GenerateValidCreate2Transfers(2)
	s.addCreate2Transfers(transfers)

	pendingTransfers, err := s.storage.GetPendingCreate2Transfers(32)
	s.NoError(err)
	s.Len(pendingTransfers, 2)

	preRoot, err := s.executionCtx.storage.StateTree.Root()
	s.NoError(err)

	commitments, err := s.rollupCtx.CreateCommitments()
	s.Nil(commitments)
	s.ErrorIs(err, ErrNotEnoughTxs)

	postRoot, err := s.executionCtx.storage.StateTree.Root()
	s.NoError(err)

	s.Equal(preRoot, postRoot)
}

func (s *Create2TransferCommitmentsTestSuite) TestCreateCommitments_StoresCorrectCommitment() {
	transfersCount := uint32(4)
	s.preparePendingCreate2Transfers(transfersCount)

	preRoot, err := s.executionCtx.storage.StateTree.Root()
	s.NoError(err)

	expectedTxsLength := encoder.Create2TransferLength * int(transfersCount)
	commitments, err := s.rollupCtx.CreateCommitments()
	s.NoError(err)
	s.Len(commitments, 1)
	s.Len(commitments[0].Transactions, expectedTxsLength)
	s.Equal(commitments[0].FeeReceiver, uint32(2))

	postRoot, err := s.executionCtx.storage.StateTree.Root()
	s.NoError(err)
	s.NotEqual(preRoot, postRoot)
	s.Equal(commitments[0].PostStateRoot, *postRoot)
}

func (s *Create2TransferCommitmentsTestSuite) TestCreateCommitments_CreatesMaximallyAsManyCommitmentsAsSpecifiedInConfig() {
	s.preparePendingCreate2Transfers(5)

	commitments, err := s.rollupCtx.CreateCommitments()
	s.NoError(err)
	s.Len(commitments, 1)
}

func (s *Create2TransferCommitmentsTestSuite) TestCreateCommitments_UpdateTransfers() {
	s.preparePendingCreate2Transfers(2)

	pendingTransfers, err := s.storage.GetPendingCreate2Transfers(2)
	s.NoError(err)
	s.Len(pendingTransfers, 2)

	commitments, err := s.rollupCtx.CreateCommitments()
	s.NoError(err)
	s.Len(commitments, 1)

	for i := range pendingTransfers {
		tx, err := s.storage.GetCreate2Transfer(pendingTransfers[i].Hash)
		s.NoError(err)
		s.Equal(commitments[0].ID, *tx.CommitmentID)
		s.Equal(uint32(i+3), *tx.ToStateID)
	}
}

func (s *Create2TransferCommitmentsTestSuite) TestRemoveTxs() {
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

	transfers := models.Create2TransferArray{transfer1, transfer2, transfer3}
	toRemove := models.Create2TransferArray{transfer2}

	s.Equal(models.Create2TransferArray{transfer1, transfer3}, removeTxs(transfers, toRemove))
}

func TestCreate2TransferCommitmentsTestSuite(t *testing.T) {
	suite.Run(t, new(Create2TransferCommitmentsTestSuite))
}

func (s *Create2TransferCommitmentsTestSuite) addCreate2Transfers(transfers []models.Create2Transfer) {
	for i := range transfers {
		err := s.storage.AddCreate2Transfer(&transfers[i])
		s.NoError(err)
	}
}

func (s *Create2TransferCommitmentsTestSuite) preparePendingCreate2Transfers(transfersAmount uint32) {
	transfers := testutils.GenerateValidCreate2Transfers(transfersAmount)
	s.addCreate2Transfers(transfers)
}
