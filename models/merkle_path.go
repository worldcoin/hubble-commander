package models

import (
	"database/sql/driver"
	"encoding/binary"
	"fmt"
	"strconv"
)

type MerklePath struct {
	Path  uint32
	Depth uint8
}

// Root is represented by empty string
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
		Path:  uint32(path),
		Depth: uint8(len(bits)),
	}
	return result, nil
}

func MakeMerklePathFromStateID(stateID uint32) MerklePath {
	return MerklePath{
		Path:  stateID,
		Depth: 32,
	}
}

// Scan implements Scanner for database/sql.
func (p *MerklePath) Scan(src interface{}) error {
	value, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("can't scan %T into MerklePath", src)
	}
	path, err := NewMerklePath(string(value[1:]))
	if err != nil {
		return err
	}
	p.Path = path.Path
	p.Depth = path.Depth
	return nil
}

// Value implements valuer for database/sql.
func (p MerklePath) Value() (driver.Value, error) {
	path := strconv.FormatInt(int64(p.Path), 2)
	return fmt.Sprintf("%0*s", p.Depth+1, path), nil
}

func (p *MerklePath) Bytes() []byte {
	bytes := make([]byte, 5)

	bytes[0] = p.Depth
	binary.LittleEndian.PutUint32(bytes[1:5], p.Path)

	return bytes
}

func (p *MerklePath) SetBytes(data []byte) error {
	if len(data) != 5 {
		return fmt.Errorf("invalid length")
	}
	p.Depth = data[0]
	p.Path = binary.LittleEndian.Uint32(data[1:5])
	return nil
}

// Move pointer left/right on the same level
func (p MerklePath) Add(value uint32) (*MerklePath, error) {
	newPath := p.Path + value
	if newPath < p.Path {
		return nil, fmt.Errorf("uint32 overflow")
	}
	maxNodeIndex := (uint32(1) << p.Depth) - 1
	if newPath > maxNodeIndex {
		return nil, fmt.Errorf("invalid index %d at depth %d", newPath, p.Depth)
	}
	p.Path = newPath
	return &p, nil
}

func (p MerklePath) Sub(value uint32) (*MerklePath, error) {
	newPath := p.Path - value
	if newPath > p.Path {
		return nil, fmt.Errorf("uint32 underflow")
	}
	p.Path = newPath
	return &p, nil
}

func (p *MerklePath) Parent() (*MerklePath, error) {
	if p.Depth == 0 {
		return nil, fmt.Errorf("cannot get parent at depth 0")
	}
	return &MerklePath{
		Path:  p.Path >> 1,
		Depth: p.Depth - 1,
	}, nil
}

func (p *MerklePath) Child(right bool) (*MerklePath, error) {
	if p.Depth >= 32 {
		return nil, fmt.Errorf("cannot have a path deeper then 32")
	}
	var bit uint32
	if right {
		bit = 1
	}
	return &MerklePath{
		Path:  p.Path<<1 + bit,
		Depth: p.Depth + 1,
	}, nil
}

func (p *MerklePath) Sibling() (*MerklePath, error) {
	if p.IsLeftNode() {
		return p.Add(1)
	}
	return p.Sub(1)
}

func (p *MerklePath) GetWitnessPaths() ([]MerklePath, error) {
	witnesses := make([]MerklePath, 0, p.Depth)
	currentPath := p
	isRoot := false

	for !isRoot {
		sibling, err := currentPath.Sibling()
		if err != nil {
			return nil, err
		}
		witnesses = append(witnesses, *sibling)
		currentPath, err = currentPath.Parent()
		if err != nil {
			return nil, err
		}
		if currentPath.Depth == 0 {
			isRoot = true
		}
	}

	return witnesses, nil
}

func (p *MerklePath) IsLeftNode() bool {
	return p.Path%2 == 0
}

func (p *MerklePath) IsRightNode() bool {
	return !p.IsLeftNode()
}

type NamespacedMerklePath struct {
	Namespace string
	Path      MerklePath
}
