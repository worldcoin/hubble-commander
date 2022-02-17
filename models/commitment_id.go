package models

const CommitmentIDDataLength = 33

type CommitmentID struct {
	// GetTransactionIDsByBatchIDs assumes BatchID is the first field
	BatchID      Uint256
	IndexInBatch uint8
}

func (c *CommitmentID) Bytes() []byte {
	encoded := make([]byte, CommitmentIDDataLength)
	copy(encoded[0:32], c.BatchID.Bytes())
	encoded[32] = c.IndexInBatch

	return encoded
}

func (c *CommitmentID) SetBytes(data []byte) error {
	if len(data) != CommitmentIDDataLength {
		return ErrInvalidLength
	}

	c.BatchID.SetBytes(data[0:32])
	c.IndexInBatch = data[32]
	return nil
}
