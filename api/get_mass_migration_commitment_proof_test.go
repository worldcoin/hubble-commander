package api

import (
	"testing"
	"time"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/encoder"
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

type GetMassMigrationCommitmentProofTestSuite struct {
	*require.Assertions
	suite.Suite
	api            *API
	storage        *st.TestStorage
	batch          *models.Batch
	commitments    []models.CommitmentWithTxs
	massMigrations []models.MassMigration
}

func (s *GetMassMigrationCommitmentProofTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *GetMassMigrationCommitmentProofTestSuite) SetupTest() {
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
			models.CommitmentID{
				BatchID:      models.MakeUint256(1),
				IndexInBatch: 0,
			},
		),
		makeMassMigration(
			common.Hash{1, 2, 3},
			0,
			1,
			models.NewTimestamp(time.Unix(150, 0).UTC()),
			models.CommitmentID{
				BatchID:      models.MakeUint256(1),
				IndexInBatch: 0,
			},
		),
		makeMassMigration(
			common.Hash{3, 4, 5},
			0,
			2,
			models.NewTimestamp(time.Unix(160, 0).UTC()),
			models.CommitmentID{
				BatchID:      models.MakeUint256(1),
				IndexInBatch: 1,
			},
		),
	}

	_, err = s.storage.StateTree.Set(0, &models.UserState{
		Balance: models.MakeUint256(100),
		TokenID: models.MakeUint256(10),
		Nonce:   models.MakeUint256(2),
	})
	s.NoError(err)

	stateRoot1, err := s.storage.StateTree.Root()
	s.NoError(err)

	_, err = s.storage.StateTree.Set(0, &models.UserState{
		Balance: models.MakeUint256(0),
		TokenID: models.MakeUint256(10),
		Nonce:   models.MakeUint256(3),
	})
	s.NoError(err)

	stateRoot2, err := s.storage.StateTree.Root()
	s.NoError(err)

	err = s.storage.BatchAddMassMigration(s.massMigrations)
	s.NoError(err)

	accountTreeRoot := utils.RandomHash()
	s.batch = &models.Batch{
		ID:                models.MakeUint256(1),
		Type:              batchtype.MassMigration,
		TransactionHash:   utils.RandomHash(),
		Hash:              ref.Hash(utils.RandomHash()),
		FinalisationBlock: ref.Uint32(10),
		AccountTreeRoot:   &accountTreeRoot,
	}

	err = s.storage.AddBatch(s.batch)
	s.NoError(err)

	s.commitments = []models.CommitmentWithTxs{
		makeMassMigrationCommitment(
			s.Assertions,
			models.CommitmentID{
				BatchID:      models.MakeUint256(1),
				IndexInBatch: 0,
			},
			*stateRoot1,
			accountTreeRoot,
			[]models.MassMigration{s.massMigrations[0], s.massMigrations[1]},
		),
		makeMassMigrationCommitment(
			s.Assertions,
			models.CommitmentID{
				BatchID:      models.MakeUint256(1),
				IndexInBatch: 1,
			},
			*stateRoot2,
			accountTreeRoot,
			[]models.MassMigration{s.massMigrations[2]},
		),
	}

	err = s.storage.AddTxCommitment(&s.commitments[0].TxCommitment)
	s.NoError(err)
	err = s.storage.AddTxCommitment(&s.commitments[1].TxCommitment)
	s.NoError(err)
}

func (s *GetMassMigrationCommitmentProofTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *GetMassMigrationCommitmentProofTestSuite) TestGetMassMigrationCommitmentProof_FirstCommitment() {
	s.testGetMassMigrationCommitmentProofEndpoint(0, s.massMigrations[:2])
}

func (s *GetMassMigrationCommitmentProofTestSuite) TestGetMassMigrationCommitmentProof_SecondCommitment() {
	s.testGetMassMigrationCommitmentProofEndpoint(1, s.massMigrations[2:])
}

func (s *GetMassMigrationCommitmentProofTestSuite) TestGetMassMigrationCommitmentProof_NonexistentBatch() {
	_, err := s.api.GetMassMigrationCommitmentProof(models.MakeUint256(10), 15)
	s.Equal(&APIError{
		Code:    50004,
		Message: "mass migration commitment inclusion proof not found",
	}, err)
}

