package encoder

import (
	"math/big"
	"testing"

	contractMassMigration "github.com/Worldcoin/hubble-commander/contracts/frontend/massmigration"
	testtx "github.com/Worldcoin/hubble-commander/contracts/test/tx"
	"github.com/Worldcoin/hubble-commander/eth/deployer"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/Worldcoin/hubble-commander/testutils/simulator"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type MassMigrationTestSuite struct {
	*require.Assertions
	suite.Suite
	sim           *simulator.Simulator
	massMigration *contractMassMigration.FrontendMassMigration
	testTx        *testtx.TestTx
}

func (s *MassMigrationTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *MassMigrationTestSuite) SetupTest() {
	sim, err := simulator.NewSimulator()
	s.NoError(err)
	s.sim = sim

	frontend, err := deployer.DeployFrontend(sim)
	s.NoError(err)
	test, err := testutils.DeployTest(sim)
	s.NoError(err)

	s.massMigration = frontend.FrontendMassMigration
	s.testTx = test.TestTx
}

func (s *MassMigrationTestSuite) TearDownTest() {
	s.sim.Close()
}

func (s *MassMigrationTestSuite) TestEncodeMassMigration() {
	massMigration := &models.MassMigration{
		TransactionBase: models.TransactionBase{
			FromStateID: 2,
			Amount:      models.MakeUint256(4),
			Fee:         models.MakeUint256(5),
			Nonce:       models.MakeUint256(6),
		},
		SpokeID: 3,
	}
	encodedMassMigration, err := EncodeMassMigration(massMigration)
	s.NoError(err)

	decodedMassMigration, err := s.massMigration.Decode(&bind.CallOpts{}, encodedMassMigration)
	s.NoError(err)
	s.EqualValues(txtype.MassMigration, decodedMassMigration.TxType.Int64())
	s.EqualValues(massMigration.FromStateID, decodedMassMigration.FromIndex.Uint64())
	s.Equal(massMigration.Amount.ToBig(), decodedMassMigration.Amount)
	s.Equal(massMigration.Fee.ToBig(), decodedMassMigration.Fee)
	s.Equal(massMigration.Nonce.ToBig(), decodedMassMigration.Nonce)
	s.EqualValues(massMigration.SpokeID, decodedMassMigration.SpokeID.Uint64())
}

func (s *MassMigrationTestSuite) TestEncodeMassMigrationForSigning() {
	massMigration := &models.MassMigration{
		TransactionBase: models.TransactionBase{
			FromStateID: 2,
			Amount:      models.MakeUint256(4),
			Fee:         models.MakeUint256(5),
			Nonce:       models.MakeUint256(6),
		},
		SpokeID: 3,
	}
	encodedMassMigration, err := EncodeMassMigration(massMigration)
	s.NoError(err)
	expected, err := s.massMigration.SignBytes(nil, encodedMassMigration)
	s.NoError(err)

	actual := EncodeMassMigrationForSigning(massMigration)
	s.Equal(expected, actual)
}

func (s *MassMigrationTestSuite) TestEncodeMassMigrationForCommitment() {
	massMigration := &models.MassMigration{
		TransactionBase: models.TransactionBase{
			FromStateID: 1,
			Amount:      models.MakeUint256(50),
			Fee:         models.MakeUint256(10),
		},
		SpokeID: 2,
	}

	encoded, err := EncodeMassMigrationForCommitment(massMigration)
	s.NoError(err)

	massMigrationCount, err := s.testTx.MassMigrationSize(nil, encoded)
	s.NoError(err)
	s.EqualValues(1, massMigrationCount.Uint64())

	decodedMassMigration, err := s.testTx.MassMigrationDecode(nil, encoded, big.NewInt(0))
	s.NoError(err)
	s.EqualValues(massMigration.FromStateID, decodedMassMigration.FromIndex.Uint64())
	s.Equal(massMigration.Amount, models.MakeUint256FromBig(*decodedMassMigration.Amount))
	s.Equal(massMigration.Fee, models.MakeUint256FromBig(*decodedMassMigration.Fee))
}

func (s *MassMigrationTestSuite) TestDecodeMassMigrationFromCommitment() {
	massMigration := &models.MassMigration{
		TransactionBase: models.TransactionBase{
			TxType:      txtype.MassMigration,
			FromStateID: 1,
			Amount:      models.MakeUint256(50),
			Fee:         models.MakeUint256(10),
		},
	}
	massMigrationHash, err := HashMassMigration(massMigration)
	s.NoError(err)
	massMigration.Hash = *massMigrationHash

	encoded, err := EncodeMassMigrationForCommitment(massMigration)
	s.NoError(err)

	decoded, err := DecodeMassMigrationFromCommitment(encoded)
	s.NoError(err)

	s.Equal(massMigration, decoded)
}

func (s *MassMigrationTestSuite) TestSerializeMassMigrations() {
	massMigrations := []models.MassMigration{
		{
			TransactionBase: models.TransactionBase{
				FromStateID: 1,
				Amount:      models.MakeUint256(50),
				Fee:         models.MakeUint256(10),
			},
			SpokeID: 2,
		},
		{
			TransactionBase: models.TransactionBase{
				FromStateID: 2,
				Amount:      models.MakeUint256(200),
				Fee:         models.MakeUint256(10),
			},
			SpokeID: 3,
		},
	}

	serialized, err := SerializeMassMigrations(massMigrations)
	s.NoError(err)

	massMigrationsCount, err := s.testTx.MassMigrationSize(nil, serialized)
	s.NoError(err)
	s.EqualValues(len(massMigrations), massMigrationsCount.Uint64())

	for i := range massMigrations {
		decodedMassMigration, err := s.testTx.MassMigrationDecode(nil, serialized, big.NewInt(int64(i)))
		s.NoError(err)
		s.EqualValues(massMigrations[i].FromStateID, decodedMassMigration.FromIndex.Uint64())
		s.Equal(massMigrations[i].Amount, models.MakeUint256FromBig(*decodedMassMigration.Amount))
		s.Equal(massMigrations[i].Fee, models.MakeUint256FromBig(*decodedMassMigration.Fee))
	}
}

func (s *MassMigrationTestSuite) TestDeserializeMassMigrations() {
	massMigrations := []models.MassMigration{
		{
			TransactionBase: models.TransactionBase{
				TxType:      txtype.MassMigration,
				FromStateID: 1,
				Amount:      models.MakeUint256(50),
				Fee:         models.MakeUint256(10),
			},
		},
		{
			TransactionBase: models.TransactionBase{
				TxType:      txtype.MassMigration,
				FromStateID: 2,
				Amount:      models.MakeUint256(200),
				Fee:         models.MakeUint256(10),
			},
		},
	}

	for i := range massMigrations {
		hash, err := HashMassMigration(&massMigrations[i])
		s.NoError(err)
		massMigrations[i].Hash = *hash
	}

	serialized, err := SerializeMassMigrations(massMigrations)
	s.NoError(err)

	deserializedMassMigrations, err := DeserializeMassMigrations(serialized)
	s.NoError(err)
	s.Len(deserializedMassMigrations, len(massMigrations))
	s.Equal(deserializedMassMigrations, massMigrations)
}

func TestMassMigrationTestSuite(t *testing.T) {
	suite.Run(t, new(MassMigrationTestSuite))
}
