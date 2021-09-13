package executor

import (
	"context"
	"testing"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/Worldcoin/hubble-commander/utils/merkletree"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type DisputeTransitionTestSuite struct {
	TestSuiteWithDisputeContext
	decodedCommitments []encoder.DecodedCommitment
	decodedBatch       eth.DecodedBatch
}

func (s *DisputeTransitionTestSuite) SetupSuite() {
	s.TestSuiteWithDisputeContext.SetupSuite()

	s.decodedCommitments = []encoder.DecodedCommitment{
		{
			StateRoot:         utils.RandomHash(),
			CombinedSignature: models.MakeRandomSignature(),
			FeeReceiver:       10,
			Transactions:      utils.RandomBytes(12),
		},
		{
			StateRoot:         utils.RandomHash(),
			CombinedSignature: models.MakeRandomSignature(),
			FeeReceiver:       10,
			Transactions:      utils.RandomBytes(12),
		},
	}
	s.decodedBatch = eth.DecodedBatch{
		Batch: models.Batch{
			ID:                models.MakeUint256(2),
			Type:              txtype.Transfer,
			TransactionHash:   utils.RandomHash(),
			Hash:              utils.NewRandomHash(),
			FinalisationBlock: ref.Uint32(10),
			AccountTreeRoot:   utils.NewRandomHash(),
		},
		Commitments: s.decodedCommitments,
	}
}

func (s *DisputeTransitionTestSuite) TestPreviousCommitmentInclusionProof_CurrentBatch() {
	expected := models.CommitmentInclusionProof{
		StateRoot: s.decodedCommitments[0].StateRoot,
		BodyRoot:  s.decodedCommitments[0].BodyHash(*s.decodedBatch.AccountTreeRoot),
		Path: &models.MerklePath{
			Path:  0,
			Depth: 2,
		},
		Witness: []common.Hash{s.decodedCommitments[1].LeafHash(*s.decodedBatch.AccountTreeRoot)},
	}

	proof, err := s.disputeCtx.previousCommitmentInclusionProof(&s.decodedBatch, 0)
	s.NoError(err)
	s.Equal(expected, *proof)
}

func (s *DisputeTransitionTestSuite) TestPreviousCommitmentInclusionProof_PreviousBatch() {
	_, err := s.storage.StateTree.Set(11, &models.UserState{
		PubKeyID: 1,
		TokenID:  models.MakeUint256(1),
		Balance:  models.MakeUint256(100),
		Nonce:    models.MakeUint256(0),
	})
	s.NoError(err)

	batch := models.Batch{
		ID:                models.MakeUint256(1),
		Type:              txtype.Transfer,
		TransactionHash:   utils.RandomHash(),
		Hash:              utils.NewRandomHash(),
		FinalisationBlock: ref.Uint32(10),
		AccountTreeRoot:   utils.NewRandomHash(),
	}
	err = s.storage.AddBatch(&batch)
	s.NoError(err)

	commitments := []models.Commitment{
		{
			ID: models.CommitmentID{
				BatchID:      batch.ID,
				IndexInBatch: 0,
			},
			Type:              txtype.Transfer,
			Transactions:      utils.RandomBytes(12),
			FeeReceiver:       11,
			CombinedSignature: models.MakeRandomSignature(),
			PostStateRoot:     utils.RandomHash(),
		},
		{
			ID: models.CommitmentID{
				BatchID:      batch.ID,
				IndexInBatch: 1,
			},
			Type:              txtype.Transfer,
			Transactions:      utils.RandomBytes(12),
			FeeReceiver:       11,
			CombinedSignature: models.MakeRandomSignature(),
			PostStateRoot:     utils.RandomHash(),
		},
	}
	for i := range commitments {
		err = s.storage.AddCommitment(&commitments[i])
		s.NoError(err)
	}

	expected := models.CommitmentInclusionProof{
		StateRoot: commitments[1].PostStateRoot,
		BodyRoot:  commitments[1].BodyHash(*batch.AccountTreeRoot),
		Path: &models.MerklePath{
			Path:  1,
			Depth: 2,
		},
		Witness: []common.Hash{commitments[0].LeafHash(*batch.AccountTreeRoot)},
	}

	proof, err := s.disputeCtx.previousCommitmentInclusionProof(&s.decodedBatch, -1)
	s.NoError(err)
	s.Equal(expected, *proof)
}

