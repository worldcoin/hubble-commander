package eth

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type SubmitBatchTestSuite struct {
	*require.Assertions
	suite.Suite
	client     *TestClient
	commitment models.Commitment
}

func (s *SubmitBatchTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *SubmitBatchTestSuite) SetupTest() {
	client, err := NewTestClient()
	s.NoError(err)
	s.client = client
	s.commitment = models.Commitment{
		Type:              txtype.Transfer,
		Transactions:      utils.RandomBytes(12),
		FeeReceiver:       uint32(1234),
		CombinedSignature: models.MakeRandomSignature(),
		PostStateRoot:     utils.RandomHash(),
	}
}

func (s *SubmitBatchTestSuite) TearDownTest() {
	s.client.Close()
}

func (s *SubmitBatchTestSuite) TestSubmitTransfersBatchAndMine_ReturnsAccountTreeRootUsed() {
	expected, err := s.client.AccountRegistry.Root(nil)
	s.NoError(err)

	_, accountRoot, err := s.client.SubmitTransfersBatchAndMine([]models.Commitment{s.commitment})
	s.NoError(err)

	s.Equal(common.BytesToHash(expected[:]), *accountRoot)
}

func (s *SubmitBatchTestSuite) TestSubmitTransfersBatchAndMine_ReturnsBatchWithCorrectHashAndType() {
	accountRoot, err := s.client.AccountRegistry.Root(nil)
	s.NoError(err)

	commitment := s.commitment
	commitment.AccountTreeRoot = ref.Hash(accountRoot)

	batch, _, err := s.client.SubmitTransfersBatchAndMine([]models.Commitment{commitment})
	s.NoError(err)

	commitmentRoot := utils.HashTwo(commitment.LeafHash(), storage.GetZeroHash(0))
	s.Equal(commitmentRoot, *batch.Hash)
	s.Equal(txtype.Transfer, batch.Type)
}

func (s *SubmitBatchTestSuite) TestSubmitCreate2TransfersBatchAndMine_ReturnsAccountTreeRootUsed() {
	expected, err := s.client.AccountRegistry.Root(nil)
	s.NoError(err)

	commitment := s.commitment
	commitment.Type = txtype.Create2Transfer

	_, accountRoot, err := s.client.SubmitCreate2TransfersBatchAndMine([]models.Commitment{s.commitment})
	s.NoError(err)

	s.Equal(common.BytesToHash(expected[:]), *accountRoot)
}

func (s *SubmitBatchTestSuite) TestSubmitCreate2TransfersBatchAndMine_ReturnsBatchWithCorrectHashAndType() {
	accountRoot, err := s.client.AccountRegistry.Root(nil)
	s.NoError(err)

	commitment := s.commitment
	commitment.Type = txtype.Create2Transfer
	commitment.AccountTreeRoot = ref.Hash(accountRoot)

	batch, _, err := s.client.SubmitCreate2TransfersBatchAndMine([]models.Commitment{commitment})
	s.NoError(err)

	commitmentRoot := utils.HashTwo(commitment.LeafHash(), storage.GetZeroHash(0))
	s.Equal(commitmentRoot, *batch.Hash)
	s.Equal(txtype.Create2Transfer, batch.Type)
}

func (s *SubmitBatchTestSuite) TestSubmitTransfersBatch_SubmitsBatchWithoutMining() {
	tx, err := s.client.SubmitTransfersBatch([]models.Commitment{s.commitment})
	s.NoError(err)
	s.NotNil(tx)

	batches, hashes, err := s.client.Client.GetBatches(ref.Uint32(0))
	s.NoError(err)
	s.Len(batches, 0)
	s.Len(hashes, 0)

	s.client.Commit()

	batches, hashes, err = s.client.Client.GetBatches(ref.Uint32(0))
	s.NoError(err)
	s.Len(batches, 1)
	s.Len(hashes, 1)
}

func (s *SubmitBatchTestSuite) TestSubmitCreate2TransfersBatch_SubmitsBatchWithoutMining() {
	commitment := s.commitment
	commitment.Type = txtype.Create2Transfer

	tx, err := s.client.SubmitCreate2TransfersBatch([]models.Commitment{s.commitment})
	s.NoError(err)
	s.NotNil(tx)

	batches, hashes, err := s.client.Client.GetBatches(ref.Uint32(0))
	s.NoError(err)
	s.Len(batches, 0)
	s.Len(hashes, 0)

	s.client.Commit()

	batches, hashes, err = s.client.Client.GetBatches(ref.Uint32(0))
	s.NoError(err)
	s.Len(batches, 1)
	s.Len(hashes, 1)
}

func TestSubmitTransferTestSuite(t *testing.T) {
	suite.Run(t, new(SubmitBatchTestSuite))
}
