package stored

import (
	"encoding/binary"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

const (
	txCommitmentBodyLength          = 4 + 64 + 33
	depositCommitmentBodyBaseLength = 32 + 32
)

type CommitmentBody interface {
	ByteEncoder
	BytesLen() int
}

func NewCommitmentBody(commitmentType batchtype.BatchType) (CommitmentBody, error) {
	// nolint:exhaustive
	switch commitmentType {
	case batchtype.Deposit:
		return new(DepositCommitmentBody), nil
	case batchtype.Transfer, batchtype.Create2Transfer:
		return new(TxCommitmentBody), nil
	default:
		return nil, errors.Errorf("unsupported commitment type: %s", commitmentType)
	}
}

type TxCommitmentBody struct {
	FeeReceiver       uint32
	CombinedSignature models.Signature
	BodyHash          *common.Hash
}

func (c *TxCommitmentBody) Bytes() []byte {
	b := make([]byte, c.BytesLen())
	binary.BigEndian.PutUint32(b[0:4], c.FeeReceiver)
	copy(b[4:68], c.CombinedSignature.Bytes())
	copy(b[68:101], EncodeHashPointer(c.BodyHash))
	return b
}

func (c *TxCommitmentBody) SetBytes(data []byte) error {
	if len(data) != txCommitmentBodyLength {
		return models.ErrInvalidLength
	}
	err := c.CombinedSignature.SetBytes(data[4:68])
	if err != nil {
		return err
	}

	c.FeeReceiver = binary.BigEndian.Uint32(data[0:4])
	c.BodyHash = decodeHashPointer(data[68:101])
	return nil
}

func (c *TxCommitmentBody) BytesLen() int {
	return txCommitmentBodyLength
}

type DepositCommitmentBody struct {
	SubTreeID   models.Uint256
	SubTreeRoot common.Hash
	Deposits    []models.PendingDeposit
}

func (c *DepositCommitmentBody) Bytes() []byte {
	b := make([]byte, c.BytesLen())
	copy(b[0:32], c.SubTreeID.Bytes())
	copy(b[32:64], c.SubTreeRoot.Bytes())

	startIndex := depositCommitmentBodyBaseLength
	for i := range c.Deposits {
		startIndex += copy(b[startIndex:startIndex+models.DepositDataLength], c.Deposits[i].Bytes())
	}

	return b
}

func (c *DepositCommitmentBody) SetBytes(data []byte) error {
	overallDepositsLength := len(data) - depositCommitmentBodyBaseLength
	if len(data) <= depositCommitmentBodyBaseLength || overallDepositsLength%models.DepositDataLength != 0 {
		return models.ErrInvalidLength
	}

	depositCount := overallDepositsLength / models.DepositDataLength
	c.Deposits = make([]models.PendingDeposit, 0, depositCount)

	startIndex := depositCommitmentBodyBaseLength
	for i := 0; i < depositCount; i++ {
		endIndex := startIndex + models.DepositDataLength
		deposit := models.PendingDeposit{}
		err := deposit.SetBytes(data[startIndex:endIndex])
		if err != nil {
			return err
		}
		c.Deposits = append(c.Deposits, deposit)
		startIndex = endIndex
	}

	c.SubTreeID.SetBytes(data[0:32])
	c.SubTreeRoot.SetBytes(data[32:64])
	return nil
}

func (c *DepositCommitmentBody) BytesLen() int {
	return depositCommitmentBodyBaseLength + len(c.Deposits)*models.DepositDataLength
}
