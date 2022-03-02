package executor

import (
	"testing"
	"time"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type FindOldestTxTimeTestSuite struct {
	*require.Assertions
	suite.Suite
	storage *st.TestStorage
	txsCtx  *TxsContext
}

func (s *FindOldestTxTimeTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *FindOldestTxTimeTestSuite) SetupTest() {
	var err error
	s.storage, err = st.NewTestStorage()
	s.NoError(err)

	_, err = s.storage.StateTree.Set(0, &models.UserState{
		PubKeyID: 0,
		TokenID:  models.MakeUint256(0),
		Balance:  models.MakeUint256(600),
		Nonce:    models.MakeUint256(0),
	})
	s.NoError(err)

	cfg := &config.RollupConfig{
		MinTxsPerCommitment:    1,
		MinCommitmentsPerBatch: 1,
	}
	executionCtx := NewTestExecutionContext(s.storage.Storage, nil, cfg)
	s.txsCtx, err = NewTestTxsContext(executionCtx, batchtype.Transfer)
	s.NoError(err)
}

func (s *FindOldestTxTimeTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *FindOldestTxTimeTestSuite) TestFindOldestTransactionTime_EmptyMempool() {
	oldest := s.txsCtx.findOldestTransactionTime()
	s.Nil(oldest)
}

func (s *FindOldestTxTimeTestSuite) TestFindOldestTransactionTime_NoTxHasTime() {
	txs := models.TransferArray{testutils.MakeTransfer(0, 1, 0, 400)}
	initTxs(s.Assertions, s.txsCtx, txs)

	oldest := s.txsCtx.findOldestTransactionTime()
	s.Nil(oldest)
}

func (s *FindOldestTxTimeTestSuite) TestFindOldestTransactionTime_FindsOldestTime() {
	oneSecondAgo := time.Now().Add(-time.Second)
	twoSecondAgo := time.Now().Add(-2 * time.Second)
	initTxs(s.Assertions, s.txsCtx, s.newTxsWithReceiveTime(oneSecondAgo, twoSecondAgo))

	oldest := s.txsCtx.findOldestTransactionTime()
	s.Equal(twoSecondAgo, oldest.Time)
}

func (s *FindOldestTxTimeTestSuite) newTxsWithReceiveTime(receiveTimes ...time.Time) models.TransferArray {
	txs := make(models.TransferArray, 0, len(receiveTimes))
	for i := range receiveTimes {
		tx := testutils.MakeTransfer(0, 1, uint64(i), 100)
		tx.ReceiveTime = models.NewTimestamp(receiveTimes[i])
		txs = append(txs, tx)
	}
	return txs
}

func TestFindOldestTxTimeTestSuite(t *testing.T) {
	suite.Run(t, new(FindOldestTxTimeTestSuite))
}
