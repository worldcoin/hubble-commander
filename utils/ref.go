package utils

import "time"

func Bool(b bool) *bool {
	return &b
}

func Int(i int) *int {
	return &i
}

func Int64(i int64) *int64 {
	return &i
}

func Uint64(u uint64) *uint64 {
	return &u
}

func String(s string) *string {
	return &s
}

func Duration(d time.Duration) *time.Duration {
	return &d
}
