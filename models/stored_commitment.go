package models

var StoredCommitmentPrefix = getBadgerHoldPrefix(StoredCommitment{})

type StoredCommitment struct {
	CommitmentBase
	Body StoredCommitmentBody
}

func MakeStoredCommitmentFromTxCommitment(c *TxCommitment) StoredCommitment {
	return StoredCommitment{
		CommitmentBase: c.CommitmentBase,
		Body: &StoredTxCommitmentBody{
			FeeReceiver:       c.FeeReceiver,
			CombinedSignature: c.CombinedSignature,
			Transactions:      c.Transactions,
		},
	}
}

func MakeStoredCommitmentFromDepositCommitment(c *DepositCommitment) StoredCommitment {
	return StoredCommitment{
		CommitmentBase: c.CommitmentBase,
		Body: &StoredDepositCommitmentBody{
			SubTreeID:   c.SubTreeID,
			SubTreeRoot: c.SubTreeRoot,
			Deposits:    c.Deposits,
		},
	}
}

func (c *StoredCommitment) Bytes() []byte {
	b := make([]byte, c.BytesLen())
	copy(b[:commitmentBaseDataLength], c.CommitmentBase.Bytes())
	copy(b[commitmentBaseDataLength:], c.Body.Bytes())
	return b
}

func (c *StoredCommitment) SetBytes(data []byte) error {
	err := c.CommitmentBase.SetBytes(data[:commitmentBaseDataLength])
	if err != nil {
		return err
	}
	return c.setBodyBytes(data[commitmentBaseDataLength:])
}

func (c *StoredCommitment) setBodyBytes(data []byte) error {
	body, err := NewStoredCommitmentBody(c.Type)
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

func (c *StoredCommitment) BytesLen() int {
	return commitmentBaseDataLength + c.Body.BytesLen()
}

func (c *StoredCommitment) ToTxCommitment() *TxCommitment {
	txCommitmentBody, ok := c.Body.(*StoredTxCommitmentBody)
	if !ok {
		panic("invalid TxCommitment body type")
	}

	return &TxCommitment{
		CommitmentBase:    c.CommitmentBase,
		FeeReceiver:       txCommitmentBody.FeeReceiver,
		CombinedSignature: txCommitmentBody.CombinedSignature,
		Transactions:      txCommitmentBody.Transactions,
	}
}

func (c *StoredCommitment) ToDepositCommitment() *DepositCommitment {
	depositCommitmentBody, ok := c.Body.(*StoredDepositCommitmentBody)
	if !ok {
		panic("invalid DepositCommitment body type")
	}

	return &DepositCommitment{
		CommitmentBase: c.CommitmentBase,
		SubTreeID:      depositCommitmentBody.SubTreeID,
		SubTreeRoot:    depositCommitmentBody.SubTreeRoot,
		Deposits:       depositCommitmentBody.Deposits,
	}
}
