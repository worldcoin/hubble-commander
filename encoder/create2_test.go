package encoder

import (
	"math/big"
	"testing"

	contractCreate2Transfer "github.com/Worldcoin/hubble-commander/contracts/frontend/create2transfer"
	testtx "github.com/Worldcoin/hubble-commander/contracts/test/tx"
	"github.com/Worldcoin/hubble-commander/eth/deployer"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/Worldcoin/hubble-commander/testutils/simulator"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type Create2TestSuite struct {
	*require.Assertions
	suite.Suite
	sim             *simulator.Simulator
	create2Transfer *contractCreate2Transfer.FrontendCreate2Transfer
	testTx          *testtx.TestTx
}

func (s *Create2TestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *Create2TestSuite) SetupTest() {
	sim, err := simulator.NewSimulator()
	s.NoError(err)
	s.sim = sim

	frontend, err := deployer.DeployFrontend(sim)
	s.NoError(err)
	test, err := testutils.DeployTest(sim)
	s.NoError(err)

	s.create2Transfer = frontend.FrontendCreate2Transfer
	s.testTx = test.TestTx
}

func (s *Create2TestSuite) TearDownTest() {
	s.sim.Close()
}

func (s *Create2TestSuite) TestEncodeCreate2TransferWithStateID() {
	encodedCreate2Transfer, err := EncodeCreate2TransferWithStateID(&models.Create2Transfer{
		TransactionBase: models.TransactionBase{
			FromStateID: 4,
			Amount:      models.MakeUint256(7),
			Fee:         models.MakeUint256(8),
			Nonce:       models.MakeUint256(9),
		},
		ToStateID:   ref.Uint32(5),
		ToPublicKey: models.PublicKey{1, 2, 3},
	}, 6)
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

func (s *Create2TestSuite) TestEncodeCreate2Transfer() {
	publicKey := models.PublicKey{1, 2, 3, 4, 5, 6}
	encodedCreate2Transfer, err := EncodeCreate2Transfer(&models.Create2Transfer{
		TransactionBase: models.TransactionBase{
			FromStateID: 4,
			Amount:      models.MakeUint256(5),
			Fee:         models.MakeUint256(6),
			Nonce:       models.MakeUint256(7),
		},
		ToPublicKey: publicKey,
	})
	s.NoError(err)
	expected, err := s.create2Transfer.EncodeWithPub(nil, contractCreate2Transfer.OffchainCreate2TransferWithPub{
		TxType:    big.NewInt(3),
		FromIndex: big.NewInt(4),
		ToPubkey:  crypto.Keccak256Hash(publicKey.Bytes()),
		Amount:    big.NewInt(5),
		Fee:       big.NewInt(6),
		Nonce:     big.NewInt(7),
	})
	s.NoError(err)
	s.Equal(expected, encodedCreate2Transfer)
}

func (s *Create2TestSuite) TestEncodeCreate2TransferForSigning() {
	tx := &models.Create2Transfer{
		TransactionBase: models.TransactionBase{
			FromStateID: 4,
			Amount:      models.MakeUint256(5),
			Fee:         models.MakeUint256(6),
			Nonce:       models.MakeUint256(7),
		},
		ToPublicKey: models.PublicKey{1, 2, 3, 4, 5, 6},
	}
	encodedCreate2Transfer, err := EncodeCreate2Transfer(tx)
	s.NoError(err)
	expected, err := s.create2Transfer.SignBytes(nil, encodedCreate2Transfer)
	s.NoError(err)

	actual, err := EncodeCreate2TransferForSigning(tx)
	s.NoError(err)
	s.Equal(expected, actual)
}

func newTxCreate2Transfer(transfer *models.Create2Transfer, toPubKeyID uint32) testtx.TxCreate2Transfer {
	return testtx.TxCreate2Transfer{
		FromIndex:  big.NewInt(int64(transfer.FromStateID)),
		ToIndex:    big.NewInt(int64(*transfer.ToStateID)),
		ToPubkeyID: big.NewInt(int64(toPubKeyID)),
		Amount:     transfer.Amount.ToBig(),
		Fee:        transfer.Fee.ToBig(),
	}
}

func (s *Create2TestSuite) TestEncodeCreate2TransferForCommitment() {
	transfer := &models.Create2Transfer{
		TransactionBase: models.TransactionBase{
			FromStateID: 1,
			Amount:      models.MakeUint256(50),
			Fee:         models.MakeUint256(10),
		},
		ToStateID: ref.Uint32(2),
	}

	expected, err := s.testTx.Create2transferSerialize(nil, []testtx.TxCreate2Transfer{newTxCreate2Transfer(transfer, 6)})
	s.NoError(err)

	encoded, err := EncodeCreate2TransferForCommitment(transfer, 6)
	s.NoError(err)

	s.Equal(expected, encoded)
}

func (s *Create2TestSuite) TestDecodeCreate2TransferForCommitment() {
	transfer := &models.Create2Transfer{
		TransactionBase: models.TransactionBase{
			TxType:      txtype.Create2Transfer,
			FromStateID: 1,
			Amount:      models.MakeUint256(50),
			Fee:         models.MakeUint256(10),
		},
		ToStateID: ref.Uint32(2),
	}

	encoded, err := EncodeCreate2TransferForCommitment(transfer, 6)
	s.NoError(err)

	decoded, pubKeyID, err := DecodeCreate2TransferFromCommitment(encoded)
	s.NoError(err)
	s.EqualValues(6, pubKeyID)

	s.Equal(transfer, decoded)
}

func (s *Create2TestSuite) TestSerializeCreate2Transfers() {
	transfer := models.Create2Transfer{
		TransactionBase: models.TransactionBase{
			FromStateID: 1,
			Amount:      models.MakeUint256(50),
			Fee:         models.MakeUint256(10),
		},
		ToStateID:   ref.Uint32(2),
		ToPublicKey: models.PublicKey{1, 2, 3},
	}
	transfer2 := models.Create2Transfer{
		TransactionBase: models.TransactionBase{
			FromStateID: 2,
			Amount:      models.MakeUint256(200),
			Fee:         models.MakeUint256(10),
		},
		ToStateID:   ref.Uint32(3),
		ToPublicKey: models.PublicKey{2, 3, 4},
	}

	expected, err := s.testTx.Create2transferSerialize(
		nil,
		[]testtx.TxCreate2Transfer{
			newTxCreate2Transfer(&transfer, 6),
			newTxCreate2Transfer(&transfer2, 5),
		},
	)
	s.NoError(err)

	serialized, err := SerializeCreate2Transfers([]models.Create2Transfer{transfer, transfer2}, []uint32{6, 5})
	s.NoError(err)

	s.Equal(expected, serialized)
}

func (s *Create2TestSuite) TestSerializeCreate2Transfers_InvalidLength() {
	transfer := models.Create2Transfer{
		TransactionBase: models.TransactionBase{
			FromStateID: 1,
			Amount:      models.MakeUint256(50),
			Fee:         models.MakeUint256(10),
		},
		ToStateID:   ref.Uint32(2),
		ToPublicKey: models.PublicKey{1, 2, 3},
	}

	serialized, err := SerializeCreate2Transfers([]models.Create2Transfer{transfer}, []uint32{})
	s.Equal(ErrInvalidSlicesLength, err)
	s.Nil(serialized)
}

func (s *Create2TestSuite) TestDeserializeCreate2Transfers() {
	txs := []models.Create2Transfer{
		{
			TransactionBase: models.TransactionBase{
				TxType:      txtype.Create2Transfer,
				FromStateID: 1,
				Amount:      models.MakeUint256(50),
				Fee:         models.MakeUint256(10),
			},
			ToStateID: ref.Uint32(2),
		},
		{
			TransactionBase: models.TransactionBase{
				TxType:      txtype.Create2Transfer,
				FromStateID: 2,
				Amount:      models.MakeUint256(200),
				Fee:         models.MakeUint256(10),
			},
			ToStateID: ref.Uint32(3),
		},
	}

	serialized, err := s.testTx.Create2transferSerialize(
		nil,
		[]testtx.TxCreate2Transfer{
			newTxCreate2Transfer(&txs[0], 6),
			newTxCreate2Transfer(&txs[1], 5),
		},
	)
	s.NoError(err)

	transfers, toPubKeyIDs, err := DeserializeCreate2Transfers(serialized)
	s.NoError(err)
	s.Equal(txs, transfers)
	s.Equal([]uint32{6, 5}, toPubKeyIDs)
}

func (s *Create2TestSuite) TestDeserializeCreate2TransferPubKeyIDs() {
	txs := []models.Create2Transfer{
		{
			TransactionBase: models.TransactionBase{
				TxType:      txtype.Create2Transfer,
				FromStateID: 1,
				Amount:      models.MakeUint256(50),
				Fee:         models.MakeUint256(10),
			},
			ToStateID: ref.Uint32(2),
		},
		{
			TransactionBase: models.TransactionBase{
				TxType:      txtype.Create2Transfer,
				FromStateID: 2,
				Amount:      models.MakeUint256(200),
				Fee:         models.MakeUint256(10),
			},
			ToStateID: ref.Uint32(3),
		},
	}

	for i := range txs {
		transferHash, err := HashCreate2Transfer(&txs[i])
		s.NoError(err)
		txs[i].Hash = *transferHash
	}

	serialized, err := s.testTx.Create2transferSerialize(
		nil,
		[]testtx.TxCreate2Transfer{
			newTxCreate2Transfer(&txs[0], 6),
			newTxCreate2Transfer(&txs[1], 5),
		},
	)
	s.NoError(err)

	toPubKeyIDs := DeserializeCreate2TransferPubKeyIDs(serialized)
	s.NoError(err)
	s.Equal([]uint32{6, 5}, toPubKeyIDs)
}

func TestCreate2TestSuite(t *testing.T) {
	suite.Run(t, new(Create2TestSuite))
}
