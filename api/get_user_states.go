package api

import (
	"context"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/o11y"
	"github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils/consts"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var getUserStatesTracer = otel.Tracer("api.getUserStates")

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

	err := a.storage.ExecuteInReadWriteTransactionWithSpan(ctx, func(txCtx context.Context, txStorage *storage.Storage) error {
		leaves, err := func() ([]models.StateLeaf, error) {
			_, innerSpan := getUserStatesTracer.Start(txCtx, "GetStateLeavesByPublicKey")
			defer innerSpan.End()

			return txStorage.GetStateLeavesByPublicKey(publicKey)
		}()
		if err != nil && !storage.IsNotFoundError(err) {
			span.SetAttributes(attribute.String("hubble.error", err.Error()))
			log.WithFields(o11y.TraceFields(txCtx)).Errorf("Error getting leaves by public key: %v", err)
			return err
		}

		for i := range leaves {
			stateID := leaves[i].StateID

			pendingState, innerErr := func() (*models.UserState, error) {
				_, innerSpan := getUserStatesTracer.Start(txCtx, "GetPendingUserState")
				defer innerSpan.End()

				return txStorage.GetPendingUserState(stateID)
			}()
			if innerErr != nil {
				return innerErr
			}

			userStates = append(userStates, dto.MakeUserStateWithID(stateID, pendingState))
		}

		pendingC2TState, err := func() (*models.UserState, error) {
			_, innerSpan := getUserStatesTracer.Start(txCtx, "GetPendingC2TState")
			defer innerSpan.End()

			return txStorage.GetPendingC2TState(publicKey)
		}()
		if err != nil {
			return err
		}
		if pendingC2TState != nil {
			userStates = append(
				userStates,
				dto.MakeUserStateWithID(consts.PendingID, pendingC2TState),
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

	totalBalance := models.NewUint256(0)
	for _, userState := range userStates {
		totalBalance = totalBalance.Add(&userState.Balance)
	}
	span.SetAttributes(attribute.String("hubble.response.totalBalance", totalBalance.String()))
	span.SetAttributes(attribute.Int("hubble.response.userStatesCount", len(userStates)))

	return userStates, nil
}
