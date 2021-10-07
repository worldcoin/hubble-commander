package executor

import (
	"time"

	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

var (
	ErrNotEnoughTxs = NewRollupError("not enough transactions")
	mockPublicKey   = models.PublicKey{1, 2, 3}
)

type FeeReceiver struct {
	StateID uint32
	TokenID models.Uint256
}

func (c *RollupContext) CreateCommitments() ([]models.Commitment, error) {
	pendingTxs, err := c.queryPendingTxs()
	if err != nil {
		return nil, err
	}

	commitmentID, err := c.NextCommitmentID()
	if err != nil {
		return nil, err
	}

	commitments := make([]models.Commitment, 0, c.cfg.MaxCommitmentsPerBatch)
	pendingAccounts := make([]models.AccountLeaf, 0)

	for i := uint8(0); len(commitments) != int(c.cfg.MaxCommitmentsPerBatch); i++ {
		var result CreateCommitmentResult
		commitmentID.IndexInBatch = i

		result, err = c.createCommitment(pendingTxs, commitmentID)
		if errors.Is(err, ErrNotEnoughTxs) {
			break
		}
		if err != nil {
			return nil, err
		}

		pendingTxs = result.PendingTxs()
		commitments = append(commitments, *result.Commitment())
		pendingAccounts = append(pendingAccounts, result.PendingAccounts()...)
	}

	if len(commitments) == 0 {
		return nil, errors.WithStack(ErrNotEnoughTxs)
	}

	err = c.registerPendingAccounts(pendingAccounts)
	if err != nil {
		return nil, err
	}

	return commitments, nil
}

func (c *RollupContext) createCommitment(pendingTxs models.GenericTransactionArray, commitmentID *models.CommitmentID) (
	CreateCommitmentResult, error,
) {
	startTime := time.Now()

	pendingTxs, err := c.refillPendingTxs(pendingTxs)
	if err != nil {
		return nil, err
	}

	feeReceiver, err := c.getCommitmentFeeReceiver()
	if err != nil {
		return nil, err
	}

	initialStateRoot, err := c.storage.StateTree.Root()
	if err != nil {
		return nil, err
	}

	executeResult, newPendingTxs, err := c.executeTxsForCommitment(pendingTxs, feeReceiver)
	if errors.Is(err, ErrNotEnoughTxs) {
		if revertErr := c.storage.StateTree.RevertTo(*initialStateRoot); revertErr != nil {
			return nil, revertErr
		}
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	commitment, err := c.BuildCommitment(executeResult, commitmentID, feeReceiver.StateID)
	if err != nil {
		return nil, err
	}

	log.Printf(
		"Created a %s commitment from %d transactions in %s",
		c.BatchType,
		executeResult.AppliedTxs().Len(),
		time.Since(startTime).Round(time.Millisecond).String(),
	)

	return c.Executor.NewCreateCommitmentResult(executeResult, commitment, newPendingTxs), nil
}

func (c *RollupContext) executeTxsForCommitment(pendingTxs models.GenericTransactionArray, feeReceiver *FeeReceiver) (
	result ExecuteTxsForCommitmentResult,
	newPendingTxs models.GenericTransactionArray,
	err error,
) {
	aggregateResult := c.Executor.NewExecuteTxsResult(c.cfg.MaxTxsPerCommitment)

	for {
		numNeededTxs := c.cfg.MaxTxsPerCommitment - uint32(aggregateResult.AppliedTxs().Len())
		executeTxsResult, err := c.ExecuteTxs(pendingTxs, numNeededTxs, feeReceiver)
		if err != nil {
			return nil, nil, err
		}

		aggregateResult.AddApplyResult(executeTxsResult)

		if aggregateResult.AppliedTxs().Len() == int(c.cfg.MaxTxsPerCommitment) {
			newPendingTxs = removeTxs(pendingTxs, aggregateResult.AllTxs())
			return c.Executor.NewExecuteTxsForCommitmentResult(aggregateResult), newPendingTxs, nil
		}

		morePendingTransfers, err := c.queryMorePendingTxs(aggregateResult.AppliedTxs())
		if errors.Is(err, ErrNotEnoughTxs) {
			newPendingTxs = removeTxs(pendingTxs, aggregateResult.AllTxs())
			return c.Executor.NewExecuteTxsForCommitmentResult(aggregateResult), newPendingTxs, nil
		}
		if err != nil {
			return nil, nil, err
		}
		pendingTxs = morePendingTransfers
	}
}

func (c *RollupContext) refillPendingTxs(pendingTxs models.GenericTransactionArray) (models.GenericTransactionArray, error) {
	if pendingTxs.Len() < int(c.cfg.MaxTxsPerCommitment) {
		return c.queryPendingTxs()
	}
	return pendingTxs, nil
}

func (c *RollupContext) queryPendingTxs() (models.GenericTransactionArray, error) {
	pendingTxs, err := c.Executor.GetPendingTxs(c.cfg.MaxCommitmentsPerBatch * c.cfg.MaxTxsPerCommitment)
	if err != nil {
		return nil, err
	}
	if pendingTxs.Len() < int(c.cfg.MinTxsPerCommitment) {
		return nil, errors.WithStack(ErrNotEnoughTxs)
	}
	return pendingTxs, nil
}

func (c *RollupContext) queryMorePendingTxs(appliedTxs models.GenericTransactionArray) (models.GenericTransactionArray, error) {
	numAppliedTransfers := uint32(appliedTxs.Len())
	pendingTransfers, err := c.Executor.GetPendingTxs(
		c.cfg.MaxCommitmentsPerBatch*c.cfg.MaxTxsPerCommitment + numAppliedTransfers,
	)
	if err != nil {
		return nil, err
	}
	pendingTransfers = removeTxs(pendingTransfers, appliedTxs)

	if pendingTransfers.Len() < int(c.cfg.MinTxsPerCommitment) {
		return nil, errors.WithStack(ErrNotEnoughTxs)
	}
	return pendingTransfers, nil
}

func (c *RollupContext) getCommitmentFeeReceiver() (*FeeReceiver, error) {
	commitmentTokenID := models.MakeUint256(0) // TODO support multiple tokens
	feeReceiverState, err := c.storage.GetFeeReceiverStateLeaf(c.cfg.FeeReceiverPubKeyID, commitmentTokenID)
	if err != nil {
		return nil, err
	}
	return &FeeReceiver{
		StateID: feeReceiverState.StateID,
		TokenID: feeReceiverState.TokenID,
	}, nil
}

func removeTxs(txList, toRemove models.GenericTransactionArray) models.GenericTransactionArray {
	outputIndex := 0
	for i := 0; i < txList.Len(); i++ {
		tx := txList.At(i)
		if !txExists(toRemove, tx) {
			txList.Set(outputIndex, tx)
			outputIndex++
		}
	}

	return txList.Slice(0, outputIndex)
}

func txExists(txList models.GenericTransactionArray, tx models.GenericTransaction) bool {
	for i := 0; i < txList.Len(); i++ {
		if txList.At(i).GetBase().Hash == tx.GetBase().Hash {
			return true
		}
	}
	return false
}

func (c *RollupContext) registerPendingAccounts(accounts []models.AccountLeaf) error {
	accounts, err := c.fillMissingAccounts(accounts)
	if err != nil {
		return err
	}
	publicKeys := make([]models.PublicKey, 0, st.AccountBatchSize)
	for i := range accounts {
		publicKeys = append(publicKeys, accounts[i].PublicKey)
		if len(publicKeys) == st.AccountBatchSize {
			_, err = c.client.RegisterBatchAccount(publicKeys)
			if err != nil {
				return err
			}
			publicKeys = make([]models.PublicKey, 0, st.AccountBatchSize)
		}
	}
	return nil
}

func (c *RollupContext) fillMissingAccounts(accounts []models.AccountLeaf) ([]models.AccountLeaf, error) {
	missingAccounts := st.AccountBatchSize - len(accounts)%st.AccountBatchSize
	if missingAccounts == st.AccountBatchSize {
		return accounts, nil
	}
	for i := 0; i < missingAccounts; i++ {
		accounts = append(accounts, models.AccountLeaf{
			PubKeyID:  accounts[len(accounts)-1].PubKeyID + 1,
			PublicKey: mockPublicKey,
		})
	}

	err := c.storage.AccountTree.SetInBatch(accounts[len(accounts)-missingAccounts:]...)
	if err != nil {
		return nil, err
	}
	return accounts, nil
}
