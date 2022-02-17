package models

import (
	"bytes"
)

type CommitmentSlot struct {
	BatchID           Uint256
	IndexInBatch      uint8
	IndexInCommitment uint8
}

const CommitmentSlotLength = 32 + 1 + 1

func NewCommitmentSlot(commitmentID CommitmentID, indexInCommitment uint8) *CommitmentSlot {
	return &CommitmentSlot{
		BatchID:           commitmentID.BatchID,
		IndexInBatch:      commitmentID.IndexInBatch,
		IndexInCommitment: indexInCommitment,
	}
}

func (k *CommitmentSlot) CommitmentID() *CommitmentID {
	return &CommitmentID{
		BatchID:      k.BatchID,
		IndexInBatch: k.IndexInBatch,
	}
}

func (k *CommitmentSlot) Bytes() []byte {
	var buf bytes.Buffer
	buf.Grow(CommitmentSlotLength)

	buf.Write(k.BatchID.Bytes())
	buf.WriteByte(k.IndexInBatch)
	buf.WriteByte(k.IndexInCommitment)

	return buf.Bytes()
}

func (k *CommitmentSlot) SetBytes(data []byte) error {
	if len(data) != CommitmentSlotLength {
		return ErrInvalidLength
	}

	k.BatchID.SetBytes(data[0:32])
	k.IndexInBatch = data[32]
	k.IndexInCommitment = data[33]
	return nil
}
