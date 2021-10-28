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
	"github.com/Worldcoin/hubble-commander/utils/merkletree"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type DepositsTestSuite struct {
	*require.Assertions
	suite.Suite
	storage        *st.TestStorage
	client         *eth.TestClient
	depositCtx     *DepositContext
	depositSubtree models.PendingDepositSubTree
}

func (s *DepositsTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())

	s.depositSubtree = models.PendingDepositSubTree{
		ID:   models.MakeUint256(1),
		Root: utils.RandomHash(),
		Deposits: []models.PendingDeposit{
			{
				ID:         models.DepositID{BlockNumber: 1, LogIndex: 0},
				ToPubKeyID: 1,
				TokenID:    models.MakeUint256(0),
				L2Amount:   models.MakeUint256(50),
			},
			{
				ID:         models.DepositID{BlockNumber: 1, LogIndex: 1},
				ToPubKeyID: 1,
				TokenID:    models.MakeUint256(0),
				L2Amount:   models.MakeUint256(50),
			},
		},
	}
}

func (s *DepositsTestSuite) SetupTest() {
	var err error
	s.storage, err = st.NewTestStorage()
	s.NoError(err)

	s.client, err = eth.NewTestClient()
	s.NoError(err)

	executionCtx := NewTestExecutionContext(s.storage.Storage, s.client.Client, nil)
	s.depositCtx = NewTestDepositContext(executionCtx)
}

func (s *DepositsTestSuite) TearDownTest() {
	s.client.Close()
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *DepositsTestSuite) TestGetVacancyProof_EmptyTree() {
	stateID, err := s.depositCtx.storage.StateTree.NextVacantSubtree(2)
	s.NoError(err)

	vacancyProof, err := s.depositCtx.GetVacancyProof(*stateID, 2)
	s.NoError(err)
	s.EqualValues(vacancyProof.PathAtDepth, 0)
	s.Len(vacancyProof.Witness, 30)
}

func (s *DepositsTestSuite) TestGetVacancyProof_SingleLeafSet() {
	_, err := s.depositCtx.storage.StateTree.Set(0, &models.UserState{})
	s.NoError(err)

	stateID, err := s.depositCtx.storage.StateTree.NextVacantSubtree(2)
	s.NoError(err)

	vacancyProof, err := s.depositCtx.GetVacancyProof(*stateID, 2)
	s.NoError(err)
	s.EqualValues(vacancyProof.PathAtDepth, 1)
	s.Len(vacancyProof.Witness, 30)
}

func (s *DepositsTestSuite) TestGetVacancyProof_TwoLeavesSet() {
	_, err := s.depositCtx.storage.StateTree.Set(0, &models.UserState{})
	s.NoError(err)
	_, err = s.depositCtx.storage.StateTree.Set(4, &models.UserState{})
	s.NoError(err)

	stateID, err := s.depositCtx.storage.StateTree.NextVacantSubtree(2)
	s.NoError(err)

	vacancyProof, err := s.depositCtx.GetVacancyProof(*stateID, 2)
	s.NoError(err)
	s.EqualValues(vacancyProof.PathAtDepth, 2)
	s.Len(vacancyProof.Witness, 30)
}

func (s *DepositsTestSuite) TestGetVacancyProof_ProducesCorrectWitness() {
	userState := &models.UserState{}
	leafWitness, err := s.depositCtx.storage.StateTree.Set(0, userState)
	s.NoError(err)

	leaf, err := st.NewStateLeaf(0, userState)
	s.NoError(err)

	currentHash := leaf.DataHash
	for i := range leafWitness[:len(leafWitness)-2] {
		currentHash = utils.HashTwo(currentHash, leafWitness[i])
	}
	firstWitness := currentHash
	secondWitness := merkletree.GetZeroHash(31)

	stateID, err := s.depositCtx.storage.StateTree.NextVacantSubtree(30)
	s.NoError(err)

	vacancyProof, err := s.depositCtx.GetVacancyProof(*stateID, 30)
	s.NoError(err)

	s.Len(vacancyProof.Witness, 2)
	s.Equal(vacancyProof.Witness[0], firstWitness)
	s.Equal(vacancyProof.Witness[1], secondWitness)
}

