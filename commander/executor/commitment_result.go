package executor

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/common"
)

type CreateCommitmentResult interface {
	AppliedTxs() models.GenericTransactionArray
	PendingAccounts() []models.AccountLeaf
	Commitment() *models.TxCommitmentWithTxs
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

func (c *CreateTransferCommitmentResult) Commitment() *models.TxCommitmentWithTxs {
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

func (c *CreateC2TCommitmentResult) Commitment() *models.TxCommitmentWithTxs {
	return c.commitment
}

type CreateMassMigrationCommitmentResult struct {
	appliedTxs models.MassMigrationArray
	commitment *models.TxCommitmentWithTxs
}

func (c *CreateMassMigrationCommitmentResult) AppliedTxs() models.GenericTransactionArray {
	return c.appliedTxs
}

func (c *CreateMassMigrationCommitmentResult) PendingAccounts() []models.AccountLeaf {
	return []models.AccountLeaf{}
}

func (c *CreateMassMigrationCommitmentResult) Commitment() *models.TxCommitmentWithTxs {
	return c.commitment
}

type BatchData interface {
	Commitments() []models.TxCommitmentWithTxs
	Metas() []models.MassMigrationMeta
	WithdrawRoots() []common.Hash
	Len() int
	AddCommitment(commitment *models.TxCommitmentWithTxs)
	AddMeta(meta *models.MassMigrationMeta)
	AddWithdrawRoot(withdrawRoot common.Hash)
}

type TxBatchData struct {
	commitments []models.TxCommitmentWithTxs
}

func (c *TxBatchData) Commitments() []models.TxCommitmentWithTxs {
	return c.commitments
}

func (c *TxBatchData) Metas() []models.MassMigrationMeta {
	panic("Meta cannot be invoked on TxBatchData")
}

func (c *TxBatchData) WithdrawRoots() []common.Hash {
	panic("WithdrawRoots cannot be invoked on TxBatchData")
}

func (c *TxBatchData) Len() int {
	return len(c.commitments)
}

func (c *TxBatchData) AddCommitment(commitment *models.TxCommitmentWithTxs) {
	c.commitments = append(c.commitments, *commitment)
}

func (c *TxBatchData) AddMeta(_ *models.MassMigrationMeta) {
	panic("AddMeta cannot be invoked on TxBatchData")
}

func (c *TxBatchData) AddWithdrawRoot(_ common.Hash) {
	panic("AddWithdrawRoot cannot be invoked on TxBatchData")
}

type MassMigrationBatchData struct {
	commitments   []models.TxCommitmentWithTxs
	metas         []models.MassMigrationMeta
	withdrawRoots []common.Hash
}

func (c *MassMigrationBatchData) Commitments() []models.TxCommitmentWithTxs {
	return c.commitments
}

func (c *MassMigrationBatchData) Metas() []models.MassMigrationMeta {
	return c.metas
}

func (c *MassMigrationBatchData) WithdrawRoots() []common.Hash {
	return c.withdrawRoots
}

func (c *MassMigrationBatchData) Len() int {
	return len(c.commitments)
}

func (c *MassMigrationBatchData) AddCommitment(commitment *models.TxCommitmentWithTxs) {
	c.commitments = append(c.commitments, *commitment)
}

func (c *MassMigrationBatchData) AddMeta(meta *models.MassMigrationMeta) {
	c.metas = append(c.metas, *meta)
}

func (c *MassMigrationBatchData) AddWithdrawRoot(withdrawRoot common.Hash) {
	c.withdrawRoots = append(c.withdrawRoots, withdrawRoot)
}