func (s *DisputeTransitionTestSuite) TestGenesisBatchCommitmentInclusionProof() {
	genesisBatch := s.addGenesisBatch()
	zeroHash := merkletree.GetZeroHash(0)

	expected := models.CommitmentInclusionProof{
		StateRoot: *genesisBatch.PrevStateRoot,
		BodyRoot:  zeroHash,
		Path: &models.MerklePath{
			Path:  0,
			Depth: 2,
		},
		Witness: []common.Hash{zeroHash},
	}

	firstBatch := s.decodedBatch
	firstBatch.ID = models.MakeUint256(1)
	proof, err := s.disputeCtx.previousCommitmentInclusionProof(&firstBatch, -1)
	s.NoError(err)
	s.Equal(expected, *proof)
}

func (s *DisputeTransitionTestSuite) TestTargetCommitmentInclusionProof() {
	expected := models.TransferCommitmentInclusionProof{
		StateRoot: s.decodedCommitments[1].StateRoot,
		Body: &models.TransferBody{
			AccountRoot:  *s.decodedBatch.AccountTreeRoot,
			Signature:    s.decodedCommitments[1].CombinedSignature,
			FeeReceiver:  s.decodedCommitments[1].FeeReceiver,
			Transactions: s.decodedCommitments[1].Transactions,
		},
		Path: &models.MerklePath{
			Path:  1,
			Depth: 2,
		},
		Witness: []common.Hash{s.decodedCommitments[0].LeafHash(*s.decodedBatch.AccountTreeRoot)},
	}

	proof, err := targetCommitmentInclusionProof(&s.decodedBatch, 1)
	s.NoError(err)
	s.Equal(expected, *proof)
}

func (s *DisputeTransitionTestSuite) TestDisputeTransition_Transfer_RemovesInvalidBatch() {
	setUserStates(s.Assertions, s.executionCtx, testDomain)

	commitmentTxs := [][]models.Transfer{
		{
			testutils.MakeTransfer(0, 2, 0, 100),
			testutils.MakeTransfer(1, 0, 0, 100),
		},
		{
			testutils.MakeTransfer(2, 0, 0, 50),
			testutils.MakeTransfer(2, 0, 1, 500),
		},
	}

	proofs := s.getTransferStateMerkleProofs(commitmentTxs)

	s.beginTransaction()
	defer s.commitTransaction()
	s.submitInvalidTransferBatch(commitmentTxs, commitmentTxs[1][1].Hash)

	remoteBatches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 1)

	err = s.disputeCtx.DisputeTransition(&remoteBatches[0], 1, proofs)
	s.NoError(err)

	s.checkBatchAfterDispute(remoteBatches[0].ID)
}

func (s *DisputeTransitionTestSuite) TestDisputeTransition_Transfer_FirstCommitment() {
	setUserStates(s.Assertions, s.executionCtx, testDomain)

	commitmentTxs := [][]models.Transfer{
		{
			testutils.MakeTransfer(0, 2, 0, 500),
		},
	}

	transfer := testutils.MakeTransfer(0, 2, 0, 50)
	submitTransferBatch(s.Assertions, s.client, s.executionCtx, &transfer)

	proofs := s.getTransferStateMerkleProofs(commitmentTxs)

	s.beginTransaction()
	defer s.commitTransaction()
	s.submitInvalidTransferBatch(commitmentTxs, commitmentTxs[0][0].Hash)

	remoteBatches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 2)

	err = s.executionCtx.storage.MarkBatchAsSubmitted(&remoteBatches[0].Batch)
	s.NoError(err)

	err = s.disputeCtx.DisputeTransition(&remoteBatches[1], 0, proofs)
	s.NoError(err)

	s.checkBatchAfterDispute(remoteBatches[1].ID)
}

