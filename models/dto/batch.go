package dto

import (
	"github.com/Worldcoin/hubble-commander/models"
)

type BatchWithCommitments struct {
	models.BatchWithAccountRoot
	Commitments []models.CommitmentWithTokenID
}
