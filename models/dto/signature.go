package dto

import (
	"encoding/hex"
	"encoding/json"
	"errors"
)

type Signature []byte

func (s *Signature) UnmarshalJSON(bytes []byte) error {
	var temp string
	err := json.Unmarshal(bytes, &temp)
	if err != nil {
		return err
	}
	if temp == "" {
		*s = Signature{}
		return nil
	}

	if len(temp) < 2 || temp[:2] != "0x" {
		return errors.New("hex string must be 0x prepended")
	}

	str, err := hex.DecodeString(temp[2:])
	if err != nil {
		return err
	}
	*s = str
	return nil
}

func (s Signature) MarshalJSON() ([]byte, error) {
	if len(s) == 0 {
		return json.Marshal("")
	}
	return json.Marshal("0x" + hex.EncodeToString(s))
}
