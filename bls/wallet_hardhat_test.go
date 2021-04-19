// +build hardhat

package bls

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/contracts/test/bls"
	"github.com/Worldcoin/hubble-commander/eth/deployer"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type WalletHardhatTestSuite struct {
	*require.Assertions
	suite.Suite
	testBLS *bls.TestBLS
}

func (s *WalletHardhatTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())

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

func (s *WalletHardhatTestSuite) TestSign_VerifySingle() {
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
		signature.BigInts(),
		wallet.PublicKey().BigInts(),
		point,
	)
	s.NoError(err)
	s.True(callSuccess)
	s.True(checkSuccess)
}

func (s *WalletHardhatTestSuite) TestSign_VerifyMultiple() {
	hexStrings := []string{"deadbeef", "cafebabe", "baadf00d"}
	signatures := make([]*Signature, 0, 3)
	publicKeys := make([][4]*big.Int, 0, 3)
	dataPoints := make([][2]*big.Int, 0, 3)

	for _, str := range hexStrings {
		bytes, err := hex.DecodeString(str)
		s.NoError(err)

		wallet, err := NewRandomWallet(testDomain)
		s.NoError(err)

		signature, err := wallet.Sign(bytes)
		s.NoError(err)

		dataPoint, err := s.testBLS.HashToPoint(nil, testDomain, bytes)
		s.NoError(err)

		signatures = append(signatures, signature)
		publicKeys = append(publicKeys, wallet.PublicKey().BigInts())
		dataPoints = append(dataPoints, dataPoint)
	}
	aggregatedSignature := NewAggregatedSignature(signatures)

	checkSuccess, callSuccess, err := s.testBLS.VerifyMultiple(
		nil,
		aggregatedSignature.BigInts(),
		publicKeys,
		dataPoints,
	)
	s.NoError(err)
	s.True(callSuccess)
	s.True(checkSuccess)
}

func TestWalletHardhatTestSuite(t *testing.T) {
	suite.Run(t, new(WalletHardhatTestSuite))
}
