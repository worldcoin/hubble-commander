package disputer

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/stretchr/testify/suite"
)

type DisputeTransferSignatureTestSuite struct {
	disputeSignatureTestSuite
}

func (s *DisputeTransferSignatureTestSuite) SetupTest() {
	s.testSuiteWithContexts.SetupTest(batchtype.Transfer, false)
	s.disputeSignatureTestSuite.setupTest()
}

func (s *DisputeTransferSignatureTestSuite) TestDisputeSignature_DisputesBatchWithInvalidSignature() {
	wallets := s.setAccounts(s.domain)

	transfer := testutils.MakeTransfer(1, 2, 0, 50)
	s.signTx(&wallets[0], &transfer)

	s.submitBatch(&transfer)

	remoteBatches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 1)

	err = s.disputeSignature(remoteBatches[0].ToDecodedTxBatch(), models.TransferArray{transfer})
	s.NoError(err)

	checkRemoteBatchAfterDispute(s.Assertions, s.client, &remoteBatches[0].GetBase().ID)
}

func (s *DisputeTransferSignatureTestSuite) TestDisputeSignature_ValidBatch() {
	wallets := s.setAccounts(s.domain)

	transfer := testutils.MakeTransfer(1, 2, 0, 50)
	s.signTx(&wallets[1], &transfer)

	s.submitBatch(&transfer)

	remoteBatches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 1)

	err = s.disputeSignature(remoteBatches[0].ToDecodedTxBatch(), models.TransferArray{transfer})
	s.NoError(err)
	_, err = s.client.GetBatch(&remoteBatches[0].GetBase().ID)
	s.NoError(err)
}

func TestDisputeTransferSignatureTestSuite(t *testing.T) {
	suite.Run(t, new(DisputeTransferSignatureTestSuite))
}
