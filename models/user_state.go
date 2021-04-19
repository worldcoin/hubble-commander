package models

// UserStateWithID for database queries
type UserStateWithID struct {
	StateID uint32
	UserState
}
