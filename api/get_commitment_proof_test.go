package api

import (
	"testing"
	"time"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/Worldcoin/hubble-commander/utils/merkletree"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type GetCommitmentProofTestSuite struct {
	*require.Assertions
	suite.Suite
	api                           *API
	storage                       *st.TestStorage
	batch                         models.Batch
	commitment                    models.Commitment
	commitmentProofNotFoundAPIErr *APIError
}

func (s *GetCommitmentProofTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *GetCommitmentProofTestSuite) SetupTest() {
	var err error
	s.storage, err = st.NewTestStorage()
	s.NoError(err)
	s.api = &API{storage: s.storage.Storage}

	s.batch = models.Batch{
		ID:                models.MakeUint256(1),
		Type:              batchtype.Transfer,
		TransactionHash:   utils.RandomHash(),
		Hash:              utils.NewRandomHash(),
		FinalisationBlock: ref.Uint32(113),
		SubmissionTime:    models.NewTimestamp(time.Unix(140, 0).UTC()),
		AccountTreeRoot:   utils.NewRandomHash(),
	}

	s.commitment = commitment
	s.commitment.ID.BatchID = s.batch.ID
	s.commitment.ID.IndexInBatch = 0

	s.commitmentProofNotFoundAPIErr = &APIError{
		Code:    20001,
		Message: "commitment not found",
	}
}

func (s *GetCommitmentProofTestSuite) addStateLeaf() {
	_, err := s.storage.StateTree.Set(uint32(1), &models.UserState{
		PubKeyID: 1,
		TokenID:  models.MakeUint256(1),
		Balance:  models.MakeUint256(420),
		Nonce:    models.MakeUint256(0),
	})
	s.NoError(err)
}

func (s *GetCommitmentProofTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *GetCommitmentProofTestSuite) TestGetCommitmentProof_TransferType() {
	err := s.storage.AddBatch(&s.batch)
	s.NoError(err)

	err = s.storage.AddCommitment(&s.commitment)
	s.NoError(err)

	s.addStateLeaf()

	transfer := models.Transfer{
		TransactionBase: models.TransactionBase{
			Hash:         utils.RandomHash(),
			TxType:       txtype.Transfer,
			FromStateID:  1,
			Amount:       models.MakeUint256(50),
			Fee:          models.MakeUint256(10),
			Nonce:        models.MakeUint256(0),
			Signature:    models.MakeRandomSignature(),
			CommitmentID: &s.commitment.ID,
		},
		ToStateID: 2,
	}
	err = s.storage.AddTransfer(&transfer)
	s.NoError(err)

	tree, err := merkletree.NewMerkleTree([]common.Hash{s.commitment.LeafHash(*s.batch.AccountTreeRoot)})
	s.NoError(err)

	path := &models.MerklePath{
		Path:  uint32(s.commitment.ID.IndexInBatch),
		Depth: tree.Depth(),
	}

	expectedCommitmentProof := &dto.TransferCommitmentInclusionProof{
		StateRoot: commitment.PostStateRoot,
		Body: &dto.TransferBody{
			AccountRoot: *s.batch.AccountTreeRoot,
			Signature:   commitment.CombinedSignature,
			FeeReceiver: commitment.FeeReceiver,
			Transactions: []dto.TransferForCommitment{{
				Hash:        transfer.Hash,
				FromStateID: transfer.FromStateID,
				Amount:      transfer.Amount,
				Fee:         transfer.Fee,
				Nonce:       transfer.Nonce,
				Signature:   transfer.Signature,
				ReceiveTime: transfer.ReceiveTime,
				ToStateID:   transfer.ToStateID,
			}},
		},
		Path:    path,
		Witness: tree.GetWitness(uint32(commitment.ID.IndexInBatch)),
	}
	commitmentProof, err := s.api.GetCommitmentProof(s.commitment.ID)

	s.NoError(err)
	s.Equal(expectedCommitmentProof, commitmentProof)
}

