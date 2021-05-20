package encoder

import (
	"encoding/binary"
	"errors"
	"math/big"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/ethereum/go-ethereum/accounts/abi"
)

var (
	tUint256Array4, _      = abi.NewType("uint256[4]", "", nil)
	ErrInvalidSlicesLength = errors.New("invalid slices length")
	ErrInvalidDataLength   = errors.New("invalid data length")
)

const create2TransferLength = 16

func EncodeCreate2TransferWithStateID(tx *models.Create2Transfer, toPubKeyID uint32) ([]byte, error) {
	arguments := abi.Arguments{
		{Name: "txType", Type: tUint256},
		{Name: "fromIndex", Type: tUint256},
		{Name: "toIndex", Type: tUint256},
		{Name: "toPubkeyID", Type: tUint256},
		{Name: "amount", Type: tUint256},
		{Name: "fee", Type: tUint256},
		{Name: "nonce", Type: tUint256},
	}
	return arguments.Pack(
		big.NewInt(int64(txtype.Create2Transfer)),
		big.NewInt(int64(tx.FromStateID)),
		big.NewInt(int64(tx.ToStateID)),
		big.NewInt(int64(toPubKeyID)),
		&tx.Amount.Int,
		&tx.Fee.Int,
		&tx.Nonce.Int,
	)
}

func EncodeCreate2Transfer(tx *models.Create2Transfer) ([]byte, error) {
	arguments := abi.Arguments{
		{Name: "txType", Type: tUint256},
		{Name: "fromIndex", Type: tUint256},
		{Name: "toPubkey", Type: tUint256Array4},
		{Name: "amount", Type: tUint256},
		{Name: "fee", Type: tUint256},
		{Name: "nonce", Type: tUint256},
	}
	return arguments.Pack(
		big.NewInt(int64(txtype.Create2Transfer)),
		big.NewInt(int64(tx.FromStateID)),
		tx.ToPublicKey.BigInts(),
		&tx.Amount.Int,
		&tx.Fee.Int,
		&tx.Nonce.Int,
	)
}

func EncodeCreate2TransferForSigning(tx *models.Create2Transfer) ([]byte, error) {
	arguments := abi.Arguments{
		{Name: "txType", Type: tUint256},
		{Name: "fromIndex", Type: tUint256},
		{Name: "toPubkey", Type: tUint256Array4},
		{Name: "nonce", Type: tUint256},
		{Name: "amount", Type: tUint256},
		{Name: "fee", Type: tUint256},
	}
	return arguments.Pack(
		big.NewInt(int64(txtype.Create2Transfer)),
		big.NewInt(int64(tx.FromStateID)),
		tx.ToPublicKey.BigInts(),
		&tx.Nonce.Int,
		&tx.Amount.Int,
		&tx.Fee.Int,
	)
}

// Encodes a create2Transfer in compact format (without signatures) for the inclusion in the commitment
func EncodeCreate2TransferForCommitment(transfer *models.Create2Transfer, toPubKeyID uint32) ([]byte, error) {
	amount, err := EncodeDecimal(transfer.Amount)
	if err != nil {
		return nil, err
	}

	fee, err := EncodeDecimal(transfer.Fee)
	if err != nil {
		return nil, err
	}

	arr := make([]byte, create2TransferLength)

	binary.BigEndian.PutUint32(arr[0:4], transfer.FromStateID)
	binary.BigEndian.PutUint32(arr[4:8], transfer.ToStateID)
	binary.BigEndian.PutUint32(arr[8:12], toPubKeyID)
	binary.BigEndian.PutUint16(arr[12:14], amount)
	binary.BigEndian.PutUint16(arr[14:16], fee)

	return arr, nil
}

func DecodeCreate2TransferFromCommitment(data []byte) (*models.Create2Transfer, uint32) {
	fromStateID := binary.BigEndian.Uint32(data[0:4])
	toStateID := binary.BigEndian.Uint32(data[4:8])
	toPubKeyID := binary.BigEndian.Uint32(data[8:12])
	amountEncoded := binary.BigEndian.Uint16(data[12:14])
	feeEncoded := binary.BigEndian.Uint16(data[14:16])

	amount := DecodeDecimal(amountEncoded)
	fee := DecodeDecimal(feeEncoded)

	return &models.Create2Transfer{
		TransactionBase: models.TransactionBase{
			TxType:      txtype.Create2Transfer,
			FromStateID: fromStateID,
			Amount:      amount,
			Fee:         fee,
		},
		ToStateID: toStateID,
	}, toPubKeyID
}

func SerializeCreate2Transfers(transfers []models.Create2Transfer, pubKeyIDs []uint32) ([]byte, error) {
	if len(transfers) != len(pubKeyIDs) {
		return nil, ErrInvalidSlicesLength
	}
	buf := make([]byte, 0, len(transfers)*create2TransferLength)

	for i := range transfers {
		encoded, err := EncodeCreate2TransferForCommitment(&transfers[i], pubKeyIDs[i])
		if err != nil {
			return nil, err
		}
		buf = append(buf, encoded...)
	}

	return buf, nil
}

func DeserializeCreate2Transfers(data []byte) ([]models.Create2Transfer, []uint32, error) {
	dataLength := len(data)
	if dataLength%create2TransferLength != 0 {
		return nil, nil, ErrInvalidDataLength
	}
	transfersCount := dataLength / create2TransferLength

	transfers := make([]models.Create2Transfer, 0, transfersCount)
	pubKeyIDs := make([]uint32, 0, transfersCount)
	for i := 0; i < transfersCount; i++ {
		transfer, pubKeyID := DecodeCreate2TransferFromCommitment(data[i*create2TransferLength : (i+1)*create2TransferLength])
		transfers = append(transfers, *transfer)
		pubKeyIDs = append(pubKeyIDs, pubKeyID)
	}

	return transfers, pubKeyIDs, nil
}
