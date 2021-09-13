package executor

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/stretchr/testify/suite"
)

type DisputeTransferSignatureTestSuite struct {
	DisputeSignatureTestSuite
}

func (s *DisputeTransferSignatureTestSuite) SetupTest() {
	s.TestSuiteWithDisputeContext.SetupTest(txtype.Transfer)
	s.DisputeSignatureTestSuite.setupTest()
}

func (s *DisputeTransferSignatureTestSuite) TestSignatureProof() {
	s.setUserStatesAndAddAccounts()

	transfers := []models.Transfer{
		testutils.MakeTransfer(1, 2, 0, 50),
		testutils.MakeTransfer(0, 2, 0, 50),
		testutils.MakeTransfer(0, 1, 1, 75),
	}

	stateProofs := make([]models.StateMerkleProof, 0, len(transfers))
	expectedPublicKeys := make([]models.PublicKey, 0, len(transfers))
	for i := range transfers {
		stateProof, err := s.executionCtx.userStateProof(transfers[i].FromStateID)
		s.NoError(err)
		stateProofs = append(stateProofs, *stateProof)

		account, err := s.storage.AccountTree.Leaf(stateProof.UserState.PubKeyID)
		s.NoError(err)
		expectedPublicKeys = append(expectedPublicKeys, account.PublicKey)
	}

	signatureProof, err := s.disputeCtx.signatureProof(stateProofs)
	s.NoError(err)
	s.Len(signatureProof.UserStates, 3)
	s.Len(signatureProof.PublicKeys, 3)

	for i := range signatureProof.UserStates {
		s.Equal(stateProofs[i].UserState, signatureProof.UserStates[i].UserState)
		s.Equal(expectedPublicKeys[i], *signatureProof.PublicKeys[i].PublicKey)
	}
}

func (s *DisputeTransferSignatureTestSuite) TestDisputeSignature_DisputesBatchWithInvalidSignature() {
	wallets := s.setUserStatesAndAddAccounts()

	transfer := testutils.MakeTransfer(1, 2, 0, 50)
	signTransfer(s.T(), &wallets[0], &transfer)

	submitTransferBatch(s.Assertions, s.client, s.rollupCtx, &transfer)

	remoteBatches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 1)

	err = s.disputeSignature(&remoteBatches[0], models.TransferArray{transfer})
	s.NoError(err)

	checkRemoteBatchAfterDispute(s.Assertions, s.client, &remoteBatches[0].ID)
}

func (s *DisputeTransferSignatureTestSuite) TestDisputeSignature_ValidBatch() {
	wallets := s.setUserStatesAndAddAccounts()

	transfer := testutils.MakeTransfer(1, 2, 0, 50)
	signTransfer(s.T(), &wallets[1], &transfer)

	submitTransferBatch(s.Assertions, s.client, s.rollupCtx, &transfer)

	remoteBatches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 1)

	err = s.disputeSignature(&remoteBatches[0], models.TransferArray{transfer})
	s.NoError(err)
	_, err = s.client.GetBatch(&remoteBatches[0].ID)
	s.NoError(err)
}

func TestDisputeTransferSignatureTestSuite(t *testing.T) {
	suite.Run(t, new(DisputeTransferSignatureTestSuite))
}
