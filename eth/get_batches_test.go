package eth

import (
	"context"
	"math/big"
	"testing"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/contracts/rollup"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type GetBatchesTestSuite struct {
	*require.Assertions
	suite.Suite
	client      *TestClient
	commitments []models.Commitment
}

func (s *GetBatchesTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
	s.commitments = []models.Commitment{
		{
			ID:                1,
			Type:              txtype.Transfer,
			Transactions:      []uint8{0, 0, 0, 0, 0, 0, 0, 1, 32, 4, 0, 0},
			FeeReceiver:       0,
			CombinedSignature: *s.mockSignature(),
		},
		{
			ID:                2,
			Type:              txtype.Transfer,
			Transactions:      []uint8{0, 0, 1, 0, 0, 0, 0, 0, 32, 1, 0, 0},
			FeeReceiver:       0,
			CombinedSignature: *s.mockSignature(),
		},
	}
}

func (s *GetBatchesTestSuite) SetupTest() {
	client, err := NewTestClient()
	s.NoError(err)
	s.client = client
}

func (s *GetBatchesTestSuite) TearDownTest() {
	s.client.Close()
}

func (s *GetBatchesTestSuite) TestGetBatches() {
	finalisationBlocks, err := s.client.GetBlocksToFinalise()
	s.NoError(err)

	batch1, err := s.client.SubmitTransfersBatchAndWait([]models.Commitment{s.commitments[0]})
	s.NoError(err)
	_, err = s.client.SubmitTransfersBatchAndWait([]models.Commitment{s.commitments[1]})
	s.NoError(err)

	rawAccountRoot, err := s.client.AccountRegistry.Root(nil)
	s.NoError(err)
	accountRoot := common.BytesToHash(rawAccountRoot[:])

	batches, err := s.client.GetBatches(&bind.FilterOpts{
		Start: uint64(*batch1.FinalisationBlock - uint32(*finalisationBlocks) + 1),
	})
	s.NoError(err)
	s.Len(batches, 1)
	s.NotEqual(common.Hash{}, batches[0].TransactionHash)
	s.Equal(accountRoot, *batches[0].AccountTreeRoot)
}

func (s *GetBatchesTestSuite) TestGetBatchIfExists_BatchExists() {
	tx, err := s.client.SubmitTransfersBatch(s.commitments)
	s.NoError(err)
	s.client.Commit()

	transaction, _, err := s.client.ChainConnection.GetBackend().TransactionByHash(context.Background(), tx.Hash())
	s.NoError(err)

	event := &rollup.RollupNewBatch{
		BatchID:     big.NewInt(1),
		AccountRoot: s.getAccountRoot(),
		BatchType:   2,
	}

	batch, err := s.client.getBatchIfExists(event, transaction)
	s.NoError(err)
	s.Equal(models.MakeUint256(1), batch.ID)
	s.Len(batch.Commitments, len(s.commitments))
	s.EqualValues(event.AccountRoot, *batch.AccountTreeRoot)
}

func (s *GetBatchesTestSuite) TestGetBatchIfExists_BatchNotExists() {
	tx, err := s.client.SubmitTransfersBatch(s.commitments)
	s.NoError(err)
	s.client.Commit()

	transaction, _, err := s.client.ChainConnection.GetBackend().TransactionByHash(context.Background(), tx.Hash())
	s.NoError(err)

	event := &rollup.RollupNewBatch{
		BatchID:     big.NewInt(5),
		AccountRoot: s.getAccountRoot(),
		BatchType:   2,
	}

	batch, err := s.client.getBatchIfExists(event, transaction)
	s.Nil(batch)
	s.ErrorIs(err, errBatchNotExists)
}

func (s *GetBatchesTestSuite) TestGetBatchIfExists_DifferentBatchHash() {
	tx, err := s.client.SubmitTransfersBatch(s.commitments)
	s.NoError(err)
	s.client.Commit()

	transaction, _, err := s.client.ChainConnection.GetBackend().TransactionByHash(context.Background(), tx.Hash())
	s.NoError(err)

	event := &rollup.RollupNewBatch{
		BatchID:     big.NewInt(1),
		AccountRoot: [32]byte{1, 2, 3},
		BatchType:   2,
	}

	batch, err := s.client.getBatchIfExists(event, transaction)
	s.Nil(batch)
	s.ErrorIs(err, errBatchNotExists)
}

func (s *GetBatchesTestSuite) getAccountRoot() common.Hash {
	rawAccountRoot, err := s.client.AccountRegistry.Root(nil)
	s.NoError(err)
	return common.BytesToHash(rawAccountRoot[:])
}

func (s *GetBatchesTestSuite) mockSignature() *models.Signature {
	wallet, err := bls.NewRandomWallet(bls.Domain{1, 2, 3, 4})
	s.NoError(err)
	signature, err := wallet.Sign(utils.RandomBytes(4))
	s.NoError(err)
	return signature.ModelsSignature()
}

func TestGetBatchesTestSuite(t *testing.T) {
	suite.Run(t, new(GetBatchesTestSuite))
}
