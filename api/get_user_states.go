package api

import (
	"context"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/o11y"
	"github.com/Worldcoin/hubble-commander/storage"
	"github.com/pkg/errors"
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
	span.SetAttributes(attribute.String("hubble.publicKey", publicKey.String()))

	userStates := make([]dto.UserStateWithID, 0)

	err := a.storage.ExecuteInReadWriteTransaction(func(txStorage *storage.Storage) error {
		leaves, err := txStorage.GetStateLeavesByPublicKey(publicKey)
		if err != nil && !storage.IsNotFoundError(err) {
			span.SetAttributes(attribute.String("hubble.error", err.Error()))
			log.WithFields(o11y.TraceFields(ctx)).Errorf("Error getting leaves by public key: %v", err)
			return err
		}

		for i := range leaves {
			stateID := leaves[i].StateID

			pendingState, innerErr := txStorage.GetPendingUserState(stateID)
			if innerErr != nil {
				return innerErr
			}

			userStates = append(userStates, dto.MakeUserStateWithID(stateID, pendingState))
		}

		pendingUserStates, err := txStorage.GetPendingUserStates(publicKey)
		if err != nil {
			return err
		}
		for i := range pendingUserStates {
			userStates = append(
				userStates,
				dto.MakePendingUserStateWithID(&pendingUserStates[i]),
			)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	if len(userStates) == 0 {
		span.SetAttributes(attribute.String("hubble.error", "user states not found"))
		return nil, errors.WithStack(
			storage.NewNotFoundError("user states"),
		)
	}

	return userStates, nil
}
