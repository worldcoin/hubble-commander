package models

type RawGenesisAccount struct {
	PrivateKey string  `yaml:"privateKey"`
	Balance    Uint256 `yaml:"balance"`
}

type GenesisAccount struct {
	PrivateKey [32]byte
	Balance    Uint256
}

type RegisteredGenesisAccount struct {
	GenesisAccount
	PublicKey PublicKey
	PubKeyID  uint32
}
