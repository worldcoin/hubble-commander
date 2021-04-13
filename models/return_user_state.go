package models

type ReturnUserState struct {
	StateID uint32
	UserState
}

type ReturnUserState2 struct {
	MerklePath MerklePath `db:"merkle_path"`
	UserState
}
