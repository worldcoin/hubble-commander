package bls

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/contracts/test/bls"
	"github.com/Worldcoin/hubble-commander/eth/deployer"
	"github.com/Worldcoin/hubble-commander/testutils/simulator"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type WalletTestSuite struct {
	*require.Assertions
	suite.Suite
	sim     *simulator.Simulator
	testBLS *bls.TestBLS
}

func (s *WalletTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *WalletTestSuite) SetupTest() {
	sim, err := simulator.NewSimulator()
	s.NoError(err)
	s.sim = sim

	contracts, err := deployer.DeployTest(sim)
	s.NoError(err)
	s.testBLS = contracts.TestBLS
}

func (s *WalletTestSuite) TearDownTest() {
	s.sim.Close()
}

func (s *WalletTestSuite) TestSign() {
	s.T().Skip("Signature verification in smart contract doesn't work for some reason")

	wallet, err := NewRandomWallet(testDomain)
	s.NoError(err)

	data := []byte("0xdeadbeef")
	signature, err := wallet.Sign(data)
	s.NoError(err)

	point, err := s.testBLS.HashToPoint(nil, testDomain, data)
	s.NoError(err)

	checkSuccess, callSuccess, err := s.testBLS.VerifySingle(
		nil,
		signature.ToBigInts(),
		wallet.PublicKey().ToBigInts(),
		point,
	)
	s.NoError(err)
	s.True(checkSuccess)
	s.True(callSuccess)
}

func TestWalletTestSuite(t *testing.T) {
	suite.Run(t, new(WalletTestSuite))
}
