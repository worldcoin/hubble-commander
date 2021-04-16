package bls

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var data = []byte("0xdeadbeef")

type SignatureTestSuite struct {
	*require.Assertions
	suite.Suite
	wallet    *Wallet
	signature *Signature
}

func (s *SignatureTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *SignatureTestSuite) SetupTest() {
	wallet, err := NewRandomWallet(testDomain)
	s.NoError(err)
	s.wallet = wallet

	s.signature, err = wallet.Sign(data)
	s.NoError(err)
}

func (s *SignatureTestSuite) TestVerify() {
	isValid, err := s.signature.Verify(data, s.wallet.PublicKey())
	s.NoError(err)
	s.True(isValid)
}

func (s *SignatureTestSuite) TestNewSignatureFromBytes() {
	bytes := s.signature.Bytes()

	signatureFromBytes, err := NewSignatureFromBytes(bytes, testDomain)
	s.NoError(err)

	isValid, err := signatureFromBytes.Verify(data, s.wallet.PublicKey())
	s.NoError(err)
	s.True(isValid)
}

func (s *SignatureTestSuite) TestNewWallet() {
	wallet, err := NewRandomWallet(testDomain)
	s.NoError(err)
	secretKey, _ := wallet.Bytes()
	fmt.Println(len(secretKey))

	wallet2, err := NewWallet(secretKey, testDomain)
	s.NoError(err)

	sig, err := wallet.Sign([]byte{1,2,3})
	s.NoError(err)
	pk := wallet2.PublicKey()

	isValid, err := sig.Verify([]byte{1,2,3}, pk)
	s.NoError(err)
	s.True(isValid)
}

func TestSignatureTestSuite(t *testing.T) {
	suite.Run(t, new(SignatureTestSuite))
}