func (s *DisputeTransitionTestSuite) TestDisputeTransition_Transfer_ValidBatch() {
	setUserStates(s.Assertions, s.executionCtx, testDomain)

	transfers := []models.Transfer{
		testutils.MakeTransfer(0, 2, 0, 50),
		testutils.MakeTransfer(0, 2, 1, 100),
	}

	submitTransferBatch(s.Assertions, s.client, s.executionCtx, &transfers[0])

	proofs := s.getTransferStateMerkleProofs([][]models.Transfer{{transfers[1]}})

	s.beginTransaction()
	defer s.commitTransaction()
	submitTransferBatch(s.Assertions, s.client, s.executionCtx, &transfers[1])

	remoteBatches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 2)

	err = s.executionCtx.storage.MarkBatchAsSubmitted(&remoteBatches[0].Batch)
	s.NoError(err)

	err = s.disputeCtx.DisputeTransition(&remoteBatches[1], 0, proofs)
	s.NoError(err)
	_, err = s.client.GetBatch(&remoteBatches[1].ID)
	s.NoError(err)
}

func (s *DisputeTransitionTestSuite) TestDisputeTransition_Create2Transfer_RemovesInvalidBatch() {
	wallets := setUserStates(s.Assertions, s.executionCtx, testDomain)

	commitmentTxs := [][]models.Create2Transfer{
		{
			testutils.MakeCreate2Transfer(0, ref.Uint32(3), 0, 100, wallets[2].PublicKey()),
			testutils.MakeCreate2Transfer(1, ref.Uint32(4), 0, 100, wallets[0].PublicKey()),
		},
		{
			testutils.MakeCreate2Transfer(2, ref.Uint32(5), 0, 50, wallets[0].PublicKey()),
			testutils.MakeCreate2Transfer(2, ref.Uint32(6), 1, 500, wallets[0].PublicKey()),
		},
	}

	pubKeyIDs := [][]uint32{{3, 4}, {5, 6}}
	proofs := s.getC2TStateMerkleProofs(commitmentTxs, pubKeyIDs)

	s.beginTransaction()
	defer s.commitTransaction()
	s.submitInvalidC2TBatch(commitmentTxs, pubKeyIDs, commitmentTxs[1][1].Hash)

	remoteBatches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 1)

	err = s.disputeCtx.DisputeTransition(&remoteBatches[0], 1, proofs)
	s.NoError(err)

	s.checkBatchAfterDispute(remoteBatches[0].ID)
}

func (s *DisputeTransitionTestSuite) TestDisputeTransition_Create2Transfer_FirstCommitment() {
	wallets := setUserStates(s.Assertions, s.executionCtx, testDomain)

	commitmentTxs := [][]models.Create2Transfer{
		{
			testutils.MakeCreate2Transfer(0, ref.Uint32(4), 0, 500, wallets[1].PublicKey()),
		},
	}
	pubKeyIDs := [][]uint32{{4}}

	transfer := testutils.MakeCreate2Transfer(0, nil, 0, 50, wallets[1].PublicKey())
	submitC2TBatch(s.Assertions, s.client, s.executionCtx, &transfer)

	registrations, unsubscribe, err := s.client.WatchRegistrations(&bind.WatchOpts{})
	s.NoError(err)
	defer unsubscribe()

	pubKeyID, err := s.client.RegisterAccount(wallets[1].PublicKey(), registrations)
	s.NoError(err)
	s.EqualValues(4, *pubKeyID)

	proofs := s.getC2TStateMerkleProofs(commitmentTxs, pubKeyIDs)

	s.beginTransaction()
	defer s.commitTransaction()
	s.submitInvalidC2TBatch(commitmentTxs, pubKeyIDs, commitmentTxs[0][0].Hash)

	remoteBatches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 2)

	err = s.executionCtx.storage.MarkBatchAsSubmitted(&remoteBatches[0].Batch)
	s.NoError(err)

	err = s.disputeCtx.DisputeTransition(&remoteBatches[1], 0, proofs)
	s.NoError(err)

	s.checkBatchAfterDispute(remoteBatches[1].ID)
}

