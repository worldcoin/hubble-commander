package disputer

import (
	"fmt"

	"github.com/Worldcoin/hubble-commander/commander/syncer"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

// Other test suites encapsulate disputeTransitionTestSuite.
// Don't add any tests on disputeTransitionTestSuite to avoid repeated runs.
type disputeTransitionTestSuite struct {
	testSuiteWithContexts
}

func (s *disputeTransitionTestSuite) getInvalidBatchStateProofs(remoteBatch eth.DecodedBatch) []models.StateMerkleProof {
	s.beginTransaction()
	defer s.rollback()

	err := s.syncCtx.SyncCommitments(remoteBatch)
	s.Error(err)

	var disputableErr *syncer.DisputableError
	s.ErrorAs(err, &disputableErr)
	s.Equal(syncer.Transition, disputableErr.Type)
	return disputableErr.Proofs
}

func (s *disputeTransitionTestSuite) getValidBatchStateProofs(syncedTxs syncer.SyncedTxs) []models.StateMerkleProof {
	feeReceiverStateID := uint32(0)

	s.beginTransaction()
	defer s.rollback()

	_, stateProofs, err := s.syncCtx.SyncTxs(syncedTxs, feeReceiverStateID)
	s.NoError(err)

	return stateProofs
}

func (s *disputeTransitionTestSuite) submitInvalidBatch(txs models.GenericTransactionArray) {
	s.beginTransaction()
	defer s.rollback()
	err := s.disputeCtx.storage.BatchAddTransaction(txs)
	s.NoError(err)

	pendingBatch, err := s.txsCtx.NewPendingBatch(s.txsCtx.BatchType)
	s.NoError(err)
	fmt.Println(*pendingBatch.PrevStateRoot)

	commitments, err := s.txsCtx.CreateCommitments()
	s.NoError(err)

	commitments[len(commitments)-1].ToCommitment().GetCommitmentBase().PostStateRoot = common.Hash{1, 2, 3}

	err = s.txsCtx.SubmitBatch(pendingBatch, commitments)
	s.NoError(err)

	s.client.GetBackend().Commit()
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
