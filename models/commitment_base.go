package models

import (
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

const commitmentBaseDataLength = commitmentIDDataLength + 1 + 32 + 33

type CommitmentBase struct {
	ID            CommitmentID
	Type          batchtype.BatchType
	PostStateRoot common.Hash
	BodyHash      *common.Hash
}

func (c *CommitmentBase) Bytes() []byte {
	b := make([]byte, commitmentBaseDataLength)
	copy(b[0:33], c.ID.Bytes())
	b[33] = byte(c.Type)
	copy(b[34:66], c.PostStateRoot.Bytes())
	copy(b[66:99], EncodeHashPointer(c.BodyHash))

	return b
}

func (c *CommitmentBase) SetBytes(data []byte) error {
	if len(data) != commitmentBaseDataLength {
		return errors.WithStack(ErrInvalidLength)
	}
	err := c.ID.SetBytes(data[0:33])
	if err != nil {
		return err
	}

	c.Type = batchtype.BatchType(data[33])
	c.PostStateRoot.SetBytes(data[34:66])
	c.BodyHash = DecodeHashPointer(data[66:99])
	return nil
}
