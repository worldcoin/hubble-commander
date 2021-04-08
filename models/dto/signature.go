package dto

import (
	"encoding/hex"
	"encoding/json"
)

type Signature []byte

func (s *Signature) UnmarshalJSON(bytes []byte) error {
	var temp string
	err := json.Unmarshal(bytes, &temp)
	if err != nil {
		return err
	}
	if len(temp) == 0 {
		*s = Signature{}
		return nil
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
