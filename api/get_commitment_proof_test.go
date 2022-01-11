package api

import (
	"testing"
	"time"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/testutils"
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
	commitment                    models.TxCommitment
	commitmentProofNotFoundAPIErr *APIError
}

func (s *GetCommitmentProofTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *GetCommitmentProofTestSuite) SetupTest() {
	var err error
	s.storage, err = st.NewTestStorage()
	s.NoError(err)
	s.api = &API{
		storage: s.storage.Storage,
		cfg:     &config.APIConfig{EnableProofMethods: true},
	}

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
	s.commitment.BodyHash = utils.NewRandomHash()

	s.commitmentProofNotFoundAPIErr = &APIError{
		Code:    50001,
		Message: "commitment inclusion proof could not be generated",
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

	err = s.storage.AddTxCommitment(&s.commitment)
	s.NoError(err)

	s.addStateLeaf()

	transfer := testutils.MakeTransfer(1, 2, 0, 50)
	transfer.CommitmentID = &s.commitment.ID
	err = s.storage.AddTransaction(&transfer)
	s.NoError(err)

	tree, err := merkletree.NewMerkleTree([]common.Hash{s.commitment.LeafHash()})
	s.NoError(err)

	path := &dto.MerklePath{
		Path:  uint32(s.commitment.ID.IndexInBatch),
		Depth: tree.Depth(),
	}

	expectedCommitmentProof := &dto.CommitmentInclusionProof{
		CommitmentInclusionProofBase: dto.CommitmentInclusionProofBase{
			StateRoot: s.commitment.PostStateRoot,
			Path:      path,
			Witness:   tree.GetWitness(uint32(s.commitment.ID.IndexInBatch)),
		},
		Body: &dto.CommitmentProofBody{
			AccountRoot: *s.batch.AccountTreeRoot,
			Signature:   s.commitment.CombinedSignature,
			FeeReceiver: s.commitment.FeeReceiver,
			Transactions: []dto.TransferForCommitment{
				dto.MakeTransferForCommitment(&transfer),
			},
		},
	}
	commitmentProof, err := s.api.GetCommitmentProof(s.commitment.ID)

	s.NoError(err)
	s.Equal(expectedCommitmentProof, commitmentProof)
}

func (s *GetCommitmentProofTestSuite) TestGetCommitmentProof_Create2TransferType() {
	err := s.storage.AddBatch(&s.batch)
	s.NoError(err)

	s.commitment.Type = batchtype.Create2Transfer
	err = s.storage.AddTxCommitment(&s.commitment)
	s.NoError(err)

	s.addStateLeaf()

	transfer := testutils.MakeCreate2Transfer(1, ref.Uint32(2), 0, 50, &models.PublicKey{2, 3, 4})
	transfer.CommitmentID = &s.commitment.ID
	err = s.storage.AddTransaction(&transfer)
	s.NoError(err)

	tree, err := merkletree.NewMerkleTree([]common.Hash{s.commitment.LeafHash()})
	s.NoError(err)

	path := &dto.MerklePath{
		Path:  uint32(s.commitment.ID.IndexInBatch),
		Depth: tree.Depth(),
	}

	expectedCommitmentProof := &dto.CommitmentInclusionProof{
		CommitmentInclusionProofBase: dto.CommitmentInclusionProofBase{
			StateRoot: s.commitment.PostStateRoot,
			Path:      path,
			Witness:   tree.GetWitness(uint32(s.commitment.ID.IndexInBatch)),
		},
		Body: &dto.CommitmentProofBody{
			AccountRoot: *s.batch.AccountTreeRoot,
			Signature:   s.commitment.CombinedSignature,
			FeeReceiver: s.commitment.FeeReceiver,
			Transactions: []dto.Create2TransferForCommitment{
				dto.MakeCreate2TransferForCommitment(&transfer),
			},
		},
	}

	commitmentProof, err := s.api.GetCommitmentProof(s.commitment.ID)
	s.NoError(err)
	s.Equal(expectedCommitmentProof, commitmentProof)
}

func (s *GetCommitmentProofTestSuite) TestGetCommitmentProof_MassMigrationType() {
	err := s.storage.AddBatch(&s.batch)
	s.NoError(err)

	s.commitment.Type = batchtype.MassMigration
	err = s.storage.AddTxCommitment(&s.commitment)
	s.NoError(err)

	s.addStateLeaf()

	massMigration := testutils.MakeMassMigration(1, 2, 0, 50)
	massMigration.CommitmentID = &s.commitment.ID
	err = s.storage.AddTransaction(&massMigration)
	s.NoError(err)

	tree, err := merkletree.NewMerkleTree([]common.Hash{s.commitment.LeafHash()})
	s.NoError(err)

	path := &dto.MerklePath{
		Path:  uint32(s.commitment.ID.IndexInBatch),
		Depth: tree.Depth(),
	}

	expectedCommitmentProof := &dto.CommitmentInclusionProof{
		CommitmentInclusionProofBase: dto.CommitmentInclusionProofBase{
			StateRoot: s.commitment.PostStateRoot,
			Path:      path,
			Witness:   tree.GetWitness(uint32(s.commitment.ID.IndexInBatch)),
		},
		Body: &dto.CommitmentProofBody{
			AccountRoot: *s.batch.AccountTreeRoot,
			Signature:   s.commitment.CombinedSignature,
			FeeReceiver: s.commitment.FeeReceiver,
			Transactions: []dto.MassMigrationForCommitment{
				dto.MakeMassMigrationForCommitment(&massMigration),
			},
		},
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

	err = s.storage.AddTxCommitment(&s.commitment)
	s.NoError(err)

	s.addStateLeaf()

	transfer := testutils.MakeTransfer(1, 2, 0, 50)
	transfer.CommitmentID = &s.commitment.ID
	err = s.storage.AddTransaction(&transfer)
	s.NoError(err)

	commitmentProof, err := s.api.GetCommitmentProof(s.commitment.ID)
	s.Equal(s.commitmentProofNotFoundAPIErr, err)
	s.Nil(commitmentProof)
}

func (s *GetCommitmentProofTestSuite) TestGetCommitmentProof_NonexistentCommitment() {
	commitmentProof, err := s.api.GetCommitmentProof(s.commitment.ID)
	s.Equal(s.commitmentProofNotFoundAPIErr, err)
	s.Nil(commitmentProof)
}

func TestGetCommitmentProofTestSuite(t *testing.T) {
	suite.Run(t, new(GetCommitmentProofTestSuite))
}
