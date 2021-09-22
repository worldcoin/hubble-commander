package syncer

import (
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/Worldcoin/hubble-commander/storage"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// Other test suites encapsulate SyncApplyTxsTestSuite.
// Don't add any tests on SyncApplyTxsTestSuite to avoid repeated runs.
type SyncApplyTxsTestSuite struct {
	*require.Assertions
	suite.Suite
	storage     *storage.TestStorage
	cfg         *config.RollupConfig
	syncCtx     *SyncContext
	feeReceiver *FeeReceiver
}

type FeeReceiver struct {
	StateID uint32
	TokenID models.Uint256
}

func (s *SyncApplyTxsTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *SyncApplyTxsTestSuite) SetupTest(batchType batchtype.BatchType) {
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

	s.syncCtx = NewTestSyncContext(s.storage.Storage, nil, s.cfg, batchType)

	s.feeReceiver = &FeeReceiver{
		StateID: 3,
		TokenID: models.MakeUint256(1),
	}
}

func (s *SyncApplyTxsTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}
