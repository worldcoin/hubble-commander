package executor

import (
	"context"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
)

var ErrNotEnoughDeposits = NewRollupError("not enough deposits")

func (c *DepositsContext) CreateAndSubmitBatch(ctx context.Context) (*models.Batch, *int, error) {
	batch, err := c.NewPendingBatch(batchtype.Deposit)
	if err != nil {
		return nil, nil, err
	}

	vacancyProof, err := c.createCommitment(ctx, batch.ID)
	if err != nil {
		return nil, nil, err
	}

	err = c.SubmitBatch(batch, vacancyProof)
	if err != nil {
		return nil, nil, err
	}

	return batch, ref.Int(1), nil
}

func (c *DepositsContext) createCommitment(ctx context.Context, batchID models.Uint256) (*models.SubtreeVacancyProof, error) {
	_, span := otel.Tracer("depositsContext").Start(ctx, "createCommitment")
	defer span.End()

	depositSubtree, err := c.storage.GetFirstPendingDepositSubtree()
	if st.IsNotFoundError(err) {
		return nil, errors.WithStack(ErrNotEnoughDeposits)
	}
	if err != nil {
		return nil, err
	}

	vacancyProof, err := c.executeDeposits(depositSubtree)
	if err != nil {
		return nil, err
	}

	err = c.addCommitment(batchID, depositSubtree)
	if err != nil {
		return nil, err
	}

	return vacancyProof, nil
}

func (c *DepositsContext) executeDeposits(depositSubtree *models.PendingDepositSubtree) (*models.SubtreeVacancyProof, error) {
	startStateID, vacancyProof, err := c.getDepositSubtreeVacancyProof()
	if err != nil {
		return nil, err
	}

	err = c.Applier.ApplyDeposits(*startStateID, depositSubtree.Deposits)
	if err != nil {
		return nil, err
	}

	err = c.storage.RemovePendingDepositSubtrees(depositSubtree.ID)
	if err != nil {
		return nil, err
	}
	return vacancyProof, nil
}

func (c *DepositsContext) addCommitment(batchID models.Uint256, depositSubtree *models.PendingDepositSubtree) error {
	commitment, err := c.newCommitment(batchID, depositSubtree)
	if err != nil {
		return errors.WithStack(err)
	}
	return c.storage.AddCommitment(commitment)
}

func (c *DepositsContext) getDepositSubtreeVacancyProof() (*uint32, *models.SubtreeVacancyProof, error) {
	subtreeDepth, err := c.client.GetMaxSubtreeDepthParam()
	if err != nil {
		return nil, nil, err
	}

	startStateID, err := c.storage.StateTree.NextVacantSubtree(*subtreeDepth)
	if err != nil {
		return nil, nil, err
	}

	vacancyProof, err := c.proverCtx.GetVacancyProof(*startStateID, *subtreeDepth)
	if err != nil {
		return nil, nil, err
	}
	return startStateID, vacancyProof, nil
}
