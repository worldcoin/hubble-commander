package models

type SignatureProof struct {
	States     []StateMerkleProof
	PublicKeys []PublicKeyProof
}

type PublicKeyProof struct {
	PublicKey *PublicKey
	Witness   Witness
}
