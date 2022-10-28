package dto

import "github.com/Worldcoin/hubble-commander/models"

type PendingState struct {
	Nonce   models.Uint256
	Balance models.Uint256
}

type PendingStateDiff struct {
	StateID    uint32
	OldNonce   *models.Uint256
	OldBalance *models.Uint256
	NewNonce   *models.Uint256
	NewBalance *models.Uint256
}
