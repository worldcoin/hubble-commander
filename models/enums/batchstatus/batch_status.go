package batchstatus

import (
	"encoding/json"

	enumerr "github.com/Worldcoin/hubble-commander/models/enums/errors"
)

type BatchStatus uint

const (
	Pending BatchStatus = iota + 1000
	InBatch
	Finalised
)

var BatchStatuses = map[BatchStatus]string{
	Pending:   "PENDING",
	InBatch:   "IN_BATCH",
	Finalised: "FINALISED",
}

func (s BatchStatus) Ref() *BatchStatus {
	return &s
}

func (s BatchStatus) String() string {
	msg, exists := BatchStatuses[s]
	if !exists {
		return "UNKNOWN"
	}
	return msg
}

func (s *BatchStatus) UnmarshalJSON(bytes []byte) error {
	var strType string
	err := json.Unmarshal(bytes, &strType)
	if err != nil {
		return err
	}

	for k, v := range BatchStatuses {
		if v == strType {
			*s = k
			return nil
		}
	}
	return enumerr.NewUnsupportedError("batch status")
}

func (s BatchStatus) MarshalJSON() ([]byte, error) {
	msg, exists := BatchStatuses[s]
	if !exists {
		return nil, enumerr.NewUnsupportedError("batch status")
	}
	return json.Marshal(msg)
}
