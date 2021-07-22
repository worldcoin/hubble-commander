package executor

import (
	"github.com/Worldcoin/hubble-commander/models"
)

func (t *TransactionExecutor) getUserStateProof(stateID uint32) (*models.StateMerkleProof, error) {
	leaf, err := t.stateTree.Leaf(stateID)
	if err != nil {
		return nil, err
	}
	witness, err := t.stateTree.GetWitness(leaf.StateID)
	if err != nil {
		return nil, err
	}
	return &models.StateMerkleProof{
		UserState: &leaf.UserState,
		Witness:   witness,
	}, nil
}
