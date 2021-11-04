package applier

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/pkg/errors"
)

func (a *Applier) ApplyDeposits(startStateID uint32, deposits []models.PendingDeposit) error {
	for i := range deposits {
		_, err := a.storage.StateTree.Set(startStateID+uint32(i), &models.UserState{
			PubKeyID: deposits[i].ToPubKeyID,
			TokenID:  deposits[i].TokenID,
			Balance:  deposits[i].L2Amount,
			Nonce:    models.MakeUint256(0),
		})
		if err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}
