package storage

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type CommitmentTestSuite struct {
	*require.Assertions
	suite.Suite
	storage *Storage
	db      *db.TestDB
}

func (s *CommitmentTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *CommitmentTestSuite) SetupTest() {
	testDB, err := db.NewTestDB()
	s.NoError(err)
	s.storage = NewTestStorage(testDB.DB)
	s.db = testDB
}

func (s *CommitmentTestSuite) TearDownTest() {
	err := s.db.Teardown()
	s.NoError(err)
}

func (s *CommitmentTestSuite) Test_AddCommitment_AddAndRetrieve() {
	commitment := models.Commitment{
		LeafHash:          common.Hash{},
		PostStateRoot:     common.Hash{},
		BodyHash:          common.Hash{},
		AccountTreeRoot:   common.Hash{},
		CombinedSignature: models.Signature{models.MakeUint256(1), models.MakeUint256(2)},
		FeeReceiver:       uint32(1),
		Transactions:      []byte{1, 2, 3},
	}

	err := s.storage.AddCommitment(&commitment)
	s.NoError(err)
}

func TestCommitmentTestSuite(t *testing.T) {
	suite.Run(t, new(CommitmentTestSuite))
}
