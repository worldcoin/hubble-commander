package executor

import (
	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/eth/rollup"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

type SyncTestSuite struct {
	TestSuiteWithExecutionContext
	transfer models.Transfer
	wallets  []bls.Wallet
	domain   *bls.Domain
}

func (s *SyncTestSuite) SetupSuite() {
	s.TestSuiteWithExecutionContext.SetupSuite()
	s.transfer = testutils.MakeTransfer(0, 1, 0, 400)
	s.setTransferHash(&s.transfer)
}

func (s *SyncTestSuite) SetupTest() {
	var err error
	s.client, err = eth.NewConfiguredTestClient(rollup.DeploymentConfig{
		Params: rollup.Params{
			MaxTxsPerCommit: models.NewUint256(1),
		},
	}, eth.ClientConfig{})
	s.NoError(err)

	s.cfg = &config.RollupConfig{
		MinCommitmentsPerBatch: 1,
		MaxCommitmentsPerBatch: 32,
		MinTxsPerCommitment:    1,
		MaxTxsPerCommitment:    1,
		DisableSignatures:      false,
	}

	s.domain, err = s.client.GetDomain()
	s.NoError(err)
	s.wallets = generateWallets(s.Assertions, s.domain, 2)
	s.setupDB()
}

func (s *SyncTestSuite) setupDB() {
	var err error
	s.storage, err = st.NewTestStorage()
	s.NoError(err)
	s.executionCtx = NewTestExecutionContext(s.storage.Storage, s.client.Client, s.cfg)

	seedDB(s.Assertions, s.storage.Storage, s.wallets)
}

func seedDB(s *require.Assertions, storage *st.Storage, wallets []bls.Wallet) {
	err := storage.AccountTree.SetSingle(&models.AccountLeaf{
		PubKeyID:  0,
		PublicKey: *wallets[0].PublicKey(),
	})
	s.NoError(err)

	err = storage.AccountTree.SetSingle(&models.AccountLeaf{
		PubKeyID:  1,
		PublicKey: *wallets[1].PublicKey(),
	})
	s.NoError(err)

	_, err = storage.StateTree.Set(0, &models.UserState{
		PubKeyID: 0,
		TokenID:  models.MakeUint256(0),
		Balance:  models.MakeUint256(1000),
		Nonce:    models.MakeUint256(0),
	})
	s.NoError(err)

	_, err = storage.StateTree.Set(1, &models.UserState{
		PubKeyID: 1,
		TokenID:  models.MakeUint256(0),
		Balance:  models.MakeUint256(0),
		Nonce:    models.MakeUint256(0),
	})
	s.NoError(err)
}

func (s *SyncTestSuite) createCommitmentWithEmptyTransactions(commitmentType txtype.TransactionType) models.Commitment {
	stateRoot, err := s.storage.StateTree.Root()
	s.NoError(err)

	feeReceiver, err := s.executionCtx.getCommitmentFeeReceiver()
	s.NoError(err)

	return models.Commitment{
		Type:              commitmentType,
		Transactions:      []byte{},
		FeeReceiver:       feeReceiver.StateID,
		CombinedSignature: models.Signature{},
		PostStateRoot:     *stateRoot,
	}
}

func (s *SyncTestSuite) syncAllBatches() {
	newRemoteBatches, err := s.client.GetAllBatches()
	s.NoError(err)

	for i := range newRemoteBatches {
		remoteBatch := &newRemoteBatches[i]
		err = s.executionCtx.SyncBatch(remoteBatch)
		s.NoError(err)
	}
}

func (s *SyncTestSuite) recreateDatabase() {
	err := s.storage.Teardown()
	s.NoError(err)
	s.setupDB()
}

func (s *SyncTestSuite) getAccountTreeRoot() common.Hash {
	rawAccountRoot, err := s.client.AccountRegistry.Root(nil)
	s.NoError(err)
	return common.BytesToHash(rawAccountRoot[:])
}

func (s *SyncTestSuite) setTransferHash(tx *models.Transfer) {
	hash, err := encoder.HashTransfer(tx)
	s.NoError(err)
	tx.Hash = *hash
}

func generateWallets(s *require.Assertions, domain *bls.Domain, walletsAmount int) []bls.Wallet {
	wallets := make([]bls.Wallet, 0, walletsAmount)
	for i := 0; i < walletsAmount; i++ {
		wallet, err := bls.NewRandomWallet(*domain)
		s.NoError(err)
		wallets = append(wallets, *wallet)
	}
	return wallets
}
