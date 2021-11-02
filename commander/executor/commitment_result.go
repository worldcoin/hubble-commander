package executor

import "github.com/Worldcoin/hubble-commander/models"

type CreateCommitmentResult interface {
	PendingTxs() models.GenericTransactionArray
	PendingAccounts() []models.AccountLeaf
	Commitment() *models.TxCommitmentWithTxs
}

type CreateTransferCommitmentResult struct {
	newPendingTxs models.GenericTransactionArray
	commitment    *models.TxCommitmentWithTxs
}

func (c *CreateTransferCommitmentResult) PendingTxs() models.GenericTransactionArray {
	return c.newPendingTxs
}

func (c *CreateTransferCommitmentResult) PendingAccounts() []models.AccountLeaf {
	return []models.AccountLeaf{}
}

func (c *CreateTransferCommitmentResult) Commitment() *models.TxCommitmentWithTxs {
	return c.commitment
}

type CreateC2TCommitmentResult struct {
	newPendingTxs   models.GenericTransactionArray
	pendingAccounts []models.AccountLeaf
	commitment      *models.TxCommitmentWithTxs
}

func (c *CreateC2TCommitmentResult) PendingTxs() models.GenericTransactionArray {
	return c.newPendingTxs
}

func (c *CreateC2TCommitmentResult) PendingAccounts() []models.AccountLeaf {
	return c.pendingAccounts
}

func (c *CreateC2TCommitmentResult) Commitment() *models.TxCommitmentWithTxs {
	return c.commitment
}
