package client

import (
	"encoding/json"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
)

type PendingCommitment struct {
	Commitment   Commitment
	Transactions []Transaction
}

func (c *PendingCommitment) ToDTO() dto.PendingCommitment {
	return dto.PendingCommitment{
		Commitment:   c.Commitment.Parsed,
		Transactions: txsToTransactionArray(c.Transactions),
	}
}

type Commitment struct {
	Parsed models.Commitment
}

func (c *Commitment) UnmarshalJSON(bytes []byte) error {
	var rawCommitment struct {
		Type *batchtype.BatchType
	}
	err := json.Unmarshal(bytes, &rawCommitment)
	if err != nil {
		return err
	}

	if rawCommitment.Type == nil {
		return ErrMissingType
	}

	switch *rawCommitment.Type {
	case batchtype.Transfer, batchtype.Create2Transfer:
		return c.unmarshalTxCommitment(bytes)
	case batchtype.MassMigration:
		return c.unmarshalMMCommitment(bytes)
	case batchtype.Deposit:
		return c.unmarshalDepositCommitment(bytes)
	default:
		return ErrNotImplemented
	}
}

func (c *Commitment) unmarshalTxCommitment(bytes []byte) error {
	var commitment models.TxCommitment
	err := json.Unmarshal(bytes, &commitment)
	if err != nil {
		return err
	}
	c.Parsed = &commitment
	return nil
}

func (c *Commitment) unmarshalMMCommitment(bytes []byte) error {
	var commitment models.MMCommitment
	err := json.Unmarshal(bytes, &commitment)
	if err != nil {
		return err
	}
	c.Parsed = &commitment
	return nil
}

func (c *Commitment) unmarshalDepositCommitment(bytes []byte) error {
	var commitment models.DepositCommitment
	err := json.Unmarshal(bytes, &commitment)
	if err != nil {
		return err
	}
	c.Parsed = &commitment
	return nil
}
