package badger

import bh "github.com/timshannon/badgerhold/v4"

func MatchAll(_ *bh.RecordAccess) (bool, error) {
	return true, nil
}
