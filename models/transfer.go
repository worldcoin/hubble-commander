package models

type Transfer struct {
	FromStateID uint32
	ToStateID   uint32
	Amount      Uint256
	Fee         Uint256
	Nonce       Uint256
	Signature   []byte
}
