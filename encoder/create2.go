package encoder

import (
	"encoding/binary"
	"errors"
	"math/big"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	tBytes32, _            = abi.NewType("bytes32", "", nil)
	ErrInvalidSlicesLength = errors.New("invalid slices length")
)

const Create2TransferLength = 16

func EncodeCreate2TransferWithStateID(create2Transfer *models.Create2Transfer, toPubKeyID uint32) ([]byte, error) {
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
		big.NewInt(int64(create2Transfer.FromStateID)),
		big.NewInt(int64(*create2Transfer.ToStateID)),
		big.NewInt(int64(toPubKeyID)),
		create2Transfer.Amount.ToBig(),
		create2Transfer.Fee.ToBig(),
		create2Transfer.Nonce.ToBig(),
	)
}

func EncodeCreate2Transfer(create2Transfer *models.Create2Transfer) ([]byte, error) {
	arguments := abi.Arguments{
		{Name: "txType", Type: tUint256},
		{Name: "fromIndex", Type: tUint256},
		{Name: "toPubkey", Type: tBytes32},
		{Name: "amount", Type: tUint256},
		{Name: "fee", Type: tUint256},
		{Name: "nonce", Type: tUint256},
	}
	return arguments.Pack(
		big.NewInt(int64(txtype.Create2Transfer)),
		big.NewInt(int64(create2Transfer.FromStateID)),
		crypto.Keccak256Hash(create2Transfer.ToPublicKey.Bytes()),
		create2Transfer.Amount.ToBig(),
		create2Transfer.Fee.ToBig(),
		create2Transfer.Nonce.ToBig(),
	)
}

func EncodeCreate2TransferForSigning(create2Transfer *models.Create2Transfer) ([]byte, error) {
	arguments := abi.Arguments{
		{Name: "txType", Type: tUint256},
		{Name: "fromIndex", Type: tUint256},
		{Name: "toPubkey", Type: tBytes32},
		{Name: "nonce", Type: tUint256},
		{Name: "amount", Type: tUint256},
		{Name: "fee", Type: tUint256},
	}
	return arguments.Pack(
		big.NewInt(int64(txtype.Create2Transfer)),
		big.NewInt(int64(create2Transfer.FromStateID)),
		crypto.Keccak256Hash(create2Transfer.ToPublicKey.Bytes()),
		create2Transfer.Nonce.ToBig(),
		create2Transfer.Amount.ToBig(),
		create2Transfer.Fee.ToBig(),
	)
}

// Encodes a create2Transfer in compact format (without signatures) for the inclusion in the commitment
func EncodeCreate2TransferForCommitment(create2Transfer *models.Create2Transfer, toPubKeyID uint32) ([]byte, error) {
	amount, err := EncodeDecimal(create2Transfer.Amount)
	if err != nil {
		return nil, err
	}

	fee, err := EncodeDecimal(create2Transfer.Fee)
	if err != nil {
		return nil, err
	}

	arr := make([]byte, Create2TransferLength)

	binary.BigEndian.PutUint32(arr[0:4], create2Transfer.FromStateID)
	binary.BigEndian.PutUint32(arr[4:8], *create2Transfer.ToStateID)
	binary.BigEndian.PutUint32(arr[8:12], toPubKeyID)
	binary.BigEndian.PutUint16(arr[12:14], amount)
	binary.BigEndian.PutUint16(arr[14:16], fee)

	return arr, nil
}

func DecodeCreate2TransferFromCommitment(data []byte) (create2Transfer *models.Create2Transfer, toPubKeyID uint32, err error) {
	fromStateID := binary.BigEndian.Uint32(data[0:4])
	toStateID := binary.BigEndian.Uint32(data[4:8])
	toPubKeyID = binary.BigEndian.Uint32(data[8:12])
	amountEncoded := binary.BigEndian.Uint16(data[12:14])
	feeEncoded := binary.BigEndian.Uint16(data[14:16])

	amount := DecodeDecimal(amountEncoded)
	fee := DecodeDecimal(feeEncoded)

	create2Transfer = &models.Create2Transfer{
		TransactionBase: models.TransactionBase{
			TxType:      txtype.Create2Transfer,
			FromStateID: fromStateID,
			Amount:      amount,
			Fee:         fee,
		},
		ToStateID: &toStateID,
	}
	return create2Transfer, toPubKeyID, nil
}

func SerializeCreate2Transfers(create2Transfers []models.Create2Transfer, pubKeyIDs []uint32) ([]byte, error) {
	if len(create2Transfers) != len(pubKeyIDs) {
		return nil, ErrInvalidSlicesLength
	}
	buf := make([]byte, 0, len(create2Transfers)*Create2TransferLength)

	for i := range create2Transfers {
		encoded, err := EncodeCreate2TransferForCommitment(&create2Transfers[i], pubKeyIDs[i])
		if err != nil {
			return nil, err
		}
		buf = append(buf, encoded...)
	}

	return buf, nil
}

func DeserializeCreate2Transfers(data []byte) ([]models.Create2Transfer, []uint32, error) {
	transfersCount := len(data) / Create2TransferLength

	transfers := make([]models.Create2Transfer, 0, transfersCount)
	pubKeyIDs := make([]uint32, 0, transfersCount)
	for i := 0; i < transfersCount; i++ {
		transfer, pubKeyID, err := DecodeCreate2TransferFromCommitment(data[i*Create2TransferLength : (i+1)*Create2TransferLength])
		if err != nil {
			return nil, nil, err
		}
		transfers = append(transfers, *transfer)
		pubKeyIDs = append(pubKeyIDs, pubKeyID)
	}

	return transfers, pubKeyIDs, nil
}

func DeserializeCreate2TransferPubKeyIDs(data []byte) []uint32 {
	transfersCount := len(data) / Create2TransferLength

	pubKeyIDs := make([]uint32, 0, transfersCount)
	for i := 0; i < transfersCount; i++ {
		pubKeyID := binary.BigEndian.Uint32(data[i*Create2TransferLength+8 : i*Create2TransferLength+12])
		pubKeyIDs = append(pubKeyIDs, pubKeyID)
	}
	return pubKeyIDs
}

func HashCreate2Transfer(create2Transfer *models.Create2Transfer) (*common.Hash, error) {
	encodedTransfer, err := EncodeCreate2Transfer(create2Transfer)
	if err != nil {
		return nil, err
	}
	hash := crypto.Keccak256Hash(encodedTransfer)
	return &hash, nil
}
