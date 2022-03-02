//go:build e2e
// +build e2e

package e2e

import (
	"fmt"
	"testing"
	"time"

	"github.com/Worldcoin/hubble-commander/api"
	admintypes "github.com/Worldcoin/hubble-commander/client"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/e2e/setup"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/models/enums/batchstatus"
	"github.com/Worldcoin/hubble-commander/models/enums/healthstatus"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/Worldcoin/hubble-commander/utils/consts"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"github.com/ybbus/jsonrpc/v2"
)

type MigrationModeE2ETestSuite struct {
	setup.E2ETestSuite

	gethRPCClient *rpc.Client
}

func (s *MigrationModeE2ETestSuite) SetupTest() {
	cfg := config.GetConfig()
	cfg.Rollup.MinTxsPerCommitment = 4
	cfg.Rollup.MaxTxsPerCommitment = 10
	cfg.Rollup.MinCommitmentsPerBatch = 1

	if cfg.Ethereum == nil {
		log.Panicf("migration mode test cannot be run on simulator")
	}

	s.SetupTestEnvironment(cfg, nil)
	s.RPCClient = s.createAdminRPCClient(cfg)

	var err error
	s.gethRPCClient, err = rpc.Dial(cfg.Ethereum.RPCURL)
	s.NoError(err)
}

func (s *MigrationModeE2ETestSuite) TearDownTest() {
	s.E2ETestSuite.TearDownTest()
	s.gethRPCClient.Close()
}

func (s *MigrationModeE2ETestSuite) TestCommanderMigrationMode() {
	s.SubmitTxBatchAndWait(func() common.Hash {
		return s.SendNTransactions(4, dto.Transfer{
			FromStateID: ref.Uint32(1),
			ToStateID:   ref.Uint32(2),
			Amount:      models.NewUint256(90),
			Fee:         models.NewUint256(10),
			Nonce:       models.NewUint256(0),
		})
	})
	s.WaitForBatchStatus(1, batchstatus.Mined)

	s.stopMiningBlocks()
	defer func() {
		s.startMiningBlocks()
	}()

	// Invalid tx
	s.SendTransaction(dto.Transfer{
		FromStateID: ref.Uint32(1),
		ToStateID:   ref.Uint32(999), // Non-existent receiver
		Amount:      models.NewUint256(90),
		Fee:         models.NewUint256(10),
		Nonce:       models.NewUint256(8),
	})
	// Some valid txs
	s.SendNTransactions(4, dto.Transfer{
		FromStateID: ref.Uint32(1),
		ToStateID:   ref.Uint32(2),
		Amount:      models.NewUint256(90),
		Fee:         models.NewUint256(10),
		Nonce:       models.NewUint256(4),
	})
	// Another invalid tx
	s.SendTransaction(dto.Transfer{
		FromStateID: ref.Uint32(2),
		ToStateID:   ref.Uint32(999), // Non-existent receiver
		Amount:      models.NewUint256(90),
		Fee:         models.NewUint256(10),
		Nonce:       models.NewUint256(0),
	})

	s.WaitForBatchStatus(2, batchstatus.Submitted)

	s.configureCommander(dto.ConfigureParams{
		CreateBatches: ref.Bool(false),
	})

	s.SendNTransactions(4, dto.Transfer{
		FromStateID: ref.Uint32(1),
		ToStateID:   ref.Uint32(2),
		Amount:      models.NewUint256(90),
		Fee:         models.NewUint256(10),
		Nonce:       models.NewUint256(8),
	})
	s.SendNTransactions(4, dto.Transfer{
		FromStateID: ref.Uint32(2),
		ToStateID:   ref.Uint32(3),
		Amount:      models.NewUint256(90),
		Fee:         models.NewUint256(10),
		Nonce:       models.NewUint256(0),
	})

	s.configureCommander(dto.ConfigureParams{
		AcceptTransactions: ref.Bool(false),
	})

	s.testValidateCommanderNotAcceptingTxs()

	migrationCommander, migrationAdminRPCClient := s.prepareMigrationCommander()

	err := migrationCommander.Start()
	s.NoError(err)
	defer func() {
		s.NoError(migrationCommander.Stop())
	}()

	// All further calls are going to be made to the new migrated commander
	s.RPCClient = migrationAdminRPCClient

	s.configureCommander(dto.ConfigureParams{
		CreateBatches: ref.Bool(false),
	})

	s.waitForCommanderReadyStatus()
	s.validateSuccessfulMigration()

	s.configureCommander(dto.ConfigureParams{
		CreateBatches: ref.Bool(true),
	})

	s.startMiningBlocks()

	s.testValidateMigratedCommanderFunctionalities()
}

func (s *MigrationModeE2ETestSuite) testValidateCommanderNotAcceptingTxs() {
	transfer, err := api.SignTransfer(&s.Wallets[0], dto.Transfer{
		FromStateID: ref.Uint32(0),
		ToStateID:   ref.Uint32(1),
		Amount:      models.NewUint256(90),
		Fee:         models.NewUint256(10),
		Nonce:       models.NewUint256(0),
	})
	s.NoError(err)

	var txHash common.Hash
	var rpcError *jsonrpc.RPCError
	err = s.RPCClient.CallFor(&txHash, "hubble_sendTransaction", []interface{}{*transfer})
	s.True(errors.As(err, &rpcError))
	s.Equal(10017, rpcError.Code)
}

