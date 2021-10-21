package models

import (
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/ethereum/go-ethereum/common"
)

const commitmentBaseDataLength = commitmentIDDataLength + 1 + 32

type CommitmentBase struct {
	ID            CommitmentID
	Type          batchtype.BatchType
	PostStateRoot common.Hash
}

func (c *CommitmentBase) Bytes() []byte {
	b := make([]byte, commitmentBaseDataLength)
	copy(b[0:33], c.ID.Bytes())
	b[33] = byte(c.Type)
	copy(b[34:66], c.PostStateRoot.Bytes())

	return b
}

func (c *CommitmentBase) SetBytes(data []byte) error {
	err := c.ID.SetBytes(data[0:33])
	if err != nil {
		return err
	}

	c.Type = batchtype.BatchType(data[33])
	c.PostStateRoot.SetBytes(data[34:66])
	return nil
}
