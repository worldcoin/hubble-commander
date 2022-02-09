package mempool

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/ethereum/go-ethereum/common"
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
	}
	s.initialNonces = map[uint32]uint{}
	s.initialNonces[0] = 10
	s.initialNonces[1] = 11
}

func (s *MempoolTestSuite) TestNewMempool() {
	NewMempool(s.initialTransactions, s.initialNonces)
}

func TestMempoolTestSuite(t *testing.T) {
	suite.Run(t, new(MempoolTestSuite))
}

func createTx(from, nonce uint32) models.GenericTransaction {
	return &models.Transfer{
		TransactionBase: models.TransactionBase{
			Hash:         common.Hash{},
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
