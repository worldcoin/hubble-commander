package prover

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/Worldcoin/hubble-commander/utils/consts"
	"github.com/Worldcoin/hubble-commander/utils/merkletree"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type CommitmentProofsTestSuite struct {
	*require.Assertions
	suite.Suite
	storage            *st.TestStorage
	proverCtx          *Context
	decodedCommitments []encoder.DecodedCommitment
	decodedBatch       eth.DecodedTxBatch
}

func (s *CommitmentProofsTestSuite) SetupSuite() {
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
	s.decodedBatch = eth.DecodedTxBatch{
		DecodedBatchBase: eth.DecodedBatchBase{
			ID:                models.MakeUint256(2),
			Type:              batchtype.Transfer,
			TransactionHash:   utils.RandomHash(),
			Hash:              utils.RandomHash(),
			FinalisationBlock: 10,
			AccountTreeRoot:   utils.RandomHash(),
		},
		Commitments: encoder.DecodedCommitmentsToCommitments(s.decodedCommitments...),
	}
}

func (s *CommitmentProofsTestSuite) SetupTest() {
	var err error
	s.storage, err = st.NewTestStorage()
	s.NoError(err)

	s.proverCtx = NewContext(s.storage.Storage)
}

func (s *CommitmentProofsTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *CommitmentProofsTestSuite) TestPreviousCommitmentInclusionProof_CurrentBatch() {
	expected := models.CommitmentInclusionProof{
		CommitmentInclusionProofBase: models.CommitmentInclusionProofBase{
			StateRoot: s.decodedCommitments[0].StateRoot,
			Path: &models.MerklePath{
				Path:  0,
				Depth: 2,
			},
			Witness: []common.Hash{s.decodedCommitments[1].LeafHash(s.decodedBatch.AccountTreeRoot)},
		},
		BodyRoot: *s.decodedCommitments[0].BodyHash(s.decodedBatch.AccountTreeRoot),
	}

	proof, err := s.proverCtx.PreviousCommitmentInclusionProof(&s.decodedBatch, 0)
	s.NoError(err)
	s.Equal(expected, *proof)
}

func (s *CommitmentProofsTestSuite) TestPreviousCommitmentInclusionProof_PreviousTransactionBatch() {
	batch := models.Batch{
		ID:                models.MakeUint256(1),
		Type:              batchtype.Transfer,
		TransactionHash:   utils.RandomHash(),
		Hash:              utils.NewRandomHash(),
		FinalisationBlock: ref.Uint32(10),
		AccountTreeRoot:   utils.NewRandomHash(),
	}
	err := s.storage.AddBatch(&batch)
	s.NoError(err)

	commitments := []models.TxCommitment{
		{
			CommitmentBase: models.CommitmentBase{
				ID: models.CommitmentID{
					BatchID:      batch.ID,
					IndexInBatch: 0,
				},
				Type:          batchtype.Transfer,
				PostStateRoot: utils.RandomHash(),
			},
			FeeReceiver:       11,
			CombinedSignature: models.MakeRandomSignature(),
			BodyHash:          utils.NewRandomHash(),
		},
		{
			CommitmentBase: models.CommitmentBase{
				ID: models.CommitmentID{
					BatchID:      batch.ID,
					IndexInBatch: 1,
				},
				Type:          batchtype.Transfer,
				PostStateRoot: utils.RandomHash(),
			},
			FeeReceiver:       11,
			CombinedSignature: models.MakeRandomSignature(),
			BodyHash:          utils.NewRandomHash(),
		},
	}
	for i := range commitments {
		err = s.storage.AddCommitment(&commitments[i])
		s.NoError(err)
	}

	expected := models.CommitmentInclusionProof{
		CommitmentInclusionProofBase: models.CommitmentInclusionProofBase{
			StateRoot: commitments[1].PostStateRoot,
			Path: &models.MerklePath{
				Path:  1,
				Depth: 2,
			},
			Witness: []common.Hash{commitments[0].LeafHash()},
		},
		BodyRoot: *commitments[1].BodyHash,
	}

	proof, err := s.proverCtx.PreviousCommitmentInclusionProof(&s.decodedBatch, -1)
	s.NoError(err)
	s.Equal(expected, *proof)
}

func (s *CommitmentProofsTestSuite) TestPreviousCommitmentInclusionProof_PreviousDepositBatch() {
	batch := models.Batch{
		ID:                models.MakeUint256(1),
		Type:              batchtype.Deposit,
		TransactionHash:   utils.RandomHash(),
		Hash:              utils.NewRandomHash(),
		FinalisationBlock: ref.Uint32(10),
		AccountTreeRoot:   utils.NewRandomHash(),
	}
	err := s.storage.AddBatch(&batch)
	s.NoError(err)

	commitment := models.DepositCommitment{
		CommitmentBase: models.CommitmentBase{
			ID: models.CommitmentID{
				BatchID:      batch.ID,
				IndexInBatch: 0,
			},
			Type: batchtype.Deposit,
		},
		SubtreeID:   models.MakeUint256(1),
		SubtreeRoot: common.Hash{1, 2, 3},
		Deposits: []models.PendingDeposit{
			{
				ID: models.DepositID{
					SubtreeID:    models.MakeUint256(1),
					DepositIndex: models.MakeUint256(0),
				},
				ToPubKeyID: 1,
				TokenID:    models.MakeUint256(0),
				L2Amount:   models.MakeUint256(10),
			},
		},
	}
	err = s.storage.AddCommitment(&commitment)
	s.NoError(err)

	expected := models.CommitmentInclusionProof{
		CommitmentInclusionProofBase: models.CommitmentInclusionProofBase{
			StateRoot: commitment.PostStateRoot,
			Path: &models.MerklePath{
				Path:  0,
				Depth: 2,
			},
			Witness: []common.Hash{consts.ZeroHash},
		},
		BodyRoot: commitment.GetBodyHash(),
	}

	proof, err := s.proverCtx.PreviousCommitmentInclusionProof(&s.decodedBatch, -1)
	s.NoError(err)
	s.Equal(expected, *proof)
}

