package stored

import "github.com/Worldcoin/hubble-commander/models"

const (
	sizeCommitmentID    = models.CommitmentIDDataLength
	sizeCommitmentSlot  = sizeCommitmentID + 1
	sizeHash            = 32
	sizeTxType          = 1
	sizeU32             = 4
	sizeU256            = 32
	sizeSignature       = 64
	sizeTimestamp       = 16
	sizePendingTxNoBody = sizeHash + sizeTxType + sizeU32 + 3*sizeU256 + sizeSignature + sizeTimestamp
	sizeBatchedTxNoBody = sizePendingTxNoBody + sizeCommitmentSlot
)
