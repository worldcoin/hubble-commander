package encoder

import (
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
	txMassMigration := &models.MassMigration{
		TransactionBase: models.TransactionBase{
			FromStateID: 2,
			Amount:      models.MakeUint256(4),
			Fee:         models.MakeUint256(5),
			Nonce:       models.MakeUint256(6),
		},
		SpokeID: models.MakeUint256(3),
	}
	encodedMassMigration, err := EncodeMassMigration(txMassMigration)
	s.NoError(err)

	decodedMassMigration, err := s.massMigration.Decode(&bind.CallOpts{}, encodedMassMigration)
	s.NoError(err)
	s.Equal(int64(txtype.MassMigration), decodedMassMigration.TxType.Int64())
	s.Equal(txMassMigration.FromStateID, uint32(decodedMassMigration.FromIndex.Uint64()))
	s.Equal(txMassMigration.Amount.ToBig(), decodedMassMigration.Amount)
	s.Equal(txMassMigration.Fee.ToBig(), decodedMassMigration.Fee)
	s.Equal(txMassMigration.Nonce.ToBig(), decodedMassMigration.Nonce)
	s.Equal(txMassMigration.SpokeID.ToBig(), decodedMassMigration.SpokeID)
}

func (s *MassMigrationTestSuite) TestEncodeMassMigrationForSigning() {
	txMassMigration := &models.MassMigration{
		TransactionBase: models.TransactionBase{
			FromStateID: 2,
			Amount:      models.MakeUint256(4),
			Fee:         models.MakeUint256(5),
			Nonce:       models.MakeUint256(6),
		},
		SpokeID: models.MakeUint256(3),
	}
	encodedMassMigration, err := EncodeMassMigration(txMassMigration)
	s.NoError(err)
	expected, err := s.massMigration.SignBytes(nil, encodedMassMigration)
	s.NoError(err)

	actual := EncodeMassMigrationForSigning(txMassMigration)
	s.Equal(expected, actual)
}

func TestMassMigrationTestSuite(t *testing.T) {
	suite.Run(t, new(MassMigrationTestSuite))
}
