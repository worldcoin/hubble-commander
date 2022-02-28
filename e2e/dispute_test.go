//go:build e2e
// +build e2e

package e2e

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/e2e/setup"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/suite"
)

type DisputesE2ETestSuite struct {
	setup.E2ETestSuite
}

func (s *DisputesE2ETestSuite) SetupTest() {
	cfg := config.GetConfig()
	cfg.Rollup.MinTxsPerCommitment = 32
	cfg.Rollup.MaxTxsPerCommitment = 32
	cfg.Rollup.MinCommitmentsPerBatch = 1

	s.SetupTestEnvironment(cfg, nil)
}

func (s *DisputesE2ETestSuite) TestDisputes() {
	s.testDisputeSignatureTransfer()
	s.testDisputeSignatureC2T()
	s.testDisputeSignatureMM()

	s.testDisputeTransitionTransfer()
	s.testDisputeTransitionC2T()
	s.testDisputeTransitionMM()
}

func (s *DisputesE2ETestSuite) TestMeasureDisputesGasUsage() {
	s.measureDisputeSignatureTransfer()
	s.measureDisputeSignatureC2T()
	s.measureDisputeSignatureMM()

	s.SubmitTxBatchAndWait(func() common.Hash {
		return s.SendNTransactions(32, dto.Transfer{
			FromStateID: ref.Uint32(1),
			ToStateID:   ref.Uint32(2),
			Amount:      models.NewUint256(90),
			Fee:         models.NewUint256(10),
			Nonce:       models.NewUint256(0),
		})
	})

	s.measureDisputeTransitionTransfer()
	s.measureDisputeTransitionC2T()
	s.measureDisputeTransitionMM()
}

func (s *DisputesE2ETestSuite) testDisputeSignatureTransfer() {
	s.requireRollbackCompleted(func() {
		s.sendNTransfersBatchWithInvalidSignature(1, 1)
	})

	s.requireBatchesCount(1)
}

func (s *DisputesE2ETestSuite) testDisputeSignatureC2T() {
	s.requireRollbackCompleted(func() {
		s.sendNC2TsBatchWithInvalidSignature(1, 1)
	})

	s.requireBatchesCount(1)
}

func (s *DisputesE2ETestSuite) testDisputeSignatureMM() {
	s.requireRollbackCompleted(func() {
		s.sendNMMsBatchWithInvalidSignature(1, 1)
	})

	s.requireBatchesCount(1)
}

func (s *DisputesE2ETestSuite) testDisputeTransitionTransfer() {
	s.SubmitTxBatchAndWait(func() common.Hash {
		return s.SendNTransactions(32, dto.Transfer{
			FromStateID: ref.Uint32(1),
			ToStateID:   ref.Uint32(2),
			Amount:      models.NewUint256(90),
			Fee:         models.NewUint256(10),
			Nonce:       models.NewUint256(0),
		})
	})

	s.requireRollbackCompleted(func() {
		s.sendNTransfersBatchWithInvalidStateRoot(2, 1)
	})

	s.requireBatchesCount(2)

	s.SubmitTxBatchAndWait(func() common.Hash {
		return s.SendNTransactions(32, dto.Transfer{
			FromStateID: ref.Uint32(1),
			ToStateID:   ref.Uint32(2),
			Amount:      models.NewUint256(90),
			Fee:         models.NewUint256(10),
			Nonce:       models.NewUint256(32),
		})
	})
}

func (s *DisputesE2ETestSuite) testDisputeTransitionC2T() {
	s.requireRollbackCompleted(func() {
		s.sendNC2TsBatchWithInvalidStateRoot(3, 1)
	})

	s.requireBatchesCount(3)

	s.SubmitTxBatchAndWait(func() common.Hash {
		firstC2TWallet := s.Wallets[len(s.Wallets)-32]
		return s.SendNTransactions(32, dto.Create2Transfer{
			FromStateID: ref.Uint32(1),
			ToPublicKey: firstC2TWallet.PublicKey(),
			Amount:      models.NewUint256(90),
			Fee:         models.NewUint256(10),
			Nonce:       models.NewUint256(64),
		})
	})
}

func (s *DisputesE2ETestSuite) testDisputeTransitionMM() {
	s.requireRollbackCompleted(func() {
		s.sendNMMsBatchWithInvalidStateRoot(4, 1)
	})

	s.requireBatchesCount(4)

	s.SubmitTxBatchAndWait(func() common.Hash {
		return s.SendNTransactions(32, dto.MassMigration{
			FromStateID: ref.Uint32(1),
			SpokeID:     ref.Uint32(1),
			Amount:      models.NewUint256(90),
			Fee:         models.NewUint256(10),
			Nonce:       models.NewUint256(96),
		})
	})
}

func (s *DisputesE2ETestSuite) measureDisputeSignatureTransfer() {
	s.requireRollbackCompleted(func() {
		s.sendNTransfersBatchWithInvalidSignature(1, s.Cfg.Rollup.MaxTxsPerCommitment)
	})

	s.requireBatchesCount(1)
}

func (s *DisputesE2ETestSuite) measureDisputeSignatureC2T() {
	s.requireRollbackCompleted(func() {
		s.sendNC2TsBatchWithInvalidSignature(1, s.Cfg.Rollup.MaxTxsPerCommitment)
	})

	s.requireBatchesCount(1)
}

func (s *DisputesE2ETestSuite) measureDisputeSignatureMM() {
	s.requireRollbackCompleted(func() {
		s.sendNMMsBatchWithInvalidSignature(1, s.Cfg.Rollup.MaxTxsPerCommitment)
	})

	s.requireBatchesCount(1)
}

func (s *DisputesE2ETestSuite) measureDisputeTransitionTransfer() {
	s.requireRollbackCompleted(func() {
		s.sendNTransfersBatchWithInvalidStateRoot(2, s.Cfg.Rollup.MaxTxsPerCommitment)
	})

	s.requireBatchesCount(2)
}

func (s *DisputesE2ETestSuite) measureDisputeTransitionC2T() {
	s.requireRollbackCompleted(func() {
		s.sendNC2TsBatchWithInvalidStateRoot(2, s.Cfg.Rollup.MaxTxsPerCommitment)
	})

	s.requireBatchesCount(2)
}

func (s *DisputesE2ETestSuite) measureDisputeTransitionMM() {
	s.requireRollbackCompleted(func() {
		s.sendNMMsBatchWithInvalidStateRoot(2, s.Cfg.Rollup.MaxTxsPerCommitment)
	})

	s.requireBatchesCount(2)
}

func TestDisputesE2ETestSuite(t *testing.T) {
	suite.Run(t, new(DisputesE2ETestSuite))
}
