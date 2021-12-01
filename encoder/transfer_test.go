package encoder

import (
	"math/big"
	"testing"

	contractTransfer "github.com/Worldcoin/hubble-commander/contracts/frontend/transfer"
	testtx "github.com/Worldcoin/hubble-commander/contracts/test/tx"
	"github.com/Worldcoin/hubble-commander/eth/deployer"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/Worldcoin/hubble-commander/testutils/simulator"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type TransferTestSuite struct {
	*require.Assertions
	suite.Suite
	sim      *simulator.Simulator
	transfer *contractTransfer.FrontendTransfer
	testTx   *testtx.TestTx
}

func (s *TransferTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *TransferTestSuite) SetupTest() {
	sim, err := simulator.NewSimulator()
	s.NoError(err)
	s.sim = sim

	frontend, err := deployer.DeployFrontend(sim)
	s.NoError(err)
	test, err := testutils.DeployTest(sim)
	s.NoError(err)

	s.transfer = frontend.FrontendTransfer
	s.testTx = test.TestTx
}

func (s *TransferTestSuite) TearDownTest() {
	s.sim.Close()
}

func (s *TransferTestSuite) TestEncodeTransfer() {
	encodedTransfer, err := EncodeTransfer(&models.Transfer{
		TransactionBase: models.TransactionBase{
			FromStateID: 2,
			Amount:      models.MakeUint256(4),
			Fee:         models.MakeUint256(5),
			Nonce:       models.MakeUint256(6),
		},
		ToStateID: 3,
	})
	s.NoError(err)
	expected, err := s.transfer.Encode(nil, contractTransfer.OffchainTransfer{
		TxType:    big.NewInt(1),
		FromIndex: big.NewInt(2),
		ToIndex:   big.NewInt(3),
		Amount:    big.NewInt(4),
		Fee:       big.NewInt(5),
		Nonce:     big.NewInt(6),
	})
	s.NoError(err)
	s.Equal(expected, encodedTransfer)
}

func (s *TransferTestSuite) TestEncodeTransferForSigning() {
	transfer := &models.Transfer{
		TransactionBase: models.TransactionBase{
			FromStateID: 2,
			Amount:      models.MakeUint256(4),
			Fee:         models.MakeUint256(5),
			Nonce:       models.MakeUint256(6),
		},
		ToStateID: 3,
	}
	encodedTransfer, err := EncodeTransfer(transfer)
	s.NoError(err)
	expected, err := s.transfer.SignBytes(nil, encodedTransfer)
	s.NoError(err)

	actual, err := EncodeTransferForSigning(transfer)
	s.NoError(err)
	s.Equal(expected, actual)
}

func newTxTransfer(transfer *models.Transfer) testtx.TxTransfer {
	return testtx.TxTransfer{
		FromIndex: big.NewInt(int64(transfer.FromStateID)),
		ToIndex:   big.NewInt(int64(transfer.ToStateID)),
		Amount:    transfer.Amount.ToBig(),
		Fee:       transfer.Fee.ToBig(),
	}
}

func (s *TransferTestSuite) TestEncodeTransferForCommitment() {
	transfer := &models.Transfer{
		TransactionBase: models.TransactionBase{
			FromStateID: 1,
			Amount:      models.MakeUint256(50),
			Fee:         models.MakeUint256(10),
		},
		ToStateID: 2,
	}

	expected, err := s.testTx.TransferSerialize(nil, []testtx.TxTransfer{newTxTransfer(transfer)})
	s.NoError(err)

	encoded, err := EncodeTransferForCommitment(transfer)
	s.NoError(err)

	s.Equal(expected, encoded)
}

func (s *TransferTestSuite) TestDecodeTransferForCommitment() {
	transfer := &models.Transfer{
		TransactionBase: models.TransactionBase{
			TxType:      txtype.Transfer,
			FromStateID: 1,
			Amount:      models.MakeUint256(50),
			Fee:         models.MakeUint256(10),
		},
		ToStateID: 2,
	}
	transferHash, err := HashTransfer(transfer)
	s.NoError(err)
	transfer.Hash = *transferHash

	encoded, err := EncodeTransferForCommitment(transfer)
	s.NoError(err)

	decoded, err := DecodeTransferFromCommitment(encoded)
	s.NoError(err)

	s.Equal(transfer, decoded)
}

func (s *TransferTestSuite) TestSerializeTransfers() {
	transfer := models.Transfer{
		TransactionBase: models.TransactionBase{
			FromStateID: 1,
			Amount:      models.MakeUint256(50),
			Fee:         models.MakeUint256(10),
		},
		ToStateID: 2,
	}
	transfer2 := models.Transfer{
		TransactionBase: models.TransactionBase{
			FromStateID: 2,
			Amount:      models.MakeUint256(200),
			Fee:         models.MakeUint256(10),
		},
		ToStateID: 3,
	}

	expected, err := s.testTx.TransferSerialize(nil, []testtx.TxTransfer{newTxTransfer(&transfer), newTxTransfer(&transfer2)})
	s.NoError(err)

	serialized, err := SerializeTransfers([]models.Transfer{transfer, transfer2})
	s.NoError(err)

	s.Equal(expected, serialized)
}

func (s *TransferTestSuite) TestDeserializeTransfers() {
	transfer := models.Transfer{
		TransactionBase: models.TransactionBase{
			TxType:      txtype.Transfer,
			FromStateID: 1,
			Amount:      models.MakeUint256(50),
			Fee:         models.MakeUint256(10),
		},
		ToStateID: 2,
	}
	transfer2 := models.Transfer{
		TransactionBase: models.TransactionBase{
			TxType:      txtype.Transfer,
			FromStateID: 2,
			Amount:      models.MakeUint256(200),
			Fee:         models.MakeUint256(10),
		},
		ToStateID: 3,
	}

	transferHash, err := HashTransfer(&transfer)
	s.NoError(err)
	transfer.Hash = *transferHash
	transferHash, err = HashTransfer(&transfer2)
	s.NoError(err)
	transfer2.Hash = *transferHash

	serialized, err := s.testTx.TransferSerialize(
		nil,
		[]testtx.TxTransfer{
			newTxTransfer(&transfer),
			newTxTransfer(&transfer2),
		},
	)
	s.NoError(err)

	transfers, err := DeserializeTransfers(serialized)
	s.NoError(err)
	s.Contains(transfers, transfer)
	s.Contains(transfers, transfer2)
}

func TestTransferTestSuite(t *testing.T) {
	suite.Run(t, new(TransferTestSuite))
}
