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
)

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

	arr := make([]byte, 16)

	binary.BigEndian.PutUint32(arr[0:4], transfer.FromStateID)
	binary.BigEndian.PutUint32(arr[4:8], transfer.ToStateID)
	binary.BigEndian.PutUint32(arr[8:12], toPubKeyID)
	binary.BigEndian.PutUint16(arr[12:14], amount)
	binary.BigEndian.PutUint16(arr[14:16], fee)

	return arr, nil
}

func SerializeCreate2Transfers(transfers []models.Create2Transfer, pubKeyIDs []uint32) ([]byte, error) {
	if len(transfers) != len(pubKeyIDs) {
		return nil, ErrInvalidSlicesLength
	}
	buf := make([]byte, 0, len(transfers)*16)

	for i := range transfers {
		encoded, err := EncodeCreate2TransferForCommitment(&transfers[i], pubKeyIDs[i])
		if err != nil {
			return nil, err
		}
		buf = append(buf, encoded...)
	}

	return buf, nil
}
