package bls

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
)

const DomainLength = 32

type Domain [DomainLength]byte

var (
	testDomain = Domain{0x00, 0x00, 0x00, 0x00}

	ErrInvalidDomainLength = errors.New("invalid domain length")
)

func (d *Domain) Bytes() []byte {
	return d[:]
}

func DomainFromBytes(data []byte) (*Domain, error) {
	if len(data) != DomainLength {
		return nil, ErrInvalidDomainLength
	}
	var domain Domain
	copy(domain[:], data)
	return &domain, nil
}

func (d Domain) MarshalText() ([]byte, error) {
	return hexutil.Bytes(d[:]).MarshalText()
}
