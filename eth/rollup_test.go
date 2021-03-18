package eth

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/eth/deployer"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/testutils/simulator"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type RollupTestSuite struct {
	*require.Assertions
	suite.Suite
	sim       *simulator.Simulator
	contracts *deployer.RollupContracts
	client    *Client
}

func (s *RollupTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *RollupTestSuite) SetupTest() {
	sim, err := simulator.NewAutominingSimulator()
	s.NoError(err)
	s.sim = sim

	contracts, err := deployer.DeployRollup(sim)
	s.NoError(err)
	s.contracts = contracts
	s.client = NewTestClient(sim.Account, contracts.Rollup, contracts.AccountRegistry)
}

func (s *RollupTestSuite) TearDownTest() {
	s.sim.Close()
}

func (s *RollupTestSuite) Test_SubmitTransfersBatch() {
	accountRoot, err := s.contracts.AccountRegistry.Root(nil)
	s.NoError(err)

	signature := models.Signature{models.MakeUint256(1), models.MakeUint256(2)}
	feeReceiver := uint32(1234)
	txs := utils.RandomBytes(12)
	bodyHash, err := encoder.GetCommitmentBodyHash(accountRoot, signature, feeReceiver, txs)
	s.NoError(err)

	postStateRoot := utils.RandomHash()
	leafHash := utils.HashTwo(postStateRoot, *bodyHash)

	commitment := models.Commitment{
		LeafHash:          leafHash,
		PostStateRoot:     postStateRoot,
		BodyHash:          *bodyHash,
		AccountTreeRoot:   accountRoot,
		CombinedSignature: signature,
		FeeReceiver:       feeReceiver,
		Transactions:      txs,
	}

	batchID, err := s.client.SubmitTransfersBatch([]*models.Commitment{&commitment})
	s.NoError(err)

	batch, err := s.contracts.Rollup.GetBatch(nil, &batchID.Int)
	s.NoError(err)

	commitmentRoot := utils.HashTwo(leafHash, storage.GetZeroHash(0))
	s.Equal(commitmentRoot, common.BytesToHash(batch.CommitmentRoot[:]))
}

func TestRollupTestSuite(t *testing.T) {
	suite.Run(t, new(RollupTestSuite))
}
