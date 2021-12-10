package disputer

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/Worldcoin/hubble-commander/utils/merkletree"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/suite"
)

type DisputeC2TSignatureTestSuite struct {
	disputeSignatureTestSuite
}

func (s *DisputeC2TSignatureTestSuite) SetupTest() {
	s.testSuiteWithContexts.SetupTest(batchtype.Create2Transfer)
	s.disputeSignatureTestSuite.setupTest()
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
		stateProof, err := s.syncCtx.UserStateProof(transfers[i].FromStateID)
		s.NoError(err)
		stateProofs = append(stateProofs, *stateProof)

		account, err := s.storage.AccountTree.Leaf(stateProof.UserState.PubKeyID)
		s.NoError(err)
		senderPublicKeys = append(senderPublicKeys, account.PublicKey)

		err = s.storage.AccountTree.SetSingle(&receiverAccounts[i])
		s.NoError(err)
		receiverPublicKeys = append(receiverPublicKeys, crypto.Keccak256Hash(transfers[i].ToPublicKey.Bytes()))
	}

	serializedTxs, err := encoder.SerializeCreate2Transfers(transfers, pubKeyIDs)
	s.NoError(err)

	commitment := &encoder.DecodedCommitment{Transactions: serializedTxs}

	signatureProof, err := s.disputeCtx.proverCtx.SignatureProofWithReceiver(commitment, stateProofs)
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

func (s *DisputeC2TSignatureTestSuite) TestDisputeSignature_DisputesBatchToNonexistentReceiver() {
	wallets := s.setUserStatesAndAddAccounts()

	transfer := testutils.MakeCreate2Transfer(0, nil, 0, 100, wallets[2].PublicKey())
	s.signC2TWithPublicKeyHash(&wallets[0], &transfer, merkletree.GetZeroHash(0))

	s.submitBatchWithNonexistentAccount(&transfer)

	remoteBatches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 1)

	err = s.disputeSignature(remoteBatches[0].ToDecodedTxBatch(), models.Create2TransferArray{transfer})
	s.NoError(err)

	checkRemoteBatchAfterDispute(s.Assertions, s.client, &remoteBatches[0].GetBase().ID)
}

func (s *DisputeC2TSignatureTestSuite) TestDisputeSignature_ValidBatch() {
	wallets := s.setUserStatesAndAddAccounts()

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

func (s *DisputeC2TSignatureTestSuite) signC2TWithPublicKeyHash(
	wallet *bls.Wallet,
	transfer *models.Create2Transfer,
	publicKeyHash common.Hash,
) {
	encodedTransfer, err := encoder.EncodeCreate2TransferForSigning(transfer)
	s.NoError(err)
	copy(encodedTransfer[32*2:32*3], publicKeyHash.Bytes())

	signature, err := wallet.Sign(encodedTransfer)
	s.NoError(err)
	transfer.Signature = *signature.ModelsSignature()
}

func (s *DisputeC2TSignatureTestSuite) submitBatchWithNonexistentAccount(tx models.GenericTransaction) {
	pendingBatch, batchData := s.createBatch(tx)
	commitments := batchData.Commitments()

	c2t, _, err := encoder.DecodeCreate2TransferFromCommitment(commitments[0].Transactions)
	s.NoError(err)

	receiverPubKeyID := uint32(1000)

	_, err = s.disputeCtx.storage.StateTree.Set(*c2t.ToStateID, &models.UserState{
		PubKeyID: receiverPubKeyID,
		TokenID:  models.MakeUint256(0),
		Balance:  tx.GetAmount(),
		Nonce:    models.MakeUint256(0),
	})
	s.NoError(err)

	root, err := s.disputeCtx.storage.StateTree.Root()
	s.NoError(err)
	commitments[0].PostStateRoot = *root

	encodedC2T, err := encoder.EncodeCreate2TransferForCommitment(c2t, receiverPubKeyID)
	s.NoError(err)
	commitments[0].Transactions = encodedC2T

	err = s.txsCtx.SubmitBatch(pendingBatch, batchData)
	s.NoError(err)

	s.client.GetBackend().Commit()
}

func TestDisputeC2TSignatureTestSuite(t *testing.T) {
	suite.Run(t, new(DisputeC2TSignatureTestSuite))
}
