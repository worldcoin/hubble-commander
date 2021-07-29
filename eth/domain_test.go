package eth

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type DomainTestSuite struct {
	*require.Assertions
	suite.Suite
	client *TestClient
}

func (s *DomainTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *DomainTestSuite) SetupTest() {
	client, err := NewTestClient()
	s.NoError(err)
	s.client = client
}

func (s *DomainTestSuite) TearDownTest() {
	s.client.Close()
}

func (s *DomainTestSuite) TestGetDomain() {
	expectedDomain, err := s.client.Rollup.DomainSeparator(&bind.CallOpts{})
	s.NoError(err)

	domain, err := s.client.GetDomain()
	s.NoError(err)
	s.Equal(bls.Domain(expectedDomain), *domain)
	s.Equal(s.client.domain, domain)
}

func TestDomainTestSuite(t *testing.T) {
	suite.Run(t, new(DomainTestSuite))
}
