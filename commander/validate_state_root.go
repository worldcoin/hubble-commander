package commander

import (
	"fmt"

	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/pkg/errors"
)

var ErrInvalidStateRoot = fmt.Errorf("current state tree root doesn't match latest commitment post state root")

func validateStateRoot(storage *st.Storage) error {
	latestCommitment, err := storage.GetLatestCommitment()
	if st.IsNotFoundError(err) {
		return nil
	}
	if err != nil {
		return err
	}
	stateRoot, err := storage.StateTree.Root()
	if err != nil {
		return err
	}
	if latestCommitment.PostStateRoot != *stateRoot {
		logLatestCommitment(latestCommitment)
		return errors.WithStack(ErrInvalidStateRoot)
	}
	return nil
}
