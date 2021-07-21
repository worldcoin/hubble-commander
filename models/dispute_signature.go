package models

import "github.com/ethereum/go-ethereum/common"

type SignatureProof struct {
	UserStates []StateMerkleProof
	PublicKeys []PublicKeyProof
}

type PublicKeyProof struct {
	PublicKey *PublicKey
	Witness   Witness
}

type SignatureProofWithReceiver struct {
	UserStates         []StateMerkleProof
	SenderPublicKeys   []PublicKeyProof
	ReceiverPublicKeys []ReceiverPublicKeyProof
}

type ReceiverPublicKeyProof struct {
	PublicKeyHash common.Hash
	Witness       Witness
}
