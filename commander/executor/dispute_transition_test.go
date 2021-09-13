package executor

import (
	"context"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

// Other test suites encapsulate DisputeTransitionTestSuite.
// Don't add any tests on DisputeTransitionTestSuite to avoid repeated runs.
type DisputeTransitionTestSuite struct {
	TestSuiteWithDisputeContext
}

func (s *DisputeTransitionTestSuite) checkBatchAfterDispute(batchID models.Uint256) {
	checkRemoteBatchAfterDispute(s.Assertions, s.client, &batchID)

	batch, err := s.storage.GetBatch(batchID)
	s.Nil(batch)
	s.True(st.IsNotFoundError(err))
}

func (s *DisputeTransitionTestSuite) beginTransaction() {
	var err error
	s.executionCtx, err = NewExecutionContext(s.storage.Storage, s.client.Client, s.cfg, context.Background())
	s.NoError(err)
	s.rollupCtx = NewTestRollupContext(s.executionCtx, s.rollupCtx.BatchType)
	s.disputeCtx = NewDisputeContext(s.executionCtx.storage, s.executionCtx.client)
}

func (s *DisputeTransitionTestSuite) applyTransfer(
	tx models.GenericTransaction,
	invalidTxHash common.Hash,
	combinedFee models.Uint256,
	receiverLeaf *models.StateLeaf,
) models.Uint256 {
	if tx.GetBase().Hash != invalidTxHash {
		transferError, appError := s.executionCtx.ApplyTransfer(tx, receiverLeaf, models.MakeUint256(0))
		s.NoError(transferError)
		s.NoError(appError)
	} else {
		senderLeaf, err := s.executionCtx.storage.StateTree.Leaf(tx.GetFromStateID())
		s.NoError(err)
		s.calculateStateAfterInvalidTransfer(senderLeaf, receiverLeaf, tx)
	}
	fee := tx.GetFee()
	return *combinedFee.Add(&fee)
}

func (s *DisputeTransitionTestSuite) calculateStateAfterInvalidTransfer(
	senderState, receiverState *models.StateLeaf,
	invalidTransfer models.GenericTransaction,
) {
	senderState.Nonce = *senderState.Nonce.AddN(1)
	amount := invalidTransfer.GetAmount()
	receiverState.Balance = *receiverState.Balance.Add(&amount)
	_, err := s.executionCtx.storage.StateTree.Set(invalidTransfer.GetFromStateID(), &senderState.UserState)
	s.NoError(err)
	_, err = s.executionCtx.storage.StateTree.Set(*invalidTransfer.GetToStateID(), &receiverState.UserState)
	s.NoError(err)
}

func (s *DisputeTransitionTestSuite) commitTransaction() {
	err := s.executionCtx.Commit()
	s.NoError(err)
}

func setUserStates(s *require.Assertions, executionCtx *ExecutionContext, domain *bls.Domain) []bls.Wallet {
	userStates := []models.UserState{
		*createUserState(0, 300, 0),
		*createUserState(1, 200, 0),
		*createUserState(2, 100, 0),
	}
	registrations, unsubscribe, err := executionCtx.client.WatchRegistrations(&bind.WatchOpts{})
	s.NoError(err)
	defer unsubscribe()

	wallets := generateWallets(s, domain, len(userStates))
	for i := range userStates {
		pubKeyID, err := executionCtx.client.RegisterAccount(wallets[i].PublicKey(), registrations)
		s.NoError(err)
		s.Equal(userStates[i].PubKeyID, *pubKeyID)

		_, err = executionCtx.storage.StateTree.Set(uint32(i), &userStates[i])
		s.NoError(err)
	}
	return wallets
}

func createUserState(pubKeyID uint32, balance, nonce uint64) *models.UserState {
	return &models.UserState{
		PubKeyID: pubKeyID,
		TokenID:  models.MakeUint256(0),
		Balance:  models.MakeUint256(balance),
		Nonce:    models.MakeUint256(nonce),
	}
}

func checkRemoteBatchAfterDispute(s *require.Assertions, client *eth.TestClient, batchID *models.Uint256) {
	_, err := client.GetBatch(batchID)
	if err == nil {
		err = client.KeepRollingBack()
		s.NoError(err)
		_, err = client.GetBatch(batchID)
	}
	s.Error(err)
	s.Equal(eth.MsgInvalidBatchID, err.Error())
}
