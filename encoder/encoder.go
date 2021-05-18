package encoder

import (
	"encoding/binary"
	"fmt"
	"math/big"

	"github.com/Worldcoin/hubble-commander/contracts/frontend/generic"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/pkg/errors"
)

var (
	tUint256, _       = abi.NewType("uint256", "", nil)
	tUint256Array4, _ = abi.NewType("uint256[4]", "", nil)

	ErrInvalidSlicesLength = errors.New("invalid slices length")
)

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
		&tx.Amount.Int,
		&tx.Fee.Int,
		&tx.Nonce.Int,
	)
}

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
		&tx.Nonce.Int,
		&tx.Amount.Int,
		&tx.Fee.Int,
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

func EncodeUserState(state generic.TypesUserState) ([]byte, error) {
	arguments := abi.Arguments{
		{Name: "pubkeyID", Type: tUint256},
		{Name: "tokenID", Type: tUint256},
		{Name: "balance", Type: tUint256},
		{Name: "nonce", Type: tUint256},
	}
	return arguments.Pack(
		state.PubkeyID,
		state.TokenID,
		state.Balance,
		state.Nonce,
	)
}

// Encodes a 256-bit integer as a number with mantissa and a decimal exponent.
// Exponent is 4 bits is packed in front of 12-bit mantissa.
// The original number can be recovered using following formula: V = M * 10^E
func EncodeDecimal(value models.Uint256) (uint16, error) {
	mantissa := new(big.Int).Set(&value.Int)
	exponent := big.NewInt(0)

	for i := 0; i < 15; i++ {
		if mantissa.Cmp(big.NewInt(0)) != 0 && big.NewInt(0).Mod(mantissa, big.NewInt(10)).Cmp(big.NewInt(0)) == 0 {
			mantissa.Div(mantissa, big.NewInt(10))
			exponent.Add(exponent, big.NewInt(1))
		} else {
			break
		}
	}

	if mantissa.Cmp(big.NewInt(0xfff)) > 0 {
		return 0, fmt.Errorf("value is not encodable as multi-precission decimal")
	}

	return uint16(exponent.Uint64())<<12 + uint16(mantissa.Uint64()), nil
}

// Encodes a transfer in compact format (without signatures) for the inclusion in the commitment
func EncodeTransferForCommitment(transfer *models.Transfer) ([]byte, error) {
	amount, err := EncodeDecimal(transfer.Amount)
	if err != nil {
		return nil, err
	}

	fee, err := EncodeDecimal(transfer.Fee)
	if err != nil {
		return nil, err
	}

	arr := make([]byte, 12)

	binary.BigEndian.PutUint32(arr[0:4], transfer.FromStateID)
	binary.BigEndian.PutUint32(arr[4:8], transfer.ToStateID)
	binary.BigEndian.PutUint16(arr[8:10], amount)
	binary.BigEndian.PutUint16(arr[10:12], fee)

	return arr, nil
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

func SerializeTransfers(transfers []models.Transfer) ([]byte, error) {
	buf := make([]byte, 0, len(transfers)*12)

	for i := range transfers {
		encoded, err := EncodeTransferForCommitment(&transfers[i])
		if err != nil {
			return nil, err
		}
		buf = append(buf, encoded...)
	}

	return buf, nil
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
