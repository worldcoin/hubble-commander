package storage

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/stored"
)

func (s *CommitmentStorage) addDepositCommitment(commitment *models.DepositCommitment) error {
	return s.database.Badger.Insert(commitment.ID, stored.MakeCommitmentFromDepositCommitment(commitment))
}