func (s *GetCommitmentProofTestSuite) TestGetCommitmentProof_Create2TransferType() {
	err := s.storage.AddBatch(&s.batch)
	s.NoError(err)

	s.commitment.Type = batchtype.Create2Transfer
	err = s.storage.AddCommitment(&s.commitment)
	s.NoError(err)

	s.addStateLeaf()

	transfer := models.Create2Transfer{
		TransactionBase: models.TransactionBase{
			Hash:         utils.RandomHash(),
			TxType:       txtype.Create2Transfer,
			FromStateID:  1,
			Amount:       models.MakeUint256(50),
			Fee:          models.MakeUint256(10),
			Nonce:        models.MakeUint256(0),
			CommitmentID: &s.commitment.ID,
		},
		ToStateID:   ref.Uint32(2),
		ToPublicKey: models.PublicKey{2, 3, 4},
	}
	err = s.storage.AddCreate2Transfer(&transfer)
	s.NoError(err)

	tree, err := merkletree.NewMerkleTree([]common.Hash{s.commitment.LeafHash(*s.batch.AccountTreeRoot)})
	s.NoError(err)

	path := &models.MerklePath{
		Path:  uint32(s.commitment.ID.IndexInBatch),
		Depth: tree.Depth(),
	}

	expectedCommitmentProof := &dto.TransferCommitmentInclusionProof{
		StateRoot: commitment.PostStateRoot,
		Body: &dto.TransferBody{
			AccountRoot: *s.batch.AccountTreeRoot,
			Signature:   commitment.CombinedSignature,
			FeeReceiver: commitment.FeeReceiver,
			Transactions: []dto.Create2TransferForCommitment{{
				Hash:        transfer.Hash,
				FromStateID: transfer.FromStateID,
				Amount:      transfer.Amount,
				Fee:         transfer.Fee,
				Nonce:       transfer.Nonce,
				Signature:   transfer.Signature,
				ReceiveTime: transfer.ReceiveTime,
				ToStateID:   transfer.ToStateID,
				ToPublicKey: transfer.ToPublicKey,
			}},
		},
		Path:    path,
		Witness: tree.GetWitness(uint32(commitment.ID.IndexInBatch)),
	}
	commitmentProof, err := s.api.GetCommitmentProof(s.commitment.ID)

	s.NoError(err)
	s.Equal(expectedCommitmentProof, commitmentProof)
}

func (s *GetCommitmentProofTestSuite) TestGetCommitmentProof_PendingBatch() {
	pendingBatch := s.batch
	pendingBatch.Hash = nil
	pendingBatch.FinalisationBlock = nil
	err := s.storage.AddBatch(&pendingBatch)
	s.NoError(err)

	err = s.storage.AddCommitment(&s.commitment)
	s.NoError(err)

	s.addStateLeaf()

	transfer := models.Transfer{
		TransactionBase: models.TransactionBase{
			Hash:         utils.RandomHash(),
			TxType:       txtype.Transfer,
			FromStateID:  1,
			Amount:       models.MakeUint256(50),
			Fee:          models.MakeUint256(10),
			Nonce:        models.MakeUint256(0),
			CommitmentID: &s.commitment.ID,
		},
		ToStateID: 2,
	}
	err = s.storage.AddTransfer(&transfer)
	s.NoError(err)

	commitmentProof, err := s.api.GetCommitmentProof(commitment.ID)
	s.Equal(s.commitmentProofNotFoundAPIErr, err)
	s.Nil(commitmentProof)
}

func (s *GetCommitmentProofTestSuite) TestGetCommitmentProof_NonexistentCommitment() {
	commitmentProof, err := s.api.GetCommitmentProof(commitment.ID)
	s.Equal(s.commitmentProofNotFoundAPIErr, err)
	s.Nil(commitmentProof)
}

func TestGetCommitmentProofTestSuite(t *testing.T) {
	suite.Run(t, new(GetCommitmentProofTestSuite))
}
