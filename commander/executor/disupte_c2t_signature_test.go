package executor

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/suite"
)

type DisputeC2TSignatureTestSuite struct {
	DisputeSignatureTestSuite
}

func (s *DisputeC2TSignatureTestSuite) SetupTest() {
	s.TestSuiteWithDisputeContext.SetupTest(txtype.Create2Transfer)
	s.DisputeSignatureTestSuite.setupTest()
}

func (s *DisputeC2TSignatureTestSuite) TestSignatureProofWithReceiver() {
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

func (s *DisputeC2TSignatureTestSuite) TestDisputeSignature_DisputesBatchWithInvalidSignature() {
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

func (s *DisputeC2TSignatureTestSuite) TestDisputeSignature_ValidBatch() {
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

func TestDisputeC2TSignatureTestSuite(t *testing.T) {
	suite.Run(t, new(DisputeC2TSignatureTestSuite))
}
