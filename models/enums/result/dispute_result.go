package result

type DisputeResult uint8

// TODO add tests checking that we handle all remaining dispute reasons

const (
	Ok DisputeResult = iota
	InvalidTokenAmount
	NotEnoughTokenBalance
	BadFromTokenID
	BadToTokenID
	BadSignature
	MismatchedAmount
	BadWithdrawRoot
	BadCompression
	TooManyTx
	BadPrecompileCall
	NonexistentReceiver
)
