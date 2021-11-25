package executor

import "github.com/Worldcoin/hubble-commander/models"

type CreateCommitmentResult interface {
	PendingAccounts() []models.AccountLeaf
	Commitment() *models.CommitmentWithTxs
}

type CreateTransferCommitmentResult struct {
	commitment *models.CommitmentWithTxs
}

func (c *CreateTransferCommitmentResult) PendingAccounts() []models.AccountLeaf {
	return []models.AccountLeaf{}
}

func (c *CreateTransferCommitmentResult) Commitment() *models.CommitmentWithTxs {
	return c.commitment
}

type CreateC2TCommitmentResult struct {
	pendingAccounts []models.AccountLeaf
	commitment      *models.CommitmentWithTxs
}

func (c *CreateC2TCommitmentResult) PendingAccounts() []models.AccountLeaf {
	return c.pendingAccounts
}

func (c *CreateC2TCommitmentResult) Commitment() *models.CommitmentWithTxs {
	return c.commitment
}
