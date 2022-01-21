package admin

import (
	"context"
	"fmt"

	"github.com/Worldcoin/hubble-commander/api/rpc"
)

var (
	errMissingAuthKey = fmt.Errorf("missing authentication key")
	errInvalidAuthKey = fmt.Errorf("invalid authentication key value")
)

func (a *API) verifyAuthKey(ctx context.Context) error {
	authKeyValue := ctx.Value(rpc.AuthKey)
	if authKeyValue == nil || authKeyValue == "" {
		return errMissingAuthKey
	}

	if authKeyValue != *a.cfg.AuthenticationKey {
		return errInvalidAuthKey
	}

	return nil
}
