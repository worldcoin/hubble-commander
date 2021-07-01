package dto

import (
	"encoding/json"
	"time"
)

type Timestamp struct {
	time.Time
}

func (t *Timestamp) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Unix())
}

func (t *Timestamp) UnmarshalJSON(data []byte) error {
	var seconds int64
	err := json.Unmarshal(data, &seconds)
	if err != nil {
		return err
	}
	*t = Timestamp{Time: time.Unix(seconds, 0)}
	return nil
}
