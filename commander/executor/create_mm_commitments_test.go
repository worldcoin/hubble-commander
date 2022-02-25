package executor

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/Worldcoin/hubble-commander/utils/merkletree"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type MMCommitmentsTestSuite struct {
	testSuiteWithTxsContext
	maxTxBytesInCommitment int
}

func (s *MMCommitmentsTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *MMCommitmentsTestSuite) SetupTest() {
	s.testSuiteWithTxsContext.SetupTestWithConfig(batchtype.MassMigration, &config.RollupConfig{
		MinTxsPerCommitment:    1,
		MaxTxsPerCommitment:    4,
		FeeReceiverPubKeyID:    2,
		MinCommitmentsPerBatch: 1,
		MaxCommitmentsPerBatch: 1,
	})
	s.maxTxBytesInCommitment = encoder.MassMigrationForCommitmentLength * int(s.cfg.MaxTxsPerCommitment)

	err := populateAccounts(s.storage.Storage, genesisBalances)
	s.NoError(err)
}

func (s *MMCommitmentsTestSuite) TestCreateCommitments_ReturnsCorrectMetaAndWithdrawRoot() {
	massMigrations := testutils.GenerateValidMassMigrations(2)
	initTxs(s.Assertions, s.txsCtx, massMigrations)

	withdrawRoot := s.generateWithdrawRoot(massMigrations)

	commitments, err := s.txsCtx.CreateCommitments()
	s.NoError(err)
	s.Len(commitments, 1)
	s.Len(commitments, 1)

	targetSpokeID := 2
	targetTokenID := models.MakeUint256(0)
	commitment := commitments[0].ToMMCommitmentWithTxs()
	targetFeeReceiver := commitment.Meta.FeeReceiver
	totalAmount := massMigrations[0].Amount.Add(&massMigrations[1].Amount)
	s.Equal(withdrawRoot, commitment.WithdrawRoot)
	s.EqualValues(targetSpokeID, commitment.Meta.SpokeID)
	s.Equal(targetTokenID, commitment.Meta.TokenID)
	s.Equal(*totalAmount, commitment.Meta.Amount)
	s.Equal(targetFeeReceiver, commitment.Meta.FeeReceiver)
}

func TestMMCommitmentsTestSuite(t *testing.T) {
	suite.Run(t, new(MMCommitmentsTestSuite))
}

func (s *MMCommitmentsTestSuite) generateWithdrawRoot(massMigrations []models.MassMigration) common.Hash {
	hashes := make([]common.Hash, 0, len(massMigrations))

	for i := range massMigrations {
		senderLeaf, err := s.storage.StateTree.Leaf(massMigrations[i].FromStateID)
		s.NoError(err)

		hash, err := encoder.HashUserState(&models.UserState{
			PubKeyID: senderLeaf.PubKeyID,
			TokenID:  senderLeaf.TokenID,
			Balance:  massMigrations[i].Amount,
			Nonce:    models.MakeUint256(0),
		})
		s.NoError(err)
		hashes = append(hashes, *hash)
	}

	merkleTree, err := merkletree.NewMerkleTree(hashes)
	s.NoError(err)

	return merkleTree.Root()
}
