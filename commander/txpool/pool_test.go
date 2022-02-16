package txpool

import (
	"context"
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type TxPoolTestSuite struct {
	*require.Assertions
	suite.Suite
	storage *st.TestStorage
	txPool  *TxPool
}

func (s *TxPoolTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *TxPoolTestSuite) SetupTest() {
	var err error
	s.storage, err = st.NewTestStorage()
	s.NoError(err)

	s.txPool, err = NewTxPool(s.storage.Storage)
	s.NoError(err)
}

func (s *TxPoolTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *TxPoolTestSuite) TestReadTxs() {
	_, err := s.storage.StateTree.Set(0, &models.UserState{
		PubKeyID: 0,
		TokenID:  models.MakeUint256(0),
		Balance:  models.MakeUint256(100),
		Nonce:    models.MakeUint256(0),
	})
	s.NoError(err)

	for i := 0; i < 5; i++ {
		s.txPool.TxChan <- s.newTransfer(0, uint64(i))
	}

	err = s.txPool.ReadTxs(context.Background())
	s.NoError(err)

	receivedTxs := s.getAllTxs(0)
	s.Len(receivedTxs, 5)
}

func (s *TxPoolTestSuite) newTransfer(from uint32, nonce uint64) *models.Transfer {
	return testutils.NewTransfer(from, 1, nonce, 100)
}

func (s *TxPoolTestSuite) getAllTxs(stateID uint32) []models.GenericTransaction {
	txs := s.txPool.Pool.GetExecutableTxs(txtype.Transfer)

	_, txMempool := s.txPool.Pool.BeginTransaction()
	for {
		tx := txMempool.GetNextExecutableTx(txtype.Transfer, stateID)
		if tx == nil {
			break
		}
		txs = append(txs, tx)
	}
	return txs
}

func TestPoolTestSuite(t *testing.T) {
	suite.Run(t, new(TxPoolTestSuite))
}
