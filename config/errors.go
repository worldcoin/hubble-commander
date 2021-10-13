package config

import "fmt"

type ErrNonMatchingKeys struct {
	publicKey string
}

func NewErrNonMatchingKeys(publicKey string) *ErrNonMatchingKeys {
	return &ErrNonMatchingKeys{publicKey: publicKey}
}

func (e *ErrNonMatchingKeys) Error() string {
	return fmt.Sprintf("public key does not match the private key of a genesis account (%s)", e.publicKey)
}

func (e *ErrNonMatchingKeys) Is(other error) bool {
	otherErr, ok := other.(*ErrNonMatchingKeys)
	if !ok {
		return false
	}
	return e.publicKey == otherErr.publicKey
}
