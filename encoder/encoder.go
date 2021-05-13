package encoder

import (
	"encoding/binary"
	"fmt"
	"math/big"

	"github.com/Worldcoin/hubble-commander/contracts/frontend/generic"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/ethereum/go-ethereum/accounts/abi"
)

var (
	tUint256, _       = abi.NewType("uint256", "", nil)
	tUint256Array4, _ = abi.NewType("uint256[4]", "", nil)
)

func EncodeCreate2Transfer(tx *models.Create2Transfer) ([]byte, error) {
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
		big.NewInt(int64(tx.ToPubKeyID)),
		&tx.Amount.Int,
		&tx.Fee.Int,
		&tx.Nonce.Int,
	)
}

func EncodeCreate2TransferWithPubKey(tx *models.Create2Transfer, publicKey *models.PublicKey) ([]byte, error) {
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
		publicKey.BigInts(),
		&tx.Amount.Int,
		&tx.Fee.Int,
		&tx.Nonce.Int,
	)
}

func EncodeCreate2TransferForSigning(tx *models.Create2Transfer, publicKey *models.PublicKey) ([]byte, error) {
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
		publicKey.BigInts(),
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

// Decodes a 256-bit integer from a number with mantissa and a decimal exponent.
// Exponent is 4 bits is packed in front of 12-bit mantissa.
// The original number can be recovered using following formula: V = M * 10^E
func DecodeDecimal(value uint16) models.Uint256 {
	exponent := value >> 12
	mantissa := value & 0x0FFF // mantissa bitmask

	m := big.NewInt(int64(mantissa))
	exp := big.NewInt(0).Exp(big.NewInt(10), big.NewInt(int64(exponent)), nil)

	return models.MakeUint256FromBig(*m.Mul(m, exp))
}

// Encodes a create2Transfer in compact format (without signatures) for the inclusion in the commitment
func EncodeCreate2TransferForCommitment(transfer *models.Create2Transfer) ([]byte, error) {
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
	binary.BigEndian.PutUint32(arr[8:12], transfer.ToPubKeyID)
	binary.BigEndian.PutUint16(arr[12:14], amount)
	binary.BigEndian.PutUint16(arr[14:16], fee)

	return arr, nil
}

func SerializeCreate2Transfers(transfers []models.Create2Transfer) ([]byte, error) {
	buf := make([]byte, 0, len(transfers)*16)

	for i := range transfers {
		encoded, err := EncodeCreate2TransferForCommitment(&transfers[i])
		if err != nil {
			return nil, err
		}
		buf = append(buf, encoded...)
	}

	return buf, nil
}
