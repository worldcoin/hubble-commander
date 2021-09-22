package executor

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/commander/applier"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type SyncApplyCreate2TransfersTestSuite struct {
	*require.Assertions
	suite.Suite
	storage     *storage.TestStorage
	client      *eth.TestClient
	cfg         *config.RollupConfig
	syncCtx     *SyncContext
	feeReceiver *FeeReceiver
}

func (s *SyncApplyCreate2TransfersTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *SyncApplyCreate2TransfersTestSuite) SetupTest() {
	var err error
	s.storage, err = storage.NewTestStorage()
	s.NoError(err)
	s.client, err = eth.NewTestClient()
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

	executionCtx := NewTestExecutionContext(s.storage.Storage, nil, s.cfg)
	s.syncCtx = NewTestSyncContext(executionCtx, batchtype.Create2Transfer)

	s.feeReceiver = &FeeReceiver{
		StateID: 3,
		TokenID: models.MakeUint256(1),
	}
}

func (s *SyncApplyCreate2TransfersTestSuite) TearDownTest() {
	s.client.Close()
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *SyncApplyCreate2TransfersTestSuite) TestApplyTxForSync_AllValid() {
	transfers, pubKeyIDs := generateValidCreate2TransfersForSync(3, 4)
	input := &SyncedC2Ts{
		txs:       transfers,
		pubKeyIDs: pubKeyIDs,
	}

	appliedTransfers, stateProofs, err := s.syncCtx.ApplyTxsForSync(input, s.feeReceiver.StateID)
	s.NoError(err)
	s.Len(appliedTransfers, 3)
	s.Len(stateProofs, 7)
}

func (s *SyncApplyCreate2TransfersTestSuite) TestApplyTxForSync_InvalidTransfer() {
	transfers, pubKeyIDs := generateValidCreate2TransfersForSync(2, 4)
	invalidTxs, invalidPubKeyIDs := generateInvalidCreate2TransfersForSync(3, 6)
	input := &SyncedC2Ts{
		txs:       append(transfers, invalidTxs...),
		pubKeyIDs: append(pubKeyIDs, invalidPubKeyIDs...),
	}

	appliedTransfers, _, err := s.syncCtx.ApplyTxsForSync(input, s.feeReceiver.StateID)
	s.Nil(appliedTransfers)

	var disputableErr *DisputableError
	s.ErrorAs(err, &disputableErr)
	s.Equal(Transition, disputableErr.Type)
	s.Len(disputableErr.Proofs, 6)
}

func (s *SyncApplyCreate2TransfersTestSuite) TestApplyTxForSync_AppliesFee() {
	generatedTransfers, pubKeyIDs := generateValidCreate2TransfersForSync(3, 4)
	input := &SyncedC2Ts{
		txs:       generatedTransfers,
		pubKeyIDs: pubKeyIDs,
	}
	_, _, err := s.syncCtx.ApplyTxsForSync(input, s.feeReceiver.StateID)
	s.NoError(err)

	feeReceiverState, err := s.syncCtx.storage.StateTree.Leaf(s.feeReceiver.StateID)
	s.NoError(err)
	s.Equal(models.MakeUint256(1003), feeReceiverState.Balance)
}

func (s *SyncApplyCreate2TransfersTestSuite) TestApplyTxForSync_ReturnsCorrectStateProofsForZeroFee() {
	transfers, pubKeyIDs := generateValidCreate2TransfersForSync(2, 5)
	for i := range transfers {
		transfers[i].Fee = models.MakeUint256(0)
	}

	input := &SyncedC2Ts{
		txs:       transfers,
		pubKeyIDs: pubKeyIDs,
	}
	_, stateProofs, err := s.syncCtx.ApplyTxsForSync(input, s.feeReceiver.StateID)
	s.NoError(err)
	s.Len(stateProofs, 5)
}

func (s *SyncApplyCreate2TransfersTestSuite) TestApplyTxForSync_InvalidFeeReceiverTokenID() {
	feeReceiver := &FeeReceiver{
		StateID: 4,
		TokenID: models.MakeUint256(4),
	}
	_, err := s.storage.StateTree.Set(feeReceiver.StateID, &models.UserState{
		PubKeyID: 4,
		TokenID:  feeReceiver.TokenID,
		Balance:  models.MakeUint256(420),
		Nonce:    models.MakeUint256(0),
	})
	s.NoError(err)

	transfers, pubKeyIDs := generateValidCreate2TransfersForSync(2, 5)
	input := &SyncedC2Ts{
		txs:       transfers,
		pubKeyIDs: pubKeyIDs,
	}

	appliedTransfers, _, err := s.syncCtx.ApplyTxsForSync(input, feeReceiver.StateID)
	s.Nil(appliedTransfers)

	var disputableErr *DisputableError
	s.ErrorAs(err, &disputableErr)
	s.Equal(Transition, disputableErr.Type)
	s.Equal(applier.ErrInvalidFeeReceiverTokenID.Error(), disputableErr.Reason)
	s.Len(disputableErr.Proofs, 5)
}

func generateValidCreate2TransfersForSync(transfersAmount, startPubKeyID uint32) (
	transfers []models.Create2Transfer,
	pubKeyIDs []uint32,
) {
	transfers = make([]models.Create2Transfer, 0, transfersAmount)
	pubKeyIDs = make([]uint32, 0, transfersAmount)

	for i := 0; i < int(transfersAmount); i++ {
		transfer := models.Create2Transfer{
			TransactionBase: models.TransactionBase{
				Hash:        utils.RandomHash(),
				TxType:      txtype.Create2Transfer,
				FromStateID: 1,
				Amount:      models.MakeUint256(1),
				Fee:         models.MakeUint256(1),
				Nonce:       models.MakeUint256(uint64(i)),
			},
			ToStateID:   ref.Uint32(startPubKeyID),
			ToPublicKey: models.PublicKey{1, 2, 3},
		}
		transfers = append(transfers, transfer)
		pubKeyIDs = append(pubKeyIDs, startPubKeyID)
		startPubKeyID++
	}
	return transfers, pubKeyIDs
}

func generateInvalidCreate2TransfersForSync(transfersAmount, startPubKeyID uint32) (
	transfers []models.Create2Transfer,
	pubKeyIDs []uint32,
) {
	transfers = make([]models.Create2Transfer, 0, transfersAmount)
	pubKeyIDs = make([]uint32, 0, transfersAmount)

	for i := 0; i < int(transfersAmount); i++ {
		transfer := models.Create2Transfer{
			TransactionBase: models.TransactionBase{
				Hash:        utils.RandomHash(),
				TxType:      txtype.Create2Transfer,
				FromStateID: 1,
				Amount:      models.MakeUint256(1_000_000),
				Fee:         models.MakeUint256(1),
				Nonce:       models.MakeUint256(0),
			},
			ToStateID:   ref.Uint32(startPubKeyID),
			ToPublicKey: models.PublicKey{1, 2, 3},
		}
		transfers = append(transfers, transfer)
		pubKeyIDs = append(pubKeyIDs, startPubKeyID)
		startPubKeyID++
	}
	return transfers, pubKeyIDs
}

func TestSyncApplyCreate2TransfersTestSuite(t *testing.T) {
	suite.Run(t, new(SyncApplyCreate2TransfersTestSuite))
}
