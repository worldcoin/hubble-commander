package executor

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/suite"
)

type DisputeSignatureTestSuite struct {
	TestSuiteWithDisputeContext
	domain *bls.Domain
}

func (s *DisputeSignatureTestSuite) SetupTest() {
	s.TestSuiteWithDisputeContext.SetupTest()

	var err error
	s.domain, err = s.client.GetDomain()
	s.NoError(err)
}

func (s *DisputeSignatureTestSuite) TestUserStateProof() {
	userState := createUserState(1, 300, 1)
	witness, err := s.executionCtx.storage.StateTree.Set(1, userState)
	s.NoError(err)

	stateProof, err := s.executionCtx.userStateProof(1)
	s.NoError(err)
	s.Equal(userState, stateProof.UserState)
	s.Equal(witness, stateProof.Witness)
}

func (s *DisputeSignatureTestSuite) TestPublicKeyProof() {
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

func (s *DisputeSignatureTestSuite) TestReceiverPublicKeyProof() {
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

func (s *DisputeSignatureTestSuite) TestSignatureProof() {
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

func (s *DisputeSignatureTestSuite) TestSignatureProofWithReceiver() {
	wallets := s.setUserStatesAndAddAccounts()

	receiverAccounts := []models.AccountLeaf{
		{PubKeyID: 3, PublicKey: *wallets[2].PublicKey()},
		{PubKeyID: 4, PublicKey: *wallets[2].PublicKey()},
		{PubKeyID: 5, PublicKey: *wallets[1].PublicKey()},
	}

	transfers := []models.Create2Transfer{
		testutils.MakeCreate2Transfer(1, ref.Uint32(3), 0, 50, &receiverAccounts[0].PublicKey),
		testutils.MakeCreate2Transfer(0, ref.Uint32(4), 0, 50, &receiverAccounts[1].PublicKey),
		testutils.MakeCreate2Transfer(0, ref.Uint32(5), 1, 75, &receiverAccounts[2].PublicKey),
	}
	pubKeyIDs := []uint32{3, 4, 5}

	stateProofs := make([]models.StateMerkleProof, 0, len(transfers))
	senderPublicKeys := make([]models.PublicKey, 0, len(transfers))
	receiverPublicKeys := make([]common.Hash, 0, len(transfers))
	for i := range transfers {
		stateProof, err := s.executionCtx.userStateProof(transfers[i].FromStateID)
		s.NoError(err)
		stateProofs = append(stateProofs, *stateProof)

		account, err := s.storage.AccountTree.Leaf(stateProof.UserState.PubKeyID)
		s.NoError(err)
		senderPublicKeys = append(senderPublicKeys, account.PublicKey)

		err = s.executionCtx.storage.AccountTree.SetSingle(&receiverAccounts[i])
		s.NoError(err)
		receiverPublicKeys = append(receiverPublicKeys, crypto.Keccak256Hash(transfers[i].ToPublicKey.Bytes()))
	}

	serializedTxs, err := encoder.SerializeCreate2Transfers(transfers, pubKeyIDs)
	s.NoError(err)

	commitment := &encoder.DecodedCommitment{Transactions: serializedTxs}

	signatureProof, err := s.disputeCtx.signatureProofWithReceiver(commitment, stateProofs)
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

func (s *DisputeSignatureTestSuite) TestDisputeSignature_DisputesTransferBatchWithInvalidSignature() {
	wallets := s.setUserStatesAndAddAccounts()

	transfer := testutils.MakeTransfer(1, 2, 0, 50)
	signTransfer(s.T(), &wallets[0], &transfer)

	submitTransferBatch(s.Assertions, s.client, s.executionCtx, &transfer)

	remoteBatches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 1)

	err = s.disputeSignature(&remoteBatches[0], models.TransferArray{transfer})
	s.NoError(err)

	checkRemoteBatchAfterDispute(s.Assertions, s.client, &remoteBatches[0].ID)
}

func (s *DisputeSignatureTestSuite) TestDisputeSignature_DisputesC2TBatchWithInvalidSignature() {
	wallets := s.setUserStatesAndAddAccounts()

	receiver := &models.AccountLeaf{
		PubKeyID:  3,
		PublicKey: *wallets[2].PublicKey(),
	}

	transfer := testutils.MakeCreate2Transfer(0, nil, 0, 100, &receiver.PublicKey)
	signCreate2Transfer(s.T(), &wallets[1], &transfer)

	submitC2TBatch(s.Assertions, s.client, s.executionCtx, &transfer)

	err := s.executionCtx.storage.AccountTree.SetSingle(receiver)
	s.NoError(err)

	remoteBatches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 1)

	err = s.disputeSignature(&remoteBatches[0], models.Create2TransferArray{transfer})
	s.NoError(err)

	checkRemoteBatchAfterDispute(s.Assertions, s.client, &remoteBatches[0].ID)
}

func (s *DisputeSignatureTestSuite) TestDisputeSignature_Transfer_ValidBatch() {
	wallets := s.setUserStatesAndAddAccounts()

	transfer := testutils.MakeTransfer(1, 2, 0, 50)
	signTransfer(s.T(), &wallets[1], &transfer)

	submitTransferBatch(s.Assertions, s.client, s.executionCtx, &transfer)

	remoteBatches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 1)

	err = s.disputeSignature(&remoteBatches[0], models.TransferArray{transfer})
	s.NoError(err)
	_, err = s.client.GetBatch(&remoteBatches[0].ID)
	s.NoError(err)
}

func (s *DisputeSignatureTestSuite) TestDisputeSignature_Create2Transfer_ValidBatch() {
	wallets := s.setUserStatesAndAddAccounts()

	receiver := &models.AccountLeaf{
		PubKeyID:  3,
		PublicKey: *wallets[2].PublicKey(),
	}

	transfer := testutils.MakeCreate2Transfer(0, nil, 0, 100, &receiver.PublicKey)
	signCreate2Transfer(s.T(), &wallets[0], &transfer)

	submitC2TBatch(s.Assertions, s.client, s.executionCtx, &transfer)

	err := s.executionCtx.storage.AccountTree.SetSingle(receiver)
	s.NoError(err)

	remoteBatches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 1)

	err = s.disputeSignature(&remoteBatches[0], models.Create2TransferArray{transfer})
	s.NoError(err)
	_, err = s.client.GetBatch(&remoteBatches[0].ID)
	s.NoError(err)
}

func (s *DisputeSignatureTestSuite) setUserStatesAndAddAccounts() []bls.Wallet {
	wallets := setUserStates(s.Assertions, s.executionCtx, s.domain)
	for i := range wallets {
		err := s.executionCtx.storage.AccountTree.SetSingle(&models.AccountLeaf{
			PubKeyID:  uint32(i),
			PublicKey: *wallets[i].PublicKey(),
		})
		s.NoError(err)
	}
	return wallets
}

func (s *DisputeSignatureTestSuite) disputeSignature(
	batch *eth.DecodedBatch,
	transfers models.GenericTransactionArray,
) error {
	proofs, err := s.executionCtx.stateMerkleProofs(transfers)
	s.NoError(err)

	return s.disputeCtx.DisputeSignature(batch, 0, proofs)
}

func TestDisputeSignatureTestSuite(t *testing.T) {
	suite.Run(t, new(DisputeSignatureTestSuite))
}