func (s *DepositsTestSuite) TestCreateCommitment_AddsCommitment() {
	err := s.storage.AddPendingDepositSubTree(&s.depositSubtree)
	s.NoError(err)

	batchID := models.MakeUint256(1)
	_, err = s.depositCtx.createCommitment(batchID)
	s.NoError(err)

	root, err := s.storage.StateTree.Root()
	s.NoError(err)

	commitment, err := s.storage.GetDepositCommitment(&models.CommitmentID{
		BatchID:      batchID,
		IndexInBatch: 0,
	})
	s.NoError(err)
	s.Equal(*root, commitment.PostStateRoot)
	s.Equal(s.depositSubtree.ID, commitment.SubTreeID)
	s.Equal(s.depositSubtree.Root, commitment.SubTreeRoot)
	s.Equal(s.depositSubtree.Deposits, commitment.Deposits)
}

func (s *DepositsTestSuite) TestCreateCommitment_NotEnoughDeposits() {
	vacancyProof, err := s.depositCtx.createCommitment(models.MakeUint256(1))
	s.ErrorIs(err, ErrNotEnoughDeposits)
	s.Nil(vacancyProof)
}

func (s *DepositsTestSuite) TestExecuteDeposits_SetsUserStates() {
	err := s.storage.AddPendingDepositSubTree(&s.depositSubtree)
	s.NoError(err)

	_, err = s.depositCtx.executeDeposits(&s.depositSubtree)
	s.NoError(err)

	for i := range s.depositSubtree.Deposits {
		stateLeaf, err := s.storage.StateTree.Leaf(uint32(i))
		s.NoError(err)
		s.Equal(s.depositSubtree.Deposits[i].L2Amount, stateLeaf.Balance)
	}
}

func (s *DepositsTestSuite) TestExecuteDeposits_RemovesDepositSubtree() {
	err := s.storage.AddPendingDepositSubTree(&s.depositSubtree)
	s.NoError(err)

	_, err = s.depositCtx.executeDeposits(&s.depositSubtree)
	s.NoError(err)

	subtree, err := s.storage.GetPendingDepositSubTree(s.depositSubtree.ID)
	s.True(st.IsNotFoundError(err))
	s.Nil(subtree)
}

func (s *DepositsTestSuite) TestExecuteDeposits_ReturnsCorrectVacancyProof() {
	_, err := s.depositCtx.storage.StateTree.Set(0, &models.UserState{})
	s.NoError(err)

	err = s.storage.AddPendingDepositSubTree(&s.depositSubtree)
	s.NoError(err)

	vacancyProof, err := s.depositCtx.executeDeposits(&s.depositSubtree)
	s.NoError(err)
	s.EqualValues(1, vacancyProof.PathAtDepth)
}

func (s *DepositsTestSuite) TestSubmitDepositBatch_SubmitsBatchOnChain() {
	s.registerToken(s.client.ExampleTokenAddress)
	s.approveTokens()
	s.queueFourDeposits()
	s.addGenesisBatch()

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

func (s *DepositsTestSuite) TestSubmitDepositBatch_StoresPendingBatch() {
	s.registerToken(s.client.ExampleTokenAddress)
	s.approveTokens()
	s.queueFourDeposits()
	s.addGenesisBatch()

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

func (s *DepositsTestSuite) addGenesisBatch() {
	root, err := s.storage.StateTree.Root()
	s.NoError(err)

	batch, err := s.client.GetBatch(models.NewUint256(0))
	s.NoError(err)

	batch.PrevStateRoot = root
	err = s.storage.AddBatch(batch)
	s.NoError(err)
}

func (s *DepositsTestSuite) registerToken(tokenAddress common.Address) *models.Uint256 {
	err := s.client.RequestRegisterTokenAndWait(tokenAddress)
	s.NoError(err)

	tokenID, err := s.client.FinalizeRegisterTokenAndWait(tokenAddress)
	s.NoError(err)

	return tokenID
}

func (s *DepositsTestSuite) approveTokens() {
	token, err := erc20.NewERC20(s.client.ExampleTokenAddress, s.client.GetBackend())
	s.NoError(err)

	_, err = token.Approve(s.client.GetAccount(), s.client.ChainState.DepositManager, utils.ParseEther("100"))
	s.NoError(err)

	s.client.GetBackend().Commit()
}

func (s *DepositsTestSuite) queueFourDeposits() []models.PendingDeposit {
	return []models.PendingDeposit{
		*s.queueDeposit(),
		*s.queueDeposit(),
		*s.queueDeposit(),
		*s.queueDeposit(),
	}
}

func (s *DepositsTestSuite) queueDeposit() *models.PendingDeposit {
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

func TestDepositsTestSuite(t *testing.T) {
	suite.Run(t, new(DepositsTestSuite))
}
