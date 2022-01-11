package stored

const (
	sizeCommitment = 34
	sizeHash       = 32
	sizeTxType     = 1
	sizeU32        = 4
	sizeU256       = 32
	sizeSignature  = 64
	sizeTimestamp  = 16
	sizePendingTx  = (sizeHash + sizeTxType + sizeU32 + 3*sizeU256 + sizeSignature + sizeTimestamp)
)
