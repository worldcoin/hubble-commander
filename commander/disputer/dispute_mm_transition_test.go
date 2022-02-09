package disputer

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/commander/syncer"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/stretchr/testify/suite"
)

type DisputeMMTransitionTestSuite struct {
	disputeTransitionTestSuite
}

func (s *DisputeMMTransitionTestSuite) SetupTest() {
	s.disputeTransitionTestSuite.SetupTest(batchtype.MassMigration, true)
}

func (s *DisputeMMTransitionTestSuite) TestDisputeTransition_RemovesInvalidBatch() {
	txs := models.MassMigrationArray{
		testutils.MakeMassMigration(0, 1, 0, 100),
		testutils.MakeMassMigration(1, 1, 0, 100),
		testutils.MakeMassMigration(2, 1, 0, 20),
		testutils.MakeMassMigration(2, 1, 1, 20),
	}
	s.submitInvalidBatch(txs)

	remoteBatches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 1)

	proofs := s.getInvalidBatchStateProofs(remoteBatches[0])
	err = s.disputeCtx.DisputeTransition(remoteBatches[0].ToDecodedTxBatch(), 1, proofs)
	s.NoError(err)

	checkRemoteBatchAfterDispute(s.Assertions, s.client, &remoteBatches[0].GetBase().ID)
}

func (s *DisputeMMTransitionTestSuite) TestDisputeTransition_ValidBatch() {
	tx := testutils.MakeMassMigration(0, 1, 0, 50)
	proofs := s.getValidBatchStateProofs(syncer.NewSyncedMMs(models.MassMigrationArray{tx}))

	s.submitBatch(&tx)

	remoteBatches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 1)

	err = s.disputeCtx.DisputeTransition(remoteBatches[0].ToDecodedTxBatch(), 0, proofs)
	s.NoError(err)
	_, err = s.client.GetContractBatch(&remoteBatches[0].GetBase().ID)
	s.NoError(err)
}

func TestDisputeMMTransitionTestSuite(t *testing.T) {
	suite.Run(t, new(DisputeMMTransitionTestSuite))
}
