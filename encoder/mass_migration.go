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

const MassMigrationLength = 77

func EncodeMassMigration(massMigration *models.MassMigration) ([]byte, error) {
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
		big.NewInt(int64(massMigration.FromStateID)),
		massMigration.Amount.ToBig(),
		massMigration.Fee.ToBig(),
		big.NewInt(int64(massMigration.SpokeID)),
		massMigration.Nonce.ToBig(),
	)
}

func EncodeMassMigrationForSigning(massMigration *models.MassMigration) []byte {
	b := make([]byte, MassMigrationLength)

	b[0] = uint8(txtype.MassMigration)
	binary.BigEndian.PutUint32(b[1:5], massMigration.FromStateID)
	copy(b[5:37], massMigration.Amount.Bytes())
	copy(b[37:69], massMigration.Fee.Bytes())
	binary.BigEndian.PutUint32(b[69:73], uint32(massMigration.Nonce.Uint64()))
	binary.BigEndian.PutUint32(b[73:77], massMigration.SpokeID)

	return b
}

func HashMassMigration(massMigration *models.MassMigration) (*common.Hash, error) {
	encodedMassMigration, err := EncodeMassMigration(massMigration)
	if err != nil {
		return nil, err
	}
	hash := crypto.Keccak256Hash(encodedMassMigration)
	return &hash, nil
}
