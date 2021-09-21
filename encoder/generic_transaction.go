package encoder

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

func HashGenericTransaction(tx models.GenericTransaction) (*common.Hash, error) {
	switch x := tx.(type) {
	case *models.Transfer:
		return HashTransfer(x)
	case *models.Create2Transfer:
		return HashCreate2Transfer(x)
	default:
		return nil, errors.Errorf("unsupported tx type: %s", tx.Type())
	}
}
