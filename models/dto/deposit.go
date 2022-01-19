package dto

import "github.com/Worldcoin/hubble-commander/models"

type DepositID struct {
	SubtreeID    models.Uint256
	DepositIndex models.Uint256
}

type PendingDeposit struct {
	ID         DepositID
	ToPubKeyID uint32
	TokenID    models.Uint256
	L2Amount   models.Uint256
}

func MakePendingDeposit(pendingDeposit *models.PendingDeposit) PendingDeposit {
	return PendingDeposit{
		ID: DepositID{
			SubtreeID:    pendingDeposit.ID.SubtreeID,
			DepositIndex: pendingDeposit.ID.DepositIndex,
		},
		ToPubKeyID: pendingDeposit.ToPubKeyID,
		TokenID:    pendingDeposit.TokenID,
		L2Amount:   pendingDeposit.L2Amount,
	}
}

func modelsPendingDepositsToDTOPendingDeposits(deposits []models.PendingDeposit) []PendingDeposit {
	dtoDeposits := make([]PendingDeposit, 0, len(deposits))

	for i := range deposits {
		dtoDeposits = append(dtoDeposits, MakePendingDeposit(&deposits[i]))
	}

	return dtoDeposits
}
