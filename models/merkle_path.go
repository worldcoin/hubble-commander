package models

import (
	"database/sql/driver"
	"fmt"
	"strconv"
)

type MerklePath struct {
	path  uint32
	depth uint8
}

func NewMerklePath(bits string) (*MerklePath, error) {
	if bits == "" {
		return &MerklePath{}, nil
	}
	if len(bits) > 32 {
		return nil, fmt.Errorf("path too long")
	}

	path, err := strconv.ParseInt(bits, 2, 33)
	if err != nil {
		return nil, err
	}
	result := &MerklePath{
		path:  uint32(path),
		depth: uint8(len(bits)),
	}
	return result, nil
}

// Scan implements Scanner for database/sql.
func (p *MerklePath) Scan(src interface{}) error {
	value, ok := src.(string)
	if !ok {
		return fmt.Errorf("can't scan %T into Uint256", src)
	}
	_, err := NewMerklePath(value)
	if err != nil {
		return err
	}
	return nil
}

// Value implements valuer for database/sql.
func (p MerklePath) Value() (driver.Value, error) {
	return strconv.FormatInt(int64(p.path), 2), nil
}

// Move pointer left/right on the same level
func (p MerklePath) Add(value uint32) (*MerklePath, error) {
	newPath := p.path + value
	if newPath < p.path {
		return nil, fmt.Errorf("uint32 overflow")
	}
	maxNodeIndex := (uint32(1) << p.depth) - 1
	if newPath > maxNodeIndex {
		return nil, fmt.Errorf("invalid index %d at depth %d", newPath, p.depth)
	}
	p.path = newPath
	return &p, nil
}

func (p MerklePath) Sub(value uint32) (*MerklePath, error) {
	newPath := p.path - value
	if newPath > p.path {
		return nil, fmt.Errorf("uint32 underflow")
	}
	p.path = newPath
	return &p, nil
}

func (p *MerklePath) Parent() (*MerklePath, error) {
	if p.depth == 0 {
		return nil, fmt.Errorf("cannot get parent at depth 0")
	}
	return &MerklePath{
		path:  p.path >> 1,
		depth: p.depth - 1,
	}, nil
}

func (p *MerklePath) Child(right bool) (*MerklePath, error) {
	if p.depth >= 32 {
		return nil, fmt.Errorf("cannot have a path deeper then 32")
	}
	var bit uint32
	if right {
		bit = 1
	}
	return &MerklePath{
		path:  p.path<<1 + bit,
		depth: p.depth + 1,
	}, nil
}

func (p *MerklePath) Sibling() (*MerklePath, error) {
	if p.path%2 == 0 {
		return p.Add(1)
	}
	return p.Sub(1)
}
