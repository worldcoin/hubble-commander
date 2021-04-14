package bls

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/kilic/bn254/bls"
)

func toBLSPublicKey(pk *models.PublicKey) *bls.PublicKey {
	blsPK, err := bls.PublicKeyFromBytes(pk[:])
	if err != nil {
		panic("failed to convert public key byte array to bls type")
	}
	return blsPK
}
