package dto

import "github.com/Worldcoin/hubble-commander/models"

type RecomputePendingState struct {
	OldNonce   models.Uint256
	OldBalance models.Uint256
	NewNonce   models.Uint256
	NewBalance models.Uint256
}