func (s *CommitmentProofsTestSuite) TestGenesisBatchCommitmentInclusionProof() {
	genesisBatch := s.addGenesisBatch()
	zeroHash := merkletree.GetZeroHash(0)

	expected := models.CommitmentInclusionProof{
		CommitmentInclusionProofBase: models.CommitmentInclusionProofBase{
			StateRoot: *genesisBatch.PrevStateRoot,
			Path: &models.MerklePath{
				Path:  0,
				Depth: 2,
			},
			Witness: []common.Hash{zeroHash},
		},
		BodyRoot: zeroHash,
	}

	firstBatch := s.decodedBatch
	firstBatch.ID = models.MakeUint256(1)
	proof, err := s.proverCtx.PreviousCommitmentInclusionProof(&firstBatch, -1)
	s.NoError(err)
	s.Equal(expected, *proof)
}

func (s *CommitmentProofsTestSuite) TestTargetTransferCommitmentInclusionProof() {
	expected := models.TransferCommitmentInclusionProof{
		CommitmentInclusionProofBase: models.CommitmentInclusionProofBase{
			StateRoot: s.decodedCommitments[1].StateRoot,
			Path: &models.MerklePath{
				Path:  1,
				Depth: 2,
			},
			Witness: []common.Hash{s.decodedCommitments[0].LeafHash(s.decodedBatch.AccountTreeRoot)},
		},
		Body: &models.TransferBody{
			AccountRoot:  s.decodedBatch.AccountTreeRoot,
			Signature:    s.decodedCommitments[1].CombinedSignature,
			FeeReceiver:  s.decodedCommitments[1].FeeReceiver,
			Transactions: s.decodedCommitments[1].Transactions,
		},
	}

	proof, err := s.proverCtx.TargetTransferCommitmentInclusionProof(&s.decodedBatch, 1)
	s.NoError(err)
	s.Equal(expected, *proof)
}

func (s *CommitmentProofsTestSuite) TestTargetMMCommitmentInclusionProof() {
	decodedMMBatch := s.decodedBatch
	decodedMMBatch.Type = batchtype.MassMigration
	decodedMMBatch.Commitments = make([]encoder.Commitment, 0, len(s.decodedCommitments))

	for i := range s.decodedCommitments {
		decodedMMBatch.Commitments = append(decodedMMBatch.Commitments,
			&encoder.DecodedMMCommitment{
				DecodedCommitment: s.decodedCommitments[i],
				Meta: &models.MassMigrationMeta{
					SpokeID:     uint32(i),
					TokenID:     models.MakeUint256(1),
					Amount:      models.MakeUint256(100),
					FeeReceiver: 0,
				},
				WithdrawRoot: utils.RandomHash(),
			},
		)
	}

	mmCommitment := decodedMMBatch.Commitments[1].(*encoder.DecodedMMCommitment)
	expected := models.MMCommitmentInclusionProof{
		CommitmentInclusionProofBase: models.CommitmentInclusionProofBase{
			StateRoot: mmCommitment.StateRoot,
			Path: &models.MerklePath{
				Path:  1,
				Depth: 2,
			},
			Witness: []common.Hash{decodedMMBatch.Commitments[0].LeafHash(decodedMMBatch.AccountTreeRoot)},
		},
		Body: &models.MMBody{
			AccountRoot:  decodedMMBatch.AccountTreeRoot,
			Signature:    mmCommitment.CombinedSignature,
			Meta:         mmCommitment.Meta,
			WithdrawRoot: mmCommitment.WithdrawRoot,
			Transactions: mmCommitment.Transactions,
		},
	}

	proof, err := s.proverCtx.TargetMMCommitmentInclusionProof(&decodedMMBatch, 1)
	s.NoError(err)
	s.Equal(expected, *proof)
}

func (s *CommitmentProofsTestSuite) addGenesisBatch() *models.Batch {
	root, err := s.storage.StateTree.Root()
	s.NoError(err)

	batch := &models.Batch{
		ID:              models.MakeUint256(0),
		Type:            batchtype.Genesis,
		TransactionHash: common.Hash{},
		Hash:            utils.NewRandomHash(),
		PrevStateRoot:   root,
	}

	err = s.storage.AddBatch(batch)
	s.NoError(err)

	return batch
}

func TestCommitmentProofsTestSuite(t *testing.T) {
	suite.Run(t, new(CommitmentProofsTestSuite))
}
