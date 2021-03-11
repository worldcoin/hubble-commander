package eth

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/testutils/deployer"
	"github.com/Worldcoin/hubble-commander/testutils/simulator"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ClientTestSuite struct {
	*require.Assertions
	suite.Suite
	sim       *simulator.Simulator
	contracts *deployer.RollupContracts
	client    *Client
}

func (s *ClientTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *ClientTestSuite) SetupTest() {
	sim, err := simulator.NewAutominingSimulator()
	s.NoError(err)
	s.sim = sim

	contracts, err := deployer.DeployRollup(sim)
	s.NoError(err)
	s.contracts = contracts
	s.client = NewTestClient(sim.Account, contracts.Rollup)
}

func (s *ClientTestSuite) TearDownTest() {
	s.sim.Close()
}

func (s *ClientTestSuite) Test_SubmitTransfersBatch() {
	s.T().Skip()

	accountRoot, err := s.contracts.AccountRegistry.Root(nil)
	s.NoError(err)

	signature := models.Signature{models.MakeUint256(1), models.MakeUint256(2)}
	feeReceiver := models.MakeUint256(1234)
	txs := utils.RandomBytes(12)
	bodyHash, err := encoder.GetCommitmentBodyHash(accountRoot, signature, feeReceiver, txs)
	s.NoError(err)

	postStateRoot := utils.RandomHash()
	leafHash := storage.HashTwo(postStateRoot, *bodyHash)

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

	commitmentRoot := storage.HashTwo(leafHash, storage.GetZeroHash(0))
	s.Equal(commitmentRoot, common.BytesToHash(batch.CommitmentRoot[:]))
}

func TestClientTestSuite(t *testing.T) {
	suite.Run(t, new(ClientTestSuite))
}
