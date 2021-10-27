package eth

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type DisputeTransitionTestSuite struct {
	*require.Assertions
	suite.Suite
	client *TestClient
}

func (s *DisputeTransitionTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *DisputeTransitionTestSuite) SetupTest() {
	var err error
	s.client, err = NewTestClient()
	s.NoError(err)
}

func (s *DisputeTransitionTestSuite) TearDownTest() {
	s.client.Close()
}

func (s *DisputeTransitionTestSuite) TestDisputeTransitionTransfer_ReturnsRevertMessage() {
	merklePath := models.MakeMerklePathFromLeafID(1)
	previousCommitmentProof := &models.CommitmentInclusionProof{
		StateRoot: utils.RandomHash(),
		BodyRoot:  utils.RandomHash(),
		Path:      &merklePath,
		Witness:   []common.Hash{},
	}
	targetCommitmentProof := &models.TransferCommitmentInclusionProof{
		StateRoot: common.Hash{},
		Body: &models.TransferBody{
			AccountRoot:  utils.RandomHash(),
			Signature:    models.Signature{},
			FeeReceiver:  0,
			Transactions: nil,
		},
		Path:    &merklePath,
		Witness: []common.Hash{},
	}

	err := s.client.DisputeTransitionTransfer(
		models.NewUint256(1),
		utils.NewRandomHash(),
		previousCommitmentProof,
		targetCommitmentProof,
		[]models.StateMerkleProof{},
	)
	s.Error(err)
	s.Equal(err.Error(), "execution reverted: Batch id greater than total number of batches, invalid batch id")
}

func TestDisputeTransitionTestSuite(t *testing.T) {
	suite.Run(t, new(DisputeTransitionTestSuite))
}