func (s *DisputeTransitionTestSuite) TestDisputeTransition_Create2Transfer_ValidBatch() {
	wallets := setUserStates(s.Assertions, s.executionCtx, testDomain)

	transfers := []models.Create2Transfer{
		testutils.MakeCreate2Transfer(0, nil, 0, 50, wallets[1].PublicKey()),
		testutils.MakeCreate2Transfer(0, ref.Uint32(4), 1, 100, wallets[1].PublicKey()),
	}
	pubKeyIDs := [][]uint32{{4}}

	submitC2TBatch(s.Assertions, s.client, s.executionCtx, &transfers[0])

	proofs := s.getC2TStateMerkleProofs([][]models.Create2Transfer{{transfers[1]}}, pubKeyIDs)

	s.beginTransaction()
	defer s.commitTransaction()

	transfers[1].ToStateID = nil
	submitC2TBatch(s.Assertions, s.client, s.executionCtx, &transfers[1])

	remoteBatches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 2)

	err = s.executionCtx.storage.MarkBatchAsSubmitted(&remoteBatches[0].Batch)
	s.NoError(err)

	err = s.disputeCtx.DisputeTransition(&remoteBatches[1], 0, proofs)
	s.NoError(err)
	_, err = s.client.GetBatch(&remoteBatches[1].ID)
	s.NoError(err)
}

func (s *DisputeTransitionTestSuite) checkBatchAfterDispute(batchID models.Uint256) {
	checkRemoteBatchAfterDispute(s.Assertions, s.client, &batchID)

	batch, err := s.storage.GetBatch(batchID)
	s.Nil(batch)
	s.True(st.IsNotFoundError(err))
}

func checkRemoteBatchAfterDispute(s *require.Assertions, client *eth.TestClient, batchID *models.Uint256) {
	_, err := client.GetBatch(batchID)
	if err == nil {
		err = client.KeepRollingBack()
		s.NoError(err)
		_, err = client.GetBatch(batchID)
	}
	s.Error(err)
	s.Equal(eth.MsgInvalidBatchID, err.Error())
}

func (s *DisputeTransitionTestSuite) beginTransaction() {
	var err error
	s.executionCtx, err = NewExecutionContext(s.storage.Storage, s.client.Client, s.cfg, context.Background())
	s.NoError(err)
	s.disputeCtx = NewDisputeContext(s.executionCtx.storage, s.executionCtx.client)
}

func (s *DisputeTransitionTestSuite) getTransferStateMerkleProofs(txs [][]models.Transfer) []models.StateMerkleProof {
	feeReceiverStateID := uint32(0)

	s.beginTransaction()
	defer s.executionCtx.Rollback(nil)

	var stateProofs []models.StateMerkleProof
	var err error
	for i := range txs {
		_, stateProofs, err = s.executionCtx.ApplyTransfersForSync(txs[i], feeReceiverStateID)
		if err != nil {
			var disputableErr *DisputableError
			s.ErrorAs(err, &disputableErr)
			s.Equal(Transition, disputableErr.Type)
			s.Len(disputableErr.Proofs, len(txs[i])*2)
			return disputableErr.Proofs
		}
	}

	return stateProofs
}

func (s *DisputeTransitionTestSuite) getC2TStateMerkleProofs(
	txs [][]models.Create2Transfer,
	pubKeyIDs [][]uint32,
) []models.StateMerkleProof {
	feeReceiverStateID := uint32(0)

	s.beginTransaction()
	defer s.executionCtx.Rollback(nil)

	var stateProofs []models.StateMerkleProof
	var err error
	for i := range txs {
		_, stateProofs, err = s.executionCtx.ApplyCreate2TransfersForSync(txs[i], pubKeyIDs[i], feeReceiverStateID)
		if err != nil {
			var disputableErr *DisputableError
			s.ErrorAs(err, &disputableErr)
			s.Equal(Transition, disputableErr.Type)
			s.Len(disputableErr.Proofs, len(txs[i])*2)
			return disputableErr.Proofs
		}
	}

	return stateProofs
}

