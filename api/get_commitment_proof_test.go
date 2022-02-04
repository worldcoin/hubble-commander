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
	api          *API
	storage      *st.TestStorage
	batch        models.Batch
	txCommitment models.TxCommitment
	mmCommitment models.MMCommitment
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
		MinedTime:         models.NewTimestamp(time.Unix(140, 0).UTC()),
		AccountTreeRoot:   utils.NewRandomHash(),
	}

	s.txCommitment = commitment
	s.txCommitment.ID.BatchID = s.batch.ID
	s.txCommitment.ID.IndexInBatch = 0
	s.txCommitment.BodyHash = utils.NewRandomHash()

	s.mmCommitment = models.MMCommitment{
		CommitmentBase: models.CommitmentBase{
			ID: models.CommitmentID{
				BatchID:      s.batch.ID,
				IndexInBatch: 0,
			},
			Type:          batchtype.MassMigration,
			PostStateRoot: utils.RandomHash(),
		},
		CombinedSignature: models.MakeRandomSignature(),
		Meta: &models.MassMigrationMeta{
			SpokeID:     1,
			TokenID:     models.MakeUint256(2),
			Amount:      models.MakeUint256(3),
			FeeReceiver: 1,
		},
		BodyHash:     utils.NewRandomHash(),
		WithdrawRoot: utils.RandomHash(),
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

	err = s.storage.AddCommitment(&s.txCommitment)
	s.NoError(err)

	s.addStateLeaf()

	transfer := testutils.MakeTransfer(1, 2, 0, 50)
	transfer.CommitmentID = &s.txCommitment.ID
	err = s.storage.AddTransaction(&transfer)
	s.NoError(err)

	tree, err := merkletree.NewMerkleTree([]common.Hash{s.txCommitment.LeafHash()})
	s.NoError(err)

	path := &dto.MerklePath{
		Path:  uint32(s.txCommitment.ID.IndexInBatch),
		Depth: tree.Depth(),
	}

	expectedCommitmentProof := &dto.CommitmentInclusionProof{
		CommitmentInclusionProofBase: dto.CommitmentInclusionProofBase{
			StateRoot: s.txCommitment.PostStateRoot,
			Path:      path,
			Witness:   tree.GetWitness(uint32(s.txCommitment.ID.IndexInBatch)),
		},
		Body: &dto.CommitmentProofBody{
			AccountRoot: *s.batch.AccountTreeRoot,
			Signature:   s.txCommitment.CombinedSignature,
			FeeReceiver: s.txCommitment.FeeReceiver,
			Transactions: []dto.TransferForCommitment{
				dto.MakeTransferForCommitment(&transfer),
			},
		},
	}
	commitmentProof, err := s.api.GetCommitmentProof(s.txCommitment.ID)

	s.NoError(err)
	s.Equal(expectedCommitmentProof, commitmentProof)
}

func (s *GetCommitmentProofTestSuite) TestGetCommitmentProof_Create2TransferType() {
	err := s.storage.AddBatch(&s.batch)
	s.NoError(err)

	s.txCommitment.Type = batchtype.Create2Transfer
	err = s.storage.AddCommitment(&s.txCommitment)
	s.NoError(err)

	s.addStateLeaf()

	transfer := testutils.MakeCreate2Transfer(1, ref.Uint32(2), 0, 50, &models.PublicKey{2, 3, 4})
	transfer.CommitmentID = &s.txCommitment.ID
	err = s.storage.AddTransaction(&transfer)
	s.NoError(err)

	tree, err := merkletree.NewMerkleTree([]common.Hash{s.txCommitment.LeafHash()})
	s.NoError(err)

	path := &dto.MerklePath{
		Path:  uint32(s.txCommitment.ID.IndexInBatch),
		Depth: tree.Depth(),
	}

	expectedCommitmentProof := &dto.CommitmentInclusionProof{
		CommitmentInclusionProofBase: dto.CommitmentInclusionProofBase{
			StateRoot: s.txCommitment.PostStateRoot,
			Path:      path,
			Witness:   tree.GetWitness(uint32(s.txCommitment.ID.IndexInBatch)),
		},
		Body: &dto.CommitmentProofBody{
			AccountRoot: *s.batch.AccountTreeRoot,
			Signature:   s.txCommitment.CombinedSignature,
			FeeReceiver: s.txCommitment.FeeReceiver,
			Transactions: []dto.Create2TransferForCommitment{
				dto.MakeCreate2TransferForCommitment(&transfer),
			},
		},
	}

	commitmentProof, err := s.api.GetCommitmentProof(s.txCommitment.ID)
	s.NoError(err)
	s.Equal(expectedCommitmentProof, commitmentProof)
}

