package testutils

import "github.com/Worldcoin/hubble-commander/models"

func GetFourDeposits() []models.PendingDeposit {
	deposits := make([]models.PendingDeposit, 4)
	for i := range deposits {
		deposits[i] = models.PendingDeposit{
			ID: models.DepositID{
				SubtreeID:    models.MakeUint256(1),
				DepositIndex: models.MakeUint256(uint64(i))},
			ToPubKeyID: 1,
			TokenID:    models.MakeUint256(0),
			L2Amount:   models.MakeUint256(10000000000),
		}
	}
	return deposits
}