func (s *DisputeTransitionTestSuite) submitInvalidTransferBatch(txs [][]models.Transfer, invalidTxHash common.Hash) *models.Batch {
	for i := range txs {
		err := s.executionCtx.storage.BatchAddTransfer(txs[i])
		s.NoError(err)
	}

	pendingBatch, err := s.executionCtx.NewPendingBatch(txtype.Transfer)
	s.NoError(err)

	commitments := s.createInvalidTransferCommitments(txs, invalidTxHash)
	s.Len(commitments, len(txs))

	err = s.executionCtx.SubmitBatch(pendingBatch, commitments)
	s.NoError(err)

	s.client.Commit()
	return pendingBatch
}

func (s *DisputeTransitionTestSuite) submitInvalidC2TBatch(
	txs [][]models.Create2Transfer,
	pubKeyIDs [][]uint32,
	invalidTxHash common.Hash,
) *models.Batch {
	for i := range txs {
		stateIDs := s.resetCreate2TransfersToStateID(txs[i])
		err := s.executionCtx.storage.BatchAddCreate2Transfer(txs[i])
		s.NoError(err)
		s.setCreate2TransfersToStateID(txs[i], stateIDs)
	}

	pendingBatch, err := s.executionCtx.NewPendingBatch(txtype.Create2Transfer)
	s.NoError(err)

	commitments := s.createInvalidC2TCommitments(txs, pubKeyIDs, invalidTxHash)
	s.Len(commitments, len(txs))

	err = s.executionCtx.SubmitBatch(pendingBatch, commitments)
	s.NoError(err)

	s.client.Commit()
	return pendingBatch
}

func (s *DisputeTransitionTestSuite) resetCreate2TransfersToStateID(txs []models.Create2Transfer) []*uint32 {
	stateIDs := make([]*uint32, 0, len(txs))
	for i := range txs {
		stateIDs = append(stateIDs, txs[i].ToStateID)
		txs[i].ToStateID = nil
	}
	return stateIDs
}

func (s *DisputeTransitionTestSuite) setCreate2TransfersToStateID(txs []models.Create2Transfer, toStateIDs []*uint32) {
	for i := range txs {
		txs[i].ToStateID = toStateIDs[i]
	}
}

func (s *DisputeTransitionTestSuite) createInvalidC2TCommitments(
	commitmentTxs [][]models.Create2Transfer,
	pubKeyIDs [][]uint32,
	invalidTxHash common.Hash,
) []models.Commitment {
	commitmentID, err := s.executionCtx.createCommitmentID()
	s.NoError(err)

	commitments := make([]models.Commitment, 0, len(commitmentTxs))
	for i := range commitmentTxs {
		commitmentID.IndexInBatch = uint8(i)
		txs := commitmentTxs[i]
		combinedFee := models.MakeUint256(0)
		for j := range txs {
			receiverLeaf := newUserLeaf(*txs[j].ToStateID, pubKeyIDs[i][j], models.MakeUint256(0))
			combinedFee = s.applyTransfer(&txs[j], invalidTxHash, combinedFee, receiverLeaf)
		}
		if combinedFee.CmpN(0) > 0 {
			_, err := s.executionCtx.ApplyFee(0, combinedFee)
			s.NoError(err)
		}

		commitment, err := s.executionCtx.buildC2TCommitment(txs, pubKeyIDs[i], commitmentID, 0, testDomain)
		s.NoError(err)
		commitments = append(commitments, *commitment)
	}

	return commitments
}

