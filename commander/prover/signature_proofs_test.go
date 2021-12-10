package prover

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/Worldcoin/hubble-commander/utils/merkletree"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type SignatureProofsTestSuite struct {
	*require.Assertions
	suite.Suite
	storage   *st.TestStorage
	proverCtx *Context
}

func (s *SignatureProofsTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *SignatureProofsTestSuite) SetupTest() {
	var err error
	s.storage, err = st.NewTestStorage()
	s.NoError(err)

	s.proverCtx = NewContext(s.storage.Storage)
}

func (s *SignatureProofsTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *SignatureProofsTestSuite) TestSignatureProof() {
	senderAccount := []models.AccountLeaf{
		{PubKeyID: 0, PublicKey: models.PublicKey{0, 2, 3}},
		{PubKeyID: 1, PublicKey: models.PublicKey{1, 2, 3}},
	}
	s.addAccounts(senderAccount...)

	transfers := []models.Transfer{
		testutils.MakeTransfer(1, 2, 0, 50),
		testutils.MakeTransfer(0, 2, 0, 50),
		testutils.MakeTransfer(0, 1, 1, 75),
	}

	stateProofs := make([]models.StateMerkleProof, 0, len(transfers))
	expectedPublicKeys := make([]models.PublicKey, 0, len(transfers))
	for i := range transfers {
		stateProofs = append(stateProofs, models.StateMerkleProof{
			UserState: &models.UserState{
				PubKeyID: transfers[i].FromStateID,
				TokenID:  models.MakeUint256(0),
				Balance:  models.MakeUint256(100),
				Nonce:    models.MakeUint256(0),
			},
		})

		account, err := s.storage.AccountTree.Leaf(transfers[i].FromStateID)
		s.NoError(err)
		expectedPublicKeys = append(expectedPublicKeys, account.PublicKey)
	}

	signatureProof, err := s.proverCtx.SignatureProof(stateProofs)
	s.NoError(err)
	s.Len(signatureProof.UserStates, len(transfers))
	s.Len(signatureProof.PublicKeys, len(transfers))

	for i := range signatureProof.UserStates {
		s.Equal(stateProofs[i].UserState, signatureProof.UserStates[i].UserState)
		s.Equal(expectedPublicKeys[i], *signatureProof.PublicKeys[i].PublicKey)
	}
}

func (s *SignatureProofsTestSuite) TestSignatureProofWithReceiver() {
	senderAccounts := []models.AccountLeaf{
		{PubKeyID: 0, PublicKey: models.PublicKey{0, 2, 3}},
		{PubKeyID: 1, PublicKey: models.PublicKey{1, 3, 4}},
	}
	receiverAccounts := []models.AccountLeaf{
		{PubKeyID: 2, PublicKey: models.PublicKey{2, 2, 3}},
		{PubKeyID: 3, PublicKey: models.PublicKey{3, 3, 4}},
		{PubKeyID: 4, PublicKey: models.PublicKey{4, 4, 5}},
	}

	s.addAccounts(senderAccounts...)
	s.addAccounts(receiverAccounts...)

	transfers := []models.Create2Transfer{
		testutils.MakeCreate2Transfer(1, ref.Uint32(3), 0, 50, &receiverAccounts[0].PublicKey),
		testutils.MakeCreate2Transfer(0, ref.Uint32(4), 0, 50, &receiverAccounts[1].PublicKey),
		testutils.MakeCreate2Transfer(0, ref.Uint32(5), 1, 75, &receiverAccounts[2].PublicKey),
	}
	pubKeyIDs := []uint32{2, 3, 4}

	stateProofs := make([]models.StateMerkleProof, 0, len(transfers))
	senderPublicKeys := make([]models.PublicKey, 0, len(transfers))
	receiverPublicKeys := make([]common.Hash, 0, len(transfers))
	for i := range transfers {
		stateProofs = append(stateProofs, models.StateMerkleProof{
			UserState: &models.UserState{
				PubKeyID: senderAccounts[transfers[i].FromStateID].PubKeyID,
				TokenID:  models.MakeUint256(0),
				Balance:  models.MakeUint256(100),
				Nonce:    models.MakeUint256(0),
			},
		})

		senderPublicKeys = append(senderPublicKeys, senderAccounts[transfers[i].FromStateID].PublicKey)
		receiverPublicKeys = append(receiverPublicKeys, crypto.Keccak256Hash(transfers[i].ToPublicKey.Bytes()))
	}

	serializedTxs, err := encoder.SerializeCreate2Transfers(transfers, pubKeyIDs)
	s.NoError(err)

	commitment := &encoder.DecodedCommitment{Transactions: serializedTxs}
	signatureProof, err := s.proverCtx.SignatureProofWithReceiver(commitment, stateProofs)
	s.NoError(err)
	s.Len(signatureProof.UserStates, 3)
	s.Len(signatureProof.SenderPublicKeys, 3)
	s.Len(signatureProof.ReceiverPublicKeys, 3)

	for i := range signatureProof.UserStates {
		s.Equal(stateProofs[i].UserState, signatureProof.UserStates[i].UserState)
		s.Equal(senderPublicKeys[i], *signatureProof.SenderPublicKeys[i].PublicKey)
		s.Equal(receiverPublicKeys[i], signatureProof.ReceiverPublicKeys[i].PublicKeyHash)
	}
}

func (s *SignatureProofsTestSuite) TestPublicKeyProof() {
	account := &models.AccountLeaf{
		PubKeyID:  1,
		PublicKey: models.PublicKey{1, 2, 3},
	}
	err := s.storage.AccountTree.SetSingle(account)
	s.NoError(err)

	publicKeyProof, err := s.proverCtx.publicKeyProof(account.PubKeyID)
	s.NoError(err)
	s.Equal(account.PublicKey, *publicKeyProof.PublicKey)
	s.Len(publicKeyProof.Witness, 32)
}

func (s *SignatureProofsTestSuite) TestReceiverPublicKeyProof() {
	account := &models.AccountLeaf{
		PubKeyID:  1,
		PublicKey: models.PublicKey{1, 2, 3},
	}
	err := s.storage.AccountTree.SetSingle(account)
	s.NoError(err)

	publicKeyProof, err := s.proverCtx.receiverPublicKeyProof(account.PubKeyID)
	s.NoError(err)
	s.Equal(crypto.Keccak256Hash(account.PublicKey.Bytes()), publicKeyProof.PublicKeyHash)
	s.Len(publicKeyProof.Witness, 32)
}

func (s *SignatureProofsTestSuite) TestReceiverPublicKeyProof_NonexistentAccount() {
	publicKeyProof, err := s.proverCtx.receiverPublicKeyProof(1)
	s.NoError(err)
	s.Equal(merkletree.GetZeroHash(0), publicKeyProof.PublicKeyHash)
	s.Len(publicKeyProof.Witness, 32)
}

func (s *SignatureProofsTestSuite) addAccounts(accounts ...models.AccountLeaf) {
	for i := range accounts {
		err := s.storage.AccountTree.SetSingle(&accounts[i])
		s.NoError(err)
	}
}

func TestSignatureProofs(t *testing.T) {
	suite.Run(t, new(SignatureProofsTestSuite))
}
