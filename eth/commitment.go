package eth

import (
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/ethereum/go-ethereum/accounts/abi"
)

type decodeCommitmentsFunc func(rollupABI *abi.ABI, calldata []byte) ([]encoder.Commitment, error)

func decodeTxCommitments(rollupABI *abi.ABI, calldata []byte) ([]encoder.Commitment, error) {
	commitments, err := encoder.DecodeTransferBatchCalldata(rollupABI, calldata)
	if err != nil {
		return nil, err
	}
	return encoder.DecodedCommitmentsToCommitments(commitments...), nil
}

func decodeMMCommitments(rollupABI *abi.ABI, calldata []byte) ([]encoder.Commitment, error) {
	commitments, err := encoder.DecodeMMBatchCalldata(rollupABI, calldata)
	if err != nil {
		return nil, err
	}
	return encoder.DecodedMMCommitmentsToCommitments(commitments...), nil
}
