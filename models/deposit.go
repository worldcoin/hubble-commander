package models

type DepositID struct {
	BlockNumber uint32
	LogIndex    uint32
}

type Deposit struct {
	ID                   DepositID
	ToPubKeyID           uint32
	TokenID              Uint256
	L2Amount             Uint256
	IncludedInCommitment *CommitmentID
}
