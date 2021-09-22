package executor

import (
	"context"
	"math/big"
	"testing"

	"github.com/Worldcoin/hubble-commander/commander/applier"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/suite"
)

type ApplyCreate2TransfersTestSuite struct {
	TestSuiteWithRollupContext
	feeReceiver *FeeReceiver
}

func (s *ApplyCreate2TransfersTestSuite) SetupTest() {
	s.TestSuiteWithRollupContext.SetupTestWithConfig(batchtype.Create2Transfer, config.RollupConfig{
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

	s.feeReceiver = &FeeReceiver{
		StateID: 3,
		TokenID: models.MakeUint256(1),
	}
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

func TestApplyCreate2TransfersTestSuite(t *testing.T) {
	suite.Run(t, new(ApplyCreate2TransfersTestSuite))
}
