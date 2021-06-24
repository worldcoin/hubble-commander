package executor

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsDisputableTransferError(t *testing.T) {
	transitionError := NewDisputableTransferError("validation duck", TransitionError)
	signatureError := NewDisputableTransferError("signature duck", SignatureError)

	require.True(t, IsDisputableTransferError(transitionError))
	require.True(t, IsDisputableTransferError(signatureError))
}
