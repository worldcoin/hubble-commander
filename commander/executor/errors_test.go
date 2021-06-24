package executor

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsDisputableTransferError(t *testing.T) {
	transitionError := NewDisputableTransferError(TransitionError, "validation duck")
	signatureError := NewDisputableTransferError(SignatureError, "signature duck")

	require.True(t, IsDisputableTransferError(transitionError))
	require.True(t, IsDisputableTransferError(signatureError))
}
