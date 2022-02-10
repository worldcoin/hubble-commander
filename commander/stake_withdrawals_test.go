package commander

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/commander/tracker"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/eth/deployer/rollup"
	"github.com/Worldcoin/hubble-commander/metrics"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type StakeWithdrawalsTestSuite struct {
	*require.Assertions
	suite.Suite
	tracker.TestSuiteWithTxsSending
	teardown   func() error
	testClient *eth.TestClient
	cmd        *Commander
}

func (s *StakeWithdrawalsTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *StakeWithdrawalsTestSuite) SetupTest() {
	testStorage, err := st.NewTestStorage()
	s.NoError(err)
	s.teardown = testStorage.Teardown

	// finalize instantly
	s.testClient, err = eth.NewConfiguredTestClient(&rollup.DeploymentConfig{
		Params: rollup.Params{BlocksToFinalise: models.NewUint256(0)},
	}, &eth.ClientConfig{})

	s.NoError(err)
	s.cmd = &Commander{
		storage:             testStorage.Storage,
		client:              s.testClient.Client,
		metrics:             metrics.NewCommanderMetrics(),
		txsTrackingChannels: s.testClient.TxsChannels,
	}

	s.StartTxsSending(s.testClient.TxsChannels.Requests)
}

func (s *StakeWithdrawalsTestSuite) TearDownTest() {
	s.StopTxsSending()
	s.testClient.Close()
	err := s.teardown()
	s.NoError(err)
}

func (s *StakeWithdrawalsTestSuite) TestSyncStakeWithdrawals() {
	batchID := models.MakeUint256(1)

	err := s.cmd.storage.AddPendingStakeWithdrawal(&models.PendingStakeWithdrawal{
		BatchID:           batchID,
		FinalisationBlock: 0,
	})
	s.NoError(err)

	commitments := s.createSingleCommitmentInSlice(batchtype.Transfer)
	_, err = s.testClient.SubmitTransfersBatchAndWait(&batchID, commitments)
	s.NoError(err)
	err = s.testClient.WithdrawStakeAndWait(&batchID)
	s.NoError(err)
	latestBlockNumber, err := s.testClient.GetLatestBlockNumber()
	s.NoError(err)
	err = s.cmd.syncStakeWithdrawals(0, *latestBlockNumber)
	s.NoError(err)

	err = s.cmd.storage.RemovePendingStakeWithdrawal(batchID)
	s.ErrorIs(err, st.NewNotFoundError("pending stake withdrawal"))
}

func (s *StakeWithdrawalsTestSuite) createSingleCommitmentInSlice(commitmentType batchtype.BatchType) []models.CommitmentWithTxs {
	stateRoot, err := s.cmd.storage.StateTree.Root()
	s.NoError(err)

	commitment := models.TxCommitmentWithTxs{
		TxCommitment: models.TxCommitment{
			CommitmentBase: models.CommitmentBase{
				Type:          commitmentType,
				PostStateRoot: *stateRoot,
			},
			FeeReceiver:       0,
			CombinedSignature: models.Signature{},
		},
		Transactions: []byte{}}

	return []models.CommitmentWithTxs{&commitment}
}

func TestStakeWithdrawalsTestSuite(t *testing.T) {
	suite.Run(t, new(StakeWithdrawalsTestSuite))
}
