package disputer

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/suite"
)

type DisputeSignatureProofsTestSuite struct {
	TestSuiteWithContexts
}

func (s *DisputeSignatureProofsTestSuite) SetupTest() {
	s.TestSuiteWithContexts.SetupTest(batchtype.Transfer)
}

func (s *DisputeSignatureProofsTestSuite) TestUserStateProof() {
	userState := createUserState(1, 300, 1)
	witness, err := s.storage.StateTree.Set(1, userState)
	s.NoError(err)

	stateProof, err := s.syncCtx.UserStateProof(1)
	s.NoError(err)
	s.Equal(userState, stateProof.UserState)
	s.Equal(witness, stateProof.Witness)
}

func (s *DisputeSignatureProofsTestSuite) TestPublicKeyProof() {
	account := &models.AccountLeaf{
		PubKeyID:  1,
		PublicKey: models.PublicKey{1, 2, 3},
	}
	err := s.storage.AccountTree.SetSingle(account)
	s.NoError(err)

	publicKeyProof, err := s.disputeCtx.publicKeyProof(account.PubKeyID)
	s.NoError(err)
	s.Equal(account.PublicKey, *publicKeyProof.PublicKey)
	s.Len(publicKeyProof.Witness, 32)
}

func (s *DisputeSignatureProofsTestSuite) TestReceiverPublicKeyProof() {
	account := &models.AccountLeaf{
		PubKeyID:  1,
		PublicKey: models.PublicKey{1, 2, 3},
	}
	err := s.storage.AccountTree.SetSingle(account)
	s.NoError(err)

	publicKeyHash := crypto.Keccak256Hash(account.PublicKey.Bytes())

	publicKeyProof, err := s.disputeCtx.receiverPublicKeyProof(account.PubKeyID)
	s.NoError(err)
	s.Equal(publicKeyHash, publicKeyProof.PublicKeyHash)
	s.Len(publicKeyProof.Witness, 32)
}

func TestDisputeSignatureProofsTestSuite(t *testing.T) {
	suite.Run(t, new(DisputeSignatureProofsTestSuite))
}
