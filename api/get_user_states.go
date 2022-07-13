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
	span.SetAttributes(attribute.String("hubble.publicKey", publicKey.String()))

	log.WithFields(o11y.TraceFields(ctx)).Infof("Getting leaves for public key: %s", publicKey.String())

	leaves, err := a.storage.GetStateLeavesByPublicKey(publicKey)
	if err != nil {
		span.SetAttributes(attribute.String("hubble.error", err.Error()))
		log.WithFields(o11y.TraceFields(ctx)).Errorf("Error getting leaves by public key: %v", err)
		return nil, err
	}

	// TODO: we're not opening a transaction so there's no guarantee this is a
	//       consistent snapshot. You might temporarily lose or gain money if you're
	//       sending money between your accounts. We should open a txn!

	userStates := make([]dto.UserStateWithID, 0, len(leaves))
	for i := range leaves {
		stateID := leaves[i].StateID

		pendingState, err := a.storage.GetPendingUserState(stateID)
		if err != nil {
			return nil, err
		}

		userStates = append(userStates, dto.MakeUserStateWithID(stateID, pendingState))
	}

	return userStates, nil
}
