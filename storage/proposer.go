package storage

func (s *Storage) SetProposer(isProposer bool) {
	if s.isProposer != isProposer {
		s.isProposer = isProposer
	}
}

func (s *Storage) IsProposer() bool {
	return s.isProposer
}
