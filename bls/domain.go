package bls

import (
	"reflect"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
)

const DomainLength = 32

type Domain [DomainLength]byte

var (
	testDomain = Domain{0x00, 0x00, 0x00, 0x00}
	domainT    = reflect.TypeOf(Domain{})

	ErrInvalidDomainLength = errors.New("invalid domain length") // TODO-API here
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

func (d *Domain) UnmarshalJSON(input []byte) error {
	return hexutil.UnmarshalFixedJSON(domainT, input, d[:])
}

func (d Domain) MarshalText() ([]byte, error) {
	return hexutil.Bytes(d[:]).MarshalText()
}
