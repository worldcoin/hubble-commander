package mempool

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
)

type MempoolTestSuite struct {
	*require.Assertions
	suite.Suite
	initialTransactions []models.GenericTransaction
	initialNonces       map[uint32]uint
}

func (s *MempoolTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())

	s.initialTransactions = []models.GenericTransaction{
		createTx(0, 10), // executable
		createTx(0, 11),
		createTx(0, 13), // gap

		createTx(1, 12), // non-executable
		createTx(1, 13),

		createTx(2, 15), // executable
		createTx(2, 16),
	}
	s.initialNonces = map[uint32]uint{}
	s.initialNonces[0] = 10
	s.initialNonces[1] = 11
	s.initialNonces[2] = 15
}

func (s *MempoolTestSuite) TestNewMempool() {
	mempool := NewMempool(s.initialTransactions, s.initialNonces)

	executable := mempool.getExecutableTxs(txtype.Transfer)
	s.Len(executable, 2)
	s.Equal(s.initialTransactions[0], executable[0])
	s.Equal(s.initialTransactions[5], executable[1])
}

func (s *MempoolTestSuite) TestAddTransaction() {
	mempool := NewMempool(s.initialTransactions, s.initialNonces)

	tx := createTx(3, 10)
	mempool.addOrReplace(tx, 10)

	executable := mempool.getExecutableTxs(txtype.Transfer)
	s.Len(executable, 3)
	s.Equal(s.initialTransactions[0], executable[0])
	s.Equal(s.initialTransactions[5], executable[1])
	s.Equal(tx, executable[2])
}

func (s *MempoolTestSuite) TestReplaceTransaction() {
	mempool := NewMempool(s.initialTransactions, s.initialNonces)

	tx := createTx(0, 10)
	mempool.addOrReplace(tx, 10)

	executable := mempool.getExecutableTxs(txtype.Transfer)
	s.Len(executable, 2)
	s.Equal(tx, executable[0])
	s.Equal(s.initialTransactions[5], executable[1])
}

func TestMempoolTestSuite(t *testing.T) {
	suite.Run(t, new(MempoolTestSuite))
}

func createTx(from, nonce uint32) models.GenericTransaction {
	return &models.Transfer{
		TransactionBase: models.TransactionBase{
			Hash:         utils.RandomHash(),
			TxType:       txtype.Transfer,
			FromStateID:  from,
			Amount:       models.Uint256{},
			Fee:          models.Uint256{},
			Nonce:        models.MakeUint256(uint64(nonce)),
			Signature:    models.Signature{},
			ReceiveTime:  nil,
			CommitmentID: nil,
			ErrorMessage: nil,
		},
		ToStateID: 0,
	}
}
