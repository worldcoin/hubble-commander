package models

type IncomingTransaction struct {
	FromIndex *Uint256
	ToIndex   *Uint256
	Amount    *Uint256
	Fee       *Uint256
	Nonce     *Uint256
	// TODO: Right now decoder expects a base64 string here, we could define a custom type with interface implementation to expect a hex string
	Signature []byte
}
