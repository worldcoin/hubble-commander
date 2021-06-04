package commander

import (
	"context"
	"testing"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type SaveBatchTestSuite struct {
	*require.Assertions
	suite.Suite
	teardown func() error
	storage  *st.Storage
	tree     *st.StateTree
	client   *eth.TestClient
	cfg      *config.RollupConfig
}

func (s *SaveBatchTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *SaveBatchTestSuite) SetupTest() {
	var err error
	s.client, err = eth.NewTestClient()
	s.NoError(err)

	s.cfg = &config.RollupConfig{
		MinCommitmentsPerBatch: 1,
		MaxCommitmentsPerBatch: 32,
		TxsPerCommitment:       1,
	}

	s.setupDB()
}

func (s *SaveBatchTestSuite) setupDB() {
	testStorage, err := st.NewTestStorageWithBadger()
	s.NoError(err)
	s.storage = testStorage.Storage
	s.teardown = testStorage.Teardown
	s.tree = st.NewStateTree(s.storage)

	s.seedDB()
}

func (s *SaveBatchTestSuite) seedDB() {
	err := s.storage.AddAccountIfNotExists(&models.Account{
		PubKeyID:  0,
		PublicKey: models.PublicKey{1, 2, 3},
	})
	s.NoError(err)

	err = s.storage.AddAccountIfNotExists(&models.Account{
		PubKeyID:  1,
		PublicKey: models.PublicKey{2, 3, 4},
	})
	s.NoError(err)

	err = s.tree.Set(0, &models.UserState{
		PubKeyID:   0,
		TokenIndex: models.MakeUint256(0),
		Balance:    models.MakeUint256(1000),
		Nonce:      models.MakeUint256(0),
	})
	s.NoError(err)

	err = s.tree.Set(1, &models.UserState{
		PubKeyID:   1,
		TokenIndex: models.MakeUint256(0),
		Balance:    models.MakeUint256(0),
		Nonce:      models.MakeUint256(0),
	})
	s.NoError(err)
}

func (s *SaveBatchTestSuite) TearDownTest() {
	s.client.Close()
	err := s.teardown()
	s.NoError(err)
}

func (s *SaveBatchTestSuite) TestSaveBatch() {
	tx := models.Transfer{
		TransactionBase: models.TransactionBase{
			TxType:      txtype.Transfer,
			FromStateID: 0,
			Amount:      models.MakeUint256(400),
			Fee:         models.MakeUint256(0),
			Nonce:       models.MakeUint256(0),
			Signature:   *s.mockSignature(),
		},
		ToStateID: 1,
	}
	err := s.storage.AddTransfer(&tx)
	s.NoError(err)

	commitments, err := createTransferCommitments([]models.Transfer{tx}, s.storage, s.cfg, testDomain)
	s.NoError(err)
	s.Len(commitments, 1)

	err = submitBatch(context.Background(), txtype.Transfer, commitments, s.storage, s.client.Client, s.cfg)
	s.NoError(err)

	//err = saveBatch(s.storage, s.client.Client)
	//s.NoError(err)

	batches, err := s.storage.GetBatchesInRange(nil, nil)
	s.NoError(err)
	s.Len(batches, 1)
	s.Nil(batches[0].Number)
	s.Nil(batches[0].Hash)

	commitment, err := s.storage.GetCommitment(commitments[0].ID)
	s.NoError(err)
	s.Equal(batches[0].ID, *commitment.IncludedInBatch)
	s.Nil(commitment.AccountTreeRoot)

	s.client.Commit()

	//err = saveBatch(s.storage, s.client.Client)
	//s.NoError(err)

	batches, err = s.storage.GetBatchesInRange(nil, nil)
	s.NoError(err)
	s.Len(batches, 1)
	s.NotNil(batches[0].Number)
	s.NotNil(batches[0].Hash)

	commitment, err = s.storage.GetCommitment(commitments[0].ID)
	s.NoError(err)
	s.Equal(batches[0].ID, *commitment.IncludedInBatch)
	s.NotNil(commitment.AccountTreeRoot)
}

func (s *SaveBatchTestSuite) mockSignature() *models.Signature {
	wallet, err := bls.NewRandomWallet(testDomain)
	s.NoError(err)
	signature, err := wallet.Sign(utils.RandomBytes(4))
	s.NoError(err)
	return signature.ModelsSignature()
}

func TestSaveBatchTestSuite(t *testing.T) {
	suite.Run(t, new(SaveBatchTestSuite))
}
