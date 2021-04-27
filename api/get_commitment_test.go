package api

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type GetCommitmentTestSuite struct {
	*require.Assertions
	suite.Suite
	api        *API
	storage    *st.Storage
	db         *db.TestDB
	tree       *st.StateTree
	commitment models.Commitment
}

func (s *GetCommitmentTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *GetCommitmentTestSuite) SetupTest() {
	testDB, err := db.NewTestDB()
	s.NoError(err)

	s.storage = st.NewTestStorage(testDB.DB)
	s.api = &API{nil, s.storage, nil}
	s.db = testDB
	s.tree = st.NewStateTree(s.storage)

	hash := utils.RandomHash()
	s.commitment = commitment
	s.commitment.AccountTreeRoot = &hash
}

func (s *GetCommitmentTestSuite) TearDownTest() {
	err := s.db.Teardown()
	s.NoError(err)
}

func (s *GetCommitmentTestSuite) TestGetCommitment() {
	commitmentID, err := s.storage.AddCommitment(&s.commitment)
	s.NoError(err)

	transfer := models.Transfer{
		TransactionBase: models.TransactionBase{
			Hash:                 utils.RandomHash(),
			TxType:               txtype.Transfer,
			FromStateID:          1,
			Amount:               models.MakeUint256(50),
			Fee:                  models.MakeUint256(10),
			Nonce:                models.MakeUint256(0),
			IncludedInCommitment: commitmentID,
		},
		ToStateID: 2,
	}
	err = s.storage.AddTransfer(&transfer)
	s.NoError(err)

	commitment, err := s.api.GetCommitment(*commitmentID)
	s.NoError(err)
	s.NotNil(commitment)
	s.Len(commitment.Transactions, 1)
	s.Contains(commitment.Transactions, transfer)
}

func TestGetCommitmentTestSuite(t *testing.T) {
	suite.Run(t, new(GetCommitmentTestSuite))
}
