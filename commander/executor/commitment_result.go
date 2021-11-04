package executor

import "github.com/Worldcoin/hubble-commander/models"

type CreateCommitmentResult interface {
	PendingTxs() models.GenericTransactionArray
	PendingAccounts() []models.AccountLeaf
	Commitment() *models.CommitmentWithTxs
}

type CreateTransferCommitmentResult struct {
	newPendingTxs models.GenericTransactionArray
	commitment    *models.CommitmentWithTxs
}

func (c *CreateTransferCommitmentResult) PendingTxs() models.GenericTransactionArray {
	return c.newPendingTxs
}

func (c *CreateTransferCommitmentResult) PendingAccounts() []models.AccountLeaf {
	return []models.AccountLeaf{}
}

func (c *CreateTransferCommitmentResult) Commitment() *models.CommitmentWithTxs {
	return c.commitment
}

type CreateC2TCommitmentResult struct {
	newPendingTxs   models.GenericTransactionArray
	pendingAccounts []models.AccountLeaf
	commitment      *models.CommitmentWithTxs
}

func (c *CreateC2TCommitmentResult) PendingTxs() models.GenericTransactionArray {
	return c.newPendingTxs
}

func (c *CreateC2TCommitmentResult) PendingAccounts() []models.AccountLeaf {
	return c.pendingAccounts
}

func (c *CreateC2TCommitmentResult) Commitment() *models.CommitmentWithTxs {
	return c.commitment
}
