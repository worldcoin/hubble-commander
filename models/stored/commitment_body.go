package stored

import (
	"encoding/binary"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

const (
	txCommitmentBodyLength          = 101 // 4 + 64 + 33
	mmCommitmentBodyLength          = 205 // 101 + 72 + 32
	depositCommitmentBodyBaseLength = 64  // 32 + 32
)

type CommitmentBody interface {
	ByteEncoder
	BytesLen() int
}

func NewCommitmentBody(commitmentType batchtype.BatchType) (CommitmentBody, error) {
	switch commitmentType {
	case batchtype.Deposit:
		return new(DepositCommitmentBody), nil
	case batchtype.Transfer, batchtype.Create2Transfer:
		return new(TxCommitmentBody), nil
	case batchtype.MassMigration:
		return new(MMCommitmentBody), nil
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

type MMCommitmentBody struct {
	TxCommitmentBody
	Meta         models.MassMigrationMeta
	WithdrawRoot common.Hash
}

func (c *MMCommitmentBody) Bytes() []byte {
	b := make([]byte, c.BytesLen())
	copy(b[0:101], c.TxCommitmentBody.Bytes())
	copy(b[101:173], c.Meta.Bytes())
	copy(b[173:205], c.WithdrawRoot.Bytes())
	return b
}

func (c *MMCommitmentBody) SetBytes(data []byte) error {
	if len(data) != mmCommitmentBodyLength {
		return models.ErrInvalidLength
	}

	err := c.TxCommitmentBody.SetBytes(data[0:101])
	if err != nil {
		return err
	}
	err = c.Meta.SetBytes(data[101:173])
	if err != nil {
		return err
	}
	c.WithdrawRoot.SetBytes(data[173:205])

	return nil
}

func (c *MMCommitmentBody) BytesLen() int {
	return mmCommitmentBodyLength
}

type DepositCommitmentBody struct {
	SubtreeID   models.Uint256
	SubtreeRoot common.Hash
	Deposits    []models.PendingDeposit
}

func (c *DepositCommitmentBody) Bytes() []byte {
	b := make([]byte, c.BytesLen())
	copy(b[0:32], c.SubtreeID.Bytes())
	copy(b[32:64], c.SubtreeRoot.Bytes())

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

	c.SubtreeID.SetBytes(data[0:32])
	c.SubtreeRoot.SetBytes(data[32:64])
	return nil
}

func (c *DepositCommitmentBody) BytesLen() int {
	return depositCommitmentBodyBaseLength + len(c.Deposits)*models.DepositDataLength
}
