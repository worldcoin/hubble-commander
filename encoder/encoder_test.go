package encoder

import (
	"math/big"
	"testing"

	contractCreate2Transfer "github.com/Worldcoin/hubble-commander/contracts/frontend/create2transfer"
	"github.com/Worldcoin/hubble-commander/contracts/frontend/generic"
	contractTransfer "github.com/Worldcoin/hubble-commander/contracts/frontend/transfer"
	testtx "github.com/Worldcoin/hubble-commander/contracts/test/tx"
	"github.com/Worldcoin/hubble-commander/contracts/test/types"
	"github.com/Worldcoin/hubble-commander/eth/deployer"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/Worldcoin/hubble-commander/testutils/simulator"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type EncoderTestSuite struct {
	*require.Assertions
	suite.Suite
	sim             *simulator.Simulator
	transfer        *contractTransfer.FrontendTransfer
	create2Transfer *contractCreate2Transfer.FrontendCreate2Transfer
	generic         *generic.FrontendGeneric
	testTx          *testtx.TestTx
	testTypes       *types.TestTypes
}

func (s *EncoderTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *EncoderTestSuite) SetupTest() {
	sim, err := simulator.NewSimulator()
	s.NoError(err)
	s.sim = sim

	frontend, err := deployer.DeployFrontend(sim)
	s.NoError(err)
	test, err := testutils.DeployTest(sim)
	s.NoError(err)

	s.transfer = frontend.FrontendTransfer
	s.create2Transfer = frontend.FrontendCreate2Transfer
	s.generic = frontend.FrontendGeneric
	s.testTx = test.TestTx
	s.testTypes = test.TestTypes
}

func (s *EncoderTestSuite) TearDownTest() {
	s.sim.Close()
}

func (s *EncoderTestSuite) TestEncodeTransfer() {
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

func (s *EncoderTestSuite) TestEncodeCreate2Transfer() {
	encodedCreate2Transfer, err := EncodeCreate2Transfer(&models.Create2Transfer{
		TransactionBase: models.TransactionBase{
			FromStateID: 4,
			Amount:      models.MakeUint256(7),
			Fee:         models.MakeUint256(8),
			Nonce:       models.MakeUint256(9),
		},
		ToStateID:  5,
		ToPubKeyID: 6,
	})
	s.NoError(err)
	expected, err := s.create2Transfer.Encode(nil, contractCreate2Transfer.OffchainCreate2Transfer{
		TxType:     big.NewInt(3),
		FromIndex:  big.NewInt(4),
		ToIndex:    big.NewInt(5),
		ToPubkeyID: big.NewInt(6),
		Amount:     big.NewInt(7),
		Fee:        big.NewInt(8),
		Nonce:      big.NewInt(9),
	})
	s.NoError(err)
	s.Equal(expected, encodedCreate2Transfer)
}

func (s *EncoderTestSuite) TestEncodeCreate2TransferWithPubKey() {
	publicKey := models.PublicKey{1, 2, 3, 4, 5, 6}
	encodedCreate2Transfer, err := EncodeCreate2TransferWithPubKey(&models.Create2Transfer{
		TransactionBase: models.TransactionBase{
			FromStateID: 4,
			Amount:      models.MakeUint256(5),
			Fee:         models.MakeUint256(6),
			Nonce:       models.MakeUint256(7),
		},
	}, &publicKey)
	s.NoError(err)
	expected, err := s.create2Transfer.EncodeWithPub(nil, contractCreate2Transfer.OffchainCreate2TransferWithPub{
		TxType:    big.NewInt(3),
		FromIndex: big.NewInt(4),
		ToPubkey:  publicKey.BigInts(),
		Amount:    big.NewInt(5),
		Fee:       big.NewInt(6),
		Nonce:     big.NewInt(7),
	})
	s.NoError(err)
	s.Equal(expected, encodedCreate2Transfer)
}

func (s *EncoderTestSuite) TestEncodeTransferForSigning() {
	tx := &models.Transfer{
		TransactionBase: models.TransactionBase{
			FromStateID: 2,
			Amount:      models.MakeUint256(4),
			Fee:         models.MakeUint256(5),
			Nonce:       models.MakeUint256(6),
		},
		ToStateID: 3,
	}
	encodedTransfer, err := EncodeTransfer(tx)
	s.NoError(err)
	expected, err := s.transfer.SignBytes(nil, encodedTransfer)
	s.NoError(err)

	actual, err := EncodeTransferForSigning(tx)
	s.NoError(err)
	s.Equal(expected, actual)
}

func (s *EncoderTestSuite) TestEncodeCreate2TransferForSigning() {
	tx := &models.Create2Transfer{
		TransactionBase: models.TransactionBase{
			FromStateID: 4,
			Amount:      models.MakeUint256(5),
			Fee:         models.MakeUint256(6),
			Nonce:       models.MakeUint256(7),
		},
	}
	publicKey := models.PublicKey{1, 2, 3, 4, 5, 6}
	encodedCreate2Transfer, err := EncodeCreate2TransferWithPubKey(tx, &publicKey)
	s.NoError(err)
	expected, err := s.create2Transfer.SignBytes(nil, encodedCreate2Transfer)
	s.NoError(err)

	actual, err := EncodeCreate2TransferForSigning(tx, &publicKey)
	s.NoError(err)
	s.Equal(expected, actual)
}

func (s *EncoderTestSuite) TestEncodeUserState() {
	state := generic.TypesUserState{
		PubkeyID: big.NewInt(1),
		TokenID:  big.NewInt(2),
		Balance:  big.NewInt(420),
		Nonce:    big.NewInt(0),
	}
	bytes, err := EncodeUserState(state)
	s.NoError(err)

	expected, err := s.generic.Encode(nil, state)
	s.NoError(err)
	s.Equal(expected, bytes)
}

func (s *EncoderTestSuite) TestEncodeDecimal() {
	num := models.MakeUint256(123400000)
	encoded, err := EncodeDecimal(num)
	s.NoError(err)

	expected, err := s.testTx.TestEncodeDecimal(nil, &num.Int)
	s.NoError(err)

	s.Equal(uint16(expected.Uint64()), encoded)
}

func newTxTransfer(transfer *models.Transfer) testtx.TxTransfer {
	return testtx.TxTransfer{
		FromIndex: big.NewInt(int64(transfer.FromStateID)),
		ToIndex:   big.NewInt(int64(transfer.ToStateID)),
		Amount:    &transfer.Amount.Int,
		Fee:       &transfer.Fee.Int,
	}
}

func (s *EncoderTestSuite) TestEncodeTransferForCommitment() {
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

func (s *EncoderTestSuite) TestSerializeTransfers() {
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

func (s *EncoderTestSuite) TestCommitmentBodyHash() {
	accountRoot := utils.RandomHash()
	signature := models.MakeRandomSignature()
	feeReceiver := models.MakeUint256(1234)
	txs := utils.RandomBytes(32)

	expectedHash, err := s.testTypes.HashTransferBody(nil, types.TypesTransferBody{
		AccountRoot: accountRoot,
		Signature:   signature.BigInts(),
		FeeReceiver: &feeReceiver.Int,
		Txs:         txs,
	})
	s.NoError(err)

	commitment := models.Commitment{
		Type:              txtype.Transfer,
		Transactions:      txs,
		FeeReceiver:       uint32(feeReceiver.Uint64()),
		CombinedSignature: signature,
		AccountTreeRoot:   &accountRoot,
	}

	s.Equal(expectedHash[:], commitment.BodyHash().Bytes())
}

func TestEncoderTestSuite(t *testing.T) {
	suite.Run(t, new(EncoderTestSuite))
}
