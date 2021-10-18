package storage

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var (
	txCommitment = models.TxCommitment{
		CommitmentBase: models.CommitmentBase{
			ID: models.CommitmentID{
				BatchID:      models.MakeUint256(1),
				IndexInBatch: 0,
			},
			Type:          batchtype.Transfer,
			PostStateRoot: utils.RandomHash(),
		},
		FeeReceiver:       uint32(1),
		CombinedSignature: models.MakeRandomSignature(),
		Transactions:      []byte{1, 2, 3},
	}
)

type TxCommitmentTestSuite struct {
	*require.Assertions
	suite.Suite
	storage *TestStorage
}

func (s *TxCommitmentTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *TxCommitmentTestSuite) SetupTest() {
	var err error
	s.storage, err = NewTestStorage()
	s.NoError(err)
}

func (s *TxCommitmentTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *TxCommitmentTestSuite) TestAddTxCommitment_AddAndRetrieve() {
	err := s.storage.AddTxCommitment(&txCommitment)
	s.NoError(err)

	actual, err := s.storage.GetTxCommitment(&txCommitment.ID)
	s.NoError(err)
	s.Equal(txCommitment, *actual)
}

func (s *TxCommitmentTestSuite) addRandomBatch() models.Uint256 {
	batch := models.Batch{
		ID:                models.MakeUint256(123),
		Type:              batchtype.Transfer,
		TransactionHash:   utils.RandomHash(),
		Hash:              utils.NewRandomHash(),
		FinalisationBlock: ref.Uint32(1234),
	}
	err := s.storage.AddBatch(&batch)
	s.NoError(err)
	return batch.ID
}

func (s *CommitmentTestSuite) TestGetCommitment_NonexistentCommitment() {
	res, err := s.storage.GetTxCommitment(&txCommitment.ID)
	s.ErrorIs(err, NewNotFoundError("commitment"))
	s.Nil(res)
}

func (s *TxCommitmentTestSuite) TestGetTxCommitmentsByBatchID() {
	err := s.storage.AddTxCommitment(&txCommitment)
	s.NoError(err)

	batchID := s.addRandomBatch()
	includedCommitment := txCommitment
	includedCommitment.ID.BatchID = batchID
	includedCommitment.FeeReceiver = 0

	expectedCommitments := make([]models.CommitmentWithTokenID, 2)
	for i := 0; i < 2; i++ {
		includedCommitment.ID.IndexInBatch = uint8(i)
		err = s.storage.AddTxCommitment(&includedCommitment)
		s.NoError(err)

		expectedCommitments[i] = models.CommitmentWithTokenID{
			ID:                 includedCommitment.ID,
			Transactions:       includedCommitment.Transactions,
			TokenID:            models.MakeUint256(1),
			FeeReceiverStateID: includedCommitment.FeeReceiver,
			CombinedSignature:  includedCommitment.CombinedSignature,
			PostStateRoot:      includedCommitment.PostStateRoot,
		}
	}

	s.addStateLeaf()

	commitments, err := s.storage.GetTxCommitmentsByBatchID(batchID)
	s.NoError(err)
	s.Len(commitments, 2)
	s.Contains(commitments, expectedCommitments[0])
	s.Contains(commitments, expectedCommitments[1])
}

func (s *CommitmentTestSuite) TestGetCommitmentsByBatchID_NonexistentCommitments() {
	batchID := s.addRandomBatch()
	commitments, err := s.storage.GetTxCommitmentsByBatchID(batchID)
	s.NoError(err)
	s.Len(commitments, 0)
}

func (s *TxCommitmentTestSuite) addStateLeaf() {
	_, err := s.storage.StateTree.Set(uint32(0), &models.UserState{
		PubKeyID: 1,
		TokenID:  models.MakeUint256(1),
		Balance:  models.MakeUint256(420),
		Nonce:    models.MakeUint256(0),
	})
	s.NoError(err)
}

func TestTxCommitmentTestSuite(t *testing.T) {
	suite.Run(t, new(TxCommitmentTestSuite))
}
