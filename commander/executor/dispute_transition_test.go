package executor

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type DisputeTransitionTestSuite struct {
	*require.Assertions
	suite.Suite
	storage             *st.Storage
	teardown            func() error
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
}

func (s *DisputeTransitionTestSuite) SetupTest() {
	testStorage, err := st.NewTestStorageWithBadger()
	s.NoError(err)
	s.storage = testStorage.Storage
	s.teardown = testStorage.Teardown

	s.transactionExecutor = NewTestTransactionExecutor(s.storage, &eth.Client{}, &config.RollupConfig{}, TransactionExecutorOpts{})
}

func (s *DisputeTransitionTestSuite) TearDownTest() {
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
		Witness: []common.Hash{s.decodedCommitments[1].BodyHash(*s.decodedBatch.AccountTreeRoot)},
	}

	proof, err := s.transactionExecutor.previousCommitmentInclusionProof(&s.decodedBatch, 0)
	s.NoError(err)
	s.Equal(expected, *proof)
}

func (s *DisputeTransitionTestSuite) TestPreviousCommitmentInclusionProof_PreviousBatch() {
	err := st.NewStateTree(s.storage).Set(11, &models.UserState{
		PubKeyID:   1,
		TokenIndex: models.MakeUint256(1),
		Balance:    models.MakeUint256(100),
		Nonce:      models.MakeUint256(0),
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
		Witness: []common.Hash{commitments[0].BodyHash(*batch.AccountTreeRoot)},
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
		Witness: []common.Hash{s.decodedCommitments[0].BodyHash(*s.decodedBatch.AccountTreeRoot)},
	}

	proof, err := targetCommitmentInclusionProof(&s.decodedBatch, 1)
	s.NoError(err)
	s.Equal(expected, *proof)
}

// nolint:funlen
func (s *DisputeTransitionTestSuite) TestRevertToForDispute() {
	s.setGenesisState()

	txs := []models.Transfer{
		{
			TransactionBase: models.TransactionBase{
				Hash:        utils.RandomHash(),
				TxType:      txtype.Transfer,
				FromStateID: 0,
				Amount:      models.MakeUint256(100),
				Fee:         models.MakeUint256(10),
				Nonce:       models.MakeUint256(0),
			},
			ToStateID: 2,
		},
		{
			TransactionBase: models.TransactionBase{
				Hash:        utils.RandomHash(),
				TxType:      txtype.Transfer,
				FromStateID: 1,
				Amount:      models.MakeUint256(100),
				Fee:         models.MakeUint256(10),
				Nonce:       models.MakeUint256(0),
			},
			ToStateID: 0,
		},
		{
			TransactionBase: models.TransactionBase{
				Hash:        utils.RandomHash(),
				TxType:      txtype.Transfer,
				FromStateID: 2,
				Amount:      models.MakeUint256(50),
				Fee:         models.MakeUint256(10),
				Nonce:       models.MakeUint256(0),
			},
			ToStateID: 0,
		},
		{
			TransactionBase: models.TransactionBase{
				Hash:        utils.RandomHash(),
				TxType:      txtype.Transfer,
				FromStateID: 2,
				Amount:      models.MakeUint256(500),
				Fee:         models.MakeUint256(10),
				Nonce:       models.MakeUint256(1),
			},
			ToStateID: 0,
		},
	}

	expectedProofs := []models.StateMerkleProof{
		{
			UserState: &models.UserState{
				PubKeyID:   0,
				TokenIndex: models.MakeUint256(0),
				Balance:    models.MakeUint256(340),
				Nonce:      models.MakeUint256(1),
			},
		},
		{
			UserState: &models.UserState{
				PubKeyID:   2,
				TokenIndex: models.MakeUint256(0),
				Balance:    models.MakeUint256(140),
				Nonce:      models.MakeUint256(1),
			},
		},
		{
			UserState: &models.UserState{
				PubKeyID:   0,
				TokenIndex: models.MakeUint256(0),
				Balance:    models.MakeUint256(290),
				Nonce:      models.MakeUint256(1),
			},
		},
		{
			UserState: &models.UserState{
				PubKeyID:   2,
				TokenIndex: models.MakeUint256(0),
				Balance:    models.MakeUint256(200),
				Nonce:      models.MakeUint256(0),
			},
		},
		{
			UserState: &models.UserState{
				PubKeyID:   0,
				TokenIndex: models.MakeUint256(0),
				Balance:    models.MakeUint256(190),
				Nonce:      models.MakeUint256(1),
			},
		},
		{
			UserState: &models.UserState{
				PubKeyID:   1,
				TokenIndex: models.MakeUint256(0),
				Balance:    models.MakeUint256(200),
				Nonce:      models.MakeUint256(0),
			},
		},
		{
			UserState: &models.UserState{
				PubKeyID:   2,
				TokenIndex: models.MakeUint256(0),
				Balance:    models.MakeUint256(100),
				Nonce:      models.MakeUint256(0),
			},
		},
		{
			UserState: &models.UserState{
				PubKeyID:   0,
				TokenIndex: models.MakeUint256(0),
				Balance:    models.MakeUint256(300),
				Nonce:      models.MakeUint256(0),
			},
		},
	}

	initialRoot, err := s.transactionExecutor.stateTree.Root()
	s.NoError(err)

	invalidTransfers := make([]models.Transfer, 0, 1)
	for i := range txs {
		senderWitness, receiverWitness := s.getTransferWitness(txs[i].FromStateID, txs[i].ToStateID)
		expectedProofs[len(expectedProofs)-1-i*2].Witness = senderWitness
		expectedProofs[len(expectedProofs)-2-i*2].Witness = receiverWitness

		var transferError error
		transferError, err = s.transactionExecutor.ApplyTransfer(&txs[i], models.MakeUint256(0))
		s.NoError(err)

		if transferError != nil {
			invalidTransfers = append(invalidTransfers, txs[i])
		}
	}
	s.Len(invalidTransfers, 1)

	proofs, err := s.transactionExecutor.stateTree.RevertToForDispute(*initialRoot, &invalidTransfers[0])
	s.NoError(err)
	s.Len(proofs, len(expectedProofs))
	s.Equal(expectedProofs, proofs)
}

func (s *DisputeTransitionTestSuite) setGenesisState() {
	userStates := []models.UserState{
		{
			PubKeyID:   0,
			TokenIndex: models.MakeUint256(0),
			Balance:    models.MakeUint256(300),
			Nonce:      models.MakeUint256(0),
		},
		{
			PubKeyID:   1,
			TokenIndex: models.MakeUint256(0),
			Balance:    models.MakeUint256(200),
			Nonce:      models.MakeUint256(0),
		},
		{
			PubKeyID:   2,
			TokenIndex: models.MakeUint256(0),
			Balance:    models.MakeUint256(100),
			Nonce:      models.MakeUint256(0),
		},
	}
	for i := range userStates {
		err := s.transactionExecutor.stateTree.Set(uint32(i), &userStates[i])
		s.NoError(err)
	}
}

func (s *DisputeTransitionTestSuite) getTransferWitness(fromStateID, toStateID uint32) (senderWitness, receiverWitness models.Witness) {
	var err error

	senderWitness, err = s.transactionExecutor.stateTree.GetWitness(models.MakeMerklePathFromStateID(fromStateID))
	s.NoError(err)

	receiverWitness, err = s.transactionExecutor.stateTree.GetWitness(models.MakeMerklePathFromStateID(toStateID))
	s.NoError(err)

	return senderWitness, receiverWitness
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(DisputeTransitionTestSuite))
}
