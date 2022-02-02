package commander

import (
	"context"
	"testing"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type MigrateTestSuite struct {
	*require.Assertions
	suite.Suite
	cmd          *Commander
	client       *eth.TestClient
	storage      *st.TestStorage
	cfg          *config.Config
	wallets      []bls.Wallet
	pendingBatch dto.PendingBatch
}

func (s *MigrateTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
	s.cfg = config.GetTestConfig()
	s.cfg.Rollup.MinTxsPerCommitment = 1
}

func (s *MigrateTestSuite) SetupTest() {
	var err error
	s.storage, err = st.NewTestStorage()
	s.NoError(err)

	s.client = newClientWithGenesisState(s.T(), s.storage)

	s.cmd = NewCommander(s.cfg, nil)
	s.cmd.client = s.client.Client
	s.cmd.storage = s.storage.Storage
	s.cmd.workersContext, s.cmd.stopWorkersContext = context.WithCancel(context.Background())

	err = s.cmd.addGenesisBatch()
	s.NoError(err)

	domain, err := s.client.GetDomain()
	s.NoError(err)
	s.wallets = testutils.GenerateWallets(s.Assertions, domain, 2)
	setAccountLeaves(s.T(), s.storage.Storage, s.wallets)

	s.pendingBatch = dto.PendingBatch{
		ID:              models.MakeUint256(1),
		Type:            batchtype.Transfer,
		TransactionHash: utils.RandomHash(),
		Commitments: []dto.PendingCommitment{
			{
				Commitment: &models.TxCommitment{
					CommitmentBase: models.CommitmentBase{
						ID: models.CommitmentID{
							BatchID:      models.MakeUint256(1),
							IndexInBatch: 0,
						},
						Type:          batchtype.Transfer,
						PostStateRoot: utils.RandomHash(),
					},
					FeeReceiver:       0,
					CombinedSignature: models.MakeRandomSignature(),
					BodyHash:          nil,
				},
				Transactions: models.TransferArray{},
			},
		},
	}
}

func (s *MigrateTestSuite) TearDownTest() {
	stopCommander(s.cmd)
	s.client.Close()
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *MigrateTestSuite) TestSyncPendingBatch_UpdatesUserBalances() {
	tx := testutils.MakeTransfer(0, 1, 0, 100)
	tx.CommitmentID = &s.pendingBatch.Commitments[0].GetCommitmentBase().ID
	s.pendingBatch.Commitments[0].Transactions = models.TransferArray{tx}

	prevSenderLeaf, err := s.storage.StateTree.Leaf(tx.FromStateID)
	s.NoError(err)
	prevReceiverLeaf, err := s.storage.StateTree.Leaf(tx.ToStateID)
	s.NoError(err)

	err = s.cmd.syncPendingBatch(&s.pendingBatch)
	s.NoError(err)

	senderLeaf, err := s.storage.StateTree.Leaf(tx.FromStateID)
	s.NoError(err)
	s.Equal(*prevSenderLeaf.Balance.Sub(&tx.Amount), senderLeaf.Balance)

	receiverLeaf, err := s.storage.StateTree.Leaf(tx.ToStateID)
	s.NoError(err)
	s.Equal(*prevReceiverLeaf.Balance.Add(&tx.Amount), receiverLeaf.Balance)
}

func (s *MigrateTestSuite) TestSyncPendingBatch_AddsPendingBatch() {
	tx := testutils.MakeTransfer(0, 1, 0, 100)
	tx.CommitmentID = &s.pendingBatch.Commitments[0].GetCommitmentBase().ID
	s.pendingBatch.Commitments[0].Transactions = models.TransferArray{tx}

	err := s.cmd.syncPendingBatch(&s.pendingBatch)
	s.NoError(err)

	expectedBatch := models.Batch{
		ID:              s.pendingBatch.ID,
		Type:            s.pendingBatch.Type,
		TransactionHash: s.pendingBatch.TransactionHash,
	}

	batch, err := s.storage.GetBatch(s.pendingBatch.ID)
	s.NoError(err)
	s.Equal(expectedBatch, *batch)

	commitments, err := s.storage.GetCommitmentsByBatchID(s.pendingBatch.ID)
	s.NoError(err)
	s.Len(commitments, 1)
	s.Equal(s.pendingBatch.Commitments[0].Commitment, commitments[0])

	txs, err := s.storage.GetTransfersByCommitmentID(commitments[0].GetCommitmentBase().ID)
	s.NoError(err)
	s.Len(txs, 1)
	s.Equal(s.pendingBatch.Commitments[0].Transactions, txs)
}

func TestMigrateTestSuite(t *testing.T) {
	suite.Run(t, new(MigrateTestSuite))
}
