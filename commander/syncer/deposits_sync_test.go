package syncer

import (
	"context"
	"testing"

	"github.com/Worldcoin/hubble-commander/commander/executor"
	"github.com/Worldcoin/hubble-commander/contracts/erc20"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/metrics"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type SyncDepositBatchTestSuite struct {
	*require.Assertions
	suite.Suite
	storage        *st.TestStorage
	client         *eth.TestClient
	syncCtx        *Context
	depositsCtx    *executor.DepositsContext
	depositSubtree models.PendingDepositSubtree
}

func (s *SyncDepositBatchTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())

	s.depositSubtree = models.PendingDepositSubtree{
		ID:       models.MakeUint256(1),
		Root:     utils.RandomHash(),
		Deposits: testutils.GetFourDeposits(),
	}
}

func (s *SyncDepositBatchTestSuite) SetupTest() {
	var err error
	s.storage, err = st.NewTestStorage()
	s.NoError(err)

	s.client, err = eth.NewTestClient()
	s.NoError(err)

	s.prepareDeposits()

	s.depositsCtx = executor.NewDepositsContext(s.storage.Storage, s.client.Client, nil, metrics.NewCommanderMetrics(), context.Background())
	s.syncCtx = NewTestContext(s.storage.Storage, s.client.Client, nil, batchtype.Deposit)
}

func (s *SyncDepositBatchTestSuite) TearDownTest() {
	s.client.Close()
	err := s.storage.Close()
	s.NoError(err)
}

func (s *SyncDepositBatchTestSuite) TestSyncBatch_SingleBatch() {
	_, _, err := s.depositsCtx.CreateAndSubmitBatch()
	s.NoError(err)
	s.client.GetBackend().Commit()
	s.depositsCtx.Rollback(nil)

	prevStateRoot, err := s.storage.StateTree.Root()
	s.NoError(err)

	s.syncBatches()

	batches, err := s.storage.GetBatchesInRange(nil, nil)
	s.NoError(err)
	s.Len(batches, 2)
	s.Equal(prevStateRoot, batches[1].PrevStateRoot)

	_, err = s.storage.GetFirstPendingDepositSubtree()
	s.True(st.IsNotFoundError(err))

	commitment, err := s.storage.GetCommitment(&models.CommitmentID{
		BatchID:      batches[1].ID,
		IndexInBatch: 0,
	})
	s.NoError(err)
	s.Equal(s.depositSubtree.Root, commitment.ToDepositCommitment().SubtreeRoot)
}

func (s *SyncDepositBatchTestSuite) TestSyncBatch_SetsUserStates() {
	_, _, err := s.depositsCtx.CreateAndSubmitBatch()
	s.NoError(err)
	s.client.GetBackend().Commit()
	s.depositsCtx.Rollback(nil)

	s.syncBatches()

	batches, err := s.storage.GetBatchesInRange(nil, nil)
	s.NoError(err)
	s.Len(batches, 2)

	for i := range s.depositSubtree.Deposits {
		stateLeaf, err := s.storage.StateTree.Leaf(uint32(i))
		s.NoError(err)
		s.Equal(s.depositSubtree.Deposits[i].ToPubKeyID, stateLeaf.PubKeyID)
		s.Equal(s.depositSubtree.Deposits[i].L2Amount, stateLeaf.Balance)
		s.Equal(s.depositSubtree.Deposits[i].TokenID, stateLeaf.TokenID)
	}
}

func (s *SyncDepositBatchTestSuite) TestSyncBatch_SyncsExistingBatch() {
	_, _, err := s.depositsCtx.CreateAndSubmitBatch()
	s.NoError(err)
	s.client.GetBackend().Commit()
	err = s.depositsCtx.Commit()
	s.NoError(err)

	s.syncBatches()

	batches, err := s.storage.GetBatchesInRange(nil, nil)
	s.NoError(err)
	s.Len(batches, 2)

	batch, err := s.storage.GetBatch(batches[1].ID)
	s.NoError(err)
	s.NotNil(batch.Hash)
	s.NotNil(batch.PrevStateRoot)
}

func (s *SyncDepositBatchTestSuite) syncBatches() {
	remoteBatches, err := s.client.GetAllBatches()
	s.NoError(err)

	for i := range remoteBatches {
		err = s.syncCtx.SyncBatch(remoteBatches[i])
		s.NoError(err)
	}
}

func (s *SyncDepositBatchTestSuite) prepareDeposits() {
	err := s.storage.AddPendingDepositSubtree(&s.depositSubtree)
	s.NoError(err)

	s.approveTokens()
	s.queueFourDeposits()
	s.addGenesisBatch()
}

func (s *SyncDepositBatchTestSuite) addGenesisBatch() {
	root, err := s.storage.StateTree.Root()
	s.NoError(err)

	contractBatch, err := s.client.GetContractBatch(models.NewUint256(0))
	s.NoError(err)

	batch := contractBatch.ToModelBatch()
	batch.PrevStateRoot = root
	err = s.storage.AddBatch(batch)
	s.NoError(err)
}

func (s *SyncDepositBatchTestSuite) approveTokens() {
	token, err := erc20.NewERC20(s.client.ExampleTokenAddress, s.client.GetBackend())
	s.NoError(err)

	_, err = token.Approve(s.client.GetAccount(), s.client.ChainState.DepositManager, utils.ParseEther("100"))
	s.NoError(err)

	s.client.GetBackend().Commit()
}

func (s *SyncDepositBatchTestSuite) queueFourDeposits() {
	for i := 0; i < 4; i++ {
		s.queueDeposit()
	}
}

func (s *SyncDepositBatchTestSuite) queueDeposit() *models.PendingDeposit {
	toPubKeyID := models.NewUint256(1)
	tokenID := models.NewUint256(0)
	l1Amount := models.NewUint256FromBig(*utils.ParseEther("10"))
	depositID, l2Amount, err := s.client.QueueDepositAndWait(toPubKeyID, l1Amount, tokenID)
	s.NoError(err)

	return &models.PendingDeposit{
		ID:         *depositID,
		ToPubKeyID: uint32(toPubKeyID.Uint64()),
		TokenID:    *tokenID,
		L2Amount:   *l2Amount,
	}
}

func TestSyncDepositBatchTestSuite(t *testing.T) {
	suite.Run(t, new(SyncDepositBatchTestSuite))
}
