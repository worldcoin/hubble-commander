package models

type GenesisAccount struct {
	PrivateKey []byte
	Balance    Uint256
}

type RegisteredGenesisAccount struct {
	GenesisAccount
	PublicKey PublicKey
	PubKeyID  uint32
}
