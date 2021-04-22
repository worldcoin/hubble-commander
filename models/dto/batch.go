package dto

import "github.com/Worldcoin/hubble-commander/models"

type BatchWithCommitments struct {
	models.Batch
	Commitments []models.Commitment
}
