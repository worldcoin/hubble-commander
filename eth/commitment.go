package eth

import (
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/ethereum/go-ethereum/common"
)

type commitment interface {
	BodyHash(accountRoot common.Hash) *common.Hash
	LeafHash(accountRoot common.Hash) common.Hash
}

func decodedTxCommitmentsToCommitments(commitments []encoder.DecodedCommitment) []commitment {
	result := make([]commitment, 0, len(commitments))
	for i := range commitments {
		result = append(result, &commitments[i])
	}
	return result
}

func decodedMMCommitmentsToCommitments(commitments []encoder.DecodedMassMigrationCommitment) []commitment {
	result := make([]commitment, 0, len(commitments))
	for i := range commitments {
		result = append(result, &commitments[i])
	}
	return result
}
