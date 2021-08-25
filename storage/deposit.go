package storage

type DepositStorage struct {
	database *Database
}

func NewDepositStorage(database *Database) *DepositStorage {
	return &DepositStorage{
		database: database,
	}
}

func (s *DepositStorage) copyWithNewDatabase(database *Database) *DepositStorage {
	newDepositStorage := *s
	newDepositStorage.database = database

	return &newDepositStorage
}
