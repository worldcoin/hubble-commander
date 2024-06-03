package executor

import (
	"context"
	"math/big"
	"testing"

	"github.com/Worldcoin/hubble-commander/contracts/erc20"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type SubmitDepositBatchTestSuite struct {
	*require.Assertions
	suite.Suite
	storage        *st.TestStorage
	client         *eth.TestClient
	depositsCtx    *DepositsContext
	depositSubtree models.PendingDepositSubtree
}

func (s *SubmitDepositBatchTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())

	s.depositSubtree = models.PendingDepositSubtree{
		ID:       models.MakeUint256(1),
		Root:     utils.RandomHash(),
		Deposits: testutils.GetFourDeposits(),
	}
}

func (s *SubmitDepositBatchTestSuite) SetupTest() {
	var err error
	s.storage, err = st.NewTestStorage()
	s.NoError(err)

	s.client, err = eth.NewTestClient()
	s.NoError(err)

	executionCtx := NewTestExecutionContext(s.storage.Storage, s.client.Client, nil)
	s.depositsCtx = NewTestDepositsContext(executionCtx)
}

func (s *SubmitDepositBatchTestSuite) TearDownTest() {
	s.client.Close()
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *SubmitDepositBatchTestSuite) TestSubmitDepositBatch_SubmitsBatchOnChain() {
	s.prepareDeposits()
	s.submitBatch()

	nextBatchID, err := s.client.Rollup.NextBatchID(nil)
	s.NoError(err)
	s.Equal(big.NewInt(2), nextBatchID)
}

func (s *SubmitDepositBatchTestSuite) TestSubmitDepositBatch_StoresPendingBatch() {
	s.prepareDeposits()
	pendingBatch := s.submitBatch()

	batch, err := s.storage.GetBatch(pendingBatch.ID)
	s.NoError(err)
	s.Equal(pendingBatch.Type, batch.Type)
	s.NotEqual(common.Hash{}, batch.TransactionHash)
	s.Equal(pendingBatch.ID, batch.ID)
	s.Equal(pendingBatch.PrevStateRoot, batch.PrevStateRoot)
	s.Nil(batch.Hash)
}

func (s *SubmitDepositBatchTestSuite) TestSubmitDepositBatch_TwoBatches() {
	s.prepareDeposits()
	s.submitBatch()

	s.queueFourDeposits()
	s.submitBatch()

	nextBatchID, err := s.client.Rollup.NextBatchID(nil)
	s.NoError(err)
	s.Equal(big.NewInt(3), nextBatchID)
}

func (s *SubmitDepositBatchTestSuite) prepareDeposits() {
	s.approveTokens()
	s.queueFourDeposits()
	s.addGenesisBatch()
}

func (s *SubmitDepositBatchTestSuite) addGenesisBatch() {
	root, err := s.storage.StateTree.Root()
	s.NoError(err)

	contractBatch, err := s.client.GetContractBatch(models.NewUint256(0))
	s.NoError(err)

	batch := contractBatch.ToModelBatch()
	batch.PrevStateRoot = *root
	err = s.storage.AddBatch(batch)
	s.NoError(err)
}

func (s *SubmitDepositBatchTestSuite) approveTokens() {
	token, err := erc20.NewERC20(s.client.ExampleTokenAddress, s.client.GetBackend())
	s.NoError(err)

	_, err = token.Approve(s.client.GetAccount(), s.client.ChainState.DepositManager, utils.ParseEther("100"))
	s.NoError(err)

	s.client.GetBackend().Commit()
}

func (s *SubmitDepositBatchTestSuite) queueFourDeposits() {
	for i := 0; i < 4; i++ {
		s.queueDeposit()
	}
}

func (s *SubmitDepositBatchTestSuite) queueDeposit() {
	toPubKeyID := models.NewUint256(1)
	tokenID := models.NewUint256(0)
	l1Amount := models.NewUint256FromBig(*utils.ParseEther("10"))
	_, _, err := s.client.QueueDepositAndWait(toPubKeyID, l1Amount, tokenID)
	s.NoError(err)
}

func (s *SubmitDepositBatchTestSuite) submitBatch() *models.Batch {
	err := s.storage.AddPendingDepositSubtree(&s.depositSubtree)
	s.NoError(err)

	pendingBatch, err := s.depositsCtx.NewPendingBatch(batchtype.Deposit)
	s.NoError(err)

	vacancyProof, err := s.depositsCtx.createCommitment(context.Background(), pendingBatch.ID)
	s.NoError(err)

	err = s.depositsCtx.SubmitBatch(pendingBatch, vacancyProof)
	s.NoError(err)

	s.client.GetBackend().Commit()

	return pendingBatch
}

func TestSubmitDepositBatchTestSuite(t *testing.T) {
	suite.Run(t, new(SubmitDepositBatchTestSuite))
}
