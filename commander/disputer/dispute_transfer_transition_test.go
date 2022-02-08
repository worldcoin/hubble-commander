package disputer

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/commander/syncer"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/stretchr/testify/suite"
)

type DisputeTransferTransitionTestSuite struct {
	disputeTransitionTestSuite
}

func (s *DisputeTransferTransitionTestSuite) SetupTest() {
	s.disputeTransitionTestSuite.SetupTest(batchtype.Transfer, true)
}

func (s *DisputeTransferTransitionTestSuite) TestDisputeTransition_RemovesInvalidBatch() {
	commitmentTxs := models.TransferArray{
		testutils.MakeTransfer(0, 2, 0, 100),
		testutils.MakeTransfer(1, 0, 0, 100),
		testutils.MakeTransfer(2, 0, 0, 20),
		testutils.MakeTransfer(2, 0, 1, 20),
	}
	s.submitInvalidBatch(commitmentTxs)

	remoteBatches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 1)

	proofs := s.getInvalidBatchStateProofs(remoteBatches[0])
	err = s.disputeCtx.DisputeTransition(remoteBatches[0].ToDecodedTxBatch(), 1, proofs)
	s.NoError(err)

	checkRemoteBatchAfterDispute(s.Assertions, s.client, &remoteBatches[0].GetBase().ID)
}

func (s *DisputeTransferTransitionTestSuite) TestDisputeTransition_FirstCommitment() {
	commitmentTxs := models.TransferArray{testutils.MakeTransfer(0, 2, 0, 100)}
	s.submitInvalidBatch(commitmentTxs)

	remoteBatches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 1)

	proofs := s.getInvalidBatchStateProofs(remoteBatches[0])
	err = s.disputeCtx.DisputeTransition(remoteBatches[0].ToDecodedTxBatch(), 0, proofs)
	s.NoError(err)

	checkRemoteBatchAfterDispute(s.Assertions, s.client, &remoteBatches[0].GetBase().ID)
}

func (s *DisputeTransferTransitionTestSuite) TestDisputeTransition_ValidBatch() {
	tx := testutils.MakeTransfer(0, 2, 0, 50)
	proofs := s.getValidBatchStateProofs(syncer.NewSyncedTransfers(models.TransferArray{tx}))

	s.submitBatch(&tx)

	remoteBatches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 1)

	err = s.disputeCtx.DisputeTransition(remoteBatches[0].ToDecodedTxBatch(), 0, proofs)
	s.NoError(err)
	_, err = s.client.GetContractBatch(&remoteBatches[0].GetBase().ID)
	s.NoError(err)
}

func TestDisputeTransferTransitionTestSuite(t *testing.T) {
	suite.Run(t, new(DisputeTransferTransitionTestSuite))
}
