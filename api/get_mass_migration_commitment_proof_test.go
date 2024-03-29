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
	commitments    []models.MMCommitmentWithTxs
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
	s.massMigrations = s.generateMassMigrationForTestSuite()

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

	err = s.storage.BatchAddTransaction(models.MakeMassMigrationArray(s.massMigrations...))
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

	s.commitments = []models.MMCommitmentWithTxs{
		makeMassMigrationCommitment(
			s.Assertions,
			s.storage,
			models.CommitmentID{
				BatchID:      s.batch.ID,
				IndexInBatch: 0,
			},
			1,
			*stateRoot1,
			accountTreeRoot,
			[]models.MassMigration{s.massMigrations[0], s.massMigrations[1]},
		),
		makeMassMigrationCommitment(
			s.Assertions,
			s.storage,
			models.CommitmentID{
				BatchID:      s.batch.ID,
				IndexInBatch: 1,
			},
			1,
			*stateRoot2,
			accountTreeRoot,
			[]models.MassMigration{s.massMigrations[2]},
		),
	}

	err = s.storage.AddCommitment(s.commitments[0].ToMMCommitmentWithTxs())
	s.NoError(err)
	err = s.storage.AddCommitment(s.commitments[1].ToMMCommitmentWithTxs())
	s.NoError(err)
}

func (s *GetMassMigrationCommitmentProofTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *GetMassMigrationCommitmentProofTestSuite) TestGetMassMigrationCommitmentProof_FirstCommitment() {
	s.testGetMassMigrationCommitmentProofEndpoint(0, 1, s.massMigrations[:2])
}

func (s *GetMassMigrationCommitmentProofTestSuite) TestGetMassMigrationCommitmentProof_SecondCommitment() {
	s.testGetMassMigrationCommitmentProofEndpoint(1, 0, s.massMigrations[2:])
}

func (s *GetMassMigrationCommitmentProofTestSuite) TestGetMassMigrationCommitmentProof_NonexistentBatch() {
	commitmentID := models.CommitmentID{
		BatchID:      models.MakeUint256(10),
		IndexInBatch: 15,
	}
	_, err := s.api.GetMassMigrationCommitmentProof(commitmentID)
	s.Equal(APIErrCannotGenerateMMCommitmentProof, err)
}

func (s *GetMassMigrationCommitmentProofTestSuite) TestGetMassMigrationCommitmentProof_NotMassMigrationCommitment() {
	err := s.storage.AddBatch(&models.Batch{
		ID:   models.MakeUint256(2),
		Type: batchtype.Transfer,
	})
	s.NoError(err)

	commitmentID := models.CommitmentID{
		BatchID:      models.MakeUint256(2),
		IndexInBatch: 0,
	}
	_, err = s.api.GetMassMigrationCommitmentProof(commitmentID)
	s.Equal(APIErrOnlyMassMigrationCommitmentsForProofing, err)
}

func (s *GetMassMigrationCommitmentProofTestSuite) testGetMassMigrationCommitmentProofEndpoint(
	commitmentIndex int,
	witnessIndex int,
	massMigrations []models.MassMigration,
) {
	withdrawTree, meta, err := prepareWithdrawTreeAndMeta(s.storage, s.commitments[commitmentIndex].Meta.FeeReceiver, massMigrations)
	s.NoError(err)

	expected := dto.MassMigrationCommitmentProof{
		CommitmentInclusionProofBase: dto.CommitmentInclusionProofBase{
			StateRoot: s.commitments[commitmentIndex].ToMMCommitmentWithTxs().PostStateRoot,
			Path: &dto.MerklePath{
				Path:  uint32(commitmentIndex),
				Depth: 2,
			},
			Witness: []common.Hash{s.commitments[witnessIndex].ToMMCommitmentWithTxs().LeafHash()},
		},
		Body: &dto.MassMigrationBody{
			AccountRoot:  *s.batch.AccountTreeRoot,
			Signature:    s.commitments[commitmentIndex].ToMMCommitmentWithTxs().CombinedSignature,
			Meta:         dto.NewMassMigrationMeta(meta),
			WithdrawRoot: withdrawTree.Root(),
			Transactions: s.commitments[commitmentIndex].ToMMCommitmentWithTxs().Transactions,
		},
	}

	commitmentID := models.CommitmentID{
		BatchID:      s.batch.ID,
		IndexInBatch: uint8(commitmentIndex),
	}

	commitmentInclusionProof, err := s.api.GetMassMigrationCommitmentProof(commitmentID)
	s.NoError(err)
	s.Equal(expected, *commitmentInclusionProof)
}

