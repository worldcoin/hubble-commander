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
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type GetWithdrawProofTestSuite struct {
	*require.Assertions
	suite.Suite
	api            *API
	storage        *st.TestStorage
	commitment     models.CommitmentWithTxs
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
		{
			TransactionBase: models.TransactionBase{
				Hash:        common.Hash{2, 3, 4},
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
				Hash:        common.Hash{1, 2, 3},
				TxType:      txtype.MassMigration,
				FromStateID: 1,
				Amount:      models.MakeUint256(90),
				Fee:         models.MakeUint256(10),
				Nonce:       models.MakeUint256(1),
				Signature:   models.MakeRandomSignature(),
				ReceiveTime: models.NewTimestamp(time.Unix(150, 0).UTC()),
				CommitmentID: &models.CommitmentID{
					BatchID:      models.MakeUint256(1),
					IndexInBatch: 0,
				},
				ErrorMessage: nil,
			},
			SpokeID: 1,
		},
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

	err = s.storage.BatchAddMassMigration(s.massMigrations)
	s.NoError(err)

	serializedMassMigrations, err := encoder.SerializeMassMigrations(s.massMigrations)
	s.NoError(err)

	s.commitment = models.CommitmentWithTxs{
		TxCommitment: models.TxCommitment{
			CommitmentBase: models.CommitmentBase{
				ID: models.CommitmentID{
					BatchID:      models.MakeUint256(1),
					IndexInBatch: 0,
				},
				Type:          batchtype.MassMigration,
				PostStateRoot: *stateRoot,
			},
			FeeReceiver:       1,
			CombinedSignature: models.MakeRandomSignature(),
			BodyHash:          nil,
		},
		Transactions: serializedMassMigrations,
	}

	accountTreeRoot := utils.RandomHash()
	s.commitment.BodyHash = s.commitment.CalcBodyHash(accountTreeRoot)

	err = s.storage.AddTxCommitment(&s.commitment.TxCommitment)
	s.NoError(err)

	err = s.storage.AddBatch(&models.Batch{
		ID: models.MakeUint256(1),
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
	_, err := s.api.GetWithdrawProof(models.MakeUint256(10), 15, utils.RandomHash())
	s.Equal(APIWithdrawProofCouldNotBeCalculated, err)
}

func (s *GetWithdrawProofTestSuite) TestGetWithdrawProof_NonexistentMassMigrationWithGivenSenderInCommitment() {
	_, err := s.api.GetWithdrawProof(models.MakeUint256(1), 0, utils.RandomHash())
	s.Equal(APIErrMassMigrationWithSenderNotFound, err)
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

	expected := dto.WithdrawProof{
		WithdrawProof: models.WithdrawProof{
			UserState: targetUserState,
			Path: models.MerklePath{
				Path:  uint32(massMigrationIndex),
				Depth: withdrawTree.Depth(),
			},
			Witness: withdrawTree.GetWitness(uint32(massMigrationIndex)),
			Root:    withdrawTree.Root(),
		},
	}

	withdrawProof, err := s.api.GetWithdrawProof(models.MakeUint256(1), 0, transactionHash)
	s.NoError(err)
	s.Equal(expected, *withdrawProof)
}

func TestGetWithdrawProofTestSuite(t *testing.T) {
	suite.Run(t, new(GetWithdrawProofTestSuite))
}
