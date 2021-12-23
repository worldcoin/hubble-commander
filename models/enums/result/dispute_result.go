package result

type DisputeResult uint8

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
