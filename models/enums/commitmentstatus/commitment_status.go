package commitmentstatus

import (
	"encoding/json"

	bs "github.com/Worldcoin/hubble-commander/models/enums/batchstatus"
	enumerr "github.com/Worldcoin/hubble-commander/models/enums/errors"
)

type CommitmentStatus uint

const (
	Pending   = CommitmentStatus(bs.Pending) // Not in use + Will be replaced in the future in favor or `Submitted`
	Mined     = CommitmentStatus(bs.Mined)
	Finalised = CommitmentStatus(bs.Finalised) // nolint:misspell

)

var CommitmentStatuses = map[CommitmentStatus]string{
	Pending:   bs.BatchStatuses[bs.Pending],
	Mined:     bs.BatchStatuses[bs.Mined],
	Finalised: bs.BatchStatuses[bs.Finalised], // nolint:misspell
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
