package models

// UserStateReceipt for database queries
type UserStateReceipt struct {
	MerklePath MerklePath `db:"merkle_path"`
	UserState
}
