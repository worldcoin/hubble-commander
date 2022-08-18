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

// Encode Remember to provide cases for both value and pointer types when adding new encoders
// TODO shorten this function by using ByteEncoder interface
//
//nolint:gocyclo, funlen
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
		return v.Bytes(), nil
	case models.CommitmentSlot:
		return v.Bytes(), nil
	case *models.CommitmentSlot:
		return v.Bytes(), nil
	case models.PendingDeposit:
		return v.Bytes(), nil
	case *models.PendingDeposit:
		return nil, errors.WithStack(errPassedByPointer)
	case models.DepositID:
		return v.Bytes(), nil
	case *models.DepositID:
		return nil, errors.WithStack(errPassedByPointer)
	case models.PendingDepositSubtree:
		return v.Bytes(), nil
	case *models.PendingDepositSubtree:
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
	case stored.FlatStateLeaf:
		return v.Bytes(), nil
	case *stored.FlatStateLeaf:
		return nil, errors.WithStack(errPassedByPointer)
	case models.StateUpdate:
		return v.Bytes(), nil
	case *models.StateUpdate:
		return nil, errors.WithStack(errPassedByPointer)
	case stored.Commitment:
		return v.Bytes(), nil
	case *stored.Commitment:
		return nil, errors.WithStack(errPassedByPointer)
	case stored.PendingTx:
		return v.Bytes(), nil
	case *stored.PendingTx:
		return nil, errors.WithStack(errPassedByPointer)
	case stored.BatchedTx:
		return v.Bytes(), nil
	case *stored.BatchedTx:
		return nil, errors.WithStack(errPassedByPointer)
	case stored.FailedTx:
		return v.Bytes(), nil
	case *stored.FailedTx:
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
	case models.PendingStakeWithdrawal:
		return v.Bytes(), nil
	case *models.PendingStakeWithdrawal:
		return nil, errors.WithStack(errPassedByPointer)
	case bh.KeyList:
		return EncodeKeyList(&v)
	case []byte:
		return v, nil
	default:
		return bh.DefaultEncode(value)
	}
}

//nolint:gocyclo, funlen
func Decode(data []byte, value interface{}) error {
	switch v := value.(type) {
	case *models.AccountNode:
		return stored.DecodeHash(data, &v.DataHash)
	case *models.AccountLeaf:
		return v.SetBytes(data)
	case *models.ChainState:
		return v.SetBytes(data)
	case *models.CommitmentID:
		return v.SetBytes(data)
	case *models.CommitmentSlot:
		return v.SetBytes(data)
	case *models.PendingDeposit:
		return v.SetBytes(data)
	case *models.DepositID:
		return v.SetBytes(data)
	case *models.PendingDepositSubtree:
		return v.SetBytes(data)
	case *models.NamespacedMerklePath:
		return v.SetBytes(data)
	case *stored.Batch:
		return v.SetBytes(data)
	case *models.MerkleTreeNode:
		return stored.DecodeHash(data, &v.DataHash)
	case *models.PublicKey:
		return v.SetBytes(data)
	case *stored.FlatStateLeaf:
		return v.SetBytes(data)
	case *models.StateUpdate:
		return v.SetBytes(data)
	case *stored.Commitment:
		return v.SetBytes(data)
	case *stored.PendingTx:
		return v.SetBytes(data)
	case *stored.BatchedTx:
		return v.SetBytes(data)
	case *stored.FailedTx:
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
	case *models.PendingStakeWithdrawal:
		return v.SetBytes(data)
	case *bh.KeyList:
		return DecodeKeyList(data, v)
	case []byte:
		copy(data, v)
		return nil
	default:
		return bh.DefaultDecode(data, value)
	}
}

//nolint: gocritic
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

//nolint: gocritic
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
