package commander

import (
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

// TODO Return a struct with 4 fields
func ApplyCreate2Transfers(
	storage *st.Storage,
	client *eth.Client,
	transfers []models.Create2Transfer,
	cfg *config.RollupConfig,
) (
	appliedTransfers []models.Create2Transfer,
	invalidTransfers []models.Create2Transfer,
	addedPubKeyIDs []uint32,
	feeReceiverStateID *uint32,
	err error,
) {
	if len(transfers) == 0 {
		return
	}
	syncedBlock, err := storage.GetSyncedBlock(client.ChainState.ChainID)
	if err != nil {
		return
	}
	events, unsubscribe, err := client.WatchRegistrations(&bind.WatchOpts{
		Start: syncedBlock,
	})
	if err != nil {
		return
	}
	defer unsubscribe()

	appliedTransfers = make([]models.Create2Transfer, 0, cfg.TxsPerCommitment)
	addedPubKeyIDs = make([]uint32, 0, cfg.TxsPerCommitment)
	combinedFee := models.NewUint256(0)

	commitmentTokenIndex, err := getTokenIndex(storage, transfers[0].FromStateID)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	var pubKeyID *uint32
	var ok bool

	for i := range transfers {
		transfer := &transfers[i]

		pubKeyID, err = getOrRegisterPubKeyID(storage, client, events, transfer, *commitmentTokenIndex)
		if err != nil {
			return nil, nil, nil, nil, err
		}

		ok, err = handleApplyC2T(storage, transfer, *pubKeyID, &appliedTransfers, &invalidTransfers, combinedFee, commitmentTokenIndex)
		if err != nil {
			return nil, nil, nil, nil, err
		}
		if !ok {
			continue
		}

		addedPubKeyIDs = append(addedPubKeyIDs, *pubKeyID)
		if uint32(len(appliedTransfers)) == cfg.TxsPerCommitment {
			break
		}
	}

	if len(appliedTransfers) > 0 {
		feeReceiverStateID, err = ApplyFee(storage, cfg.FeeReceiverPubKeyID, *commitmentTokenIndex, *combinedFee)
		if err != nil {
			return nil, nil, nil, nil, err
		}
	}

	return appliedTransfers, invalidTransfers, addedPubKeyIDs, feeReceiverStateID, nil
}

func ApplyCreate2TransfersForSync(
	storage *st.Storage,
	transfers []models.Create2Transfer,
	pubKeyIDs []uint32,
	cfg *config.RollupConfig,
) (
	appliedTransfers []models.Create2Transfer,
	invalidTransfers []models.Create2Transfer,
	err error,
) {
	if len(transfers) == 0 {
		return
	}
	if len(transfers) != len(pubKeyIDs) {
		return nil, nil, ErrInvalidSliceLength
	}

	appliedTransfers = make([]models.Create2Transfer, 0, cfg.TxsPerCommitment)
	combinedFee := models.NewUint256(0)

	commitmentTokenIndex, err := getTokenIndex(storage, transfers[0].FromStateID)
	if err != nil {
		return nil, nil, err
	}

	for i := range transfers {
		transfer := &transfers[i]

		_, err = handleApplyC2T(storage, transfer, pubKeyIDs[i], &appliedTransfers, &invalidTransfers, combinedFee, commitmentTokenIndex)
		if err != nil {
			return nil, nil, err
		}

		if uint32(len(appliedTransfers)) == cfg.TxsPerCommitment {
			break
		}
	}

	if len(appliedTransfers) > 0 {
		_, err = ApplyFee(storage, cfg.FeeReceiverPubKeyID, *commitmentTokenIndex, *combinedFee)
		if err != nil {
			return nil, nil, err
		}
	}

	return appliedTransfers, invalidTransfers, nil
}

func getTokenIndex(storage *st.Storage, stateID uint32) (*models.Uint256, error) {
	senderLeaf, err := storage.GetStateLeaf(stateID)
	if err != nil {
		return nil, err
	}
	return &senderLeaf.TokenIndex, nil
}

func handleApplyC2T(
	storage *st.Storage,
	transfer *models.Create2Transfer,
	pubKeyID uint32,
	appliedTxs, invalidTxs *[]models.Create2Transfer,
	combinedFee, tokenIndex *models.Uint256,
) (bool, error) {
	transferError, appError := ApplyCreate2Transfer(storage, transfer, pubKeyID, *tokenIndex)
	if appError != nil {
		return false, appError
	}
	if transferError != nil {
		logAndSaveTransactionError(storage, &transfer.TransactionBase, transferError)
		*invalidTxs = append(*invalidTxs, *transfer)
		return false, nil
	}

	*appliedTxs = append(*appliedTxs, *transfer)
	*combinedFee = *combinedFee.Add(&transfer.Fee)
	return true, nil
}
