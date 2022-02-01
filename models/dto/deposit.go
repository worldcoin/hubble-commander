package dto

import "github.com/Worldcoin/hubble-commander/models"

type DepositID struct {
	SubtreeID    models.Uint256
	DepositIndex models.Uint256
}

type Deposit struct {
	ID         DepositID
	ToPubKeyID uint32
	TokenID    models.Uint256
	L2Amount   models.Uint256
}

func MakeDeposit(pendingDeposit *models.PendingDeposit) Deposit {
	return Deposit{
		ID: DepositID{
			SubtreeID:    pendingDeposit.ID.SubtreeID,
			DepositIndex: pendingDeposit.ID.DepositIndex,
		},
		ToPubKeyID: pendingDeposit.ToPubKeyID,
		TokenID:    pendingDeposit.TokenID,
		L2Amount:   pendingDeposit.L2Amount,
	}
}

func MakeDeposits(deposits []models.PendingDeposit) []Deposit {
	dtoDeposits := make([]Deposit, 0, len(deposits))

	for i := range deposits {
		dtoDeposits = append(dtoDeposits, MakeDeposit(&deposits[i]))
	}

	return dtoDeposits
}
