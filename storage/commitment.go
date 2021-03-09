package storage

import "github.com/Worldcoin/hubble-commander/models"

func (s *Storage) AddCommitment(commitment *models.Commitment) error {
	_, err := s.DB.ExecBuilder(
		s.QB.Insert("commitment").
			Values(
				commitment.LeafHash,
				commitment.PostStateRoot,
				commitment.BodyHash,
				commitment.AccountTreeRoot,
				commitment.CombinedSignature,
				commitment.FeeReceiver,
				commitment.Transactions,
			),
	)
	return err
}
