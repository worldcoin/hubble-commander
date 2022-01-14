package storage

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/stored"
)

func (s *CommitmentStorage) AddTxCommitment(commitment *models.TxCommitment) error {
	return s.database.Badger.Insert(commitment.ID, stored.MakeCommitmentFromTxCommitment(commitment))
}
