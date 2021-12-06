package eth

import (
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/ethereum/go-ethereum/common"
)

type GenericCommitment interface {
	BodyHash(accountRoot common.Hash) *common.Hash
	LeafHash(accountRoot common.Hash) common.Hash
}

//TODO-sync: move to encoder package and remove above interface
func decodedTxCommitmentsToCommitments(commitments []encoder.DecodedCommitment) []encoder.GenericCommitment {
	result := make([]encoder.GenericCommitment, 0, len(commitments))
	for i := range commitments {
		result = append(result, &commitments[i])
	}
	return result
}

func decodedMMCommitmentsToCommitments(commitments []encoder.DecodedMMCommitment) []encoder.GenericCommitment {
	result := make([]encoder.GenericCommitment, 0, len(commitments))
	for i := range commitments {
		result = append(result, &commitments[i])
	}
	return result
}
