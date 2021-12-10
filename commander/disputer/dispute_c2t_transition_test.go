package disputer

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/commander/syncer"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/stretchr/testify/suite"
)

type DisputeCT2TransitionTestSuite struct {
	disputeTransitionTestSuite
}

func (s *DisputeCT2TransitionTestSuite) SetupTest() {
	s.SetupTestWithConfig(batchtype.Create2Transfer, &config.RollupConfig{
		MinCommitmentsPerBatch: 1,
		MaxCommitmentsPerBatch: 32,
		MinTxsPerCommitment:    1,
		MaxTxsPerCommitment:    2,
		DisableSignatures:      true,
	})
}

func (s *DisputeCT2TransitionTestSuite) TestDisputeTransition_RemovesInvalidBatch() {
	wallets := s.setUserStatesAndAddAccounts()

	commitmentTxs := []models.GenericTransactionArray{
		models.Create2TransferArray{
			testutils.MakeCreate2Transfer(0, nil, 0, 100, wallets[2].PublicKey()),
			testutils.MakeCreate2Transfer(1, nil, 0, 100, wallets[0].PublicKey()),
		},
		models.Create2TransferArray{
			testutils.MakeCreate2Transfer(2, nil, 0, 20, wallets[0].PublicKey()),
			testutils.MakeCreate2Transfer(2, nil, 1, 20, wallets[0].PublicKey()),
		},
	}

	s.submitInvalidBatch(commitmentTxs)

	remoteBatches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 1)

	proofs := s.getInvalidBatchStateProofs(remoteBatches[0])
	err = s.disputeCtx.DisputeTransition(remoteBatches[0].ToDecodedTxBatch(), 1, proofs)
	s.NoError(err)

	checkRemoteBatchAfterDispute(s.Assertions, s.client, &remoteBatches[0].GetBase().ID)
}

func (s *DisputeCT2TransitionTestSuite) TestDisputeTransition_FirstCommitment() {
	wallets := s.setUserStatesAndAddAccounts()

	commitmentTxs := []models.GenericTransactionArray{
		models.Create2TransferArray{testutils.MakeCreate2Transfer(0, nil, 0, 100, wallets[1].PublicKey())},
	}

	s.submitInvalidBatch(commitmentTxs)

	remoteBatches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 1)

	proofs := s.getInvalidBatchStateProofs(remoteBatches[0])
	err = s.disputeCtx.DisputeTransition(remoteBatches[0].ToDecodedTxBatch(), 0, proofs)
	s.NoError(err)

	checkRemoteBatchAfterDispute(s.Assertions, s.client, &remoteBatches[0].GetBase().ID)
}

func (s *DisputeCT2TransitionTestSuite) TestDisputeTransition_ValidBatch() {
	wallets := s.setUserStatesAndAddAccounts()
	tx := testutils.MakeCreate2Transfer(0, ref.Uint32(3), 0, 50, wallets[1].PublicKey())

	proofs := s.getValidBatchStateProofs(syncer.NewSyncedC2Ts(models.Create2TransferArray{tx}, []uint32{1}))

	tx.ToStateID = nil
	s.beginTransaction()
	defer s.commitTransaction()
	s.submitBatch(&tx)

	remoteBatches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 1)

	err = s.disputeCtx.DisputeTransition(remoteBatches[0].ToDecodedTxBatch(), 0, proofs)
	s.NoError(err)
	_, err = s.client.GetBatch(&remoteBatches[0].GetBase().ID)
	s.NoError(err)
}

func (s *DisputeCT2TransitionTestSuite) setUserStatesAndAddAccounts() []bls.Wallet {
	wallets := setUserStates(s.Assertions, s.disputeCtx, &bls.TestDomain)
	for i := range wallets {
		err := s.storage.AccountTree.SetSingle(&models.AccountLeaf{
			PubKeyID:  uint32(i),
			PublicKey: *wallets[i].PublicKey(),
		})
		s.NoError(err)
	}
	return wallets
}

func TestDisputeCT2TransitionTestSuite(t *testing.T) {
	suite.Run(t, new(DisputeCT2TransitionTestSuite))
}
