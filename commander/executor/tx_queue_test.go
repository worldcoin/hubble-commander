package executor

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type TxQueueTestSuite struct {
	*require.Assertions
	suite.Suite
}

func (s *TxQueueTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *TxQueueTestSuite) TestTxQueue_RemoveFromQueue_Transfer() {
	transfer1 := createRandomTransferWithHash()
	transfer2 := createRandomTransferWithHash()
	transfer3 := createRandomTransferWithHash()

	transfers := NewTxQueue(models.TransferArray{transfer1, transfer2, transfer3})
	toRemove := models.TransferArray{transfer2}

	transfers.RemoveFromQueue(toRemove)

	s.Equal(models.TransferArray{transfer1, transfer3}, transfers.PickTxsForCommitment())
}

func (s *TxQueueTestSuite) TestTxQueue_RemoveFromQueue_C2T() {
	transfer1 := createRandomC2TWithHash()
	transfer2 := createRandomC2TWithHash()
	transfer3 := createRandomC2TWithHash()

	transfers := NewTxQueue(models.Create2TransferArray{transfer1, transfer2, transfer3})
	toRemove := models.Create2TransferArray{transfer2}

	transfers.RemoveFromQueue(toRemove)

	s.Equal(models.Create2TransferArray{transfer1, transfer3}, transfers.PickTxsForCommitment())
}

func createRandomTransferWithHash() models.Transfer {
	return models.Transfer{
		TransactionBase: models.TransactionBase{
			Hash: utils.RandomHash(),
		},
	}
}

func createRandomC2TWithHash() models.Create2Transfer {
	return models.Create2Transfer{
		TransactionBase: models.TransactionBase{
			Hash: utils.RandomHash(),
		},
	}
}

func TestTxQueueTestSuite(t *testing.T) {
	suite.Run(t, new(TxQueueTestSuite))
}
