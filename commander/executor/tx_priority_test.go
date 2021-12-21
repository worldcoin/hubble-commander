package executor

import (
	"sort"
	"testing"
	"time"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type TxPriorityTestSuite struct {
	*require.Assertions
	suite.Suite
}

func (s *TxPriorityTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *TxPriorityTestSuite) TestHigherPriority_LowerNonceFirst() {
	txs := models.TransferArray{
		s.makeTx(3),
		s.makeTx(1),
		s.makeTx(2),
	}
	sort.Slice(byPriority(txs))

	expectedTxs := models.TransferArray{
		s.makeTx(1),
		s.makeTx(2),
		s.makeTx(3),
	}
	s.Equal(expectedTxs, txs)
}

func (s *TxPriorityTestSuite) TestHigherPriority_HigherFeeFirst() {
	txs := models.TransferArray{
		s.makeTxWithFee(3, 10),
		s.makeTxWithFee(2, 0),
		s.makeTxWithFee(1, 20),
		s.makeTxWithFee(1, 10),
		s.makeTxWithFee(1, 30),
	}
	sort.Slice(byPriority(txs))

	expectedTxs := models.TransferArray{
		s.makeTxWithFee(1, 30),
		s.makeTxWithFee(1, 20),
		s.makeTxWithFee(1, 10),
		s.makeTxWithFee(2, 0),
		s.makeTxWithFee(3, 10),
	}
	s.Equal(expectedTxs, txs)
}

func (s *TxPriorityTestSuite) TestHigherPriority_EarlierReceiveTimeFirst() {
	now := models.Timestamp{Time: time.Now()}
	earlier := now.Add(-10 * time.Second)
	later := now.Add(10 * time.Second)
	txs := models.TransferArray{
		s.makeTxWithReceiveTime(3, 10, &later),
		s.makeTxWithReceiveTime(3, 10, nil),
		s.makeTxWithReceiveTime(3, 11, &later),
		s.makeTxWithReceiveTime(3, 10, &earlier),
		s.makeTxWithReceiveTime(1, 0, &now),
		s.makeTxWithReceiveTime(3, 10, &now),
	}
	sort.Slice(byPriority(txs))

	expectedTxs := models.TransferArray{
		s.makeTxWithReceiveTime(1, 0, &now),
		s.makeTxWithReceiveTime(3, 11, &later),
		s.makeTxWithReceiveTime(3, 10, &earlier),
		s.makeTxWithReceiveTime(3, 10, &now),
		s.makeTxWithReceiveTime(3, 10, &later),
		s.makeTxWithReceiveTime(3, 10, nil),
	}
	s.Equal(expectedTxs, txs)
}

func (s *TxPriorityTestSuite) makeTx(nonce uint64) models.Transfer {
	return s.makeTxWithFee(nonce, 0)
}

func (s *TxPriorityTestSuite) makeTxWithFee(nonce, fee uint64) models.Transfer {
	return s.makeTxWithReceiveTime(nonce, fee, nil)
}

func (s *TxPriorityTestSuite) makeTxWithReceiveTime(nonce, fee uint64, receiveTime *models.Timestamp) models.Transfer {
	return models.Transfer{
		TransactionBase: models.TransactionBase{
			Nonce:       models.MakeUint256(nonce),
			Fee:         models.MakeUint256(fee),
			ReceiveTime: receiveTime,
		},
	}
}

func TestTxPriorityTestSuite(t *testing.T) {
	suite.Run(t, new(TxPriorityTestSuite))
}
