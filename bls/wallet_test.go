// +build hardhat

package bls

import (
	"encoding/hex"
	"testing"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/contracts/test/bls"
	"github.com/Worldcoin/hubble-commander/eth/deployer"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type WalletTestSuite struct {
	*require.Assertions
	suite.Suite
	testBLS *bls.TestBLS
}

func (s *WalletTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *WalletTestSuite) SetupTest() {
	cfg := config.GetTestConfig()
	s.NotNil(cfg.Ethereum, "This test must be run against hardhat node instance with gas estimator contract deployed")

	dep, err := deployer.NewRPCDeployer(cfg.Ethereum)
	s.NoError(err)

	opts := *dep.GetAccount()
	opts.GasLimit = 3_000_000
	_, _, testBLS, err := bls.DeployTestBLS(&opts, dep.GetBackend())
	s.NoError(err)
	s.testBLS = testBLS
}

func (s *WalletTestSuite) TestSign() {
	wallet, err := NewRandomWallet(testDomain)
	s.NoError(err)

	data, err := hex.DecodeString("deadbeef")
	s.NoError(err)

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
	s.True(callSuccess)
	s.True(checkSuccess)
}

func TestWalletTestSuite(t *testing.T) {
	suite.Run(t, new(WalletTestSuite))
}
