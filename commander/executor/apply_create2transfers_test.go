package executor

import (
	"context"
	"math/big"
	"testing"

	"github.com/Worldcoin/hubble-commander/commander/applier"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/contracts/accountregistry"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/suite"
)

type ApplyCreate2TransfersTestSuite struct {
	TestSuiteWithRollupContext
	feeReceiver *FeeReceiver
	events      chan *accountregistry.AccountRegistrySinglePubkeyRegistered
	unsubscribe func()
}

func (s *ApplyCreate2TransfersTestSuite) SetupTest() {
	s.TestSuiteWithRollupContext.SetupTestWithConfig(txtype.Create2Transfer, config.RollupConfig{
		FeeReceiverPubKeyID: 3,
		MaxTxsPerCommitment: 6,
	})

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

	_, err := s.storage.StateTree.Set(1, &senderState)
	s.NoError(err)
	_, err = s.storage.StateTree.Set(2, &receiverState)
	s.NoError(err)
	_, err = s.storage.StateTree.Set(3, &feeReceiverState)
	s.NoError(err)

	s.events, s.unsubscribe, err = s.client.WatchRegistrations(&bind.WatchOpts{})
	s.NoError(err)

	s.feeReceiver = &FeeReceiver{
		StateID: 3,
		TokenID: models.MakeUint256(1),
	}
}

func (s *ApplyCreate2TransfersTestSuite) TearDownTest() {
	s.unsubscribe()
	s.TestSuiteWithRollupContext.TearDownTest()
}

func (s *ApplyCreate2TransfersTestSuite) TestApplyTxs_AllValid() {
	generatedTransfers := testutils.GenerateValidCreate2Transfers(3)

	transfers, err := s.rollupCtx.ApplyTxs(generatedTransfers, s.cfg.MaxTxsPerCommitment, s.feeReceiver)
	s.NoError(err)

	s.Len(transfers.AppliedTxs(), 3)
	s.Len(transfers.InvalidTxs(), 0)
	s.Len(transfers.AddedPubKeyIDs(), 3)
}

func (s *ApplyCreate2TransfersTestSuite) TestApplyTxs_SomeValid() {
	generatedTransfers := testutils.GenerateValidCreate2Transfers(2)
	generatedTransfers = append(generatedTransfers, testutils.GenerateInvalidCreate2Transfers(3)...)

	transfers, err := s.rollupCtx.ApplyTxs(generatedTransfers, s.cfg.MaxTxsPerCommitment, s.feeReceiver)
	s.NoError(err)

	s.Len(transfers.AppliedTxs(), 2)
	s.Len(transfers.InvalidTxs(), 3)
	s.Len(transfers.AddedPubKeyIDs(), 2)
}

func (s *ApplyCreate2TransfersTestSuite) TestApplyTxs_AppliesNoMoreThanLimit() {
	generatedTransfers := testutils.GenerateValidCreate2Transfers(7)

	transfers, err := s.rollupCtx.ApplyTxs(generatedTransfers, s.cfg.MaxTxsPerCommitment, s.feeReceiver)
	s.NoError(err)

	s.Len(transfers.AppliedTxs(), 6)
	s.Len(transfers.InvalidTxs(), 0)
	s.Len(transfers.AddedPubKeyIDs(), 6)
}

func (s *ApplyCreate2TransfersTestSuite) TestApplyTxs_SavesTransferErrors() {
	generatedTransfers := testutils.GenerateValidCreate2Transfers(3)
	generatedTransfers = append(generatedTransfers, testutils.GenerateInvalidCreate2Transfers(2)...)

	for i := range generatedTransfers {
		err := s.storage.AddCreate2Transfer(&generatedTransfers[i])
		s.NoError(err)
	}

	transfers, err := s.rollupCtx.ApplyTxs(generatedTransfers, s.cfg.MaxTxsPerCommitment, s.feeReceiver)
	s.NoError(err)

	s.Len(transfers.AppliedTxs(), 3)
	s.Len(transfers.InvalidTxs(), 2)
	s.Len(transfers.AddedPubKeyIDs(), 3)

	for i := range generatedTransfers {
		transfer, err := s.storage.GetCreate2Transfer(generatedTransfers[i].Hash)
		s.NoError(err)
		if i < 3 {
			s.Nil(transfer.ErrorMessage)
		} else {
			s.Equal(*transfer.ErrorMessage, applier.ErrNonceTooLow.Error())
		}
	}
}

func (s *ApplyCreate2TransfersTestSuite) TestApplyTxs_AppliesFee() {
	generatedTransfers := testutils.GenerateValidCreate2Transfers(3)

	_, err := s.rollupCtx.ApplyTxs(generatedTransfers, s.cfg.MaxTxsPerCommitment, s.feeReceiver)
	s.NoError(err)

	feeReceiverState, err := s.executionCtx.storage.StateTree.Leaf(s.feeReceiver.StateID)
	s.NoError(err)
	s.Equal(models.MakeUint256(1003), feeReceiverState.Balance)
}