func (s *GetMassMigrationCommitmentProofTestSuite) testGetMassMigrationCommitmentProofEndpoint(
	commitmentIndex int,
	massMigrations []models.MassMigration,
) {
	withdrawTree, meta := s.prepareWithdrawTreeAndMeta(commitmentIndex, massMigrations)

	witnessIndex := 0
	if commitmentIndex == 0 {
		witnessIndex = 1
	}

	expected := dto.MassMigrationCommitmentProof{
		CommitmentInclusionProofBase: dto.CommitmentInclusionProofBase{
			StateRoot: s.commitments[commitmentIndex].PostStateRoot,
			Path: &dto.MerklePath{
				Path:  uint32(commitmentIndex),
				Depth: 2,
			},
			Witness: []common.Hash{s.commitments[witnessIndex].LeafHash()},
		},
		Body: &dto.MassMigrationBody{
			AccountRoot:  *s.batch.AccountTreeRoot,
			Signature:    s.commitments[commitmentIndex].CombinedSignature,
			Meta:         meta,
			WithdrawRoot: withdrawTree.Root(),
			Transactions: s.commitments[commitmentIndex].Transactions,
		},
	}

	commitmentInclusionProof, err := s.api.GetMassMigrationCommitmentProof(s.batch.ID, uint8(commitmentIndex))
	s.NoError(err)
	s.Equal(expected, *commitmentInclusionProof)
}

func (s *GetMassMigrationCommitmentProofTestSuite) prepareWithdrawTreeAndMeta(
	commitmentIndex int,
	massMigrations []models.MassMigration,
) (*merkletree.MerkleTree, *dto.MassMigrationMeta) {
	hashes := make([]common.Hash, 0, len(massMigrations))
	meta := &dto.MassMigrationMeta{
		FeeReceiver: s.commitments[commitmentIndex].FeeReceiver,
	}

	for i := range massMigrations {
		senderLeaf, err := s.storage.StateTree.Leaf(s.massMigrations[i].FromStateID)
		s.NoError(err)

		if i == 0 {
			meta.SpokeID = massMigrations[i].SpokeID
			meta.TokenID = senderLeaf.TokenID
		}

		meta.Amount = *meta.Amount.Add(&massMigrations[i].Amount)

		hash, err := encoder.HashUserState(&models.UserState{
			PubKeyID: senderLeaf.PubKeyID,
			TokenID:  senderLeaf.TokenID,
			Balance:  s.massMigrations[i].Amount,
			Nonce:    models.MakeUint256(0),
		})
		s.NoError(err)

		hashes = append(hashes, *hash)
	}

	withdrawTree, err := merkletree.NewMerkleTree(hashes)
	s.NoError(err)

	return withdrawTree, meta
}

func makeMassMigration(
	hash common.Hash,
	from uint32,
	nonce uint64,
	receiveTime *models.Timestamp,
	commitmentID models.CommitmentID,
) models.MassMigration {
	return models.MassMigration{
		TransactionBase: models.TransactionBase{
			Hash:         hash,
			TxType:       txtype.MassMigration,
			FromStateID:  from,
			Amount:       models.MakeUint256(90),
			Fee:          models.MakeUint256(10),
			Nonce:        models.MakeUint256(nonce),
			Signature:    models.MakeRandomSignature(),
			ReceiveTime:  receiveTime,
			CommitmentID: &commitmentID,
			ErrorMessage: nil,
		},
		SpokeID: 1,
	}
}

func makeMassMigrationCommitment(
	assertion *require.Assertions,
	commitmentID models.CommitmentID,
	stateRoot common.Hash,
	accountRoot common.Hash,
	massMigrations []models.MassMigration,
) models.CommitmentWithTxs {
	serializedMassMigrations, err := encoder.SerializeMassMigrations(massMigrations)
	assertion.NoError(err)

	massMigrationCommitment := models.CommitmentWithTxs{
		TxCommitment: models.TxCommitment{
			CommitmentBase: models.CommitmentBase{
				ID:            commitmentID,
				Type:          batchtype.MassMigration,
				PostStateRoot: stateRoot,
			},
			FeeReceiver:       1,
			CombinedSignature: models.MakeRandomSignature(),
			BodyHash:          nil,
		},
		Transactions: serializedMassMigrations,
	}

	massMigrationCommitment.BodyHash = massMigrationCommitment.CalcBodyHash(accountRoot)

	return massMigrationCommitment
}

func TestGetMassMigrationCommitmentProofTestSuite(t *testing.T) {
	suite.Run(t, new(GetMassMigrationCommitmentProofTestSuite))
}
