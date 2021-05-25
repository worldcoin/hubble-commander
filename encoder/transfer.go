package encoder

import (
	"encoding/binary"
	"math/big"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/ethereum/go-ethereum/accounts/abi"
)

const transferLength = 12

func EncodeTransfer(tx *models.Transfer) ([]byte, error) {
	arguments := abi.Arguments{
		{Name: "txType", Type: tUint256},
		{Name: "fromIndex", Type: tUint256},
		{Name: "toIndex", Type: tUint256},
		{Name: "amount", Type: tUint256},
		{Name: "fee", Type: tUint256},
		{Name: "nonce", Type: tUint256},
	}
	return arguments.Pack(
		big.NewInt(int64(txtype.Transfer)),
		big.NewInt(int64(tx.FromStateID)),
		big.NewInt(int64(tx.ToStateID)),
		tx.Amount.ToBig(),
		tx.Fee.ToBig(),
		tx.Nonce.ToBig(),
	)
}

func EncodeTransferForSigning(tx *models.Transfer) ([]byte, error) {
	arguments := abi.Arguments{
		{Name: "txType", Type: tUint256},
		{Name: "fromIndex", Type: tUint256},
		{Name: "toIndex", Type: tUint256},
		{Name: "nonce", Type: tUint256},
		{Name: "amount", Type: tUint256},
		{Name: "fee", Type: tUint256},
	}
	return arguments.Pack(
		big.NewInt(int64(txtype.Transfer)),
		big.NewInt(int64(tx.FromStateID)),
		big.NewInt(int64(tx.ToStateID)),
		tx.Nonce.ToBig(),
		tx.Amount.ToBig(),
		tx.Fee.ToBig(),
	)
}

// EncodeTransferForCommitment Encodes a transfer in compact format (without signatures) for the inclusion in the commitment
func EncodeTransferForCommitment(transfer *models.Transfer) ([]byte, error) {
	amount, err := EncodeDecimal(transfer.Amount)
	if err != nil {
		return nil, err
	}

	fee, err := EncodeDecimal(transfer.Fee)
	if err != nil {
		return nil, err
	}

	arr := make([]byte, transferLength)

	binary.BigEndian.PutUint32(arr[0:4], transfer.FromStateID)
	binary.BigEndian.PutUint32(arr[4:8], transfer.ToStateID)
	binary.BigEndian.PutUint16(arr[8:10], amount)
	binary.BigEndian.PutUint16(arr[10:12], fee)

	return arr, nil
}

func DecodeTransferFromCommitment(data []byte) *models.Transfer {
	fromStateID := binary.BigEndian.Uint32(data[0:4])
	toStateID := binary.BigEndian.Uint32(data[4:8])
	amountEncoded := binary.BigEndian.Uint16(data[8:10])
	feeEncoded := binary.BigEndian.Uint16(data[10:12])

	amount := DecodeDecimal(amountEncoded)
	fee := DecodeDecimal(feeEncoded)

	return &models.Transfer{
		TransactionBase: models.TransactionBase{
			TxType:      txtype.Transfer,
			FromStateID: fromStateID,
			Amount:      amount,
			Fee:         fee,
		},
		ToStateID: toStateID,
	}
}

func SerializeTransfers(transfers []models.Transfer) ([]byte, error) {
	buf := make([]byte, 0, len(transfers)*transferLength)

	for i := range transfers {
		encoded, err := EncodeTransferForCommitment(&transfers[i])
		if err != nil {
			return nil, err
		}
		buf = append(buf, encoded...)
	}

	return buf, nil
}

func DeserializeTransfers(data []byte) ([]models.Transfer, error) {
	dataLength := len(data)
	if dataLength%transferLength != 0 {
		return nil, ErrInvalidDataLength
	}
	transfersCount := dataLength / transferLength

	res := make([]models.Transfer, 0, transfersCount)
	for i := 0; i < transfersCount; i++ {
		transfer := DecodeTransferFromCommitment(data[i*transferLength : (i+1)*transferLength])
		res = append(res, *transfer)
	}

	return res, nil
}
