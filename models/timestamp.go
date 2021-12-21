package models

import (
	"encoding/json"
	"time"
)

type Timestamp struct {
	time.Time
}

func NewTimestamp(t time.Time) *Timestamp {
	return &Timestamp{Time: t}
}

func (t Timestamp) Add(d time.Duration) Timestamp {
	return Timestamp{t.Time.Add(d)}
}

func (t Timestamp) Before(other Timestamp) bool {
	return t.Time.Before(other.Time)
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

func (t *Timestamp) Bytes() []byte {
	utcTimestamp := t.Time.UTC()
	bytes, _ := utcTimestamp.MarshalBinary()
	return bytes
}

func (t *Timestamp) SetBytes(data []byte) error {
	return t.UnmarshalBinary(data)
}
