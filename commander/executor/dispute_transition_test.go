package executor

import (
	"context"
	"testing"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type DisputeTransitionTestSuite struct {
	*require.Assertions
	suite.Suite
	storage             *st.Storage
	teardown            func() error
	client              *eth.TestClient
	cfg                 *config.RollupConfig
	transactionExecutor *TransactionExecutor
	decodedCommitments  []encoder.DecodedCommitment
	decodedBatch        eth.DecodedBatch
}

func (s *DisputeTransitionTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
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
	s.cfg = &config.RollupConfig{
		MinCommitmentsPerBatch: 1,
		MaxCommitmentsPerBatch: 32,
		TxsPerCommitment:       2,
		DevMode:                false,
	}
}

func (s *DisputeTransitionTestSuite) SetupTest() {
	testStorage, err := st.NewTestStorageWithBadger()
	s.NoError(err)
	s.storage = testStorage.Storage
	s.teardown = testStorage.Teardown

	s.client, err = eth.NewTestClient()
	s.NoError(err)

	s.transactionExecutor = NewTestTransactionExecutor(s.storage, s.client.Client, s.cfg, context.Background())
}

func (s *DisputeTransitionTestSuite) TearDownTest() {
	s.client.Close()
	err := s.teardown()
	s.NoError(err)
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

	proof, err := s.transactionExecutor.previousCommitmentInclusionProof(&s.decodedBatch, 0)
	s.NoError(err)
	s.Equal(expected, *proof)
}

func (s *DisputeTransitionTestSuite) TestPreviousCommitmentInclusionProof_PreviousBatch() {
	_, err := st.NewStateTree(s.storage).Set(11, &models.UserState{
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
			Type:              txtype.Transfer,
			Transactions:      utils.RandomBytes(12),
			FeeReceiver:       11,
			CombinedSignature: models.MakeRandomSignature(),
			PostStateRoot:     utils.RandomHash(),
			IncludedInBatch:   &batch.ID,
		},
		{
			Type:              txtype.Transfer,
			Transactions:      utils.RandomBytes(12),
			FeeReceiver:       11,
			CombinedSignature: models.MakeRandomSignature(),
			PostStateRoot:     utils.RandomHash(),
			IncludedInBatch:   &batch.ID,
		},
	}
	for i := range commitments {
		_, err = s.storage.AddCommitment(&commitments[i])
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

	proof, err := s.transactionExecutor.previousCommitmentInclusionProof(&s.decodedBatch, -1)
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

func (s *DisputeTransitionTestSuite) TestDisputeTransition_RemovesInvalidBatch() {
	s.setUserStates()

	commitmentTxs := [][]models.Transfer{
		{
			s.createTransfer(0, 2, 0, 100),
			s.createTransfer(1, 0, 0, 100),
		},
		{
			s.createTransfer(2, 0, 0, 50),
			s.createTransfer(2, 0, 1, 500),
		},
	}

	proofs := s.getStateMerkleProofs(commitmentTxs)

	s.beginExecutorTransaction()
	s.createAndSubmitInvalidTransferBatch(commitmentTxs, commitmentTxs[1][1].Hash)

	remoteBatches, err := s.client.GetBatches(&bind.FilterOpts{})
	s.NoError(err)
	s.Len(remoteBatches, 1)

	err = s.transactionExecutor.disputeTransition(&remoteBatches[0], 1, proofs)
	s.NoError(err)

	_, err = s.client.GetBatch(&remoteBatches[0].ID)
	s.Error(err)
	s.Equal("execution reverted: Batch id greater than total number of batches, invalid batch id", err.Error())

	batch, err := s.storage.GetBatch(remoteBatches[0].ID)
	s.Nil(batch)
	s.True(st.IsNotFoundError(err))
}

func (s *DisputeTransitionTestSuite) TestDisputeTransition_FirstCommitment() {
	s.setUserStates()

	commitmentTxs := [][]models.Transfer{
		{
			s.createTransfer(0, 2, 0, 100),
			s.createTransfer(1, 0, 0, 100),
		},
	}

	proofs := s.getStateMerkleProofs(commitmentTxs)

	s.beginExecutorTransaction()
	s.createAndSubmitInvalidTransferBatch(commitmentTxs, commitmentTxs[1][1].Hash)

	remoteBatches, err := s.client.GetBatches(&bind.FilterOpts{})
	s.NoError(err)
	s.Len(remoteBatches, 1)

	err = s.transactionExecutor.disputeTransition(&remoteBatches[0], 1, proofs)
	s.NoError(err)

	_, err = s.client.GetBatch(&remoteBatches[0].ID)
	s.Error(err)
	s.Equal("execution reverted: Batch id greater than total number of batches, invalid batch id", err.Error())

	batch, err := s.storage.GetBatch(remoteBatches[0].ID)
	s.Nil(batch)
	s.True(st.IsNotFoundError(err))
}

func (s *DisputeTransitionTestSuite) beginExecutorTransaction() {
	var err error
	s.transactionExecutor, err = NewTransactionExecutor(s.storage, s.client.Client, s.cfg, context.Background())
	s.NoError(err)
}

func (s *DisputeTransitionTestSuite) getStateMerkleProofs(txs [][]models.Transfer) []models.StateMerkleProof {
	s.beginExecutorTransaction()

	feeReceiver := &FeeReceiver{
		StateID: 0,
		TokenID: models.MakeUint256(0),
	}

	var disputableTransferError *DisputableTransferError
	for i := range txs {
		_, err := s.transactionExecutor.ApplyTransfersForSync(txs[i], feeReceiver)
		if err != nil {
			s.ErrorAs(err, &disputableTransferError)
			s.Len(disputableTransferError.Proofs, len(txs[i])*2)
		}
	}
	s.NotNil(disputableTransferError)

	s.transactionExecutor.Rollback(nil)
	return disputableTransferError.Proofs
}

func (s *DisputeTransitionTestSuite) createAndSubmitInvalidTransferBatch(txs [][]models.Transfer, invalidTxHash common.Hash) *models.Batch {
	for i := range txs {
		err := s.storage.BatchAddTransfer(txs[i])
		s.NoError(err)
	}

	pendingBatch, err := s.transactionExecutor.NewPendingBatch(txtype.Transfer)
	s.NoError(err)

	commitments := s.createInvalidTransferCommitments(txs, invalidTxHash)
	s.Len(commitments, len(txs))

	err = s.transactionExecutor.SubmitBatch(pendingBatch, commitments)
	s.NoError(err)

	s.client.Commit()
	return pendingBatch
}

func (s *DisputeTransitionTestSuite) createInvalidTransferCommitments(
	commitmentTxs [][]models.Transfer,
	invalidTxHash common.Hash,
) []models.Commitment {
	commitments := make([]models.Commitment, 0, len(commitmentTxs))
	for i := range commitmentTxs {
		txs := commitmentTxs[i]
		combinedFee := models.MakeUint256(0)
		for j := range txs {
			senderState, receiverState, err := s.transactionExecutor.getParticipantsStates(&transfer)
			s.NoError(err)

			if txs[j].Hash != invalidTxHash {
				transferError, appError := s.transactionExecutor.ApplyTransfer(&txs[j], models.MakeUint256(0))
				s.NoError(transferError)
				s.NoError(appError)
			} else {
				s.calculateStateAfterInvalidTransfer(senderState, receiverState, &txs[j])
			}
			combinedFee = *combinedFee.Add(&txs[j].Fee)
		}
		if combinedFee.CmpN(0) > 0 {
			err := s.transactionExecutor.ApplyFee(0, combinedFee)
			s.NoError(err)
		}
		commitment, err := s.transactionExecutor.buildTransferCommitment(txs, 0, testDomain)
		s.NoError(err)
		commitments = append(commitments, *commitment)
	}

	return commitments
}

func (s *DisputeTransitionTestSuite) calculateStateAfterInvalidTransfer(
	senderState, receiverState *models.StateLeaf,
	invalidTransfer *models.Transfer,
) {
	senderState.Nonce = *senderState.Nonce.AddN(1)
	receiverState.Balance = *receiverState.Balance.Add(&invalidTransfer.Amount)
	_, err := s.transactionExecutor.stateTree.Set(invalidTransfer.GetFromStateID(), &senderState.UserState)
	s.NoError(err)
	_, err = s.transactionExecutor.stateTree.Set(invalidTransfer.ToStateID, &receiverState.UserState)
	s.NoError(err)
}

func (s *DisputeTransitionTestSuite) setUserStates() {
	userStates := []models.UserState{
		*s.createUserState(0, 300, 0),
		*s.createUserState(1, 200, 0),
		*s.createUserState(2, 100, 0),
	}
	registrations, unsubscribe, err := s.client.WatchRegistrations(&bind.WatchOpts{})
	s.NoError(err)
	defer unsubscribe()

	for i := range userStates {
		wallet, err := bls.NewRandomWallet(*testDomain)
		s.NoError(err)
		pubKeyID, err := s.client.RegisterAccount(wallet.PublicKey(), registrations)
		s.NoError(err)
		s.Equal(userStates[i].PubKeyID, *pubKeyID)

		_, err = s.transactionExecutor.stateTree.Set(uint32(i), &userStates[i])
		s.NoError(err)
	}
}

func (s *DisputeTransitionTestSuite) createUserState(pubKeyID uint32, balance, nonce uint64) *models.UserState {
	return &models.UserState{
		PubKeyID: pubKeyID,
		TokenID:  models.MakeUint256(0),
		Balance:  models.MakeUint256(balance),
		Nonce:    models.MakeUint256(nonce),
	}
}

func (s *DisputeTransitionTestSuite) createTransfer(from, to uint32, nonce, amount uint64) models.Transfer {
	return models.Transfer{
		TransactionBase: models.TransactionBase{
			Hash:        utils.RandomHash(),
			TxType:      txtype.Transfer,
			FromStateID: from,
			Amount:      models.MakeUint256(amount),
			Fee:         models.MakeUint256(10),
			Nonce:       models.MakeUint256(nonce),
		},
		ToStateID: to,
	}
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(DisputeTransitionTestSuite))
}
