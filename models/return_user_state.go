package models

// ReturnUserState for database queries
type ReturnUserState struct {
	MerklePath MerklePath `db:"merkle_path"`
	UserState
}
