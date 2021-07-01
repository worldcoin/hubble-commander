package bls

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/kilic/bn254/bls"
	log "github.com/sirupsen/logrus"
)

func fromBLSPublicKey(blsPK *bls.PublicKey) *models.PublicKey {
	var pk models.PublicKey
	pkBytes := blsPK.ToBytes()
	copy(pk[:32], pkBytes[32:64])
	copy(pk[32:64], pkBytes[:32])
	copy(pk[64:96], pkBytes[96:])
	copy(pk[96:], pkBytes[64:96])
	return &pk
}

func toBLSPublicKey(pk *models.PublicKey) *bls.PublicKey {
	blsBytes := make([]byte, 128)
	copy(blsBytes[:32], pk[32:64])
	copy(blsBytes[32:64], pk[:32])
	copy(blsBytes[64:96], pk[96:])
	copy(blsBytes[96:], pk[64:96])
	blsPK, err := bls.PublicKeyFromBytes(blsBytes)
	if err != nil {
		log.Panicf("failed to convert public key byte array to bls type: %v", err)
	}
	return blsPK
}

func PrivateToPublicKey(privateKey [32]byte) (*models.PublicKey, error) {
	keyPair, err := bls.NewKeyPairFromSecret(privateKey[:])
	if err != nil {
		return nil, err
	}
	return fromBLSPublicKey(keyPair.Public), nil
}
