package api

import (
	"context"

	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/o11y"
	"github.com/Worldcoin/hubble-commander/storage"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var getUserStateAPIErrors = map[error]*APIError{
	storage.AnyNotFoundError: NewAPIError(99002, "user state not found"),
}

func (a *API) GetUserState(ctx context.Context, id uint32) (*dto.UserStateWithID, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.Int64("hubble.stateID", int64(id)))

	log.WithFields(o11y.TraceFields(ctx)).Infof("Getting state for id: %d", id)

	userState, err := a.unsafeGetUserState(id)
	if err != nil {
		// WithFields allows APM to associate these lines with the trace
		log.WithFields(o11y.TraceFields(ctx)).Errorf("Error getting user state: %v", err)

		return nil, sanitizeError(err, getUserStateAPIErrors)
	}

	return userState, nil
}

func (a *API) unsafeGetUserState(id uint32) (*dto.UserStateWithID, error) {
	userState, err := a.storage.GetPendingUserState(id)
	if err != nil {
		return nil, err
	}

	addressibleValue := dto.MakeUserStateWithID(id, userState)
	return &addressibleValue, nil
}
