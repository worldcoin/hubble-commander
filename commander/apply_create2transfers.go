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
	events, unsubscribe, err := client.WatchRegistrations(&bind.WatchOpts{})
	if err != nil {
		return
	}
	defer unsubscribe()

	appliedTransfers = make([]models.Create2Transfer, 0, cfg.TxsPerCommitment)
	addedPubKeyIDs = make([]uint32, 0, cfg.TxsPerCommitment)
	combinedFee := models.MakeUint256(0)

	senderLeaf, err := storage.GetStateLeaf(transfers[0].FromStateID)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	commitmentTokenIndex := senderLeaf.TokenIndex
	var pubKeyID *uint32

	for i := range transfers {
		transfer := transfers[i]

		pubKeyID, err = getOrRegisterPubKeyID(storage, client, events, &transfer, commitmentTokenIndex)
		if err != nil {
			return nil, nil, nil, nil, err
		}

		transferError, appError := ApplyCreate2Transfer(storage, &transfer, *pubKeyID, commitmentTokenIndex)
		if appError != nil {
			return nil, nil, nil, nil, appError
		}
		if transferError != nil {
			logAndSaveTransactionError(storage, &transfer.TransactionBase, transferError)
			invalidTransfers = append(invalidTransfers, transfer)
			continue
		}

		addedPubKeyIDs = append(addedPubKeyIDs, *pubKeyID)
		appliedTransfers = append(appliedTransfers, transfer)
		combinedFee = *combinedFee.Add(&transfer.Fee)

		if uint32(len(appliedTransfers)) == cfg.TxsPerCommitment {
			break
		}
	}

	if len(appliedTransfers) > 0 {
		feeReceiverStateID, err = ApplyFee(storage, cfg.FeeReceiverPubKeyID, commitmentTokenIndex, combinedFee)
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

	appliedTransfers = make([]models.Create2Transfer, 0, cfg.TxsPerCommitment)
	combinedFee := models.MakeUint256(0)

	senderLeaf, err := storage.GetStateLeaf(transfers[0].FromStateID)
	if err != nil {
		return nil, nil, err
	}
	commitmentTokenIndex := senderLeaf.TokenIndex

	for i := range transfers {
		transfer := transfers[i]
		pubKeyID := pubKeyIDs[i]

		transferError, appError := ApplyCreate2Transfer(storage, &transfer, pubKeyID, commitmentTokenIndex)
		if appError != nil {
			return nil, nil, appError
		}
		if transferError != nil {
			logAndSaveTransactionError(storage, &transfer.TransactionBase, transferError)
			invalidTransfers = append(invalidTransfers, transfer)
			continue
		}

		appliedTransfers = append(appliedTransfers, transfer)
		combinedFee = *combinedFee.Add(&transfer.Fee)

		if uint32(len(appliedTransfers)) == cfg.TxsPerCommitment {
			break
		}
	}

	if len(appliedTransfers) > 0 {
		_, err = ApplyFee(storage, cfg.FeeReceiverPubKeyID, commitmentTokenIndex, combinedFee)
		if err != nil {
			return nil, nil, err
		}
	}

	return appliedTransfers, invalidTransfers, nil
}