func (s *MigrationModeE2ETestSuite) testValidateMigratedCommanderFunctionalities() {
	s.waitForMinedBatchAndValidate(2, 4)

	s.waitForMinedBatchAndValidate(3, 8)

	s.SendNTransactions(7, dto.Transfer{
		FromStateID: ref.Uint32(0),
		ToStateID:   ref.Uint32(1),
		Amount:      models.NewUint256(90),
		Fee:         models.NewUint256(10),
		Nonce:       models.NewUint256(0),
	})

	s.waitForMinedBatchAndValidate(4, 7)
}

func (s *MigrationModeE2ETestSuite) createAdminRPCClient(cfg *config.Config) jsonrpc.RPCClient {
	url := fmt.Sprint("http://localhost:", cfg.API.Port)
	adminRPCClient := jsonrpc.NewClientWithOpts(url, &jsonrpc.RPCClientOpts{
		CustomHeaders: map[string]string{
			consts.AuthKeyHeader: cfg.API.AuthenticationKey,
		},
	})

	return adminRPCClient
}

func (s *MigrationModeE2ETestSuite) stopMiningBlocks() {
	err := s.gethRPCClient.Call(nil, "miner_stop")
	s.NoError(err)
	log.Printf("Stopping mining new blocks")
}

func (s *MigrationModeE2ETestSuite) startMiningBlocks() {
	err := s.gethRPCClient.Call(nil, "miner_start")
	s.NoError(err)
	log.Printf("Starting mining new block")
}

func (s *MigrationModeE2ETestSuite) configureCommander(params dto.ConfigureParams) {
	_, err := s.RPCClient.Call("admin_configure", []interface{}{params})
	s.NoError(err)
}

func (s *MigrationModeE2ETestSuite) waitForCommanderReadyStatus() {
	s.Eventually(func() bool {
		var status string
		err := s.RPCClient.CallFor(&status, "hubble_getStatus")
		s.NoError(err)
		return status == healthstatus.Ready
	}, 30*time.Second, testutils.TryInterval)
}

func (s *MigrationModeE2ETestSuite) prepareMigrationCommander() (*setup.InProcessCommander, jsonrpc.RPCClient) {
	newCommanderCfg := *s.Cfg
	newCommanderCfg.Bootstrap.Prune = true
	newCommanderCfg.Bootstrap.Migrate = true
	bootstrapURL := fmt.Sprint("http://localhost:", newCommanderCfg.API.Port)
	newCommanderCfg.Bootstrap.BootstrapNodeURL = &bootstrapURL
	newCommanderCfg.API.Port = "5555"
	newCommanderCfg.Metrics.Port = "2222"
	newCommanderCfg.Badger.Path += "_migration"
	migrationCommander, err := setup.CreateInProcessCommander(&newCommanderCfg, nil)
	s.NoError(err)

	migrationAdminRPCClient := s.createAdminRPCClient(&newCommanderCfg)

	return migrationCommander, migrationAdminRPCClient
}

func (s *MigrationModeE2ETestSuite) validateSuccessfulMigration() {
	s.validateFailedTxs()
	s.validatePendingTxs()
	s.validatePendingBatches()
}

func (s *MigrationModeE2ETestSuite) validatePendingTxs() {
	pendingTxs := make([]admintypes.Transaction, 0)
	err := s.RPCClient.CallFor(&pendingTxs, "admin_getPendingTransactions")
	s.NoError(err)
	s.Len(pendingTxs, 8)
}

func (s *MigrationModeE2ETestSuite) validateFailedTxs() {
	failedTxs := make([]admintypes.Transaction, 0)
	err := s.RPCClient.CallFor(&failedTxs, "admin_getFailedTransactions")
	s.NoError(err)
	s.Len(failedTxs, 2)
}

func (s *MigrationModeE2ETestSuite) validatePendingBatches() {
	pendingBatches := make([]admintypes.Batch, 0)
	err := s.RPCClient.CallFor(&pendingBatches, "admin_getPendingBatches")
	s.NoError(err)

	s.Len(pendingBatches, 1)
	s.Equal(models.MakeUint256(2), pendingBatches[0].ID)
	s.Len(pendingBatches[0].Commitments, 1)
	s.Len(pendingBatches[0].Commitments[0].Transactions, 4)
}

func (s *MigrationModeE2ETestSuite) validateBatch(batchID uint64, txCount int) {
	batch := s.GetBatchByID(batchID)

	s.Len(batch.Commitments, 1)

	commitmentID := models.CommitmentID{
		BatchID:      batch.ID,
		IndexInBatch: 0,
	}

	commitment := s.GetCommitment(commitmentID)
	s.Len(commitment.Transactions, txCount)
}

func (s *MigrationModeE2ETestSuite) waitForMinedBatchAndValidate(batchID uint64, txCount int) {
	s.WaitForBatchStatus(batchID, batchstatus.Mined)
	s.validateBatch(batchID, txCount)
}

func TestMigrationModeE2ETestSuite(t *testing.T) {
	suite.Run(t, new(MigrationModeE2ETestSuite))
}
