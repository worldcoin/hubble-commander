package disputer

import (
	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
)

// Other test suites encapsulate DisputeSignatureTestSuite.
// Don't add any tests on DisputeSignatureTestSuite to avoid repeated runs.
type DisputeSignatureTestSuite struct {
	TestSuiteWithContexts
	domain *bls.Domain
}

func (s *DisputeSignatureTestSuite) setupTest() {
	s.NotNil(s.client) // make sure TestSuiteWithContexts.SetupTest was called before

	var err error
	s.domain, err = s.client.GetDomain()
	s.NoError(err)
}

func (s *DisputeSignatureTestSuite) setUserStatesAndAddAccounts() []bls.Wallet {
	wallets := setUserStates(s.Assertions, s.disputeCtx, s.domain)
	for i := range wallets {
		err := s.storage.AccountTree.SetSingle(&models.AccountLeaf{
			PubKeyID:  uint32(i),
			PublicKey: *wallets[i].PublicKey(),
		})
		s.NoError(err)
	}
	return wallets
}

func (s *DisputeSignatureTestSuite) disputeSignature(
	batch *eth.DecodedBatch,
	transfers models.GenericTransactionArray,
) error {
	proofs, err := s.syncCtx.StateMerkleProofs(transfers)
	s.NoError(err)

	return s.disputeCtx.DisputeSignature(batch, 0, proofs)
}
