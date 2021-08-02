package commander

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type RollupTestSuite struct {
	*require.Assertions
	suite.Suite
	storage   *st.TestStorage
	stateTree *st.StoredMerkleTree
}

func (s *RollupTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *RollupTestSuite) SetupTest() {
	var err error
	s.storage, err = st.NewTestStorageWithBadger()
	s.NoError(err)
	s.stateTree = st.NewStoredMerkleTree("state", s.storage.Database) // Must be the same state tree as in Storage object
}

func (s *RollupTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *RollupTestSuite) TestValidateStateRoot_SameStateRootHash() {
	commitment := models.Commitment{
		Type:              txtype.Transfer,
		Transactions:      []byte{1, 2, 3},
		FeeReceiver:       0,
		CombinedSignature: models.Signature{},
		PostStateRoot:     common.Hash{1, 2, 3},
	}
	_, err := s.storage.AddCommitment(&commitment)
	s.NoError(err)

	_, _, err = s.stateTree.SetNode(&models.MerklePath{Path: 0, Depth: 0}, commitment.PostStateRoot)
	s.NoError(err)

	err = validateStateRoot(s.storage.Storage)
	s.NoError(err)
}

func (s *RollupTestSuite) TestValidateStateRoot_DifferentStateRootHash() {
	commitment := models.Commitment{
		Type:              txtype.Transfer,
		Transactions:      []byte{1, 2, 3},
		FeeReceiver:       0,
		CombinedSignature: models.Signature{},
		PostStateRoot:     common.Hash{1, 2, 3},
	}
	_, err := s.storage.AddCommitment(&commitment)
	s.NoError(err)

	err = validateStateRoot(s.storage.Storage)
	s.Equal(ErrInvalidStateRoot, err)
}

func (s *RollupTestSuite) TestValidateStateRoot_FirstCommitment() {
	_, _, err := s.stateTree.SetNode(&models.MerklePath{Path: 0, Depth: 0}, common.Hash{1, 2, 3})
	s.NoError(err)

	err = validateStateRoot(s.storage.Storage)
	s.NoError(err)
}

func TestRollupTestSuite(t *testing.T) {
	suite.Run(t, new(RollupTestSuite))
}
