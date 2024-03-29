package eth

import (
	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

func (c *Client) GetDomain() (*bls.Domain, error) {
	if c.domain != nil {
		return c.domain, nil
	}

	domainSeparator, err := c.Rollup.DomainSeparator(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	domain := bls.Domain(domainSeparator)
	c.domain = &domain
	return &domain, nil
}
