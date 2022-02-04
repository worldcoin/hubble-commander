package commander

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const (
	getPendingBatchesMethod     = "GetPendingBatches"
	getFailedTransactionsMethod = "GetFailedTransactions"
)

type MockHubble struct {
	mock.Mock
}

func (m *MockHubble) GetPendingBatches() ([]dto.PendingBatch, error) {
	args := m.Called()
	return args.Get(0).([]dto.PendingBatch), args.Error(1)
}

func (m *MockHubble) GetFailedTransactions() (models.GenericTransactionArray, error) {
	args := m.Called()
	return args.Get(0).(models.GenericTransactionArray), args.Error(1)
}

type MigrateTestSuite struct {
	*require.Assertions
	suite.Suite
	storage        *st.TestStorage
	cmd            *Commander
	cfg            *config.Config
	pendingBatches []dto.PendingBatch
	failedTxs      models.GenericTransactionArray
}

func (s *MigrateTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *MigrateTestSuite) SetupTest() {
	var err error
	s.storage, err = st.NewTestStorage()
	s.NoError(err)
	s.cfg = config.GetTestConfig()

	s.cmd = NewCommander(s.cfg, nil)
	s.cmd.storage = s.storage.Storage

	setStateLeaves(s.T(), s.storage.Storage)

	s.pendingBatches = []dto.PendingBatch{
		makePendingBatch(2, models.TransferArray{testutils.MakeTransfer(0, 1, 1, 100)}),
		makePendingBatch(1, models.TransferArray{testutils.MakeTransfer(0, 1, 0, 100)}),
	}

	s.failedTxs = models.MakeTransferArray(
		makeFailedTransfer(0),
		makeFailedTransfer(1),
	)
}

func (s *MigrateTestSuite) TearDownTest() {
	stopCommander(s.cmd)
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *MigrateTestSuite) TestMigrate_MissingBootstrapNodeURL() {
	err := s.cmd.migrate()
	s.ErrorIs(err, errMissingBootstrapNodeURL)
}

func (s *MigrateTestSuite) TestMigrateCommanderData_SetsMigrateToFalse() {
	hubble := new(MockHubble)
	hubble.On(getPendingBatchesMethod).
		Return([]dto.PendingBatch{}, nil)
	hubble.On(getFailedTransactionsMethod).
		Return(models.TransferArray{}, nil)

	err := s.cmd.migrateCommanderData(hubble)
	s.NoError(err)

	s.False(s.cmd.isMigrating())
}

func (s *MigrateTestSuite) TestMigrateCommanderData_SyncsFailedTxs() {
	hubble := new(MockHubble)
	hubble.On(getPendingBatchesMethod).Return(s.pendingBatches, nil)
	hubble.On(getFailedTransactionsMethod).Return(s.failedTxs, nil)

	err := s.cmd.migrateCommanderData(hubble)
	s.NoError(err)

	for i := 0; i < s.failedTxs.Len(); i++ {
		tx, err := s.cmd.storage.GetTransfer(s.failedTxs.At(i).GetBase().Hash)
		s.NoError(err)
		s.Equal(*s.failedTxs.At(i).ToTransfer(), *tx)
	}
}

func (s *MigrateTestSuite) TestMigrateCommanderData_SyncsBatches() {
	batches := []dto.PendingBatch{
		makePendingBatch(1, models.TransferArray{testutils.MakeTransfer(0, 1, 0, 100)}),
		makePendingBatch(2, models.TransferArray{testutils.MakeTransfer(0, 1, 1, 100)}),
	}

	hubble := new(MockHubble)
	hubble.On(getPendingBatchesMethod).Return(batches, nil)
	hubble.On(getFailedTransactionsMethod).Return(models.TransferArray{}, nil)

	err := s.cmd.migrateCommanderData(hubble)
	s.NoError(err)

	leaf, err := s.storage.StateTree.Leaf(1)
	s.NoError(err)
	s.EqualValues(200, leaf.Balance.Uint64())
}

func makePendingBatch(batchID uint64, txs models.GenericTransactionArray) dto.PendingBatch {
	return dto.PendingBatch{
		ID:              models.MakeUint256(batchID),
		Type:            batchtype.Transfer,
		TransactionHash: utils.RandomHash(),
		Commitments: []dto.PendingCommitment{
			{
				Commitment: &models.TxCommitment{
					CommitmentBase: models.CommitmentBase{
						ID: models.CommitmentID{
							BatchID:      models.MakeUint256(batchID),
							IndexInBatch: 0,
						},
						Type:          batchtype.Transfer,
						PostStateRoot: utils.RandomHash(),
					},
					FeeReceiver:       0,
					CombinedSignature: models.MakeRandomSignature(),
				},
				Transactions: txs,
			},
		},
	}
}

func makeFailedTransfer(nonce uint64) models.Transfer {
	transfer := testutils.MakeTransfer(0, 1, nonce, 100)
	transfer.ErrorMessage = ref.String("failed quack")
	return transfer
}

func TestMigrateTestSuite(t *testing.T) {
	suite.Run(t, new(MigrateTestSuite))
}
