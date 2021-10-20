package prover

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils"
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
	decodedBatch       eth.DecodedBatch
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
	s.decodedBatch = eth.DecodedBatch{
		Batch: models.Batch{
			ID:                models.MakeUint256(2),
			Type:              batchtype.Transfer,
			TransactionHash:   utils.RandomHash(),
			Hash:              utils.NewRandomHash(),
			FinalisationBlock: ref.Uint32(10),
			AccountTreeRoot:   utils.NewRandomHash(),
		},
		Commitments: s.decodedCommitments,
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
		StateRoot: s.decodedCommitments[0].StateRoot,
		BodyRoot:  s.decodedCommitments[0].BodyHash(*s.decodedBatch.AccountTreeRoot),
		Path: &models.MerklePath{
			Path:  0,
			Depth: 2,
		},
		Witness: []common.Hash{s.decodedCommitments[1].LeafHash(*s.decodedBatch.AccountTreeRoot)},
	}

	proof, err := s.proverCtx.PreviousCommitmentInclusionProof(&s.decodedBatch, 0)
	s.NoError(err)
	s.Equal(expected, *proof)
}

func (s *CommitmentProofsTestSuite) TestPreviousCommitmentInclusionProof_PreviousBatch() {
	_, err := s.storage.StateTree.Set(11, &models.UserState{
		PubKeyID: 1,
		TokenID:  models.MakeUint256(1),
		Balance:  models.MakeUint256(100),
		Nonce:    models.MakeUint256(0),
	})
	s.NoError(err)

	batch := models.Batch{
		ID:                models.MakeUint256(1),
		Type:              batchtype.Transfer,
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
			Type:              batchtype.Transfer,
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
			Type:              batchtype.Transfer,
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

	proof, err := s.proverCtx.PreviousCommitmentInclusionProof(&s.decodedBatch, -1)
	s.NoError(err)
	s.Equal(expected, *proof)
}

func (s *CommitmentProofsTestSuite) TestGenesisBatchCommitmentInclusionProof() {
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
	proof, err := s.proverCtx.PreviousCommitmentInclusionProof(&firstBatch, -1)
	s.NoError(err)
	s.Equal(expected, *proof)
}

func (s *CommitmentProofsTestSuite) TestTargetCommitmentInclusionProof() {
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

	proof, err := s.proverCtx.TargetCommitmentInclusionProof(&s.decodedBatch, 1)
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
