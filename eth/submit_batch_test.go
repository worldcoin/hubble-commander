package eth

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
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

func (s *SubmitBatchTestSuite) TestSubmitTransfersBatchAndWait_ReturnsBatchWithCorrectAccountTreeRoot() {
	expected, err := s.client.AccountRegistry.Root(nil)
	s.NoError(err)

	batch, err := s.client.SubmitTransfersBatchAndWait([]models.Commitment{s.commitment})
	s.NoError(err)

	s.Equal(common.BytesToHash(expected[:]), *batch.AccountTreeRoot)
}

func (s *SubmitBatchTestSuite) TestSubmitTransfersBatchAndWait_ReturnsBatchWithCorrectHashAndType() {
	accountRoot, err := s.client.AccountRegistry.Root(nil)
	s.NoError(err)

	commitment := s.commitment

	batch, err := s.client.SubmitTransfersBatchAndWait([]models.Commitment{commitment})
	s.NoError(err)

	commitmentRoot := utils.HashTwo(commitment.LeafHash(accountRoot), storage.GetZeroHash(0))
	s.Equal(commitmentRoot, *batch.Hash)
	s.Equal(txtype.Transfer, batch.Type)
}

func (s *SubmitBatchTestSuite) TestSubmitCreate2TransfersBatchAndWait_ReturnsBatchWithCorrectAccountTreeRoot() {
	expected, err := s.client.AccountRegistry.Root(nil)
	s.NoError(err)

	commitment := s.commitment
	commitment.Type = txtype.Create2Transfer

	batch, err := s.client.SubmitCreate2TransfersBatchAndWait([]models.Commitment{s.commitment})
	s.NoError(err)

	s.Equal(common.BytesToHash(expected[:]), *batch.AccountTreeRoot)
}

func (s *SubmitBatchTestSuite) TestSubmitCreate2TransfersBatchAndWait_ReturnsBatchWithCorrectHashAndType() {
	accountRoot, err := s.client.AccountRegistry.Root(nil)
	s.NoError(err)

	commitment := s.commitment
	commitment.Type = txtype.Create2Transfer

	batch, err := s.client.SubmitCreate2TransfersBatchAndWait([]models.Commitment{commitment})
	s.NoError(err)

	commitmentRoot := utils.HashTwo(commitment.LeafHash(accountRoot), storage.GetZeroHash(0))
	s.Equal(commitmentRoot, *batch.Hash)
	s.Equal(txtype.Create2Transfer, batch.Type)
}

func (s *SubmitBatchTestSuite) TestSubmitTransfersBatch_SubmitsBatchWithoutWaitingForItToBeMined() {
	tx, err := s.client.SubmitTransfersBatch([]models.Commitment{s.commitment})
	s.NoError(err)
	s.NotNil(tx)

	batches, err := s.client.Client.GetBatches(&bind.FilterOpts{})
	s.NoError(err)
	s.Len(batches, 0)

	s.client.Commit()

	batches, err = s.client.Client.GetBatches(&bind.FilterOpts{})
	s.NoError(err)
	s.Len(batches, 1)
}

func (s *SubmitBatchTestSuite) TestSubmitCreate2TransfersBatch_SubmitsBatchWithoutWaitingForItToBeMined() {
	commitment := s.commitment
	commitment.Type = txtype.Create2Transfer

	tx, err := s.client.SubmitCreate2TransfersBatch([]models.Commitment{s.commitment})
	s.NoError(err)
	s.NotNil(tx)

	batches, err := s.client.Client.GetBatches(&bind.FilterOpts{})
	s.NoError(err)
	s.Len(batches, 0)

	s.client.Commit()

	batches, err = s.client.Client.GetBatches(&bind.FilterOpts{})
	s.NoError(err)
	s.Len(batches, 1)
}

func TestSubmitTransferTestSuite(t *testing.T) {
	suite.Run(t, new(SubmitBatchTestSuite))
}
