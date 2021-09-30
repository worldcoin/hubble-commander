package executor

import (
	"context"
	"testing"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ApplyCreate2TransfersTestSuite struct {
	*require.Assertions
	suite.Suite
	storage             *storage.TestStorage
	cfg                 *config.RollupConfig
	client              *eth.TestClient
	transactionExecutor *TransactionExecutor
	feeReceiver         *FeeReceiver
}

func (s *ApplyCreate2TransfersTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *ApplyCreate2TransfersTestSuite) SetupTest() {
	var err error
	s.storage, err = storage.NewTestStorageWithBadger()
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

	s.transactionExecutor = NewTestTransactionExecutor(s.storage.Storage, s.client.Client, s.cfg, context.Background())
	s.feeReceiver = &FeeReceiver{
		StateID: 3,
		TokenID: models.MakeUint256(1),
	}
}

func (s *ApplyCreate2TransfersTestSuite) TearDownTest() {
	s.client.Close()
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *ApplyCreate2TransfersTestSuite) TestApplyCreate2Transfers_AllValid() {
	pending := &PendingC2Ts{
		Txs: generateValidCreate2Transfers(3),
	}

	transfers, err := s.transactionExecutor.ApplyCreate2Transfers(pending, s.cfg.MaxTxsPerCommitment, s.feeReceiver)
	s.NoError(err)

	s.Len(transfers.appliedTransfers, 3)
	s.Len(transfers.invalidTransfers, 0)
	s.Len(transfers.pendingAccounts, 3)
	s.Len(transfers.addedPubKeyIDs, 3)
}

func (s *ApplyCreate2TransfersTestSuite) TestApplyCreate2Transfers_SomeValid() {
	generatedTransfers := generateValidCreate2Transfers(2)
	pending := &PendingC2Ts{
		Txs: append(generatedTransfers, generateInvalidCreate2Transfers(3)...),
	}

	transfers, err := s.transactionExecutor.ApplyCreate2Transfers(pending, s.cfg.MaxTxsPerCommitment, s.feeReceiver)
	s.NoError(err)

	s.Len(transfers.appliedTransfers, 2)
	s.Len(transfers.invalidTransfers, 3)
	s.Len(transfers.pendingAccounts, 2)
	s.Len(transfers.addedPubKeyIDs, 2)
}

func (s *ApplyCreate2TransfersTestSuite) TestApplyCreate2Transfers_AppliesNoMoreThanLimit() {
	pending := &PendingC2Ts{
		Txs: generateValidCreate2Transfers(7),
	}

	transfers, err := s.transactionExecutor.ApplyCreate2Transfers(pending, s.cfg.MaxTxsPerCommitment, s.feeReceiver)
	s.NoError(err)

	s.Len(transfers.appliedTransfers, 6)
	s.Len(transfers.invalidTransfers, 0)
	s.Len(transfers.pendingAccounts, 6)
	s.Len(transfers.addedPubKeyIDs, 6)
}

func (s *ApplyCreate2TransfersTestSuite) TestApplyCreate2Transfers_SavesTransferErrors() {
	generatedTransfers := generateValidCreate2Transfers(3)
	pending := &PendingC2Ts{
		Txs: append(generatedTransfers, generateInvalidCreate2Transfers(2)...),
	}

	for i := range pending.Txs {
		_, err := s.storage.AddCreate2Transfer(&pending.Txs[i])
		s.NoError(err)
	}

	transfers, err := s.transactionExecutor.ApplyCreate2Transfers(pending, s.cfg.MaxTxsPerCommitment, s.feeReceiver)
	s.NoError(err)

	s.Len(transfers.appliedTransfers, 3)
	s.Len(transfers.invalidTransfers, 2)
	s.Len(transfers.pendingAccounts, 3)
	s.Len(transfers.addedPubKeyIDs, 3)

	for i := range pending.Txs {
		transfer, err := s.storage.GetCreate2Transfer(pending.Txs[i].Hash)
		s.NoError(err)
		if i < 3 {
			s.Nil(transfer.ErrorMessage)
		} else {
			s.Equal(*transfer.ErrorMessage, ErrNonceTooLow.Error())
		}
	}
}

func (s *ApplyCreate2TransfersTestSuite) TestApplyCreate2Transfers_AppliesFee() {
	pending := &PendingC2Ts{
		Txs: generateValidCreate2Transfers(3),
	}

	_, err := s.transactionExecutor.ApplyCreate2Transfers(pending, s.cfg.MaxTxsPerCommitment, s.feeReceiver)
	s.NoError(err)

	feeReceiverState, err := s.transactionExecutor.storage.StateTree.Leaf(s.feeReceiver.StateID)
	s.NoError(err)
	s.Equal(models.MakeUint256(1003), feeReceiverState.Balance)
}

func (s *ApplyCreate2TransfersTestSuite) TestApplyCreate2Transfers_AddsPendingAccounts() {
	pending := &PendingC2Ts{
		Txs: generateValidCreate2Transfers(3),
	}
	pending.Txs[0].ToPublicKey = models.PublicKey{1, 1, 1}
	pending.Txs[1].ToPublicKey = models.PublicKey{2, 2, 2}
	pending.Txs[2].ToPublicKey = models.PublicKey{3, 3, 3}

	transfers, err := s.transactionExecutor.ApplyCreate2Transfers(pending, s.cfg.MaxTxsPerCommitment, s.feeReceiver)
	s.NoError(err)

	s.Len(transfers.appliedTransfers, 3)
	s.Len(transfers.invalidTransfers, 0)
	s.Len(transfers.pendingAccounts, 3)
	s.Len(transfers.addedPubKeyIDs, 3)

	for i := range pending.Txs {
		s.Equal(pending.Txs[i].ToPublicKey, transfers.pendingAccounts[i].PublicKey)
	}
}

func (s *ApplyCreate2TransfersTestSuite) TestApplyCreate2TransfersForSync_AllValid() {
	transfers, pubKeyIDs := generateValidCreate2TransfersForSync(3, 4)

	appliedTransfers, stateProofs, err := s.transactionExecutor.ApplyCreate2TransfersForSync(transfers, pubKeyIDs, s.feeReceiver.StateID)
	s.NoError(err)
	s.Len(appliedTransfers, 3)
	s.Len(stateProofs, 7)
}

func (s *ApplyCreate2TransfersTestSuite) TestApplyCreate2TransfersForSync_InvalidTransfer() {
	transfers, pubKeyIDs := generateValidCreate2TransfersForSync(2, 4)
	invalidTxs, invalidPubKeyIDs := generateInvalidCreate2TransfersForSync(3, 6)

	transfers = append(transfers, invalidTxs...)
	pubKeyIDs = append(pubKeyIDs, invalidPubKeyIDs...)

	appliedTransfers, _, err := s.transactionExecutor.ApplyCreate2TransfersForSync(transfers, pubKeyIDs, s.feeReceiver.StateID)
	s.Nil(appliedTransfers)

	var disputableErr *DisputableError
	s.ErrorAs(err, &disputableErr)
	s.Equal(Transition, disputableErr.Type)
	s.Len(disputableErr.Proofs, 6)
}

func (s *ApplyCreate2TransfersTestSuite) TestApplyCreate2TransfersForSync_InvalidSlicesLength() {
	generatedTransfers := generateValidCreate2Transfers(3)
	_, _, err := s.transactionExecutor.ApplyCreate2TransfersForSync(generatedTransfers, []uint32{1, 2}, s.feeReceiver.StateID)
	s.Equal(ErrInvalidSlicesLength, err)
}

func (s *ApplyCreate2TransfersTestSuite) TestApplyCreate2TransfersForSync_AppliesFee() {
	generatedTransfers, pubKeyIDs := generateValidCreate2TransfersForSync(3, 4)

	_, _, err := s.transactionExecutor.ApplyCreate2TransfersForSync(generatedTransfers, pubKeyIDs, s.feeReceiver.StateID)
	s.NoError(err)

	feeReceiverState, err := s.transactionExecutor.storage.StateTree.Leaf(s.feeReceiver.StateID)
	s.NoError(err)
	s.Equal(models.MakeUint256(1003), feeReceiverState.Balance)
}

func (s *ApplyTransfersTestSuite) TestApplyCreate2TransfersForSync_ReturnsCorrectStateProofsForZeroFee() {
	transfers, pubKeyIDs := generateValidCreate2TransfersForSync(2, 5)
	for i := range transfers {
		transfers[i].Fee = models.MakeUint256(0)
	}

	_, stateProofs, err := s.transactionExecutor.ApplyCreate2TransfersForSync(transfers, pubKeyIDs, s.feeReceiver.StateID)
	s.NoError(err)
	s.Len(stateProofs, 5)
}

func (s *ApplyTransfersTestSuite) TestApplyCreate2TransfersForSync_InvalidFeeReceiverTokenID() {
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

	appliedTransfers, _, err := s.transactionExecutor.ApplyCreate2TransfersForSync(transfers, pubKeyIDs, feeReceiver.StateID)
	s.Nil(appliedTransfers)

	var disputableErr *DisputableError
	s.ErrorAs(err, &disputableErr)
	s.Equal(Transition, disputableErr.Type)
	s.Equal(ErrInvalidFeeReceiverTokenID.Error(), disputableErr.Reason)
	s.Len(disputableErr.Proofs, 5)
}

func (s *ApplyCreate2TransfersTestSuite) TestGetPubKeyID_GetsPubKeyIDFromAccountTree() {
	pubKeyID, isPending, err := s.transactionExecutor.getPubKeyID(PendingAccounts{}, &create2Transfer, models.MakeUint256(1))
	s.NoError(err)
	s.True(isPending)
	s.Equal(uint32(2147483648), *pubKeyID)
}

func (s *ApplyCreate2TransfersTestSuite) TestGetPubKeyID_PredictsPubKeyIDInCaseThereIsNoUnusedOne() {
	pendingAccounts := PendingAccounts{
		{
			PubKeyID:  2147483650,
			PublicKey: models.PublicKey{1, 2, 3},
		},
	}
	pubKeyID, isPending, err := s.transactionExecutor.getPubKeyID(pendingAccounts, &create2Transfer, models.MakeUint256(1))
	s.NoError(err)
	s.True(isPending)
	s.Equal(uint32(2147483651), *pubKeyID)
}

func (s *ApplyCreate2TransfersTestSuite) TestGetPubKeyID_ReturnsUnusedPubKeyID() {
	for i := 1; i <= 10; i++ {
		err := s.storage.AccountTree.SetSingle(&models.AccountLeaf{
			PubKeyID:  uint32(i),
			PublicKey: models.PublicKey{1, 2, 3},
		})
		s.NoError(err)
	}

	c2T := create2Transfer
	c2T.ToPublicKey = models.PublicKey{1, 2, 3}

	pubKeyID, isPending, err := s.transactionExecutor.getPubKeyID(PendingAccounts{}, &c2T, models.MakeUint256(1))
	s.NoError(err)
	s.False(isPending)
	s.Equal(uint32(4), *pubKeyID)
}

func generateValidCreate2Transfers(transfersAmount uint32) []models.Create2Transfer {
	transfers := make([]models.Create2Transfer, 0, transfersAmount)
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
			ToStateID:   nil,
			ToPublicKey: models.PublicKey{1, 2, 3},
		}
		transfers = append(transfers, transfer)
	}
	return transfers
}

func generateInvalidCreate2Transfers(transfersAmount uint64) []models.Create2Transfer {
	transfers := make([]models.Create2Transfer, 0, transfersAmount)
	for i := uint64(0); i < transfersAmount; i++ {
		transfer := models.Create2Transfer{
			TransactionBase: models.TransactionBase{
				Hash:        utils.RandomHash(),
				TxType:      txtype.Create2Transfer,
				FromStateID: 1,
				Amount:      models.MakeUint256(1),
				Fee:         models.MakeUint256(1),
				Nonce:       models.MakeUint256(0),
			},
			ToStateID:   nil,
			ToPublicKey: models.PublicKey{1, 2, 3},
		}
		transfers = append(transfers, transfer)
	}
	return transfers
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

func TestApplyCreate2TransfersTestSuite(t *testing.T) {
	suite.Run(t, new(ApplyCreate2TransfersTestSuite))
}
