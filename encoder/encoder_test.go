package encoder

import (
	"math/big"
	"testing"

	"github.com/Worldcoin/hubble-commander/contracts/frontend/generic"
	"github.com/Worldcoin/hubble-commander/contracts/frontend/transfer"
	testtx "github.com/Worldcoin/hubble-commander/contracts/test/tx"
	"github.com/Worldcoin/hubble-commander/contracts/test/types"
	"github.com/Worldcoin/hubble-commander/eth/deployer"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/Worldcoin/hubble-commander/testutils/simulator"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type EncoderTestSuite struct {
	*require.Assertions
	suite.Suite
	sim       *simulator.Simulator
	transfer  *transfer.FrontendTransfer
	generic   *generic.FrontendGeneric
	testTx    *testtx.TestTx
	testTypes *types.TestTypes
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
	s.generic = frontend.FrontendGeneric
	s.testTx = test.TestTx
	s.testTypes = test.TestTypes
}

func (s *EncoderTestSuite) TearDownTest() {
	s.sim.Close()
}

func (s *EncoderTestSuite) TestEncodeTransfer() {
	encodedTransfer, err := EncodeTransfer(&models.Transfer{
		FromStateID: 2,
		ToStateID:   3,
		Amount:      models.MakeUint256(4),
		Fee:         models.MakeUint256(5),
		Nonce:       models.MakeUint256(6),
	})
	s.NoError(err)
	expected, err := s.transfer.Encode(nil, transfer.OffchainTransfer{
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

func newTxTransfer(tx *models.Transaction) testtx.TxTransfer {
	return testtx.TxTransfer{
		FromIndex: big.NewInt(int64(tx.FromIndex)),
		ToIndex:   big.NewInt(int64(tx.ToIndex)),
		Amount:    &tx.Amount.Int,
		Fee:       &tx.Fee.Int,
	}
}

func (s *EncoderTestSuite) TestEncodeTransaction() {
	tx := &models.Transaction{
		FromIndex: 1,
		ToIndex:   2,
		Amount:    models.MakeUint256(50),
		Fee:       models.MakeUint256(10),
	}

	expected, err := s.testTx.TransferSerialize(nil, []testtx.TxTransfer{newTxTransfer(tx)})
	s.NoError(err)

	encoded, err := EncodeTransaction(tx)
	s.NoError(err)

	s.Equal(expected, encoded)
}

func (s *EncoderTestSuite) TestSerializeTransactions() {
	tx := models.Transaction{
		FromIndex: 1,
		ToIndex:   2,
		Amount:    models.MakeUint256(50),
		Fee:       models.MakeUint256(10),
	}
	tx2 := models.Transaction{
		FromIndex: 2,
		ToIndex:   3,
		Amount:    models.MakeUint256(200),
		Fee:       models.MakeUint256(10),
	}

	expected, err := s.testTx.TransferSerialize(nil, []testtx.TxTransfer{newTxTransfer(&tx), newTxTransfer(&tx2)})
	s.NoError(err)

	serialized, err := SerializeTransactions([]models.Transaction{tx, tx2})
	s.NoError(err)

	s.Equal(expected, serialized)
}

func (s *EncoderTestSuite) TestCommitmentBodyHash() {
	accountRoot := utils.RandomHash()
	signature := models.MakeSignature(1, 2)
	feeReceiver := models.MakeUint256(1234)
	txs := utils.RandomBytes(32)

	expectedHash, err := s.testTypes.HashTransferBody(nil, types.TypesTransferBody{
		AccountRoot: accountRoot,
		Signature:   [2]*big.Int{&signature[0].Int, &signature[1].Int},
		FeeReceiver: &feeReceiver.Int,
		Txs:         txs,
	})
	s.NoError(err)

	commitment := models.Commitment{
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
