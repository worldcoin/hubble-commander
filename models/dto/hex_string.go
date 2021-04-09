package dto

import (
	"encoding/hex"
	"encoding/json"
	"errors"
)

var ErrHexStringNotPrepended = errors.New("hex string must be 0x prepended")

type HexString []byte

func (s *HexString) UnmarshalJSON(bytes []byte) error {
	var temp string
	err := json.Unmarshal(bytes, &temp)
	if err != nil {
		return err
	}
	if temp == "" {
		*s = HexString{}
		return nil
	}

	if len(temp) < 2 || temp[:2] != "0x" {
		return ErrHexStringNotPrepended
	}

	str, err := hex.DecodeString(temp[2:])
	if err != nil {
		return err
	}
	*s = str
	return nil
}

func (s HexString) MarshalJSON() ([]byte, error) {
	if len(s) == 0 {
		return json.Marshal("")
	}
	return json.Marshal("0x" + hex.EncodeToString(s))
}
