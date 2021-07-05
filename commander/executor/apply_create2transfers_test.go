package executor

import (
	"context"
	"math/big"
	"testing"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/contracts/accountregistry"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ApplyCreate2TransfersTestSuite struct {
	*require.Assertions
	suite.Suite
	teardown            func() error
	storage             *storage.Storage
	tree                *storage.StateTree
	cfg                 *config.RollupConfig
	client              *eth.TestClient
	transactionExecutor *TransactionExecutor
	feeReceiver         *FeeReceiver
	events              chan *accountregistry.AccountRegistrySinglePubkeyRegistered
	unsubscribe         func()
}

func (s *ApplyCreate2TransfersTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *ApplyCreate2TransfersTestSuite) SetupTest() {
	testStorage, err := storage.NewTestStorageWithBadger()
	s.NoError(err)
	s.storage = testStorage.Storage
	s.teardown = testStorage.Teardown
	s.tree = storage.NewStateTree(s.storage)
	s.client, err = eth.NewTestClient()
	s.NoError(err)
	s.cfg = &config.RollupConfig{
		FeeReceiverPubKeyID: 3,
		TxsPerCommitment:    6,
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

	for i := 1; i <= 10; i++ {
		err = s.storage.AddAccountIfNotExists(&models.Account{
			PubKeyID:  uint32(i),
			PublicKey: models.PublicKey{1, 2, 3},
		})
		s.NoError(err)
	}

	_, err = s.tree.Set(1, &senderState)
	s.NoError(err)
	_, err = s.tree.Set(2, &receiverState)
	s.NoError(err)
	_, err = s.tree.Set(3, &feeReceiverState)
	s.NoError(err)

	s.events, s.unsubscribe, err = s.client.WatchRegistrations(&bind.WatchOpts{})
	s.NoError(err)

	s.transactionExecutor = NewTestTransactionExecutor(s.storage, s.client.Client, s.cfg, TransactionExecutorOpts{})
	s.feeReceiver = &FeeReceiver{
		StateID: 3,
		TokenID: models.MakeUint256(1),
	}
}

func (s *ApplyCreate2TransfersTestSuite) TearDownTest() {
	s.unsubscribe()
	s.client.Close()
	err := s.teardown()
	s.NoError(err)
}

func (s *ApplyCreate2TransfersTestSuite) TestApplyCreate2Transfers_AllValid() {
	generatedTransfers := generateValidCreate2Transfers(3)

	transfers, err := s.transactionExecutor.ApplyCreate2Transfers(generatedTransfers, s.cfg.TxsPerCommitment, s.feeReceiver)
	s.NoError(err)

	s.Len(transfers.appliedTransfers, 3)
	s.Len(transfers.invalidTransfers, 0)
	s.Len(transfers.addedPubKeyIDs, 3)
}

func (s *ApplyCreate2TransfersTestSuite) TestApplyCreate2Transfers_SomeValid() {
	generatedTransfers := generateValidCreate2Transfers(2)
	generatedTransfers = append(generatedTransfers, generateInvalidCreate2Transfers(3)...)

	transfers, err := s.transactionExecutor.ApplyCreate2Transfers(generatedTransfers, s.cfg.TxsPerCommitment, s.feeReceiver)
	s.NoError(err)

	s.Len(transfers.appliedTransfers, 2)
	s.Len(transfers.invalidTransfers, 3)
	s.Len(transfers.addedPubKeyIDs, 2)
}

func (s *ApplyCreate2TransfersTestSuite) TestApplyCreate2Transfers_AppliesNoMoreThanLimit() {
	generatedTransfers := generateValidCreate2Transfers(7)

	transfers, err := s.transactionExecutor.ApplyCreate2Transfers(generatedTransfers, s.cfg.TxsPerCommitment, s.feeReceiver)
	s.NoError(err)

	s.Len(transfers.appliedTransfers, 6)
	s.Len(transfers.invalidTransfers, 0)
	s.Len(transfers.addedPubKeyIDs, 6)
}

func (s *ApplyCreate2TransfersTestSuite) TestApplyCreate2Transfers_SavesTransferErrors() {
	generatedTransfers := generateValidCreate2Transfers(3)
	generatedTransfers = append(generatedTransfers, generateInvalidCreate2Transfers(2)...)

	for i := range generatedTransfers {
		_, err := s.storage.AddCreate2Transfer(&generatedTransfers[i])
		s.NoError(err)
	}

	transfers, err := s.transactionExecutor.ApplyCreate2Transfers(generatedTransfers, s.cfg.TxsPerCommitment, s.feeReceiver)
	s.NoError(err)

	s.Len(transfers.appliedTransfers, 3)
	s.Len(transfers.invalidTransfers, 2)
	s.Len(transfers.addedPubKeyIDs, 3)

	for i := range generatedTransfers {
		transfer, err := s.storage.GetCreate2Transfer(generatedTransfers[i].Hash)
		s.NoError(err)
		if i < 3 {
			s.Nil(transfer.ErrorMessage)
		} else {
			s.Equal(*transfer.ErrorMessage, ErrNonceTooLow.Error())
		}
	}
}

func (s *ApplyCreate2TransfersTestSuite) TestApplyCreate2Transfers_AppliesFee() {
	generatedTransfers := generateValidCreate2Transfers(3)

	_, err := s.transactionExecutor.ApplyCreate2Transfers(generatedTransfers, s.cfg.TxsPerCommitment, s.feeReceiver)
	s.NoError(err)

	feeReceiverState, err := s.transactionExecutor.storage.GetStateLeaf(s.feeReceiver.StateID)
	s.NoError(err)
	s.Equal(models.MakeUint256(1003), feeReceiverState.Balance)
}

func (s *ApplyCreate2TransfersTestSuite) TestApplyCreate2Transfers_RegistersPublicKeys() {
	generatedTransfers := generateValidCreate2Transfers(3)
	generatedTransfers[0].ToPublicKey = models.PublicKey{1, 1, 1}
	generatedTransfers[1].ToPublicKey = models.PublicKey{2, 2, 2}
	generatedTransfers[2].ToPublicKey = models.PublicKey{3, 3, 3}

	latestBlockNumber, err := s.client.GetLatestBlockNumber()
	s.NoError(err)

	transfers, err := s.transactionExecutor.ApplyCreate2Transfers(generatedTransfers, s.cfg.TxsPerCommitment, s.feeReceiver)
	s.NoError(err)

	s.Len(transfers.appliedTransfers, 3)
	s.Len(transfers.invalidTransfers, 0)
	s.Len(transfers.addedPubKeyIDs, 3)

	registeredAccounts := s.getRegisteredAccounts(*latestBlockNumber)
	for i := range generatedTransfers {
		s.Equal(registeredAccounts[i], models.Account{
			PubKeyID:  transfers.addedPubKeyIDs[i],
			PublicKey: generatedTransfers[i].ToPublicKey,
		})
	}
}

func (s *ApplyCreate2TransfersTestSuite) TestApplyCreate2TransfersForSync_AllValid() {
	generatedTransfers, pubKeyIDs := generateValidCreate2TransfersForSync(3, 4)

	appliedTransfers, err := s.transactionExecutor.ApplyCreate2TransfersForSync(generatedTransfers, pubKeyIDs, s.feeReceiver)
	s.NoError(err)

	s.Len(appliedTransfers, 3)
}

func (s *ApplyCreate2TransfersTestSuite) TestApplyCreate2TransfersForSync_InvalidTransfer() {
	generatedTransfers, pubKeyIDs := generateValidCreate2TransfersForSync(2, 4)
	invalidTxs, invalidPubKeyIDs := generateInvalidCreate2TransfersForSync(3, 6)

	generatedTransfers = append(generatedTransfers, invalidTxs...)
	pubKeyIDs = append(pubKeyIDs, invalidPubKeyIDs...)

	appliedTransfers, err := s.transactionExecutor.ApplyCreate2TransfersForSync(generatedTransfers, pubKeyIDs, s.feeReceiver)
	s.Nil(appliedTransfers)

	var disputableTransferError *DisputableTransferError
	s.ErrorAs(err, &disputableTransferError)
	s.Len(disputableTransferError.Proofs, 6)
}

func (s *ApplyCreate2TransfersTestSuite) TestApplyCreate2TransfersForSync_InvalidSlicesLength() {
	generatedTransfers := generateValidCreate2Transfers(3)
	_, err := s.transactionExecutor.ApplyCreate2TransfersForSync(generatedTransfers, []uint32{1, 2}, s.feeReceiver)
	s.Equal(ErrInvalidSliceLength, err)
}

func (s *ApplyCreate2TransfersTestSuite) TestApplyCreate2TransfersForSync_AppliesFee() {
	generatedTransfers, pubKeyIDs := generateValidCreate2TransfersForSync(3, 4)

	_, err := s.transactionExecutor.ApplyCreate2TransfersForSync(generatedTransfers, pubKeyIDs, s.feeReceiver)
	s.NoError(err)

	feeReceiverState, err := s.transactionExecutor.storage.GetStateLeaf(s.feeReceiver.StateID)
	s.NoError(err)
	s.Equal(models.MakeUint256(1003), feeReceiverState.Balance)
}

func (s *ApplyCreate2TransfersTestSuite) TestGetOrRegisterPubKeyID_RegistersPubKeyIDInCaseThereIsNoUnusedOne() {
	pubKeyID, err := s.transactionExecutor.getOrRegisterPubKeyID(s.events, &create2Transfer, models.MakeUint256(1))
	s.NoError(err)
	s.Equal(uint32(0), *pubKeyID)
}

func (s *ApplyCreate2TransfersTestSuite) TestGetOrRegisterPubKeyID_ReturnsUnusedPubKeyID() {
	c2T := create2Transfer
	c2T.ToPublicKey = models.PublicKey{1, 2, 3}

	pubKeyID, err := s.transactionExecutor.getOrRegisterPubKeyID(s.events, &c2T, models.MakeUint256(1))
	s.NoError(err)
	s.Equal(uint32(4), *pubKeyID)
}

func (s *ApplyCreate2TransfersTestSuite) getRegisteredAccounts(startBlockNumber uint64) []models.Account {
	it, err := s.client.AccountRegistry.FilterSinglePubkeyRegistered(&bind.FilterOpts{Start: startBlockNumber})
	s.NoError(err)

	registeredAccounts := make([]models.Account, 0)
	for it.Next() {
		tx, _, err := s.client.ChainConnection.GetBackend().TransactionByHash(context.Background(), it.Event.Raw.TxHash)
		s.NoError(err)

		unpack, err := s.client.AccountRegistryABI.Methods["register"].Inputs.Unpack(tx.Data()[4:])
		s.NoError(err)

		pubkey := unpack[0].([4]*big.Int)
		registeredAccounts = append(registeredAccounts, models.Account{
			PubKeyID:  uint32(it.Event.PubkeyID.Uint64()),
			PublicKey: models.MakePublicKeyFromInts(pubkey),
		})
	}
	return registeredAccounts
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
