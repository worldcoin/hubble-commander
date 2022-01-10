package models

type PendingStakeWithdrawal struct {
	BatchID           Uint256
	FinalisationBlock uint32 `badgerhold:"index"`
}