func (s *DisputeTransitionTestSuite) createInvalidTransferCommitments(
	commitmentTxs [][]models.Transfer,
	invalidTxHash common.Hash,
) []models.Commitment {
	commitmentID, err := s.executionCtx.createCommitmentID()
	s.NoError(err)

	commitments := make([]models.Commitment, 0, len(commitmentTxs))
	for i := range commitmentTxs {
		commitmentID.IndexInBatch = uint8(i)
		txs := commitmentTxs[i]
		combinedFee := models.MakeUint256(0)
		for j := range txs {
			receiverLeaf, err := s.executionCtx.storage.StateTree.Leaf(txs[j].ToStateID)
			s.NoError(err)
			combinedFee = s.applyTransfer(&txs[j], invalidTxHash, combinedFee, receiverLeaf)
		}
		if combinedFee.CmpN(0) > 0 {
			_, err := s.executionCtx.ApplyFee(0, combinedFee)
			s.NoError(err)
		}

		commitment, err := s.executionCtx.buildTransferCommitment(txs, commitmentID, 0, testDomain)
		s.NoError(err)
		commitments = append(commitments, *commitment)
	}

	return commitments
}

func (s *DisputeTransitionTestSuite) applyTransfer(
	tx models.GenericTransaction,
	invalidTxHash common.Hash,
	combinedFee models.Uint256,
	receiverLeaf *models.StateLeaf,
) models.Uint256 {
	if tx.GetBase().Hash != invalidTxHash {
		transferError, appError := s.executionCtx.ApplyTransfer(tx, receiverLeaf, models.MakeUint256(0))
		s.NoError(transferError)
		s.NoError(appError)
	} else {
		senderLeaf, err := s.executionCtx.storage.StateTree.Leaf(tx.GetFromStateID())
		s.NoError(err)
		s.calculateStateAfterInvalidTransfer(senderLeaf, receiverLeaf, tx)
	}
	fee := tx.GetFee()
	return *combinedFee.Add(&fee)
}

func (s *DisputeTransitionTestSuite) calculateStateAfterInvalidTransfer(
	senderState, receiverState *models.StateLeaf,
	invalidTransfer models.GenericTransaction,
) {
	senderState.Nonce = *senderState.Nonce.AddN(1)
	amount := invalidTransfer.GetAmount()
	receiverState.Balance = *receiverState.Balance.Add(&amount)
	_, err := s.executionCtx.storage.StateTree.Set(invalidTransfer.GetFromStateID(), &senderState.UserState)
	s.NoError(err)
	_, err = s.executionCtx.storage.StateTree.Set(*invalidTransfer.GetToStateID(), &receiverState.UserState)
	s.NoError(err)
}

func (s *DisputeTransitionTestSuite) addGenesisBatch() *models.Batch {
	root, err := s.storage.StateTree.Root()
	s.NoError(err)

	batch, err := s.client.GetBatch(models.NewUint256(0))
	s.NoError(err)
	batch.PrevStateRoot = root

	err = s.storage.AddBatch(batch)
	s.NoError(err)

	return batch
}

func (s *DisputeTransitionTestSuite) commitTransaction() {
	err := s.executionCtx.Commit()
	s.NoError(err)
}

func setUserStates(s *require.Assertions, executionCtx *ExecutionContext, domain *bls.Domain) []bls.Wallet {
	userStates := []models.UserState{
		*createUserState(0, 300, 0),
		*createUserState(1, 200, 0),
		*createUserState(2, 100, 0),
	}
	registrations, unsubscribe, err := executionCtx.client.WatchRegistrations(&bind.WatchOpts{})
	s.NoError(err)
	defer unsubscribe()

	wallets := generateWallets(s, domain, len(userStates))
	for i := range userStates {
		pubKeyID, err := executionCtx.client.RegisterAccount(wallets[i].PublicKey(), registrations)
		s.NoError(err)
		s.Equal(userStates[i].PubKeyID, *pubKeyID)

		_, err = executionCtx.storage.StateTree.Set(uint32(i), &userStates[i])
		s.NoError(err)
	}
	return wallets
}

func createUserState(pubKeyID uint32, balance, nonce uint64) *models.UserState {
	return &models.UserState{
		PubKeyID: pubKeyID,
		TokenID:  models.MakeUint256(0),
		Balance:  models.MakeUint256(balance),
		Nonce:    models.MakeUint256(nonce),
	}
}

func TestDisputeTransitionTestSuite(t *testing.T) {
	suite.Run(t, new(DisputeTransitionTestSuite))
}
