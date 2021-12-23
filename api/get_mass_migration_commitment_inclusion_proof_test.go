package api

import (
	"testing"
	"time"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/models"
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

type GetMMCommitmentInclusionProofTestSuite struct {
	*require.Assertions
	suite.Suite
	api            *API
	storage        *st.TestStorage
	batch          *models.Batch
	commitments    []models.CommitmentWithTxs
	massMigrations []models.MassMigration
}

func (s *GetMMCommitmentInclusionProofTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *GetMMCommitmentInclusionProofTestSuite) SetupTest() {
	var err error
	s.storage, err = st.NewTestStorage()
	s.NoError(err)
	s.api = &API{
		storage: s.storage.Storage,
		cfg:     &config.APIConfig{EnableProofMethods: true},
	}

	s.massMigrations = []models.MassMigration{
		{
			TransactionBase: models.TransactionBase{
				Hash:        utils.RandomHash(),
				TxType:      txtype.MassMigration,
				FromStateID: 0,
				Amount:      models.MakeUint256(90),
				Fee:         models.MakeUint256(10),
				Nonce:       models.MakeUint256(0),
				Signature:   models.MakeRandomSignature(),
				ReceiveTime: models.NewTimestamp(time.Unix(140, 0).UTC()),
				CommitmentID: &models.CommitmentID{
					BatchID:      models.MakeUint256(1),
					IndexInBatch: 0,
				},
				ErrorMessage: nil,
			},
			SpokeID: 1,
		},
		{
			TransactionBase: models.TransactionBase{
				Hash:        utils.RandomHash(),
				TxType:      txtype.MassMigration,
				FromStateID: 0,
				Amount:      models.MakeUint256(90),
				Fee:         models.MakeUint256(10),
				Nonce:       models.MakeUint256(1),
				Signature:   models.MakeRandomSignature(),
				ReceiveTime: models.NewTimestamp(time.Unix(150, 0).UTC()),
				CommitmentID: &models.CommitmentID{
					BatchID:      models.MakeUint256(1),
					IndexInBatch: 1,
				},
				ErrorMessage: nil,
			},
			SpokeID: 1,
		},
	}

	_, err = s.storage.StateTree.Set(0, &models.UserState{
		Balance: models.MakeUint256(100),
		TokenID: models.MakeUint256(10),
		Nonce:   models.MakeUint256(1),
	})
	s.NoError(err)

	stateRoot1, err := s.storage.StateTree.Root()
	s.NoError(err)

	_, err = s.storage.StateTree.Set(0, &models.UserState{
		Balance: models.MakeUint256(0),
		TokenID: models.MakeUint256(10),
		Nonce:   models.MakeUint256(2),
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

	serializedMassMigrations1, err := encoder.SerializeMassMigrations([]models.MassMigration{s.massMigrations[0]})
	s.NoError(err)
	serializedMassMigrations2, err := encoder.SerializeMassMigrations([]models.MassMigration{s.massMigrations[1]})
	s.NoError(err)

	s.commitments = []models.CommitmentWithTxs{
		{
			TxCommitment: models.TxCommitment{
				CommitmentBase: models.CommitmentBase{
					ID: models.CommitmentID{
						BatchID:      models.MakeUint256(1),
						IndexInBatch: 0,
					},
					Type:          batchtype.MassMigration,
					PostStateRoot: *stateRoot1,
				},
				FeeReceiver:       1,
				CombinedSignature: models.MakeRandomSignature(),
				BodyHash:          nil,
			},
			Transactions: serializedMassMigrations1,
		},
		{
			TxCommitment: models.TxCommitment{
				CommitmentBase: models.CommitmentBase{
					ID: models.CommitmentID{
						BatchID:      models.MakeUint256(1),
						IndexInBatch: 1,
					},
					Type:          batchtype.MassMigration,
					PostStateRoot: *stateRoot2,
				},
				FeeReceiver:       1,
				CombinedSignature: models.MakeRandomSignature(),
				BodyHash:          nil,
			},
			Transactions: serializedMassMigrations2,
		},
	}

	s.commitments[0].BodyHash = s.commitments[0].CalcBodyHash(accountTreeRoot)
	s.commitments[1].BodyHash = s.commitments[1].CalcBodyHash(accountTreeRoot)

	err = s.storage.AddTxCommitment(&s.commitments[0].TxCommitment)
	s.NoError(err)
	err = s.storage.AddTxCommitment(&s.commitments[1].TxCommitment)
	s.NoError(err)
}

func (s *GetMMCommitmentInclusionProofTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *GetMMCommitmentInclusionProofTestSuite) TestGetMassMigrationCommitmentInclusionProof_FirstCommitment() {
	s.testGetMassMigrationCommitmentInclusionProofEndpoint(0)
}

func (s *GetMMCommitmentInclusionProofTestSuite) TestGetMassMigrationCommitmentInclusionProof_SecondCommitment() {
	s.testGetMassMigrationCommitmentInclusionProofEndpoint(1)
}

func (s *GetMMCommitmentInclusionProofTestSuite) TestGetMassMigrationCommitmentInclusionProof_NonexistentBatch() {
	_, err := s.api.GetMassMigrationCommitmentInclusionProof(models.MakeUint256(10), 15)
	s.Equal(&APIError{
		Code:    50004,
		Message: "mass migration commitment inclusion proof not found",
	}, err)
}

func (s *GetMMCommitmentInclusionProofTestSuite) testGetMassMigrationCommitmentInclusionProofEndpoint(commitmentIndex int) {
	senderLeaf, err := s.storage.StateTree.Leaf(s.massMigrations[commitmentIndex].FromStateID)
	s.NoError(err)

	hash, err := encoder.HashUserState(&models.UserState{
		PubKeyID: senderLeaf.PubKeyID,
		TokenID:  senderLeaf.TokenID,
		Balance:  s.massMigrations[commitmentIndex].Amount,
		Nonce:    models.MakeUint256(0),
	})
	s.NoError(err)

	withdrawTree, err := merkletree.NewMerkleTree([]common.Hash{*hash})
	s.NoError(err)

	witnessIndex := 0
	if commitmentIndex == 0 {
		witnessIndex = 1
	}

	expected := models.MMCommitmentInclusionProof{
		CommitmentInclusionProofBase: models.CommitmentInclusionProofBase{
			StateRoot: s.commitments[commitmentIndex].PostStateRoot,
			Path: &models.MerklePath{
				Path:  uint32(commitmentIndex),
				Depth: 2,
			},
			Witness: []common.Hash{s.commitments[witnessIndex].LeafHash()},
		},
		Body: &models.MMBody{
			AccountRoot: *s.batch.AccountTreeRoot,
			Signature:   s.commitments[commitmentIndex].CombinedSignature,
			Meta: &models.MassMigrationMeta{
				SpokeID:     s.massMigrations[commitmentIndex].SpokeID,
				TokenID:     senderLeaf.TokenID,
				Amount:      s.massMigrations[commitmentIndex].Amount,
				FeeReceiver: s.commitments[commitmentIndex].FeeReceiver,
			},
			WithdrawRoot: withdrawTree.Root(),
			Transactions: s.commitments[commitmentIndex].Transactions,
		},
	}

	commitmentInclusionProof, err := s.api.GetMassMigrationCommitmentInclusionProof(s.batch.ID, uint8(commitmentIndex))
	s.NoError(err)
	s.Equal(expected, *commitmentInclusionProof)
}

func TestGetMMCommitmentInclusionProofTestSuite(t *testing.T) {
	suite.Run(t, new(GetMMCommitmentInclusionProofTestSuite))
}
