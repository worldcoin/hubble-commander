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

//TODO-ref: move to correct package
func (s *DisputeTransferSignatureTestSuite) TestSignatureProof() {
	s.setAccounts(s.domain)

	transfers := []models.Transfer{
		testutils.MakeTransfer(1, 2, 0, 50),
		testutils.MakeTransfer(0, 2, 0, 50),
		testutils.MakeTransfer(0, 1, 1, 75),
	}

	stateProofs := make([]models.StateMerkleProof, 0, len(transfers))
	expectedPublicKeys := make([]models.PublicKey, 0, len(transfers))
	for i := range transfers {
		stateProof, err := s.syncCtx.UserStateProof(transfers[i].FromStateID)
		s.NoError(err)
		stateProofs = append(stateProofs, *stateProof)

		account, err := s.storage.AccountTree.Leaf(stateProof.UserState.PubKeyID)
		s.NoError(err)
		expectedPublicKeys = append(expectedPublicKeys, account.PublicKey)
	}

	signatureProof, err := s.disputeCtx.proverCtx.SignatureProof(stateProofs)
	s.NoError(err)
	s.Len(signatureProof.UserStates, 3)
	s.Len(signatureProof.PublicKeys, 3)

	for i := range signatureProof.UserStates {
		s.Equal(stateProofs[i].UserState, signatureProof.UserStates[i].UserState)
		s.Equal(expectedPublicKeys[i], *signatureProof.PublicKeys[i].PublicKey)
	}
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
