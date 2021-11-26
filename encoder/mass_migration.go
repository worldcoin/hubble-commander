package encoder

import (
	"encoding/binary"
	"math/big"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/ethereum/go-ethereum/accounts/abi"
)

const MassMigrationLength = 77

func EncodeMassMigration(tx *models.MassMigration) ([]byte, error) {
	arguments := abi.Arguments{
		{Name: "txType", Type: tUint256},
		{Name: "fromIndex", Type: tUint256},
		{Name: "amount", Type: tUint256},
		{Name: "fee", Type: tUint256},
		{Name: "spokeID", Type: tUint256},
		{Name: "nonce", Type: tUint256},
	}
	return arguments.Pack(
		big.NewInt(int64(txtype.MassMigration)),
		big.NewInt(int64(tx.FromStateID)),
		tx.Amount.ToBig(),
		tx.Fee.ToBig(),
		tx.SpokeID.ToBig(),
		tx.Nonce.ToBig(),
	)
}

func EncodeMassMigrationForSigning(tx *models.MassMigration) ([]byte, error) {
	b := make([]byte, MassMigrationLength)

	b[0] = uint8(txtype.MassMigration)
	binary.BigEndian.PutUint32(b[1:5], tx.FromStateID)
	copy(b[5:37], tx.Amount.Bytes())
	copy(b[37:69], tx.Fee.Bytes())
	binary.BigEndian.PutUint32(b[69:73], uint32(tx.Nonce.Uint64()))
	binary.BigEndian.PutUint32(b[73:77], uint32(tx.SpokeID.Uint64()))

	return b, nil
}
