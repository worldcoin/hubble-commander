package eth

import (
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/ethereum/go-ethereum/accounts/abi"
)

type decodeCommitmentsFunc func(rollupABI *abi.ABI, calldata []byte) ([]encoder.GenericCommitment, error)

func decodedTxCommitments(rollupABI *abi.ABI, calldata []byte) ([]encoder.GenericCommitment, error) {
	commitments, err := encoder.DecodeBatchCalldata(rollupABI, calldata)
	if err != nil {
		return nil, err
	}

	result := make([]encoder.GenericCommitment, 0, len(commitments))
	for i := range commitments {
		result = append(result, &commitments[i])
	}
	return result, nil
}

func decodedMMCommitments(rollupABI *abi.ABI, calldata []byte) ([]encoder.GenericCommitment, error) {
	commitments, err := encoder.DecodeMMBatchCalldata(rollupABI, calldata)
	if err != nil {
		return nil, err
	}

	result := make([]encoder.GenericCommitment, 0, len(commitments))
	for i := range commitments {
		result = append(result, &commitments[i])
	}
	return result, nil
}
