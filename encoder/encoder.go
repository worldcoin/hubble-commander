package encoder

import (
	"encoding/binary"
	"fmt"
	"math/big"

	"github.com/Worldcoin/hubble-commander/contracts/frontend/generic"
	"github.com/Worldcoin/hubble-commander/contracts/frontend/transfer"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/accounts/abi"
)

var (
	tUint256, _ = abi.NewType("uint256", "", nil)
)

func EncodeTransfer(tx transfer.OffchainTransfer) ([]uint8, error) {
	arguments := abi.Arguments{
		{Name: "txType", Type: tUint256},
		{Name: "fromIndex", Type: tUint256},
		{Name: "toIndex", Type: tUint256},
		{Name: "amount", Type: tUint256},
		{Name: "fee", Type: tUint256},
		{Name: "nonce", Type: tUint256},
	}
	encodedBytes, err := arguments.Pack(
		tx.TxType,
		tx.FromIndex,
		tx.ToIndex,
		tx.Amount,
		tx.Fee,
		tx.Nonce,
	)
	if err != nil {
		return nil, err
	}
	return encodedBytes, nil
}

func EncodeUserState(state generic.TypesUserState) ([]uint8, error) {
	arguments := abi.Arguments{
		{Name: "pubkeyID", Type: tUint256},
		{Name: "tokenID", Type: tUint256},
		{Name: "balance", Type: tUint256},
		{Name: "nonce", Type: tUint256},
	}
	encodedBytes, err := arguments.Pack(
		state.PubkeyID,
		state.TokenID,
		state.Balance,
		state.Nonce,
	)
	if err != nil {
		return nil, err
	}
	return encodedBytes, nil
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

// Encodes a transaction in compact format (without signatures) for the inclusion in the commitment
func EncodeTransaction(transaction *models.Transaction) ([]uint8, error) {
	amount, err := EncodeDecimal(transaction.Amount)
	if err != nil {
		return nil, err
	}

	fee, err := EncodeDecimal(transaction.Fee)
	if err != nil {
		return nil, err
	}

	arr := make([]byte, 12)

	binary.BigEndian.PutUint32(arr[0:4], transaction.FromIndex)
	binary.BigEndian.PutUint32(arr[4:8], transaction.ToIndex)
	binary.BigEndian.PutUint16(arr[8:10], amount)
	binary.BigEndian.PutUint16(arr[10:12], fee)

	return arr, nil
}

func SerializeTransactions(txs []models.Transaction) ([]byte, error) {
	buf := make([]byte, 0, len(txs)*12)

	for i := range txs {
		encoded, err := EncodeTransaction(&txs[i])
		if err != nil {
			return nil, err
		}
		buf = append(buf, encoded...)
	}

	return buf, nil
}
