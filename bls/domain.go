package bls

import "github.com/pkg/errors"

const DomainLength = 32

type Domain = [DomainLength]byte

var (
	testDomain = Domain{0x00, 0x00, 0x00, 0x00}

	ErrInvalidDomainLength = errors.New("invalid domain length")
)

func DomainFromBytes(data []byte) (*Domain, error) {
	if len(data) != DomainLength {
		return nil, ErrInvalidDomainLength
	}
	var domain [DomainLength]byte
	copy(domain[:], data)
	return &domain, nil
}
