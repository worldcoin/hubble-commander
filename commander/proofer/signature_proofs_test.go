package proofer

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils/merkletree"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type SingatureProofsTestSuite struct {
	*require.Assertions
	suite.Suite
	storage            *st.TestStorage
	prooferCtx         *Context
	decodedCommitments []encoder.DecodedCommitment
	decodedBatch       eth.DecodedBatch
}

func (s *SingatureProofsTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *SingatureProofsTestSuite) SetupTest() {
	var err error
	s.storage, err = st.NewTestStorage()
	s.NoError(err)

	s.prooferCtx = NewContext(s.storage.Storage)
}

func (s *SingatureProofsTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *SingatureProofsTestSuite) TestPublicKeyProof() {
	account := &models.AccountLeaf{
		PubKeyID:  1,
		PublicKey: models.PublicKey{1, 2, 3},
	}
	err := s.storage.AccountTree.SetSingle(account)
	s.NoError(err)

	publicKeyProof, err := s.prooferCtx.publicKeyProof(account.PubKeyID)
	s.NoError(err)
	s.Equal(account.PublicKey, *publicKeyProof.PublicKey)
	s.Len(publicKeyProof.Witness, 32)
}

func (s *SingatureProofsTestSuite) TestReceiverPublicKeyProof() {
	account := &models.AccountLeaf{
		PubKeyID:  1,
		PublicKey: models.PublicKey{1, 2, 3},
	}
	err := s.storage.AccountTree.SetSingle(account)
	s.NoError(err)

	publicKeyProof, err := s.prooferCtx.receiverPublicKeyProof(account.PubKeyID)
	s.NoError(err)
	s.Equal(crypto.Keccak256Hash(account.PublicKey.Bytes()), publicKeyProof.PublicKeyHash)
	s.Len(publicKeyProof.Witness, 32)
}

func (s *SingatureProofsTestSuite) TestReceiverPublicKeyProof_NonexistentAccount() {
	publicKeyProof, err := s.prooferCtx.receiverPublicKeyProof(1)
	s.NoError(err)
	s.Equal(merkletree.GetZeroHash(0), publicKeyProof.PublicKeyHash)
	s.Len(publicKeyProof.Witness, 32)
}

func TestSignatureProofs(t *testing.T) {
	suite.Run(t, new(SingatureProofsTestSuite))
}
