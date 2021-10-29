package executor

import (
	"math/big"
	"testing"

	"github.com/Worldcoin/hubble-commander/contracts/erc20"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	st "github.com/Worldcoin/hubble-commander/storage"
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
	depositCtx     *DepositContext
	depositSubtree models.PendingDepositSubTree
}

func (s *SubmitDepositBatchTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())

	s.depositSubtree = models.PendingDepositSubTree{
		ID:       models.MakeUint256(1),
		Root:     utils.RandomHash(),
		Deposits: getFourDeposits(),
	}
}

func (s *SubmitDepositBatchTestSuite) SetupTest() {
	var err error
	s.storage, err = st.NewTestStorage()
	s.NoError(err)

	s.client, err = eth.NewTestClient()
	s.NoError(err)

	executionCtx := NewTestExecutionContext(s.storage.Storage, s.client.Client, nil)
	s.depositCtx = NewTestDepositContext(executionCtx)
}

func (s *SubmitDepositBatchTestSuite) TearDownTest() {
	s.client.Close()
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *SubmitDepositBatchTestSuite) TestSubmitDepositBatch_SubmitsBatchOnChain() {
	s.prepareDeposits()

	pendingBatch, err := s.depositCtx.NewPendingBatch(batchtype.Deposit)
	s.NoError(err)

	_, vacancyProof, err := s.depositCtx.getDepositSubtreeVacancyProof()
	s.NoError(err)

	err = s.depositCtx.SubmitBatch(pendingBatch, vacancyProof)
	s.NoError(err)

	s.client.GetBackend().Commit()

	nextBatchID, err := s.client.Rollup.NextBatchID(nil)
	s.NoError(err)
	s.Equal(big.NewInt(2), nextBatchID)
}

func (s *SubmitDepositBatchTestSuite) TestSubmitDepositBatch_StoresPendingBatch() {
	s.prepareDeposits()

	pendingBatch, err := s.depositCtx.NewPendingBatch(batchtype.Deposit)
	s.NoError(err)

	_, vacancyProof, err := s.depositCtx.getDepositSubtreeVacancyProof()
	s.NoError(err)

	err = s.depositCtx.SubmitBatch(pendingBatch, vacancyProof)
	s.NoError(err)

	s.client.GetBackend().Commit()

	batch, err := s.storage.GetBatch(pendingBatch.ID)
	s.NoError(err)
	s.Equal(pendingBatch.Type, batch.Type)
	s.NotEqual(common.Hash{}, batch.TransactionHash)
	s.Equal(pendingBatch.ID, batch.ID)
	s.Equal(pendingBatch.PrevStateRoot, batch.PrevStateRoot)
	s.Nil(batch.Hash)
}

func (s *SubmitDepositBatchTestSuite) prepareDeposits() {
	s.registerToken(s.client.ExampleTokenAddress)
	s.approveTokens()
	s.queueFourDeposits()
	s.addGenesisBatch()
}

func (s *SubmitDepositBatchTestSuite) addGenesisBatch() {
	root, err := s.storage.StateTree.Root()
	s.NoError(err)

	batch, err := s.client.GetBatch(models.NewUint256(0))
	s.NoError(err)

	batch.PrevStateRoot = root
	err = s.storage.AddBatch(batch)
	s.NoError(err)
}

func (s *SubmitDepositBatchTestSuite) registerToken(tokenAddress common.Address) *models.Uint256 {
	err := s.client.RequestRegisterTokenAndWait(tokenAddress)
	s.NoError(err)

	tokenID, err := s.client.FinalizeRegisterTokenAndWait(tokenAddress)
	s.NoError(err)

	return tokenID
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

func getFourDeposits() []models.PendingDeposit {
	deposits := make([]models.PendingDeposit, 4)
	for i := range deposits {
		deposits[i] = models.PendingDeposit{
			ID:         models.DepositID{BlockNumber: 1, LogIndex: uint32(i)},
			ToPubKeyID: 1,
			TokenID:    models.MakeUint256(0),
			L2Amount:   models.MakeUint256(10000000000),
		}
	}
	return deposits
}

func (s *SubmitDepositBatchTestSuite) queueDeposit() *models.PendingDeposit {
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

func TestSubmitDepositBatchTestSuite(t *testing.T) {
	suite.Run(t, new(SubmitDepositBatchTestSuite))
}
