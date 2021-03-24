package eth

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/eth/deployer"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/testutils/simulator"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/Worldcoin/hubble-commander/utils/ref"
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
	s.client, err = NewClient(sim.Account, NewClientParams{
		Rollup:          contracts.Rollup,
		AccountRegistry: contracts.AccountRegistry,
	})
	s.NoError(err)
}

func (s *RollupTestSuite) TearDownTest() {
	s.sim.Close()
}

func (s *RollupTestSuite) Test_SubmitTransfersBatch_ReturnsAccountTreeRootUsed() {
	expected, err := s.contracts.AccountRegistry.Root(nil)
	s.NoError(err)

	commitment := models.Commitment{
		Transactions:      utils.RandomBytes(12),
		FeeReceiver:       uint32(1234),
		CombinedSignature: models.MakeSignature(1, 2),
		PostStateRoot:     utils.RandomHash(),
	}

	_, accountRoot, err := s.client.SubmitTransfersBatch([]models.Commitment{commitment})
	s.NoError(err)

	s.Equal(common.BytesToHash(expected[:]), *accountRoot)
}

func (s *RollupTestSuite) Test_SubmitTransfersBatch_ReturnsBatchWithCorrectHash() {
	accountRoot, err := s.contracts.AccountRegistry.Root(nil)
	s.NoError(err)

	commitment := models.Commitment{
		Transactions:      utils.RandomBytes(12),
		FeeReceiver:       uint32(1234),
		CombinedSignature: models.MakeSignature(1, 2),
		PostStateRoot:     utils.RandomHash(),
		AccountTreeRoot:   ref.Hash(accountRoot),
	}

	batch, _, err := s.client.SubmitTransfersBatch([]models.Commitment{commitment})
	s.NoError(err)

	commitmentRoot := utils.HashTwo(commitment.LeafHash(), storage.GetZeroHash(0))
	s.Equal(commitmentRoot, batch.Hash)
}

func TestRollupTestSuite(t *testing.T) {
	suite.Run(t, new(RollupTestSuite))
}
