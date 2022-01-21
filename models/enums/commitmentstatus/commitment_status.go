package commitmentstatus

import (
	"encoding/json"

	"github.com/Worldcoin/hubble-commander/models/enums/batchstatus"
	enumerr "github.com/Worldcoin/hubble-commander/models/enums/errors"
)

type CommitmentStatus uint

const (
	Pending   = CommitmentStatus(batchstatus.Pending)
	InBatch   = CommitmentStatus(batchstatus.Submitted)
	Finalised = CommitmentStatus(batchstatus.Finalised) // nolint:misspell

)

var CommitmentStatuses = map[CommitmentStatus]string{
	Pending:   "PENDING",
	InBatch:   "IN_BATCH",
	Finalised: "FINALISED", // nolint:misspell
}

func (s CommitmentStatus) Ref() *CommitmentStatus {
	return &s
}

func (s CommitmentStatus) String() string {
	msg, exists := CommitmentStatuses[s]
	if !exists {
		return "UNKNOWN"
	}
	return msg
}

func (s *CommitmentStatus) UnmarshalJSON(bytes []byte) error {
	var strType string
	err := json.Unmarshal(bytes, &strType)
	if err != nil {
		return err
	}

	for k, v := range CommitmentStatuses {
		if v == strType {
			*s = k
			return nil
		}
	}
	return enumerr.NewUnsupportedError("commitment status")
}

func (s CommitmentStatus) MarshalJSON() ([]byte, error) {
	msg, exists := CommitmentStatuses[s]
	if !exists {
		return nil, enumerr.NewUnsupportedError("commitment status")
	}
	return json.Marshal(msg)
}
