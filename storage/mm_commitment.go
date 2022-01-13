package storage

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/stored"
)

func (s *CommitmentStorage) AddMMCommitment(commitment *models.MMCommitment) error {
	return s.database.Badger.Insert(commitment.ID, stored.MakeCommitmentFromMMCommitment(commitment))
}
