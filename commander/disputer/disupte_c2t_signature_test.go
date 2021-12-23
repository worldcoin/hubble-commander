package disputer

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/stretchr/testify/suite"
)

type DisputeC2TSignatureTestSuite struct {
	disputeSignatureTestSuite
}

func (s *DisputeC2TSignatureTestSuite) SetupTest() {
	s.testSuiteWithContexts.SetupTest(batchtype.Create2Transfer, false)
	s.disputeSignatureTestSuite.setupTest()
}

func (s *DisputeC2TSignatureTestSuite) TestDisputeSignature_DisputesBatchWithInvalidSignature() {
	wallets := s.setAccounts(s.domain)

	transfer := testutils.MakeCreate2Transfer(0, nil, 0, 100, wallets[2].PublicKey())
	s.signTx(&wallets[1], &transfer)

	s.submitBatch(&transfer)

	remoteBatches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 1)

	err = s.disputeSignature(remoteBatches[0].ToDecodedTxBatch(), models.Create2TransferArray{transfer})
	s.NoError(err)

	checkRemoteBatchAfterDispute(s.Assertions, s.client, &remoteBatches[0].GetBase().ID)
}

func (s *DisputeC2TSignatureTestSuite) TestDisputeSignature_ValidBatch() {
	wallets := s.setAccounts(s.domain)

	transfer := testutils.MakeCreate2Transfer(0, nil, 0, 100, wallets[2].PublicKey())
	s.signTx(&wallets[0], &transfer)

	s.submitBatch(&transfer)

	remoteBatches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 1)

	err = s.disputeSignature(remoteBatches[0].ToDecodedTxBatch(), models.Create2TransferArray{transfer})
	s.NoError(err)
	_, err = s.client.GetBatch(&remoteBatches[0].GetBase().ID)
	s.NoError(err)
}

func TestDisputeC2TSignatureTestSuite(t *testing.T) {
	suite.Run(t, new(DisputeC2TSignatureTestSuite))
}
