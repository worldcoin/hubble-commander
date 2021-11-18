package executor

import (
	"strings"
	"time"

	"github.com/Worldcoin/hubble-commander/metrics"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func (c *TxsContext) CreateAndSubmitBatch() error {
	startTime := time.Now()
	batch, err := c.NewPendingBatch(c.BatchType)
	if err != nil {
		return err
	}

	commitments, err := c.CreateCommitments()
	if err != nil {
		return err
	}

	err = c.SubmitBatch(batch, commitments)
	if err != nil {
		return err
	}

	duration := measureBatchBuildAndSubmissionTime(startTime, c.commanderMetrics, batch.Type)
	logNewBatch(batch, len(commitments), duration)

	return nil
}

func measureBatchBuildAndSubmissionTime(
	start time.Time,
	commanderMetrics *metrics.CommanderMetrics,
	batchType batchtype.BatchType,
) time.Duration {
	duration := time.Since(start).Round(time.Millisecond)

	lowercaseBatchType := strings.ToLower(batchType.String())
	commanderMetrics.BatchBuildAndSubmissionTimes.WithLabelValues(lowercaseBatchType).Observe(float64(duration.Milliseconds()))
	return duration
}

func logNewBatch(batch *models.Batch, commitmentsCount int, duration time.Duration) {
	log.Printf(
		"Submitted a %s batch with %d commitment(s) on chain in %s. Batch ID: %d. Transaction hash: %v",
		batch.Type.String(),
		commitmentsCount,
		duration,
		batch.ID.Uint64(),
		batch.TransactionHash,
	)
}

func (c *ExecutionContext) NewPendingBatch(batchType batchtype.BatchType) (*models.Batch, error) {
	prevStateRoot, err := c.storage.StateTree.Root()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	batchID, err := c.storage.GetNextBatchID()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &models.Batch{
		ID:            *batchID,
		Type:          batchType,
		PrevStateRoot: prevStateRoot,
	}, nil
}
