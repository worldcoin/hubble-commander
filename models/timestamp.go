package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type Timestamp struct {
	time.Time
}

func NewTimestamp(t time.Time) *Timestamp {
	return &Timestamp{Time: t}
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

// Scan implements Scanner for database/sql.
func (t *Timestamp) Scan(src interface{}) error {
	value, ok := src.(time.Time)
	if !ok {
		return fmt.Errorf("can't scan %T into Timestamp", src)
	}
	t.Time = value.UTC()
	return nil
}

// Value implements valuer for database/sql.
func (t Timestamp) Value() (driver.Value, error) {
	return t.Time.UTC(), nil
}
