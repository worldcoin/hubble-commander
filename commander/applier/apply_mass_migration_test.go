package applier

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ApplyMassMigrationTestSuite struct {
	*require.Assertions
	suite.Suite
	storage       *st.TestStorage
	applier       *Applier
	massMigration models.MassMigration
	tokenID       models.Uint256
}

func (s *ApplyMassMigrationTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
	s.massMigration = models.MassMigration{
		TransactionBase: models.TransactionBase{
			FromStateID: 1,
			Amount:      models.MakeUint256(100),
			Fee:         models.MakeUint256(10),
			Nonce:       models.MakeUint256(0),
		},
		SpokeID: 2,
	}
}

func (s *ApplyMassMigrationTestSuite) SetupTest() {
	var err error
	s.storage, err = st.NewTestStorage()
	s.NoError(err)
	s.applier = NewApplier(s.storage.Storage)

	_, err = s.storage.StateTree.Set(senderState.PubKeyID, &senderState)
	s.NoError(err)

	s.tokenID = senderState.TokenID
}

func (s *ApplyMassMigrationTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *ApplyMassMigrationTestSuite) TestApplyMassMigration() {
	_, txError, appError := s.applier.ApplyMassMigration(&s.massMigration, s.tokenID)
	s.NoError(txError)
	s.NoError(appError)

	senderLeaf, err := s.storage.StateTree.Leaf(1)
	s.NoError(err)

	s.Equal(uint64(290), senderLeaf.Balance.Uint64())
	s.Equal(*s.massMigration.Nonce.AddN(1), senderLeaf.Nonce)
}

func (s *ApplyMassMigrationTestSuite) TestApplyMassMigration_ValidatesSenderTokenID() {
	setUserStatesInTree(s.Assertions, s.storage)

	_, txError, appError := s.applier.ApplyMassMigration(&s.massMigration, models.MakeUint256(3))
	s.NoError(txError)
	s.ErrorIs(appError, ErrInvalidSenderTokenID)
}

func (s *ApplyMassMigrationTestSuite) TestApplyMassMigration_ValidatesNonce() {
	massMigrationWithBadNonce := s.massMigration
	massMigrationWithBadNonce.Nonce = models.MakeUint256(1)
	setUserStatesInTree(s.Assertions, s.storage)

	_, txError, appError := s.applier.ApplyMassMigration(&massMigrationWithBadNonce, models.MakeUint256(1))
	s.ErrorIs(txError, ErrNonceTooHigh)
	s.NoError(appError)
}

func (s *ApplyMassMigrationTestSuite) TestApplyMassMigrationForSync_ValidatesNonce() {
	setUserStatesInTree(s.Assertions, s.storage)

	bigMassMigration := s.massMigration
	bigMassMigration.Amount = models.MakeUint256(1_000_000)

	synced, txError, appError := s.applier.ApplyMassMigrationForSync(&bigMassMigration, models.MakeUint256(1))
	s.NotNil(synced)
	s.ErrorIs(txError, ErrBalanceTooLow)
	s.NoError(appError)

	s.Equal(&bigMassMigration, synced.Tx.ToMassMigration())
	s.Equal(senderState, *synced.SenderStateProof.UserState)
	s.Len(synced.SenderStateProof.Witness, st.StateTreeDepth)
}

func (s *ApplyMassMigrationTestSuite) TestApplyTransferForSync_ValidatesSenderTokenID() {
	setUserStatesInTree(s.Assertions, s.storage)

	synced, txError, appError := s.applier.ApplyMassMigrationForSync(&s.massMigration, models.MakeUint256(3))
	s.NotNil(synced)
	s.ErrorIs(txError, ErrInvalidSenderTokenID)
	s.NoError(appError)

	s.Equal(&s.massMigration, synced.Tx.ToMassMigration())
	s.Equal(senderState, *synced.SenderStateProof.UserState)
	s.Len(synced.SenderStateProof.Witness, st.StateTreeDepth)
}

func (s *ApplyMassMigrationTestSuite) TestApplyTransferForSync_ReturnsTransferWithUpdatedNonce() {
	setUserStatesInTree(s.Assertions, s.storage)

	transferWithModifiedNonce := s.massMigration
	transferWithModifiedNonce.Nonce = models.MakeUint256(1234)

	synced, txError, appError := s.applier.ApplyMassMigrationForSync(&transferWithModifiedNonce, models.MakeUint256(1))
	s.NoError(appError)
	s.NoError(txError)

	s.Equal(models.MakeUint256(1234), transferWithModifiedNonce.Nonce)
	s.Equal(models.MakeUint256(0), synced.Tx.ToMassMigration().GetNonce())
}

func TestApplyMassMigrationTestSuite(t *testing.T) {
	suite.Run(t, new(ApplyMassMigrationTestSuite))
}
