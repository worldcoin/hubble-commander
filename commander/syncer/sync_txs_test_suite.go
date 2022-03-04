package syncer

import (
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/Worldcoin/hubble-commander/storage"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// Other test suites encapsulate syncTxsTestSuite.
// Don't add any tests on syncTxsTestSuite to avoid repeated runs.
type syncTxsTestSuite struct {
	*require.Assertions
	suite.Suite
	storage            *storage.TestStorage
	cfg                *config.RollupConfig
	syncCtx            *TxsContext
	feeReceiverStateID uint32
}

func (s *syncTxsTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *syncTxsTestSuite) SetupTest(txType txtype.TransactionType) {
	var err error
	s.storage, err = storage.NewTestStorage()
	s.NoError(err)
	s.cfg = &config.RollupConfig{
		FeeReceiverPubKeyID: 3,
		MaxTxsPerCommitment: 6,
	}

	senderState := models.UserState{
		PubKeyID: 1,
		TokenID:  models.MakeUint256(1),
		Balance:  models.MakeUint256(420),
		Nonce:    models.MakeUint256(0),
	}
	receiverState := models.UserState{
		PubKeyID: 2,
		TokenID:  models.MakeUint256(1),
		Balance:  models.MakeUint256(0),
		Nonce:    models.MakeUint256(0),
	}
	feeReceiverState := models.UserState{
		PubKeyID: 3,
		TokenID:  models.MakeUint256(1),
		Balance:  models.MakeUint256(1000),
		Nonce:    models.MakeUint256(0),
	}

	_, err = s.storage.StateTree.Set(1, &senderState)
	s.NoError(err)
	_, err = s.storage.StateTree.Set(2, &receiverState)
	s.NoError(err)
	_, err = s.storage.StateTree.Set(3, &feeReceiverState)
	s.NoError(err)

	s.syncCtx, err = NewTestTxsContext(s.storage.Storage, nil, s.cfg, txType)
	s.NoError(err)

	s.feeReceiverStateID = 3
}

func (s *syncTxsTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}