func prepareWithdrawTreeAndMeta(
	storage *st.TestStorage,
	feeReceiver uint32,
	massMigrations []models.MassMigration,
) (*merkletree.MerkleTree, *models.MassMigrationMeta, error) {
	hashes := make([]common.Hash, 0, len(massMigrations))
	meta := &models.MassMigrationMeta{
		FeeReceiver: feeReceiver,
	}

	for i := range massMigrations {
		senderLeaf, err := storage.StateTree.Leaf(massMigrations[i].FromStateID)
		if err != nil {
			return nil, nil, err
		}

		if i == 0 {
			meta.SpokeID = massMigrations[i].SpokeID
			meta.TokenID = senderLeaf.TokenID
		}

		meta.Amount = *meta.Amount.Add(&massMigrations[i].Amount)

		hash, err := encoder.HashUserState(&models.UserState{
			PubKeyID: senderLeaf.PubKeyID,
			TokenID:  senderLeaf.TokenID,
			Balance:  massMigrations[i].Amount,
			Nonce:    models.MakeUint256(0),
		})
		if err != nil {
			return nil, nil, err
		}

		hashes = append(hashes, *hash)
	}

	withdrawTree, err := merkletree.NewMerkleTree(hashes)
	if err != nil {
		return nil, nil, err
	}

	return withdrawTree, meta, nil
}

func (s *GetMassMigrationCommitmentProofTestSuite) generateMassMigrationForTestSuite() []models.MassMigration {
	return []models.MassMigration{
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
			0,
			1,
			models.NewTimestamp(time.Unix(150, 0).UTC()),
			models.CommitmentSlot{
				BatchID:           models.MakeUint256(1),
				IndexInBatch:      0,
				IndexInCommitment: 1,
			},
		),
		makeMassMigration(
			common.Hash{3, 4, 5},
			0,
			2,
			models.NewTimestamp(time.Unix(160, 0).UTC()),
			models.CommitmentSlot{
				BatchID:           models.MakeUint256(1),
				IndexInBatch:      1,
				IndexInCommitment: 2,
			},
		),
	}
}

func makeMassMigration(
	hash common.Hash,
	from uint32,
	nonce uint64,
	receiveTime *models.Timestamp,
	commitmentSlot models.CommitmentSlot,
) models.MassMigration {
	return models.MassMigration{
		TransactionBase: models.TransactionBase{
			Hash:           hash,
			TxType:         txtype.MassMigration,
			FromStateID:    from,
			Amount:         models.MakeUint256(90),
			Fee:            models.MakeUint256(10),
			Nonce:          models.MakeUint256(nonce),
			Signature:      models.MakeRandomSignature(),
			ReceiveTime:    receiveTime,
			CommitmentSlot: &commitmentSlot,
			ErrorMessage:   nil,
		},
		SpokeID: 1,
	}
}

func makeMassMigrationCommitment(
	assertion *require.Assertions,
	storage *st.TestStorage,
	commitmentID models.CommitmentID,
	feeReceiver uint32,
	stateRoot common.Hash,
	accountRoot common.Hash,
	massMigrations []models.MassMigration,
) models.MMCommitmentWithTxs {
	serializedMassMigrations, err := encoder.SerializeMassMigrations(massMigrations)
	assertion.NoError(err)

	withdrawTree, meta, err := prepareWithdrawTreeAndMeta(storage, feeReceiver, massMigrations)
	assertion.NoError(err)

	massMigrationCommitment := models.MMCommitmentWithTxs{
		MMCommitment: models.MMCommitment{
			CommitmentBase: models.CommitmentBase{
				ID:            commitmentID,
				Type:          batchtype.MassMigration,
				PostStateRoot: stateRoot,
			},
			CombinedSignature: models.MakeRandomSignature(),
			Meta:              meta,
			WithdrawRoot:      withdrawTree.Root(),
		},
		Transactions: serializedMassMigrations,
	}

	massMigrationCommitment.CalcAndSetBodyHash(accountRoot)
	return massMigrationCommitment
}

func TestGetMassMigrationCommitmentProofTestSuite(t *testing.T) {
	suite.Run(t, new(GetMassMigrationCommitmentProofTestSuite))
}
