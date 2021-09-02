package api

import "github.com/Worldcoin/hubble-commander/models"

func (a *API) GetPublicKeyByIDProof(id uint32) (*models.PublicKeyProof, error) {
	leaf, err := a.storage.AccountTree.Leaf(id)

	if err != nil {
		return nil, err
	}

	witness, err := a.storage.AccountTree.GetWitness(id)

	if err != nil {
		return nil, err
	}

	return &models.PublicKeyProof{
		PublicKey: &leaf.PublicKey,
		Witness:   witness,
	}, nil
}
