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
	"github.com/Worldcoin/hubble-commander/testutils/simulator"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/ethereum/go-ethereum/common"
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
	test, err := deployer.DeployTest(sim)
	s.NoError(err)

	s.transfer = frontend.FrontendTransfer
	s.generic = frontend.FrontendGeneric
	s.testTx = test.TestTx
	s.testTypes = test.TestTypes
}

func (s *EncoderTestSuite) TearDownTest() {
	s.sim.Close()
}

func (s *EncoderTestSuite) TestEncodeTransferZero() {
	tx := transfer.OffchainTransfer{
		TxType:    big.NewInt(0),
		FromIndex: big.NewInt(0),
		ToIndex:   big.NewInt(0),
		Amount:    big.NewInt(0),
		Fee:       big.NewInt(0),
		Nonce:     big.NewInt(0),
	}
	bytes, err := EncodeTransfer(tx)
	s.NoError(err)
	expected, err := s.transfer.Encode(nil, tx)
	s.NoError(err)
	s.Equal(expected, bytes)
}

func (s *EncoderTestSuite) TestEncodeTransferNonZero() {
	tx := transfer.OffchainTransfer{
		TxType:    big.NewInt(1),
		FromIndex: big.NewInt(2),
		ToIndex:   big.NewInt(3),
		Amount:    big.NewInt(4),
		Fee:       big.NewInt(5),
		Nonce:     big.NewInt(6),
	}
	bytes, err := EncodeTransfer(tx)
	s.NoError(err)
	expected, err := s.transfer.Encode(nil, tx)
	s.NoError(err)
	s.Equal(expected, bytes)
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

func (s *EncoderTestSuite) TestDecimalEncoding() {
	num := models.MakeUint256(123400000)
	encoded, err := EncodeDecimal(num)
	s.NoError(err)

	expected, err := s.testTx.TestEncodeDecimal(nil, &num.Int)
	s.NoError(err)

	s.Equal(uint16(expected.Uint64()), encoded)
}

func (s *EncoderTestSuite) TestTransactionEncoding() {
	tx := models.Transaction{
		Hash:      common.Hash{},
		FromIndex: models.MakeUint256(1),
		ToIndex:   models.MakeUint256(2),
		Amount:    models.MakeUint256(50),
		Fee:       models.MakeUint256(10),
		Nonce:     models.MakeUint256(22),
		Signature: nil,
	}

	encoded, err := EncodeTransaction(&tx)
	s.NoError(err)

	txTransfer := testtx.TxTransfer{
		FromIndex: &tx.FromIndex.Int,
		ToIndex:   &tx.ToIndex.Int,
		Amount:    &tx.Amount.Int,
		Fee:       &tx.Fee.Int,
	}
	expected, err := s.testTx.TransferSerialize(nil, []testtx.TxTransfer{txTransfer})
	s.NoError(err)

	s.Equal(expected, encoded)
}

func (s *EncoderTestSuite) TestGetCommitmentBodyHash() {
	accountRoot := utils.RandomHash()
	signature := models.Signature{models.MakeUint256(1), models.MakeUint256(2)}
	feeReceiver := models.MakeUint256(1234)
	txs := utils.RandomBytes(32)

	expectedHash, err := s.testTypes.HashTransferBody(nil, types.TypesTransferBody{
		AccountRoot: accountRoot,
		Signature:   [2]*big.Int{&signature[0].Int, &signature[1].Int},
		FeeReceiver: &feeReceiver.Int,
		Txs:         txs,
	})
	s.NoError(err)

	bodyHash, err := GetCommitmentBodyHash(
		accountRoot,
		signature,
		uint32(feeReceiver.Uint64()),
		txs,
	)
	s.NoError(err)

	s.Equal(expectedHash[:], bodyHash.Bytes())
}

func TestEncoderTestSuite(t *testing.T) {
	suite.Run(t, new(EncoderTestSuite))
}
