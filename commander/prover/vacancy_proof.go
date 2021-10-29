package prover

import (
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
)

func (c *Context) GetVacancyProof(startStateID uint32, subtreeDepth uint8) (*models.SubtreeVacancyProof, error) {
	path := models.MerklePath{
		Path:  startStateID >> subtreeDepth,
		Depth: st.StateTreeDepth - subtreeDepth,
	}
	witness, err := c.storage.StateTree.GetNodeWitness(path)
	if err != nil {
		return nil, err
	}

	return &models.SubtreeVacancyProof{
		PathAtDepth: path.Path,
		Witness:     witness,
	}, nil
}
