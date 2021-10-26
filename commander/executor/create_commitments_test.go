package executor

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type CreateCommitmentsTestSuite struct {
	testSuiteWithRollupContext
	wallets []bls.Wallet
}

func (s *CreateCommitmentsTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *CreateCommitmentsTestSuite) SetupTest() {
	s.testSuiteWithRollupContext.SetupTestWithConfig(batchtype.Transfer, config.RollupConfig{
		MinTxsPerCommitment:    2,
		MaxTxsPerCommitment:    32,
		MinCommitmentsPerBatch: 1,
		MaxCommitmentsPerBatch: 32,
	})

	domain, err := s.client.GetDomain()
	s.NoError(err)
	s.wallets = testutils.GenerateWallets(s.Assertions, domain, 2)

	s.addUserStates()
}

func (s *CreateCommitmentsTestSuite) TestCreateCommitments_DoesNotCreateCommitmentsWithLessTxsThenRequired() {
	validTransfer := testutils.MakeTransfer(1, 2, 0, 100)
	s.hashSignAndAddTransfer(&s.wallets[0], &validTransfer)
	invalidTransfer := testutils.MakeTransfer(2, 1, 1234, 100)
	s.hashSignAndAddTransfer(&s.wallets[1], &invalidTransfer)

	preStateRoot, err := s.storage.StateTree.Root()
	s.NoError(err)

	commitments, err := s.rollupCtx.CreateCommitments()
	s.Nil(commitments)
	s.ErrorIs(err, ErrNotEnoughTxs)

	postStateRoot, err := s.storage.StateTree.Root()
	s.NoError(err)

	s.Equal(preStateRoot, postStateRoot)
}

func (s *CreateCommitmentsTestSuite) TestCreateCommitments_ReturnsErrorIfCouldNotCreateEnoughCommitments() {
	s.cfg.MinTxsPerCommitment = 1
	s.cfg.MaxTxsPerCommitment = 1
	s.cfg.MinCommitmentsPerBatch = 2

	validTransfer := testutils.MakeTransfer(1, 2, 0, 100)
	s.hashSignAndAddTransfer(&s.wallets[0], &validTransfer)
	invalidTransfer := testutils.MakeTransfer(2, 1, 1234, 100)
	s.hashSignAndAddTransfer(&s.wallets[1], &invalidTransfer)

	commitments, err := s.rollupCtx.CreateCommitments()
	s.Nil(commitments)
	s.ErrorIs(err, ErrNotEnoughCommitments)
}

func (s *CreateCommitmentsTestSuite) hashSignAndAddTransfer(wallet *bls.Wallet, transfer *models.Transfer) {
	hash, err := encoder.HashTransfer(transfer)
	s.NoError(err)
	transfer.Hash = *hash

	encodedTransfer, err := encoder.EncodeTransferForSigning(transfer)
	s.NoError(err)
	signature, err := wallet.Sign(encodedTransfer)
	s.NoError(err)
	transfer.Signature = *signature.ModelsSignature()

	err = s.storage.AddTransfer(transfer)
	s.NoError(err)
}

func (s *CreateCommitmentsTestSuite) addUserStates() {
	feeReceiver := &models.UserState{
		PubKeyID: 0,
		TokenID:  models.MakeUint256(0),
		Balance:  models.MakeUint256(0),
		Nonce:    models.MakeUint256(0),
	}
	_, err := s.storage.StateTree.Set(0, feeReceiver)
	s.NoError(err)

	_, err = s.storage.StateTree.Set(1, &models.UserState{
		PubKeyID: 1,
		TokenID:  models.MakeUint256(0),
		Balance:  models.MakeUint256(1000),
		Nonce:    models.MakeUint256(0),
	})
	s.NoError(err)

	_, err = s.storage.StateTree.Set(2, &models.UserState{
		PubKeyID: 2,
		TokenID:  models.MakeUint256(0),
		Balance:  models.MakeUint256(1000),
		Nonce:    models.MakeUint256(0),
	})
	s.NoError(err)
}

func TestCreateCommitmentsTestSuite(t *testing.T) {
	suite.Run(t, new(CreateCommitmentsTestSuite))
}
