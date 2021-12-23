package db

import (
	"fmt"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/stored"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	bh "github.com/timshannon/badgerhold/v4"
)

var errPassedByPointer = fmt.Errorf("pointer was passed to Encode, pass by value instead")

// nolint:gocyclo, funlen
// Encode Remember to provide cases for both value and pointer types when adding new encoders
// TODO shorten this function by using ByteEncoder interface
func Encode(value interface{}) ([]byte, error) {
	switch v := value.(type) {
	case models.AccountNode:
		return stored.EncodeHash(&v.DataHash)
	case *models.AccountNode:
		return nil, errors.WithStack(errPassedByPointer)
	case models.AccountLeaf:
		return v.Bytes(), nil
	case *models.AccountLeaf:
		return nil, errors.WithStack(errPassedByPointer)
	case stored.Batch:
		return v.Bytes(), nil
	case *stored.Batch:
		return nil, errors.WithStack(errPassedByPointer)
	case models.ChainState:
		return v.Bytes(), nil
	case *models.ChainState:
		return nil, errors.WithStack(errPassedByPointer)
	case models.CommitmentID:
		return v.Bytes(), nil
	case *models.CommitmentID:
		return stored.EncodeCommitmentIDPointer(v), nil
	case models.PendingDeposit:
		return v.Bytes(), nil
	case *models.PendingDeposit:
		return nil, errors.WithStack(errPassedByPointer)
	case models.DepositID:
		return v.Bytes(), nil
	case *models.DepositID:
		return nil, errors.WithStack(errPassedByPointer)
	case models.PendingDepositSubTree:
		return v.Bytes(), nil
	case *models.PendingDepositSubTree:
		return nil, errors.WithStack(errPassedByPointer)
	case models.NamespacedMerklePath:
		return v.Bytes(), nil
	case *models.NamespacedMerklePath:
		return nil, errors.WithStack(errPassedByPointer)
	case models.MerkleTreeNode:
		return stored.EncodeHash(&v.DataHash)
	case *models.MerkleTreeNode:
		return nil, errors.WithStack(errPassedByPointer)
	case models.PublicKey:
		return v.Bytes(), nil
	case *models.PublicKey:
		return nil, errors.WithStack(errPassedByPointer)
	case stored.StateLeaf:
		return v.Bytes(), nil
	case *stored.StateLeaf:
		return nil, errors.WithStack(errPassedByPointer)
	case models.StateUpdate:
		return v.Bytes(), nil
	case *models.StateUpdate:
		return nil, errors.WithStack(errPassedByPointer)
	case stored.Commitment:
		return v.Bytes(), nil
	case *stored.Commitment:
		return nil, errors.WithStack(errPassedByPointer)
	case stored.Tx:
		return v.Bytes(), nil
	case *stored.Tx:
		return nil, errors.WithStack(errPassedByPointer)
	case stored.TxReceipt:
		return v.Bytes(), nil
	case *stored.TxReceipt:
		return nil, errors.WithStack(errPassedByPointer)
	case models.Uint256:
		return v.Bytes(), nil
	case *models.Uint256:
		return nil, errors.WithStack(errPassedByPointer)
	case common.Hash:
		return v.Bytes(), nil
	case *common.Hash:
		return stored.EncodeHashPointer(v), nil
	case string:
		return stored.EncodeString(v), nil
	case *string:
		return nil, errors.WithStack(errPassedByPointer)
	case uint32:
		return stored.EncodeUint32(v), nil
	case *uint32:
		return stored.EncodeUint32Pointer(v), nil
	case uint64:
		return stored.EncodeUint64(v), nil
	case *uint64:
		return nil, errors.WithStack(errPassedByPointer)
	case models.RegisteredToken:
		return v.Contract.Bytes(), nil
	case *models.RegisteredToken:
		return nil, errors.WithStack(errPassedByPointer)
	case models.RegisteredSpoke:
		return v.Contract.Bytes(), nil
	case *models.RegisteredSpoke:
		return nil, errors.WithStack(errPassedByPointer)
	case bh.KeyList:
		return EncodeKeyList(&v)
	default:
		return bh.DefaultEncode(value)
	}
}

// nolint:gocyclo, funlen
func Decode(data []byte, value interface{}) error {
	switch v := value.(type) {
	case *models.AccountNode:
		return stored.DecodeHash(data, &v.DataHash)
	case *models.AccountLeaf:
		return v.SetBytes(data)
	case *models.ChainState:
		return v.SetBytes(data)
	case *models.CommitmentID:
		return decodeCommitmentIDPointer(data, &value, v)
	case *models.PendingDeposit:
		return v.SetBytes(data)
	case *models.DepositID:
		return v.SetBytes(data)
	case *models.PendingDepositSubTree:
		return v.SetBytes(data)
	case *models.NamespacedMerklePath:
		return v.SetBytes(data)
	case *stored.Batch:
		return v.SetBytes(data)
	case *models.MerkleTreeNode:
		return stored.DecodeHash(data, &v.DataHash)
	case *models.PublicKey:
		return v.SetBytes(data)
	case *stored.StateLeaf:
		return v.SetBytes(data)
	case *models.StateUpdate:
		return v.SetBytes(data)
	case *stored.Commitment:
		return v.SetBytes(data)
	case *stored.Tx:
		return v.SetBytes(data)
	case *stored.TxReceipt:
		return v.SetBytes(data)
	case *models.Uint256:
		v.SetBytes(data)
		return nil
	case *common.Hash:
		return decodeHashPointer(data, &value, v)
	case *string:
		return stored.DecodeString(data, v)
	case *uint32:
		return decodeUint32Pointer(data, &value, v)
	case *uint64:
		return stored.DecodeUint64(data, v)
	case *models.RegisteredToken:
		v.Contract.SetBytes(data)
		return nil
	case *models.RegisteredSpoke:
		v.Contract.SetBytes(data)
		return nil
	case *bh.KeyList:
		return DecodeKeyList(data, v)
	default:
		return bh.DefaultDecode(data, value)
	}
}

// nolint: gocritic
func decodeHashPointer(data []byte, value *interface{}, dst *common.Hash) error {
	if len(data) == 32 {
		return stored.DecodeHash(data, dst)
	}
	if data[0] == 1 {
		return stored.DecodeHash(data[1:], dst)
	}
	*value = nil
	return nil
}

// nolint: gocritic
func decodeCommitmentIDPointer(data []byte, value *interface{}, dst *models.CommitmentID) error {
	if len(data) == 33 {
		return dst.SetBytes(data)
	}
	if data[0] == 1 {
		return dst.SetBytes(data[1:])
	}
	*value = nil
	return nil
}

// nolint: gocritic
func decodeUint32Pointer(data []byte, value *interface{}, dst *uint32) error {
	if len(data) == 4 {
		return stored.DecodeUint32(data, dst)
	}
	if data[0] == 1 {
		return stored.DecodeUint32(data[1:], dst)
	}
	*value = nil
	return nil
}

func DecodeKey(data []byte, key interface{}, prefix []byte) error {
	return Decode(data[len(prefix):], key)
}
