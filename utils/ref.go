package utils

import (
	"time"

	"github.com/Worldcoin/hubble-commander/models"
)

// TODO move this to ref package for better code readability

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

func Uint256(u models.Uint256) *models.Uint256 {
	return &u
}
