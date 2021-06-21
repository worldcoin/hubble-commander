package commander

import (
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/pkg/errors"
)

var ErrInvalidCommitmentRoot = errors.New("invalid commitment root of batch #0")

func verifyCommitmentRoot(storage *st.Storage, client *eth.Client) error {
	firstBatch, err := client.GetBatch(models.NewUint256(0))
	if err != nil {
		return err
	}
	stateRoot, err := st.NewStateTree(storage).Root()
	if err != nil {
		return err
	}

	zeroHash := st.GetZeroHash(0)
	commitmentRoot := utils.HashTwo(utils.HashTwo(*stateRoot, zeroHash), zeroHash)
	if *firstBatch.Hash != commitmentRoot {
		return ErrInvalidCommitmentRoot
	}
	return nil
}
