package commander

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ValidateStateRootTestSuite struct {
	*require.Assertions
	suite.Suite
	storage *st.TestStorage
}

func (s *ValidateStateRootTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *ValidateStateRootTestSuite) SetupTest() {
	var err error
	s.storage, err = st.NewTestStorage()
	s.NoError(err)
}

func (s *ValidateStateRootTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *ValidateStateRootTestSuite) TestValidateStateRoot_SameStateRootHash() {
	root, err := s.storage.StateTree.Root()
	s.NoError(err)

	commitment := models.TxCommitment{
		CommitmentBase: models.CommitmentBase{
			Type:          batchtype.Transfer,
			PostStateRoot: *root,
		},
		FeeReceiver:       0,
		CombinedSignature: models.Signature{},
	}
	err = s.storage.AddCommitment(&commitment)
	s.NoError(err)

	err = validateStateRoot(s.storage.Storage)
	s.NoError(err)
}

func (s *ValidateStateRootTestSuite) TestValidateStateRoot_DifferentStateRootHash() {
	commitment := models.TxCommitment{
		CommitmentBase: models.CommitmentBase{
			Type:          batchtype.Transfer,
			PostStateRoot: common.Hash{1, 2, 3},
		},
		FeeReceiver:       0,
		CombinedSignature: models.Signature{},
	}
	err := s.storage.AddCommitment(&commitment)
	s.NoError(err)

	err = validateStateRoot(s.storage.Storage)
	s.ErrorIs(err, ErrInvalidStateRoot)
}

func (s *ValidateStateRootTestSuite) TestValidateStateRoot_FirstCommitment() {
	err := validateStateRoot(s.storage.Storage)
	s.NoError(err)
}

func TestValidateStateRootTestSuite(t *testing.T) {
	suite.Run(t, new(ValidateStateRootTestSuite))
}
