package bls

import "github.com/pkg/errors"

type Domain = [32]byte

var (
	testDomain = Domain{0x00, 0x00, 0x00, 0x00}

	ErrInvalidDomainLength = errors.New("invalid domain length")
)

func DomainFromBytes(data []byte) (*Domain, error) {
	if len(data) != 32 {
		return nil, ErrInvalidDomainLength
	}
	var domain [32]byte
	copy(domain[:], data)
	return &domain, nil
}
