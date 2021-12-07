package encoder

import (
	"encoding/binary"
	"math/big"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

const TransferLength = 12

func EncodeTransfer(transfer *models.Transfer) ([]byte, error) {
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
		big.NewInt(int64(transfer.FromStateID)),
		big.NewInt(int64(transfer.ToStateID)),
		transfer.Amount.ToBig(),
		transfer.Fee.ToBig(),
		transfer.Nonce.ToBig(),
	)
}

func EncodeTransferForSigning(transfer *models.Transfer) ([]byte, error) {
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
		big.NewInt(int64(transfer.FromStateID)),
		big.NewInt(int64(transfer.ToStateID)),
		transfer.Nonce.ToBig(),
		transfer.Amount.ToBig(),
		transfer.Fee.ToBig(),
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

	arr := make([]byte, TransferLength)

	binary.BigEndian.PutUint32(arr[0:4], transfer.FromStateID)
	binary.BigEndian.PutUint32(arr[4:8], transfer.ToStateID)
	binary.BigEndian.PutUint16(arr[8:10], amount)
	binary.BigEndian.PutUint16(arr[10:12], fee)

	return arr, nil
}

func DecodeTransferFromCommitment(data []byte) (*models.Transfer, error) {
	fromStateID := binary.BigEndian.Uint32(data[0:4])
	toStateID := binary.BigEndian.Uint32(data[4:8])
	amountEncoded := binary.BigEndian.Uint16(data[8:10])
	feeEncoded := binary.BigEndian.Uint16(data[10:12])

	amount := DecodeDecimal(amountEncoded)
	fee := DecodeDecimal(feeEncoded)

	transfer := &models.Transfer{
		TransactionBase: models.TransactionBase{
			TxType:      txtype.Transfer,
			FromStateID: fromStateID,
			Amount:      amount,
			Fee:         fee,
		},
		ToStateID: toStateID,
	}
	return transfer, nil
}

func SerializeTransfers(transfers []models.Transfer) ([]byte, error) {
	buf := make([]byte, 0, len(transfers)*TransferLength)

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
	transfersCount := len(data) / TransferLength

	res := make([]models.Transfer, 0, transfersCount)
	for i := 0; i < transfersCount; i++ {
		transfer, err := DecodeTransferFromCommitment(data[i*TransferLength : (i+1)*TransferLength])
		if err != nil {
			return nil, err
		}
		res = append(res, *transfer)
	}

	return res, nil
}

func HashTransfer(transfer *models.Transfer) (*common.Hash, error) {
	encodedTransfer, err := EncodeTransfer(transfer)
	if err != nil {
		return nil, err
	}
	hash := crypto.Keccak256Hash(encodedTransfer)
	return &hash, nil
}
