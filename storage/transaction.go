
func (s *Storage) GetLatestTransactionNonce(accountIndex uint32) (*models.Uint256, error) {
	res := make([]models.Uint256, 0, 1)
	err := s.DB.Query(
		s.QB.Select("transaction.nonce").
			From("transaction").
			Where(squirrel.Eq{"from_index": accountIndex}).
			OrderBy("nonce DESC").
			Limit(1),
	).Into(&res)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, ErrTransactionNotFound
	}
	return &res[0], nil
}
