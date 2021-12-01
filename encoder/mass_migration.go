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

const (
	MassMigrationForSignatureLength  = 77
	MassMigrationForCommitmentLength = 8
)

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
	b := make([]byte, MassMigrationForSignatureLength)

	b[0] = uint8(txtype.MassMigration)
	binary.BigEndian.PutUint32(b[1:5], massMigration.FromStateID)
	copy(b[5:37], massMigration.Amount.Bytes())
	copy(b[37:69], massMigration.Fee.Bytes())
	binary.BigEndian.PutUint32(b[69:73], uint32(massMigration.Nonce.Uint64()))
	binary.BigEndian.PutUint32(b[73:77], massMigration.SpokeID)

	return b
}

func EncodeMassMigrationForCommitment(massMigration *models.MassMigration) ([]byte, error) {
	amount, err := EncodeDecimal(massMigration.Amount)
	if err != nil {
		return nil, err
	}

	fee, err := EncodeDecimal(massMigration.Fee)
	if err != nil {
		return nil, err
	}

	b := make([]byte, MassMigrationForCommitmentLength)

	binary.BigEndian.PutUint32(b[0:4], massMigration.FromStateID)
	binary.BigEndian.PutUint16(b[4:6], amount)
	binary.BigEndian.PutUint16(b[6:8], fee)

	return b, nil
}

func DecodeMassMigrationFromCommitment(data []byte) (*models.MassMigration, error) {
	fromStateID := binary.BigEndian.Uint32(data[0:4])
	amountEncoded := binary.BigEndian.Uint16(data[4:6])
	feeEncoded := binary.BigEndian.Uint16(data[6:8])

	amount := DecodeDecimal(amountEncoded)
	fee := DecodeDecimal(feeEncoded)

	massMigration := &models.MassMigration{
		TransactionBase: models.TransactionBase{
			TxType:      txtype.MassMigration,
			FromStateID: fromStateID,
			Amount:      amount,
			Fee:         fee,
		},
	}
	return massMigration, nil
}

func SerializeMassMigrations(massMigrations []models.MassMigration) ([]byte, error) {
	buf := make([]byte, 0, len(massMigrations)*MassMigrationForCommitmentLength)

	for i := range massMigrations {
		encoded, err := EncodeMassMigrationForCommitment(&massMigrations[i])
		if err != nil {
			return nil, err
		}
		buf = append(buf, encoded...)
	}

	return buf, nil
}

func DeserializeMassMigrations(data []byte) ([]models.MassMigration, error) {
	massMigrationsCount := len(data) / MassMigrationForCommitmentLength

	res := make([]models.MassMigration, 0, massMigrationsCount)
	for i := 0; i < massMigrationsCount; i++ {
		massMigration, err := DecodeMassMigrationFromCommitment(data[i*MassMigrationForCommitmentLength : (i+1)*MassMigrationForCommitmentLength])
		if err != nil {
			return nil, err
		}
		res = append(res, *massMigration)
	}

	return res, nil
}

func HashMassMigration(massMigration *models.MassMigration) (*common.Hash, error) {
	encodedMassMigration, err := EncodeMassMigration(massMigration)
	if err != nil {
		return nil, err
	}
	hash := crypto.Keccak256Hash(encodedMassMigration)
	return &hash, nil
}
