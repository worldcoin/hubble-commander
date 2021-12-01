package executor

import "github.com/Worldcoin/hubble-commander/models"

type CreateCommitmentResult interface {
	AppliedTxs() models.GenericTransactionArray
	PendingAccounts() []models.AccountLeaf
	Commitment() *models.CommitmentWithTxs
}

type CreateTransferCommitmentResult struct {
	appliedTxs models.TransferArray
	commitment *models.CommitmentWithTxs
}

func (c *CreateTransferCommitmentResult) AppliedTxs() models.GenericTransactionArray {
	return c.appliedTxs
}

func (c *CreateTransferCommitmentResult) PendingAccounts() []models.AccountLeaf {
	return []models.AccountLeaf{}
}

func (c *CreateTransferCommitmentResult) Commitment() *models.CommitmentWithTxs {
	return c.commitment
}

type CreateC2TCommitmentResult struct {
	appliedTxs      models.Create2TransferArray
	pendingAccounts []models.AccountLeaf
	commitment      *models.CommitmentWithTxs
}

func (c *CreateC2TCommitmentResult) AppliedTxs() models.GenericTransactionArray {
	return c.appliedTxs
}

func (c *CreateC2TCommitmentResult) PendingAccounts() []models.AccountLeaf {
	return c.pendingAccounts
}

func (c *CreateC2TCommitmentResult) Commitment() *models.CommitmentWithTxs {
	return c.commitment
}

type CreateMassMigrationCommitmentResult struct {
	appliedTxs models.MassMigrationArray
	commitment *models.CommitmentWithTxs
}

func (c *CreateMassMigrationCommitmentResult) AppliedTxs() models.GenericTransactionArray {
	return c.appliedTxs
}

func (c *CreateMassMigrationCommitmentResult) PendingAccounts() []models.AccountLeaf {
	return []models.AccountLeaf{}
}

func (c *CreateMassMigrationCommitmentResult) Commitment() *models.CommitmentWithTxs {
	return c.commitment
}

type CreateCommitmentsResult interface {
	Commitments() []models.CommitmentWithTxs
	Metas() []models.MassMigrationMeta
	Len() int
	AddCommitment(commitment *models.CommitmentWithTxs)
	AddResult(createCommitmentResult CreateCommitmentResult)
}

type CreateTxCommitmentsResult struct {
	commitments []models.CommitmentWithTxs
}

func (c *CreateTxCommitmentsResult) Commitments() []models.CommitmentWithTxs {
	return c.commitments
}

func (c *CreateTxCommitmentsResult) Metas() []models.MassMigrationMeta {
	panic("Meta cannot be invoked on CreateTxCommitmentsResult")
}

func (c *CreateTxCommitmentsResult) Len() int {
	return len(c.commitments)
}

func (c *CreateTxCommitmentsResult) AddCommitment(commitment *models.CommitmentWithTxs) {
	c.commitments = append(c.commitments, *commitment)
}

func (c *CreateTxCommitmentsResult) AddResult(result CreateCommitmentResult) {
	c.AddCommitment(result.Commitment())
}

type CreateMassMigrationCommitmentsResult struct {
	commitments []models.CommitmentWithTxs
	metas       []models.MassMigrationMeta
}

func (c *CreateMassMigrationCommitmentsResult) Commitments() []models.CommitmentWithTxs {
	return c.commitments
}

func (c *CreateMassMigrationCommitmentsResult) Metas() []models.MassMigrationMeta {
	return c.metas
}

func (c *CreateMassMigrationCommitmentsResult) Len() int {
	return len(c.commitments)
}

func (c *CreateMassMigrationCommitmentsResult) AddCommitment(commitment *models.CommitmentWithTxs) {
	c.commitments = append(c.commitments, *commitment)
}

func (c *CreateMassMigrationCommitmentsResult) AddResult(result CreateCommitmentResult) {
	c.AddCommitment(result.Commitment())

	txs := result.AppliedTxs().ToMassMigrationArray()
	totalAmount := models.NewUint256(0)
	for i := range txs {
		txAmount := txs.At(i).GetAmount()
		totalAmount = totalAmount.Add(&txAmount)
	}

	c.metas = append(c.metas, models.MassMigrationMeta{
		SpokeID:     uint32(txs.At(0).ToMassMigration().SpokeID.Uint64()),
		TokenID:     commitmentTokenID, // TODO: support multiple tokens
		Amount:      *totalAmount,
		FeeReceiver: result.Commitment().FeeReceiver,
	})
}
