// a test suite specifically for the interface between the syncer and the
// mempool, ensures the syncer makes the correct changes to the mempool when
// new batches arrive
package syncer

import (
	"context"
	"testing"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/Worldcoin/hubble-commander/models/stored"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/stretchr/testify/suite"
)

type MempoolSyncTestSuite struct {
	syncTestSuite
}

func (s *MempoolSyncTestSuite) SetupTest() {
	var config = config.RollupConfig{
		MinCommitmentsPerBatch: 1,
		MaxCommitmentsPerBatch: 32,
		MinTxsPerCommitment:    1,
		MaxTxsPerCommitment:    32,
		DisableSignatures:      true,
	}

	s.testSuiteWithSyncAndRollupContext.SetupTestWithConfig(batchtype.Transfer, &config)
	s.syncTestSuite.setupTest()
}

// this is the simple case: we have some transfers which make it into the batch and
// then those same transfers are still in the mempool when we apply the batch. These
// should be removed from the mempool without emitting any log messages
func (s *MempoolSyncTestSuite) Test_BatchContainsMempoolTxs() {
	// stateID=0 starts with balance 1000
	// stateID=1 starts with balance 0
	// NewTransfer(from, to, nonce, amount) bundles a fee of 10

	// 1. insert some transactions into the mempool
	s.addTransfer(0, 1, 0, 20)
	s.addTransfer(1, 0, 0, 10)

	s.assertPendingState(0, 1, 1000-30+10) // TODO: this is incorrect b/c stateID=0 is the fee receiver
	s.assertPendingState(1, 1, 0)

	savedMempoolState, err := s.storage.GetAllMempoolTransactions()
	s.NoError(err)

	// 2. create a batch using them
	s.mineBatch()

	// 3. reset the database but put the mempool back how it was:
	s.resetDatabaseWithMempool(savedMempoolState)

	// 4. now... sync this batch from the chain!
	s.syncAllBatches()

	// 5. those txns should have been evicted from the mempool and the mempool
	//    should be correct:
	s.assertMempoolSize(txtype.Transfer, 0)
	s.assertPendingState(0, 1, 1000-30+10)
	s.assertPendingState(1, 1, 0)
}

// this case is more complicated: the batch contains some unexpected transactions and
// the mempool should evict the conflicting transactions as well as update its idea
// of the current state to reflect the new reality
func (s *MempoolSyncTestSuite) Test_BatchContainsUnknownCompatibleTxs() {
	s.addAccount(2, 0, 1000)

	s.addTransfer(2, 1, 0, 20)
	s.addTransfer(2, 1, 1, 20)
	s.addTransfer(1, 2, 0, 10)

	s.mineBatch()

	s.assertPendingState(2, 2, 1000-30-30+10)
	s.assertPendingState(1, 1, 0+20+20-20)

	// now we throw some mismatched transactions into the mempool:
	s.recreateDatabase()
	s.addAccount(2, 0, 1000)

	s.addTransfer(2, 1, 0, 30)
	s.addTransfer(1, 2, 0, 5)

	s.syncAllBatches()

	// our mismatched transactions should have been thrown away:

	s.assertMempoolSize(txtype.Transfer, 0)
	s.assertPendingState(2, 2, 1000-30-30+10)
	s.assertPendingState(1, 1, 0+20+20-20)
}

// this case is unimplemented:
func (s *MempoolSyncTestSuite) Test_BatchContainsUnknownIncompatibleTxs() {

}

func (s *MempoolSyncTestSuite) addAccount(stateID, nonce, balance uint32) {
	wallet, err := bls.NewRandomWallet(*s.domain)
	s.NoError(err)

	err = s.storage.AccountTree.SetSingle(&models.AccountLeaf{
		PubKeyID:  stateID,
		PublicKey: *wallet.PublicKey(),
	})
	s.NoError(err)

	_, err = s.storage.StateTree.Set(stateID, &models.UserState{
		PubKeyID: 1,
		TokenID:  models.MakeUint256(0),
		Balance:  models.MakeUint256(uint64(balance)),
		Nonce:    models.MakeUint256(uint64(nonce)),
	})
	s.NoError(err)
}

func (s *MempoolSyncTestSuite) addTransfer(from, to, nonce, amount uint32) {
	err := s.storage.AddMempoolTx(testutils.NewTransfer(from, to, uint64(nonce), uint64(amount)))
	s.NoError(err)
}

func (s *MempoolSyncTestSuite) assertMempoolSize(txType txtype.TransactionType, count uint32) {
	size, err := s.storage.CountPendingTxsOfType(txType)
	s.NoError(err)
	s.Equal(count, size)
}

func (s *MempoolSyncTestSuite) assertPendingState(stateID, nonce, balance uint32) {
	bal, err := s.storage.GetPendingBalance(stateID)
	s.NoError(err)
	s.Equal(models.MakeUint256(uint64(balance)), *bal)

	non, err := s.storage.GetPendingNonce(stateID)
	s.NoError(err)
	s.Equal(models.MakeUint256(uint64(nonce)), *non)
}

func (s *MempoolSyncTestSuite) resetDatabaseWithMempool(mempool []stored.PendingTx) {
	s.recreateDatabase()

	for i := range mempool {
		tx := mempool[i]
		err := s.storage.AddMempoolTx(tx.ToGenericTransaction())
		s.NoError(err)
	}
}

func (s *MempoolSyncTestSuite) mineBatch() {
	pendingBatch, err := s.txsCtx.NewPendingBatch(s.txsCtx.BatchType)
	s.NoError(err)

	commitments, err := s.txsCtx.CreateCommitments(context.Background())
	s.NoError(err)
	s.Len(commitments, 1)

	err = s.txsCtx.SubmitBatch(context.Background(), pendingBatch, commitments)
	s.NoError(err)

	s.client.GetBackend().Commit()
}

func TestMempoolSyncTestSuite(t *testing.T) {
	suite.Run(t, new(MempoolSyncTestSuite))
}
