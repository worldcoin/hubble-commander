package models

type ReturnUserState struct {
	MerklePath MerklePath `db:"merkle_path"`
	UserState
}
