package disputer

import (
	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

// Other test suites encapsulate disputeTransitionTestSuite.
// Don't add any tests on disputeTransitionTestSuite to avoid repeated runs.
type disputeTransitionTestSuite struct {
	testSuiteWithContexts
}

func (s *disputeTransitionTestSuite) applyTransfer(
	tx models.GenericTransaction,
	invalidTxHash common.Hash,
	combinedFee models.Uint256,
	receiverLeaf *models.StateLeaf,
) models.Uint256 {
	if tx.GetBase().Hash != invalidTxHash {
		txError, appError := s.txsCtx.Applier.ApplyTx(tx, receiverLeaf, models.MakeUint256(0))
		s.NoError(txError)
		s.NoError(appError)
	} else {
		senderLeaf, err := s.disputeCtx.storage.StateTree.Leaf(tx.GetFromStateID())
		s.NoError(err)
		s.calculateStateAfterInvalidTransfer(senderLeaf, receiverLeaf, tx)
	}

	fee := tx.GetFee()
	return *combinedFee.Add(&fee)
}

func (s *disputeTransitionTestSuite) calculateStateAfterInvalidTransfer(
	senderState, receiverState *models.StateLeaf,
	invalidTransfer models.GenericTransaction,
) {
	senderState.Nonce = *senderState.Nonce.AddN(1)
	amount := invalidTransfer.GetAmount()
	receiverState.Balance = *receiverState.Balance.Add(&amount)
	_, err := s.disputeCtx.storage.StateTree.Set(invalidTransfer.GetFromStateID(), &senderState.UserState)
	s.NoError(err)
	_, err = s.disputeCtx.storage.StateTree.Set(*invalidTransfer.GetToStateID(), &receiverState.UserState)
	s.NoError(err)
}

func setUserStates(s *require.Assertions, disputeCtx *Context, domain *bls.Domain) []bls.Wallet {
	userStates := []models.UserState{
		*createUserState(0, 300),
		*createUserState(1, 200),
		*createUserState(2, 100),
	}

	wallets := testutils.GenerateWallets(s, domain, len(userStates))
	for i := range userStates {
		pubKeyID, err := disputeCtx.client.RegisterAccountAndWait(wallets[i].PublicKey())
		s.NoError(err)
		s.Equal(userStates[i].PubKeyID, *pubKeyID)

		_, err = disputeCtx.storage.StateTree.Set(uint32(i), &userStates[i])
		s.NoError(err)
	}
	return wallets
}

func createUserState(pubKeyID uint32, balance uint64) *models.UserState {
	return &models.UserState{
		PubKeyID: pubKeyID,
		TokenID:  models.MakeUint256(0),
		Balance:  models.MakeUint256(balance),
		Nonce:    models.MakeUint256(0),
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