func (s *GetCommitmentProofTestSuite) TestGetCommitmentProof_MassMigrationType() {
	err := s.storage.AddBatch(&s.batch)
	s.NoError(err)

	err = s.storage.AddCommitment(&s.mmCommitment)
	s.NoError(err)

	s.addStateLeaf()

	massMigration := testutils.MakeMassMigration(1, 2, 0, 50)
	massMigration.CommitmentID = &s.mmCommitment.ID
	err = s.storage.AddTransaction(&massMigration)
	s.NoError(err)

	tree, err := merkletree.NewMerkleTree([]common.Hash{s.mmCommitment.LeafHash()})
	s.NoError(err)

	path := &dto.MerklePath{
		Path:  uint32(s.mmCommitment.ID.IndexInBatch),
		Depth: tree.Depth(),
	}

	expectedCommitmentProof := &dto.CommitmentInclusionProof{
		CommitmentInclusionProofBase: dto.CommitmentInclusionProofBase{
			StateRoot: s.mmCommitment.PostStateRoot,
			Path:      path,
			Witness:   tree.GetWitness(uint32(s.mmCommitment.ID.IndexInBatch)),
		},
		Body: &dto.CommitmentProofBody{
			AccountRoot: *s.batch.AccountTreeRoot,
			Signature:   s.mmCommitment.CombinedSignature,
			FeeReceiver: s.mmCommitment.Meta.FeeReceiver,
			Transactions: []dto.MassMigrationForCommitment{
				dto.MakeMassMigrationForCommitment(&massMigration),
			},
		},
	}

	commitmentProof, err := s.api.GetCommitmentProof(s.mmCommitment.ID)
	s.NoError(err)
	s.Equal(expectedCommitmentProof, commitmentProof)
}

func (s *GetCommitmentProofTestSuite) TestGetCommitmentProof_DepositType() {
	s.batch.Type = batchtype.Deposit
	err := s.storage.AddBatch(&s.batch)
	s.NoError(err)

	commitmentProof, err := s.api.GetCommitmentProof(s.txCommitment.ID)
	s.Equal(APIErrUnsupportedCommitmentTypeForProof, err)
	s.Nil(commitmentProof)
}

func (s *GetCommitmentProofTestSuite) TestGetCommitmentProof_PendingBatch() {
	pendingBatch := s.batch
	pendingBatch.Hash = nil
	pendingBatch.FinalisationBlock = nil
	err := s.storage.AddBatch(&pendingBatch)
	s.NoError(err)

	err = s.storage.AddCommitment(&s.txCommitment)
	s.NoError(err)

	s.addStateLeaf()

	transfer := testutils.MakeTransfer(1, 2, 0, 50)
	transfer.CommitmentID = &s.txCommitment.ID
	err = s.storage.AddTransaction(&transfer)
	s.NoError(err)

	commitmentProof, err := s.api.GetCommitmentProof(s.txCommitment.ID)
	s.Equal(APIErrCannotGenerateCommitmentProof, err)
	s.Nil(commitmentProof)
}

func (s *GetCommitmentProofTestSuite) TestGetCommitmentProof_NonexistentCommitment() {
	commitmentProof, err := s.api.GetCommitmentProof(s.txCommitment.ID)
	s.Equal(APIErrCannotGenerateCommitmentProof, err)
	s.Nil(commitmentProof)
}

func TestGetCommitmentProofTestSuite(t *testing.T) {
	suite.Run(t, new(GetCommitmentProofTestSuite))
}
