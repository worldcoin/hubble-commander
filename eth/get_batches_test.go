package eth

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type GetBatchesTestSuite struct {
	*require.Assertions
	suite.Suite
	client *TestClient
}

func (s *GetBatchesTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *GetBatchesTestSuite) SetupTest() {
	client, err := NewTestClient()
	s.NoError(err)
	s.client = client
}

func (s *GetBatchesTestSuite) TearDownTest() {
	s.client.Close()
}

func (s *GetBatchesTestSuite) TestGetBatches() {
	commitment1 := models.Commitment{
		ID:                1,
		Type:              txtype.Transfer,
		Transactions:      []uint8{0, 0, 0, 0, 0, 0, 0, 1, 32, 4, 0, 0},
		FeeReceiver:       0,
		CombinedSignature: *s.mockSignature(),
	}

	commitment2 := models.Commitment{
		ID:                2,
		Type:              txtype.Transfer,
		Transactions:      []uint8{0, 0, 1, 0, 0, 0, 0, 0, 32, 1, 0, 0},
		FeeReceiver:       0,
		CombinedSignature: *s.mockSignature(),
	}

	finalisationBlocks, err := s.client.GetBlocksToFinalise()
	s.NoError(err)

	batch1, _, err := s.client.SubmitTransfersBatch([]models.Commitment{commitment1})
	s.NoError(err)
	_, _, err = s.client.SubmitTransfersBatch([]models.Commitment{commitment2})
	s.NoError(err)

	submissionBlockBatch1 := *batch1.FinalisationBlock - uint32(*finalisationBlocks)

	batches, err := s.client.GetBatches(&submissionBlockBatch1)
	s.NoError(err)
	s.Len(batches, 1)
}

func (s *GetBatchesTestSuite) mockSignature() *models.Signature {
	wallet, err := bls.NewRandomWallet(bls.Domain{1, 2, 3, 4})
	s.NoError(err)
	signature, err := wallet.Sign(utils.RandomBytes(4))
	s.NoError(err)
	return signature.ModelsSignature()
}

func TestGetBatchesTestSuite(t *testing.T) {
	suite.Run(t, new(GetBatchesTestSuite))
}
