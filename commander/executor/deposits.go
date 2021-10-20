package executor

import (
	"time"

	"github.com/Worldcoin/hubble-commander/commander/prover"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func (c *ExecutionContext) GetVacancyProof(startStateID uint32, subtreeDepth uint8) (*models.SubtreeVacancyProof, error) {
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

func (c *ExecutionContext) CreateAndSubmitDepositBatch() error {
	startTime := time.Now()
	batch, err := c.NewPendingBatch(batchtype.Deposit)
	if err != nil {
		return errors.WithStack(err)
	}

	vacancyProof, err := c.ExecuteDeposits()
	if err != nil {
		return errors.WithStack(err)
	}

	err = c.SubmitDepositBatch(batch, vacancyProof)
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

func (c *ExecutionContext) ExecuteDeposits() (*models.SubtreeVacancyProof, error) {
	depositSubTree, err := c.storage.GetFirstPendingDepositSubTree()
	if st.IsNotFoundError(err) {
		//TODO-dep: return error that will be omitted
		return nil, nil
	}
	if err != nil {
		return nil, errors.WithStack(err)
	}

	startStateID, vacancyProof, err := c.getDepositSubtreeVacancyProof()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	err = c.ApplyDeposits(*startStateID, depositSubTree.Deposits)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	//TODO-dep: create and insert commitment after merging branch that introduce it
	err = c.storage.DeletePendingDepositSubTrees(depositSubTree.ID)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return vacancyProof, nil
}

func (c *ExecutionContext) getDepositSubtreeVacancyProof() (*uint32, *models.SubtreeVacancyProof, error) {
	subtreeDepth, err := c.client.GetDepositSubtreeDepth()
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

func (c *ExecutionContext) SubmitDepositBatch(batch *models.Batch, vacancyProof *models.SubtreeVacancyProof) error {
	proverCtx := prover.NewContext(c.storage)
	commitmentInclusionProof, err := proverCtx.PreviousBatchCommitmentInclusionProof(batch.ID)
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
