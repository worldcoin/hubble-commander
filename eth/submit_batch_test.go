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
		CombinedSignature: models.MakeSignature(1, 2),
		PostStateRoot:     utils.RandomHash(),
	}
}

func (s *SubmitBatchTestSuite) TearDownTest() {
	s.client.Close()
}

func (s *SubmitBatchTestSuite) TestSubmitTransfersBatch_ReturnsAccountTreeRootUsed() {
	expected, err := s.client.AccountRegistry.Root(nil)
	s.NoError(err)

	_, accountRoot, err := s.client.SubmitTransfersBatch([]models.Commitment{s.commitment})
	s.NoError(err)

	s.Equal(common.BytesToHash(expected[:]), *accountRoot)
}

func (s *SubmitBatchTestSuite) TestSubmitTransfersBatch_ReturnsBatchWithCorrectHashAndType() {
	accountRoot, err := s.client.AccountRegistry.Root(nil)
	s.NoError(err)

	commitment := s.commitment
	commitment.AccountTreeRoot = ref.Hash(accountRoot)

	batch, _, err := s.client.SubmitTransfersBatch([]models.Commitment{commitment})
	s.NoError(err)

	commitmentRoot := utils.HashTwo(commitment.LeafHash(), storage.GetZeroHash(0))
	s.Equal(commitmentRoot, batch.Hash)
	s.Equal(txtype.Transfer, batch.Type)
}

func (s *SubmitBatchTestSuite) TestSubmitTransfersBatch_Create2TransferReturnsAccountTreeRootUsed() {
	expected, err := s.client.AccountRegistry.Root(nil)
	s.NoError(err)

	commitment := s.commitment
	commitment.Type = txtype.Create2Transfer

	_, accountRoot, err := s.client.SubmitCreate2TransfersBatch([]models.Commitment{s.commitment})
	s.NoError(err)

	s.Equal(common.BytesToHash(expected[:]), *accountRoot)
}

func (s *SubmitBatchTestSuite) TestSubmitTransfersBatch_Create2TransferReturnsBatchWithCorrectHashAndType() {
	accountRoot, err := s.client.AccountRegistry.Root(nil)
	s.NoError(err)

	commitment := s.commitment
	commitment.Type = txtype.Create2Transfer
	commitment.AccountTreeRoot = ref.Hash(accountRoot)

	batch, _, err := s.client.SubmitCreate2TransfersBatch([]models.Commitment{commitment})
	s.NoError(err)

	commitmentRoot := utils.HashTwo(commitment.LeafHash(), storage.GetZeroHash(0))
	s.Equal(commitmentRoot, batch.Hash)
	s.Equal(txtype.Create2Transfer, batch.Type)
}

func TestSubmitTransferTestSuite(t *testing.T) {
	suite.Run(t, new(SubmitBatchTestSuite))
}
