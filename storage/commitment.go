package storage

import "github.com/Worldcoin/hubble-commander/models"

func (s *Storage) AddCommitment(commitment *models.Commitment) error {
	_, err := s.QB.Insert("commitment").
		Values(
			commitment.LeafHash,
			commitment.BodyHash,
			commitment.AccountTreeRoot,
			commitment.CombinedSignature,
			commitment.FeeReceiver,
			commitment.Transactions,
		).
		RunWith(s.DB).
		Exec()

	return err
}