func (s *ApplyCreate2TransfersTestSuite) TestApplyTxs_RegistersPublicKeys() {
	generatedTransfers := testutils.GenerateValidCreate2Transfers(3)
	generatedTransfers[0].ToPublicKey = models.PublicKey{1, 1, 1}
	generatedTransfers[1].ToPublicKey = models.PublicKey{2, 2, 2}
	generatedTransfers[2].ToPublicKey = models.PublicKey{3, 3, 3}

	latestBlockNumber, err := s.client.GetLatestBlockNumber()
	s.NoError(err)

	transfers, err := s.rollupCtx.ApplyTxs(generatedTransfers, s.cfg.MaxTxsPerCommitment, s.feeReceiver)
	s.NoError(err)

	s.Len(transfers.AppliedTxs(), 3)
	s.Len(transfers.InvalidTxs(), 0)
	s.Len(transfers.AddedPubKeyIDs(), 3)

	registeredAccounts := s.getRegisteredAccounts(*latestBlockNumber)
	for i := range generatedTransfers {
		s.Equal(registeredAccounts[i], models.AccountLeaf{
			PubKeyID:  transfers.AddedPubKeyIDs()[i],
			PublicKey: generatedTransfers[i].ToPublicKey,
		})
	}
}

func (s *ApplyCreate2TransfersTestSuite) TestApplyCreate2TransfersForSync_AllValid() {
	transfers, pubKeyIDs := generateValidCreate2TransfersForSync(3, 4)

	appliedTransfers, stateProofs, err := s.rollupCtx.ApplyCreate2TransfersForSync(transfers, pubKeyIDs, s.feeReceiver.StateID)
	s.NoError(err)
	s.Len(appliedTransfers, 3)
	s.Len(stateProofs, 7)
}

func (s *ApplyCreate2TransfersTestSuite) TestApplyCreate2TransfersForSync_InvalidTransfer() {
	transfers, pubKeyIDs := generateValidCreate2TransfersForSync(2, 4)
	invalidTxs, invalidPubKeyIDs := generateInvalidCreate2TransfersForSync(3, 6)

	transfers = append(transfers, invalidTxs...)
	pubKeyIDs = append(pubKeyIDs, invalidPubKeyIDs...)

	appliedTransfers, _, err := s.rollupCtx.ApplyCreate2TransfersForSync(transfers, pubKeyIDs, s.feeReceiver.StateID)
	s.Nil(appliedTransfers)

	var disputableErr *DisputableError
	s.ErrorAs(err, &disputableErr)
	s.Equal(Transition, disputableErr.Type)
	s.Len(disputableErr.Proofs, 6)
}

func (s *ApplyCreate2TransfersTestSuite) TestApplyCreate2TransfersForSync_InvalidSlicesLength() {
	generatedTransfers := testutils.GenerateValidCreate2Transfers(3)
	_, _, err := s.rollupCtx.ApplyCreate2TransfersForSync(generatedTransfers, []uint32{1, 2}, s.feeReceiver.StateID)
	s.Equal(applier.ErrInvalidSlicesLength, err)
}

func (s *ApplyCreate2TransfersTestSuite) TestApplyCreate2TransfersForSync_AppliesFee() {
	generatedTransfers, pubKeyIDs := generateValidCreate2TransfersForSync(3, 4)

	_, _, err := s.rollupCtx.ApplyCreate2TransfersForSync(generatedTransfers, pubKeyIDs, s.feeReceiver.StateID)
	s.NoError(err)

	feeReceiverState, err := s.executionCtx.storage.StateTree.Leaf(s.feeReceiver.StateID)
	s.NoError(err)
	s.Equal(models.MakeUint256(1003), feeReceiverState.Balance)
}

func (s *ApplyCreate2TransfersTestSuite) TestApplyCreate2TransfersForSync_ReturnsCorrectStateProofsForZeroFee() {
	transfers, pubKeyIDs := generateValidCreate2TransfersForSync(2, 5)
	for i := range transfers {
		transfers[i].Fee = models.MakeUint256(0)
	}

	_, stateProofs, err := s.rollupCtx.ApplyCreate2TransfersForSync(transfers, pubKeyIDs, s.feeReceiver.StateID)
	s.NoError(err)
	s.Len(stateProofs, 5)
}

func (s *ApplyCreate2TransfersTestSuite) TestApplyCreate2TransfersForSync_InvalidFeeReceiverTokenID() {
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

	appliedTransfers, _, err := s.rollupCtx.ApplyCreate2TransfersForSync(transfers, pubKeyIDs, feeReceiver.StateID)
	s.Nil(appliedTransfers)

	var disputableErr *DisputableError
	s.ErrorAs(err, &disputableErr)
	s.Equal(Transition, disputableErr.Type)
	s.Equal(applier.ErrInvalidFeeReceiverTokenID.Error(), disputableErr.Reason)
	s.Len(disputableErr.Proofs, 5)
}

func (s *ApplyCreate2TransfersTestSuite) getRegisteredAccounts(startBlockNumber uint64) []models.AccountLeaf {
	it, err := s.client.AccountRegistry.FilterSinglePubkeyRegistered(&bind.FilterOpts{Start: startBlockNumber})
	s.NoError(err)

	registeredAccounts := make([]models.AccountLeaf, 0)
	for it.Next() {
		tx, _, err := s.client.ChainConnection.GetBackend().TransactionByHash(context.Background(), it.Event.Raw.TxHash)
		s.NoError(err)

		unpack, err := s.client.AccountRegistryABI.Methods["register"].Inputs.Unpack(tx.Data()[4:])
		s.NoError(err)

		pubkey := unpack[0].([4]*big.Int)
		registeredAccounts = append(registeredAccounts, models.AccountLeaf{
			PubKeyID:  uint32(it.Event.PubkeyID.Uint64()),
			PublicKey: models.MakePublicKeyFromInts(pubkey),
		})
	}
	return registeredAccounts
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
