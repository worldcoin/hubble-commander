package executor

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/Worldcoin/hubble-commander/utils/merkletree"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type MassMigrationCommitmentsTestSuite struct {
	*require.Assertions
	suite.Suite
	storage                *st.TestStorage
	cfg                    *config.RollupConfig
	txsCtx                 *TxsContext
	maxTxBytesInCommitment int
}

func (s *MassMigrationCommitmentsTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *MassMigrationCommitmentsTestSuite) SetupTest() {
	var err error
	s.storage, err = st.NewTestStorage()
	s.NoError(err)
	s.cfg = &config.RollupConfig{
		MinTxsPerCommitment:    1,
		MaxTxsPerCommitment:    4,
		FeeReceiverPubKeyID:    2,
		MinCommitmentsPerBatch: 1,
		MaxCommitmentsPerBatch: 1,
	}
	s.maxTxBytesInCommitment = encoder.MassMigrationForCommitmentLength * int(s.cfg.MaxTxsPerCommitment)

	err = populateAccounts(s.storage.Storage, genesisBalances)
	s.NoError(err)

	executionCtx := NewTestExecutionContext(s.storage.Storage, eth.DomainOnlyTestClient, s.cfg)
	s.txsCtx = NewTestTxsContext(executionCtx, batchtype.MassMigration)
}

func (s *MassMigrationCommitmentsTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *MassMigrationCommitmentsTestSuite) TestCreateCommitments_ReturnsCorrectMeta() {
	massMigrations := testutils.GenerateValidMassMigrations(2)
	s.addMassMigrations(massMigrations)

	batchData, err := s.txsCtx.CreateCommitments()
	s.NoError(err)
	s.Len(batchData.Commitments(), 1)
	s.Len(batchData.Metas(), 1)

	targetSpokeID := 2
	targetTokenID := models.MakeUint256(0)
	targetFeeReceiver := batchData.Commitments()[0].FeeReceiver
	totalAmount := massMigrations[0].Amount.Add(&massMigrations[1].Amount)
	s.EqualValues(targetSpokeID, batchData.Metas()[0].SpokeID)
	s.Equal(targetTokenID, batchData.Metas()[0].TokenID)
	s.Equal(*totalAmount, batchData.Metas()[0].Amount)
	s.Equal(targetFeeReceiver, batchData.Metas()[0].FeeReceiver)
}

func (s *MassMigrationCommitmentsTestSuite) TestCreateCommitments_ReturnsCorrectWithdrawRoot() {
	massMigrations := testutils.GenerateValidMassMigrations(2)
	s.addMassMigrations(massMigrations)

	withdrawRoot := s.generateWithdrawRoot(massMigrations)

	batchData, err := s.txsCtx.CreateCommitments()
	s.NoError(err)
	s.Len(batchData.Commitments(), 1)
	s.Len(batchData.WithdrawRoots(), 1)
	s.Equal(withdrawRoot, batchData.WithdrawRoots()[0])
}

func TestMassMigrationCommitmentsTestSuite(t *testing.T) {
	suite.Run(t, new(MassMigrationCommitmentsTestSuite))
}

func (s *MassMigrationCommitmentsTestSuite) addMassMigrations(massMigrations []models.MassMigration) {
	err := s.storage.BatchAddMassMigration(massMigrations)
	s.NoError(err)
}

func (s *MassMigrationCommitmentsTestSuite) generateWithdrawRoot(massMigrations []models.MassMigration) common.Hash {
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
