package disputer

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/stretchr/testify/suite"
)

type DisputeMMSignatureTestSuite struct {
	disputeSignatureTestSuite
}

func (s *DisputeMMSignatureTestSuite) SetupTest() {
	s.testSuiteWithContexts.SetupTest(batchtype.MassMigration, false)
	s.disputeSignatureTestSuite.setupTest()
}

func (s *DisputeMMSignatureTestSuite) TestDisputeSignature_DisputesBatchWithInvalidSignature() {
	wallets := s.setAccounts(s.domain)

	massMigration := testutils.MakeMassMigration(1, 2, 0, 50)
	s.signTx(&wallets[0], &massMigration)

	s.submitBatch(&massMigration)

	remoteBatches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 1)

	err = s.disputeSignature(remoteBatches[0].ToDecodedTxBatch(), models.MassMigrationArray{massMigration})
	s.NoError(err)

	checkRemoteBatchAfterDispute(s.Assertions, s.client, &remoteBatches[0].GetBase().ID)
}

func (s *DisputeMMSignatureTestSuite) TestDisputeSignature_ValidBatch() {
	wallets := s.setAccounts(s.domain)

	massMigration := testutils.MakeMassMigration(1, 2, 0, 50)
	s.signTx(&wallets[1], &massMigration)

	s.submitBatch(&massMigration)

	remoteBatches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 1)

	err = s.disputeSignature(remoteBatches[0].ToDecodedTxBatch(), models.MassMigrationArray{massMigration})
	s.NoError(err)
	_, err = s.client.GetContractBatch(&remoteBatches[0].GetBase().ID)
	s.NoError(err)
}

func TestDisputeMMSignatureTestSuite(t *testing.T) {
	suite.Run(t, new(DisputeMMSignatureTestSuite))
}
