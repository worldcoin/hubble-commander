package eth

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/Worldcoin/hubble-commander/utils/merkletree"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type SubmitBatchTestSuite struct {
	*require.Assertions
	suite.Suite
	client     *TestClient
	commitment models.TxCommitmentWithTxs
}

func (s *SubmitBatchTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *SubmitBatchTestSuite) SetupTest() {
	client, err := NewTestClient()
	s.NoError(err)
	s.client = client
	s.commitment = models.TxCommitmentWithTxs{
		TxCommitment: models.TxCommitment{
			CommitmentBase: models.CommitmentBase{
				Type:          batchtype.Transfer,
				PostStateRoot: utils.RandomHash(),
			},
			FeeReceiver:       uint32(1234),
			CombinedSignature: models.MakeRandomSignature(),
		},
		Transactions: utils.RandomBytes(12),
	}
}

func (s *SubmitBatchTestSuite) TearDownTest() {
	s.client.Close()
}

func (s *SubmitBatchTestSuite) TestSubmitTransfersBatchAndWait_ReturnsCorrectBatch() {
	commitment := s.commitment

	accountRoot, err := s.client.AccountRegistry.Root(nil)
	s.NoError(err)
	commitment.BodyHash = commitment.CalcBodyHash(accountRoot)
	commitmentRoot := utils.HashTwo(commitment.LeafHash(), merkletree.GetZeroHash(0))
	minFinalisationBlock := s.getMinFinalisationBlock()

	batch, err := s.client.SubmitTransfersBatchAndWait([]models.TxCommitmentWithTxs{commitment})
	s.NoError(err)

	s.Equal(models.MakeUint256(1), batch.ID)
	s.Equal(batchtype.Transfer, batch.Type)
	s.Equal(commitmentRoot, *batch.Hash)
	s.GreaterOrEqual(*batch.FinalisationBlock, minFinalisationBlock)
	s.Equal(common.BytesToHash(accountRoot[:]), *batch.AccountTreeRoot)
}

func (s *SubmitBatchTestSuite) TestSubmitCreate2TransfersBatchAndWait_ReturnsCorrectBatch() {
	commitment := s.commitment
	commitment.Type = batchtype.Create2Transfer

	accountRoot, err := s.client.AccountRegistry.Root(nil)
	s.NoError(err)
	commitment.BodyHash = commitment.CalcBodyHash(accountRoot)
	commitmentRoot := utils.HashTwo(commitment.LeafHash(), merkletree.GetZeroHash(0))
	minFinalisationBlock := s.getMinFinalisationBlock()

	batch, err := s.client.SubmitCreate2TransfersBatchAndWait([]models.TxCommitmentWithTxs{commitment})
	s.NoError(err)

	s.Equal(models.MakeUint256(1), batch.ID)
	s.Equal(batchtype.Create2Transfer, batch.Type)
	s.Equal(commitmentRoot, *batch.Hash)
	s.GreaterOrEqual(*batch.FinalisationBlock, minFinalisationBlock)
	s.Equal(common.BytesToHash(accountRoot[:]), *batch.AccountTreeRoot)
}

func (s *SubmitBatchTestSuite) getMinFinalisationBlock() uint32 {
	latestBlockNumber, err := s.client.Blockchain.GetLatestBlockNumber()
	s.NoError(err)
	blocksToFinalise, err := s.client.GetBlocksToFinalise()
	s.NoError(err)
	return uint32(*latestBlockNumber) + uint32(*blocksToFinalise)
}

func (s *SubmitBatchTestSuite) TestSubmitTransfersBatch_SubmitsBatchWithoutWaitingForItToBeMined() {
	tx, err := s.client.SubmitTransfersBatch([]models.TxCommitmentWithTxs{s.commitment})
	s.NoError(err)
	s.NotNil(tx)

	batches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(batches, 0)

	s.client.GetBackend().Commit()

	batches, err = s.client.GetAllBatches()
	s.NoError(err)
	s.Len(batches, 1)
}

func (s *SubmitBatchTestSuite) TestSubmitCreate2TransfersBatch_SubmitsBatchWithoutWaitingForItToBeMined() {
	commitment := s.commitment
	commitment.Type = batchtype.Create2Transfer

	tx, err := s.client.SubmitCreate2TransfersBatch([]models.TxCommitmentWithTxs{s.commitment})
	s.NoError(err)
	s.NotNil(tx)

	batches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(batches, 0)

	s.client.GetBackend().Commit()

	batches, err = s.client.GetAllBatches()
	s.NoError(err)
	s.Len(batches, 1)
}

func TestSubmitTransferTestSuite(t *testing.T) {
	suite.Run(t, new(SubmitBatchTestSuite))
}
