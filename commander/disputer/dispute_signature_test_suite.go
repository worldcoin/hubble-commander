package disputer

import (
	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
)

// Other test suites encapsulate disputeSignatureTestSuite.
// Don't add any tests on disputeSignatureTestSuite to avoid repeated runs.
type disputeSignatureTestSuite struct {
	testSuiteWithContexts
	domain *bls.Domain
}

func (s *disputeSignatureTestSuite) setupTest() {
	s.NotNil(s.client) // make sure testSuiteWithContexts.SetupTest was called before

	var err error
	s.domain, err = s.client.GetDomain()
	s.NoError(err)
}

func (s *disputeSignatureTestSuite) setUserStatesAndAddAccounts() []bls.Wallet {
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

func (s *disputeSignatureTestSuite) disputeSignature(
	batch *eth.DecodedTxBatch,
	transfers models.GenericTransactionArray,
) error {
	proofs, err := s.syncCtx.StateMerkleProofs(transfers)
	s.NoError(err)

	return s.disputeCtx.DisputeSignature(batch, 0, proofs)
}
