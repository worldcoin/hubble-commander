package dto

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTimestamp_MarshalJSON(t *testing.T) {
	timestamp := Timestamp{Time: time.Now()}
	json, err := timestamp.MarshalJSON()
	require.NoError(t, err)

	seconds := timestamp.Unix()
	stringValue := fmt.Sprintf("%d", seconds)
	expected := []byte(stringValue)

	require.Equal(t, expected, json)
}

func TestTimestamp_UnmarshalJSON(t *testing.T) {
	seconds := time.Now().Unix()
	stringValue := fmt.Sprintf("%d", seconds)
	json := []byte(stringValue)

	timestamp := Timestamp{}
	err := timestamp.UnmarshalJSON(json)
	require.NoError(t, err)

	require.Equal(t, seconds, timestamp.Unix())
}
