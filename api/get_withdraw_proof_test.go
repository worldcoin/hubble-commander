package api

import (
	"testing"
	"time"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/Worldcoin/hubble-commander/utils/merkletree"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type GetWithdrawProofTestSuite struct {
	*require.Assertions
	suite.Suite
	api            *API
	storage        *st.TestStorage
	commitment     models.MMCommitmentWithTxs
	massMigrations []models.MassMigration
}

func (s *GetWithdrawProofTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *GetWithdrawProofTestSuite) SetupTest() {
	var err error
	s.storage, err = st.NewTestStorage()
	s.NoError(err)
	s.api = &API{
		storage: s.storage.Storage,
		cfg:     &config.APIConfig{EnableProofMethods: true},
	}

	// unsorted mass migrations
	s.massMigrations = []models.MassMigration{
		makeMassMigration(
			common.Hash{2, 3, 4},
			0,
			0,
			models.NewTimestamp(time.Unix(140, 0).UTC()),
			models.CommitmentSlot{
				BatchID:           models.MakeUint256(1),
				IndexInBatch:      0,
				IndexInCommitment: 0,
			},
		),
		makeMassMigration(
			common.Hash{1, 2, 3},
			1,
			1,
			models.NewTimestamp(time.Unix(150, 0).UTC()),
			models.CommitmentSlot{
				BatchID:           models.MakeUint256(1),
				IndexInBatch:      0,
				IndexInCommitment: 1,
			},
		),
	}

	_, err = s.storage.StateTree.Set(0, &models.UserState{
		Balance: models.MakeUint256(0),
		TokenID: models.MakeUint256(10),
		Nonce:   models.MakeUint256(1),
	})
	s.NoError(err)

	_, err = s.storage.StateTree.Set(1, &models.UserState{
		Balance: models.MakeUint256(0),
		TokenID: models.MakeUint256(10),
		Nonce:   models.MakeUint256(2),
	})
	s.NoError(err)

	stateRoot, err := s.storage.StateTree.Root()
	s.NoError(err)

	err = s.storage.BatchAddTransaction(models.MakeMassMigrationArray(s.massMigrations...))
	s.NoError(err)

	accountTreeRoot := utils.RandomHash()
	s.commitment = makeMassMigrationCommitment(
		s.Assertions,
		s.storage,
		models.CommitmentID{
			BatchID:      models.MakeUint256(1),
			IndexInBatch: 0,
		},
		1,
		*stateRoot,
		accountTreeRoot,
		s.massMigrations,
	)

	err = s.storage.AddCommitment(&s.commitment)
	s.NoError(err)

	err = s.storage.AddBatch(&models.Batch{
		ID:   models.MakeUint256(1),
		Type: batchtype.MassMigration,
	})
	s.NoError(err)
}

func (s *GetWithdrawProofTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *GetWithdrawProofTestSuite) TestGetWithdrawProof_FirstMassMigrationInCommitment() {
	s.testGetWithdrawProofEndpoint(s.massMigrations[0].Hash)
}

func (s *GetWithdrawProofTestSuite) TestGetWithdrawProof_SecondMassMigrationInCommitment() {
	s.testGetWithdrawProofEndpoint(s.massMigrations[1].Hash)
}

func (s *GetWithdrawProofTestSuite) TestGetWithdrawProof_NonexistentBatch() {
	commitmentID := models.CommitmentID{
		BatchID:      models.MakeUint256(10),
		IndexInBatch: 15,
	}
	_, err := s.api.GetWithdrawProof(commitmentID, utils.RandomHash())
	s.Equal(APIWithdrawProofCouldNotBeCalculated, err)
}

func (s *GetWithdrawProofTestSuite) TestGetWithdrawProof_InvalidBatchType() {
	commitmentID := models.CommitmentID{
		BatchID:      models.MakeUint256(2),
		IndexInBatch: 0,
	}

	err := s.storage.AddBatch(&models.Batch{
		ID:   commitmentID.BatchID,
		Type: batchtype.Transfer,
	})
	s.NoError(err)

	_, err = s.api.GetWithdrawProof(commitmentID, utils.RandomHash())
	s.Equal(APIErrOnlyMassMigrationBatches, err)
}

func (s *GetWithdrawProofTestSuite) TestGetWithdrawProof_NonexistentMassMigrationWithGivenSenderInCommitment() {
	commitmentID := models.CommitmentID{
		BatchID:      models.MakeUint256(1),
		IndexInBatch: 0,
	}
	_, err := s.api.GetWithdrawProof(commitmentID, utils.RandomHash())
	s.Equal(APIErrMassMigrationWithTxHashNotFound, err)
}

func (s *GetWithdrawProofTestSuite) testGetWithdrawProofEndpoint(transactionHash common.Hash) {
	var (
		targetUserState    *models.UserState
		massMigrationIndex int
	)

	hashes := make([]common.Hash, 0, len(s.massMigrations))

	for i := range s.massMigrations {
		senderLeaf, err := s.storage.StateTree.Leaf(s.massMigrations[i].FromStateID)
		s.NoError(err)

		massMigrationUserState := &models.UserState{
			PubKeyID: senderLeaf.PubKeyID,
			TokenID:  senderLeaf.TokenID,
			Balance:  s.massMigrations[i].Amount,
			Nonce:    models.MakeUint256(0),
		}

		hash, err := encoder.HashUserState(massMigrationUserState)
		s.NoError(err)

		hashes = append(hashes, *hash)

		if s.massMigrations[i].Hash == transactionHash {
			targetUserState = massMigrationUserState
			massMigrationIndex = i
		}
	}

	withdrawTree, err := merkletree.NewMerkleTree(hashes)
	s.NoError(err)

	dtoTargetUserState := dto.MakeUserState(targetUserState)
	expected := dto.WithdrawProof{
		UserState: &dtoTargetUserState,
		Path: dto.MerklePath{
			Path:  uint32(massMigrationIndex),
			Depth: withdrawTree.Depth(),
		},
		Witness: withdrawTree.GetWitness(uint32(massMigrationIndex)),
		Root:    withdrawTree.Root(),
	}

	commitmentID := models.CommitmentID{
		BatchID:      models.MakeUint256(1),
		IndexInBatch: 0,
	}
	withdrawProof, err := s.api.GetWithdrawProof(commitmentID, transactionHash)
	s.NoError(err)
	s.Equal(expected, *withdrawProof)
}

func TestGetWithdrawProofTestSuite(t *testing.T) {
	suite.Run(t, new(GetWithdrawProofTestSuite))
}
