package api

import (
	"context"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/o11y"
	"github.com/Worldcoin/hubble-commander/storage"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var getUserStatesAPIErrors = map[error]*APIError{
	storage.AnyNotFoundError: NewAPIError(99003, "user states not found"),
}

func (a *API) GetUserStates(ctx context.Context, publicKey *models.PublicKey) ([]dto.UserStateWithID, error) {
	batch, err := a.unsafeGetUserStates(ctx, publicKey)
	if err != nil {
		return nil, sanitizeError(err, getUserStatesAPIErrors)
	}

	return batch, nil
}

func (a *API) unsafeGetUserStates(ctx context.Context, publicKey *models.PublicKey) ([]dto.UserStateWithID, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("publicKey", publicKey.String()))

	log.WithFields(o11y.TraceFields(ctx)).Infof("Getting leaves for public key: %s", publicKey.String())

	leaves, err := a.storage.GetStateLeavesByPublicKey(publicKey)
	if err != nil {
		span.SetAttributes(attribute.String("error", err.Error()))
		log.WithFields(o11y.TraceFields(ctx)).Errorf("Error getting leaves by public key: %v", err)
		return nil, err
	}

	userStates := make([]dto.UserStateWithID, 0, len(leaves))
	for i := range leaves {
		userStates = append(userStates, dto.MakeUserStateWithID(&leaves[i]))
	}

	return userStates, nil
}
