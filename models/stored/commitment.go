package stored

import (
	"github.com/Worldcoin/hubble-commander/models"
)

var CommitmentPrefix = models.GetBadgerHoldPrefix(Commitment{})

type Commitment struct {
	models.CommitmentBase
	Body CommitmentBody
}

func MakeCommitmentFromTxCommitment(c *models.TxCommitment) Commitment {
	return Commitment{
		CommitmentBase: c.CommitmentBase,
		Body: &TxCommitmentBody{
			FeeReceiver:       c.FeeReceiver,
			CombinedSignature: c.CombinedSignature,
			BodyHash:          c.BodyHash,
		},
	}
}

func MakeCommitmentFromMMCommitment(c *models.MMCommitment) Commitment {
	return Commitment{
		CommitmentBase: c.CommitmentBase,
		Body: &MMCommitmentBody{
			CombinedSignature: c.CombinedSignature,
			BodyHash:          c.BodyHash,
			Meta:              *c.Meta,
			WithdrawRoot:      c.WithdrawRoot,
		},
	}
}

func MakeCommitmentFromDepositCommitment(c *models.DepositCommitment) Commitment {
	return Commitment{
		CommitmentBase: c.CommitmentBase,
		Body: &DepositCommitmentBody{
			SubtreeID:   c.SubtreeID,
			SubtreeRoot: c.SubtreeRoot,
			Deposits:    c.Deposits,
		},
	}
}

func (c *Commitment) Bytes() []byte {
	b := make([]byte, c.BytesLen())
	copy(b[:models.CommitmentBaseDataLength], c.CommitmentBase.Bytes())
	copy(b[models.CommitmentBaseDataLength:], c.Body.Bytes())
	return b
}

func (c *Commitment) SetBytes(data []byte) error {
	err := c.CommitmentBase.SetBytes(data[:models.CommitmentBaseDataLength])
	if err != nil {
		return err
	}
	return c.setBodyBytes(data[models.CommitmentBaseDataLength:])
}

func (c *Commitment) setBodyBytes(data []byte) error {
	body, err := NewCommitmentBody(c.Type)
	if err != nil {
		return err
	}
	err = body.SetBytes(data)
	if err != nil {
		return err
	}
	c.Body = body
	return nil
}

func (c *Commitment) BytesLen() int {
	return models.CommitmentBaseDataLength + c.Body.BytesLen()
}

func (c *Commitment) ToTxCommitment() *models.TxCommitment {
	txCommitmentBody, ok := c.Body.(*TxCommitmentBody)
	if !ok {
		panic("invalid TxCommitment body type")
	}

	return &models.TxCommitment{
		CommitmentBase:    c.CommitmentBase,
		FeeReceiver:       txCommitmentBody.FeeReceiver,
		CombinedSignature: txCommitmentBody.CombinedSignature,
		BodyHash:          txCommitmentBody.BodyHash,
	}
}

func (c *Commitment) ToMMCommitment() *models.MMCommitment {
	mmCommitmentBody, ok := c.Body.(*MMCommitmentBody)
	if !ok {
		panic("invalid MMCommitment body type")
	}

	return &models.MMCommitment{
		CommitmentBase:    c.CommitmentBase,
		CombinedSignature: mmCommitmentBody.CombinedSignature,
		BodyHash:          mmCommitmentBody.BodyHash,
		Meta:              &mmCommitmentBody.Meta,
		WithdrawRoot:      mmCommitmentBody.WithdrawRoot,
	}
}

func (c *Commitment) ToDepositCommitment() *models.DepositCommitment {
	depositCommitmentBody, ok := c.Body.(*DepositCommitmentBody)
	if !ok {
		panic("invalid DepositCommitment body type")
	}

	return &models.DepositCommitment{
		CommitmentBase: c.CommitmentBase,
		SubtreeID:      depositCommitmentBody.SubtreeID,
		SubtreeRoot:    depositCommitmentBody.SubtreeRoot,
		Deposits:       depositCommitmentBody.Deposits,
	}
}
