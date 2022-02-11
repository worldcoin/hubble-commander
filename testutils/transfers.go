package testutils

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/Worldcoin/hubble-commander/utils"
)

func MakeTransfer(from, to uint32, nonce, amount uint64) models.Transfer {
	return *NewTransfer(from, to, nonce, amount)
}

func NewTransfer(from, to uint32, nonce, amount uint64) *models.Transfer {
	return &models.Transfer{
		TransactionBase: models.TransactionBase{
			Hash:        utils.RandomHash(),
			TxType:      txtype.Transfer,
			FromStateID: from,
			Amount:      models.MakeUint256(amount),
			Fee:         models.MakeUint256(10),
			Nonce:       models.MakeUint256(nonce),
		},
		ToStateID: to,
	}
}

func MakeCreate2Transfer(from uint32, to *uint32, nonce, amount uint64, publicKey *models.PublicKey) models.Create2Transfer {
	return *NewCreate2Transfer(from, to, nonce, amount, publicKey)
}

func NewCreate2Transfer(from uint32, to *uint32, nonce, amount uint64, publicKey *models.PublicKey) *models.Create2Transfer {
	c2t := &models.Create2Transfer{
		TransactionBase: models.TransactionBase{
			Hash:        utils.RandomHash(),
			TxType:      txtype.Create2Transfer,
			FromStateID: from,
			Amount:      models.MakeUint256(amount),
			Fee:         models.MakeUint256(10),
			Nonce:       models.MakeUint256(nonce),
		},
		ToStateID: to,
	}
	if publicKey != nil {
		c2t.ToPublicKey = *publicKey
	}
	return c2t
}

func MakeMassMigration(from, spokeID uint32, nonce, amount uint64) models.MassMigration {
	return *NewMassMigration(from, spokeID, nonce, amount)
}

func NewMassMigration(from, spokeID uint32, nonce, amount uint64) *models.MassMigration {
	return &models.MassMigration{
		TransactionBase: models.TransactionBase{
			Hash:        utils.RandomHash(),
			TxType:      txtype.MassMigration,
			FromStateID: from,
			Amount:      models.MakeUint256(amount),
			Fee:         models.MakeUint256(10),
			Nonce:       models.MakeUint256(nonce),
		},
		SpokeID: spokeID,
	}
}

func GenerateValidCreate2Transfers(transfersAmount uint32) models.Create2TransferArray {
	transfers := make(models.Create2TransferArray, 0, transfersAmount)
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

func GenerateInvalidCreate2Transfers(transfersAmount uint64) models.Create2TransferArray {
	transfers := make(models.Create2TransferArray, 0, transfersAmount)
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

func GenerateValidTransfers(transfersAmount uint32) models.TransferArray {
	transfers := make([]models.Transfer, 0, transfersAmount)
	for i := 0; i < int(transfersAmount); i++ {
		transfer := models.Transfer{
			TransactionBase: models.TransactionBase{
				Hash:        utils.RandomHash(),
				TxType:      txtype.Transfer,
				FromStateID: 1,
				Amount:      models.MakeUint256(1),
				Fee:         models.MakeUint256(1),
				Nonce:       models.MakeUint256(uint64(i)),
			},
			ToStateID: 2,
		}
		transfers = append(transfers, transfer)
	}
	return transfers
}

func GenerateInvalidTransfers(transfersAmount uint64) []models.Transfer {
	transfers := make([]models.Transfer, 0, transfersAmount)
	for i := uint64(0); i < transfersAmount; i++ {
		transfer := models.Transfer{
			TransactionBase: models.TransactionBase{
				Hash:        utils.RandomHash(),
				TxType:      txtype.Transfer,
				FromStateID: 1,
				Amount:      models.MakeUint256(1_000_000),
				Fee:         models.MakeUint256(1),
				Nonce:       models.MakeUint256(0),
			},
			ToStateID: 2,
		}
		transfers = append(transfers, transfer)
	}
	return transfers
}

func GenerateValidMassMigrations(massMigrationsAmount uint32) models.MassMigrationArray {
	massMigrations := make([]models.MassMigration, 0, massMigrationsAmount)
	for i := 0; i < int(massMigrationsAmount); i++ {
		massMigration := models.MassMigration{
			TransactionBase: models.TransactionBase{
				Hash:        utils.RandomHash(),
				TxType:      txtype.MassMigration,
				FromStateID: 1,
				Amount:      models.MakeUint256(1),
				Fee:         models.MakeUint256(1),
				Nonce:       models.MakeUint256(uint64(i)),
			},
			SpokeID: 2,
		}
		massMigrations = append(massMigrations, massMigration)
	}
	return massMigrations
}
