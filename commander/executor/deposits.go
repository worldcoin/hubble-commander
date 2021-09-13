package executor

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/storage"
)

func (c *ExecutionContext) GetVacancyProof(startStateID uint32, subtreeDepth uint8) (*models.SubtreeVacancyProof, error) {
	path := models.MerklePath{
		Path:  startStateID >> subtreeDepth,
		Depth: storage.StateTreeDepth - subtreeDepth,
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
