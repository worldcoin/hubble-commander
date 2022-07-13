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
	cfg := config.GetCommanderConfigAndSetupLogger()
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
		postStateRoot := common.Hash{67, 176, 141, 0, 2, 50, 199, 28, 151, 99, 125, 127, 174, 183, 63, 187, 57, 94, 96, 143, 40, 65, 71, 129,
			140, 109, 96, 138, 47, 64, 15, 207}
		s.sendNTransfersBatchWithInvalidSignature(1, 1, postStateRoot)
	})

	s.requireBatchesCount(1)
}

func (s *DisputesE2ETestSuite) testDisputeSignatureC2T() {
	s.requireRollbackCompleted(func() {
		postStateRoot := common.Hash{148, 161, 17, 219, 203, 227, 135, 229, 169, 243, 249, 27, 68, 196, 71, 192, 229, 71, 220, 254, 108, 164,
			8, 136, 125, 30, 61, 226, 146, 163, 116, 94}
		s.sendNC2TsBatchWithInvalidSignature(1, 1, postStateRoot)
	})

	s.requireBatchesCount(1)
}

func (s *DisputesE2ETestSuite) testDisputeSignatureMM() {
	s.requireRollbackCompleted(func() {
		postStateRoot := common.Hash{161, 160, 52, 208, 71, 192, 169, 92, 255, 70, 215, 143, 101, 201, 129, 143, 184, 94, 24, 48, 137, 254,
			133, 57, 232, 166, 27, 150, 225, 125, 207, 107}
		s.sendNMMsBatchWithInvalidSignature(1, 1, postStateRoot)
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
		postStateRoot := common.Hash{218, 39, 24, 204, 77, 219, 251, 62, 246, 186, 68, 119, 175, 83, 17, 234, 153, 92, 73, 42, 200, 2, 98,
			27, 106, 20, 147, 37, 147, 100, 145, 233}
		s.sendNTransfersBatchWithInvalidSignature(1, s.Cfg.Rollup.MaxTxsPerCommitment, postStateRoot)
	})

	s.requireBatchesCount(1)
}

func (s *DisputesE2ETestSuite) measureDisputeSignatureC2T() {
	s.requireRollbackCompleted(func() {
		postStateRoot := common.Hash{42, 173, 143, 147, 46, 6, 72, 244, 166, 38, 240, 98, 63, 249, 81, 151, 53, 8, 196, 185, 69, 71, 5, 166,
			172, 117, 126, 138, 192, 163, 99, 127}
		s.sendNC2TsBatchWithInvalidSignature(1, s.Cfg.Rollup.MaxTxsPerCommitment, postStateRoot)
	})

	s.requireBatchesCount(1)
}

func (s *DisputesE2ETestSuite) measureDisputeSignatureMM() {
	s.requireRollbackCompleted(func() {
		postStateRoot := common.Hash{157, 167, 18, 112, 146, 192, 201, 9, 143, 73, 16, 150, 43, 69, 28, 236, 185, 148, 237, 111, 173, 121,
			11, 169, 195, 70, 158, 54, 157, 100, 207, 38}
		s.sendNMMsBatchWithInvalidSignature(1, s.Cfg.Rollup.MaxTxsPerCommitment, postStateRoot)
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
