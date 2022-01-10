package executor

import (
	"github.com/Worldcoin/hubble-commander/models"
)

type CreateCommitmentResult interface {
	AppliedTxs() models.GenericTransactionArray
	PendingAccounts() []models.AccountLeaf
	Commitment() models.CommitmentWithTxs
}

type CreateTransferCommitmentResult struct {
	appliedTxs models.TransferArray
	commitment *models.TxCommitmentWithTxs
}

func (c *CreateTransferCommitmentResult) AppliedTxs() models.GenericTransactionArray {
	return c.appliedTxs
}

func (c *CreateTransferCommitmentResult) PendingAccounts() []models.AccountLeaf {
	return []models.AccountLeaf{}
}

func (c *CreateTransferCommitmentResult) Commitment() models.CommitmentWithTxs {
	return c.commitment
}

type CreateC2TCommitmentResult struct {
	appliedTxs      models.Create2TransferArray
	pendingAccounts []models.AccountLeaf
	commitment      *models.TxCommitmentWithTxs
}

func (c *CreateC2TCommitmentResult) AppliedTxs() models.GenericTransactionArray {
	return c.appliedTxs
}

func (c *CreateC2TCommitmentResult) PendingAccounts() []models.AccountLeaf {
	return c.pendingAccounts
}

func (c *CreateC2TCommitmentResult) Commitment() models.CommitmentWithTxs {
	return c.commitment
}

type CreateMassMigrationCommitmentResult struct {
	appliedTxs models.MassMigrationArray
	commitment *models.MMCommitmentWithTxs
}

func (c *CreateMassMigrationCommitmentResult) AppliedTxs() models.GenericTransactionArray {
	return c.appliedTxs
}

func (c *CreateMassMigrationCommitmentResult) PendingAccounts() []models.AccountLeaf {
	return []models.AccountLeaf{}
}

func (c *CreateMassMigrationCommitmentResult) Commitment() models.CommitmentWithTxs {
	return c.commitment
}
