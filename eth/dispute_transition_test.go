package eth

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type DisputeTransitionTestSuite struct {
	*require.Assertions
	suite.Suite
	testSuiteWithRequestsSending
	client *TestClient
}

func (s *DisputeTransitionTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *DisputeTransitionTestSuite) SetupTest() {
	var err error
	s.client, err = NewTestClient()
	s.NoError(err)

	s.StartTxsSending(s.T(), s.client.TxsChannels.Requests)
}

func (s *DisputeTransitionTestSuite) TearDownTest() {
	s.StopTxsSending()
	s.client.Close()
}

func (s *DisputeTransitionTestSuite) TestDisputeTransitionTransfer_ReturnsRevertMessage() {
	commitment := &models.TxCommitmentWithTxs{
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

	batch, err := s.client.SubmitTransfersBatchAndWait(models.NewUint256(1), []models.CommitmentWithTxs{commitment})
	s.NoError(err)

	merklePath := models.MakeMerklePathFromLeafID(1)
	previousCommitmentProof := &models.CommitmentInclusionProof{
		CommitmentInclusionProofBase: models.CommitmentInclusionProofBase{
			StateRoot: utils.RandomHash(),
			Path:      &merklePath,
			Witness:   []common.Hash{},
		},
		BodyRoot: utils.RandomHash(),
	}
	targetCommitmentProof := &models.TransferCommitmentInclusionProof{
		CommitmentInclusionProofBase: models.CommitmentInclusionProofBase{
			StateRoot: common.Hash{},
			Path:      &merklePath,
			Witness:   []common.Hash{},
		},
		Body: &models.TransferBody{
			AccountRoot:  utils.RandomHash(),
			Signature:    models.Signature{},
			FeeReceiver:  0,
			Transactions: nil,
		},
	}

	err = s.client.DisputeTransitionTransfer(
		models.NewUint256(1),
		batch.Hash,
		previousCommitmentProof,
		targetCommitmentProof,
		[]models.StateMerkleProof{},
	)
	var disputeError *DisputeTxRevertedError
	s.ErrorAs(err, &disputeError)
	s.Equal("dispute of batch #1 failed: execution reverted: previous commitment has wrong path", err.Error())
}

func TestDisputeTransitionTestSuite(t *testing.T) {
	suite.Run(t, new(DisputeTransitionTestSuite))
}
