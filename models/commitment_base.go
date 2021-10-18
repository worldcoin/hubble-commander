package models

import (
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/ethereum/go-ethereum/common"
)

const commitmentIDDataLength = 33

type CommitmentBase struct {
	ID            CommitmentID
	Type          batchtype.BatchType
	PostStateRoot common.Hash
}

type CommitmentID struct {
	BatchID      Uint256
	IndexInBatch uint8
}

func (c *CommitmentID) Bytes() []byte {
	encoded := make([]byte, commitmentIDDataLength)
	copy(encoded[0:32], utils.PadLeft(c.BatchID.Bytes(), 32))
	encoded[32] = c.IndexInBatch

	return encoded
}

func (c *CommitmentID) SetBytes(data []byte) error {
	if len(data) != commitmentIDDataLength {
		return ErrInvalidLength
	}

	c.BatchID.SetBytes(data[0:32])
	c.IndexInBatch = data[32]
	return nil
}
