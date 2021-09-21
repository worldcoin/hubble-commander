package executor

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/stretchr/testify/suite"
)

type RevertBatchesTestSuite struct {
	TestSuiteWithRollupContext
	transfer models.Transfer
	wallets  []bls.Wallet
}

func (s *RevertBatchesTestSuite) SetupTest() {
	s.TestSuiteWithRollupContext.SetupTestWithConfig(batchtype.Transfer, config.RollupConfig{
		MinCommitmentsPerBatch: 1,
		MaxCommitmentsPerBatch: 32,
		MinTxsPerCommitment:    1,
		MaxTxsPerCommitment:    1,
	})

	s.transfer = testutils.MakeTransfer(0, 1, 0, 400)
	s.setTransferHash(&s.transfer)

	domain, err := s.client.GetDomain()
	s.NoError(err)

	s.wallets = generateWallets(s.Assertions, domain, 2)

	seedDB(s.Assertions, s.storage.Storage, s.wallets)
}

func (s *RevertBatchesTestSuite) TestRevertBatches_RevertsState() {
	initialStateRoot, err := s.storage.StateTree.Root()
	s.NoError(err)

	signTransfer(s.T(), &s.wallets[s.transfer.FromStateID], &s.transfer)
	pendingBatch := submitTransferBatch(s.Assertions, s.client, s.rollupCtx, &s.transfer)

	err = s.executionCtx.RevertBatches(pendingBatch)
	s.NoError(err)

	stateRoot, err := s.storage.StateTree.Root()
	s.NoError(err)
	s.Equal(*initialStateRoot, *stateRoot)

	state0, err := s.storage.StateTree.Leaf(s.transfer.FromStateID)
	s.NoError(err)
	s.Equal(uint64(1000), state0.Balance.Uint64())

	state1, err := s.storage.StateTree.Leaf(s.transfer.ToStateID)
	s.NoError(err)
	s.Equal(uint64(0), state1.Balance.Uint64())
}

func (s *RevertBatchesTestSuite) TestRevertBatches_ExcludesTransactionsFromCommitments() {
	signTransfer(s.T(), &s.wallets[s.transfer.FromStateID], &s.transfer)
	pendingBatch := submitTransferBatch(s.Assertions, s.client, s.rollupCtx, &s.transfer)

	err := s.executionCtx.RevertBatches(pendingBatch)
	s.NoError(err)

	transfer, err := s.storage.GetTransfer(s.transfer.Hash)
	s.NoError(err)
	s.Nil(transfer.CommitmentID)
}

func (s *RevertBatchesTestSuite) TestRevertBatches_DeletesCommitmentsAndBatches() {
	transfers := make([]models.Transfer, 2)
	transfers[0] = s.transfer
	transfers[1] = testutils.MakeTransfer(0, 1, 1, 200)

	pendingBatches := make([]models.Batch, 2)
	for i := range pendingBatches {
		signTransfer(s.T(), &s.wallets[transfers[i].FromStateID], &transfers[i])
		pendingBatches[i] = *submitTransferBatch(s.Assertions, s.client, s.rollupCtx, &transfers[i])
	}

	latestCommitment, err := s.executionCtx.storage.GetLatestCommitment()
	s.NoError(err)
	s.Equal(models.MakeUint256(2), latestCommitment.ID.BatchID)

	err = s.executionCtx.RevertBatches(&pendingBatches[0])
	s.NoError(err)

	_, err = s.executionCtx.storage.GetLatestCommitment()
	s.ErrorIs(err, st.NewNotFoundError("commitment"))

	batches, err := s.storage.GetBatchesInRange(nil, nil)
	s.NoError(err)
	s.Len(batches, 0)
}

func (s *RevertBatchesTestSuite) setTransferHash(tx *models.Transfer) {
	hash, err := encoder.HashTransfer(tx)
	s.NoError(err)
	tx.Hash = *hash
}

func TestRevertBatchesTestSuite(t *testing.T) {
	suite.Run(t, new(RevertBatchesTestSuite))
}
