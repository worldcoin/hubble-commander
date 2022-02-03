package stored

const (
	sizeCommitment      = 34 // See {Encode,decode}CommitmentIDPointer
	sizeHash            = 32
	sizeTxType          = 1
	sizeU32             = 4
	sizeU256            = 32
	sizeSignature       = 64
	sizeTimestamp       = 16
	sizePendingTxNoBody = (sizeHash + sizeTxType + sizeU32 + 3*sizeU256 + sizeSignature + sizeTimestamp)
	sizeBatchedTxNoBody = sizePendingTxNoBody + sizeCommitment
)
