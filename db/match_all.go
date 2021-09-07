package db

import bh "github.com/timshannon/badgerhold/v3"

func MatchAll(_ *bh.RecordAccess) (bool, error) {
	return true, nil
}
