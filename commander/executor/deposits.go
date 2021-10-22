package executor

import (
	"time"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

var ErrNotEnoughDeposits = NewRollupError("not enough deposits")

func (c *DepositContext) GetVacancyProof(startStateID uint32, subtreeDepth uint8) (*models.SubtreeVacancyProof, error) {
	path := models.MerklePath{
		Path:  startStateID >> subtreeDepth,
		Depth: st.StateTreeDepth - subtreeDepth,
	}
	witness, err := c.storage.StateTree.GetNodeWitness(path)
	if err != nil {
		return nil, err
	}

	return &models.SubtreeVacancyProof{
		PathAtDepth: path.Path,
		Witness:     witness,
	}, nil
}

func (c *DepositContext) CreateAndSubmitBatch() error {
	startTime := time.Now()
	batch, err := c.NewPendingBatch(batchtype.Deposit)
	if err != nil {
		return errors.WithStack(err)
	}

	vacancyProof, err := c.createCommitment(batch.ID)
	if err != nil {
		return errors.WithStack(err)
	}

	err = c.SubmitBatch(batch, vacancyProof)
	if err != nil {
		return errors.WithStack(err)
	}

	log.Printf(
		"Submitted a %s batch with %d commitment(s) on chain in %s. Batch ID: %d. Transaction hash: %v",
		batchtype.Deposit.String(),
		1,
		time.Since(startTime).Round(time.Millisecond).String(),
		batch.ID.Uint64(),
		batch.TransactionHash,
	)
	return nil
}

func (c *DepositContext) createCommitment(batchID models.Uint256) (*models.SubtreeVacancyProof, error) {
	depositSubtree, err := c.storage.GetFirstPendingDepositSubTree()
	if st.IsNotFoundError(err) {
		return nil, errors.WithStack(ErrNotEnoughDeposits)
	}
	if err != nil {
		return nil, errors.WithStack(err)
	}

	vacancyProof, err := c.ExecuteDeposits(depositSubtree)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	err = c.addCommitment(batchID, depositSubtree)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return vacancyProof, nil
}

func (c *DepositContext) ExecuteDeposits(depositSubtree *models.PendingDepositSubTree) (*models.SubtreeVacancyProof, error) {
	startStateID, vacancyProof, err := c.getDepositSubtreeVacancyProof()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	err = c.ApplyDeposits(*startStateID, depositSubtree.Deposits)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	err = c.storage.DeletePendingDepositSubTrees(depositSubtree.ID)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return vacancyProof, nil
}

func (c *DepositContext) addCommitment(batchID models.Uint256, depositSubtree *models.PendingDepositSubTree) error {
	commitment, err := c.newCommitment(batchID, depositSubtree)
	if err != nil {
		return err
	}
	return c.storage.AddDepositCommitment(commitment)
}

func (c *DepositContext) getDepositSubtreeVacancyProof() (*uint32, *models.SubtreeVacancyProof, error) {
	subtreeDepth, err := c.client.GetMaxSubTreeDepthParam()
	if err != nil {
		return nil, nil, err
	}

	startStateID, err := c.storage.StateTree.NextVacantSubtree(*subtreeDepth)
	if err != nil {
		return nil, nil, err
	}

	vacancyProof, err := c.GetVacancyProof(*startStateID, *subtreeDepth)
	if err != nil {
		return nil, nil, err
	}
	return startStateID, vacancyProof, nil
}

func (c *DepositContext) SubmitBatch(batch *models.Batch, vacancyProof *models.SubtreeVacancyProof) error {
	commitmentInclusionProof, err := c.proverCtx.PreviousBatchCommitmentInclusionProof(batch.ID)
	if err != nil {
		return err
	}

	tx, err := c.client.SubmitDeposits(commitmentInclusionProof, vacancyProof)
	if err != nil {
		return err
	}

	batch.TransactionHash = tx.Hash()
	return c.storage.AddBatch(batch)
}
