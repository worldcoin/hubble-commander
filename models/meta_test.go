package models

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/stretchr/testify/require"
)

func TestMassMigrationMeta_Bytes_ReturnsACopy(t *testing.T) {
	meta := MassMigrationMeta{
		SpokeID:     1,
		TokenID:     MakeUint256(2),
		Amount:      MakeUint256(3),
		FeeReceiver: 4,
	}

	expected := meta

	bytes := meta.Bytes()
	bytes[0] = 9

	require.Equal(t, expected, meta)
}

func TestMassMigrationMeta_SetBytes(t *testing.T) {
	meta := MassMigrationMeta{
		SpokeID:     1,
		TokenID:     MakeUint256(2),
		Amount:      MakeUint256(3),
		FeeReceiver: 4,
	}

	bytes := meta.Bytes()
	newMeta := MassMigrationMeta{}
	err := newMeta.SetBytes(bytes)
	require.NoError(t, err)
	require.Equal(t, meta, newMeta)
}

func TestMassMigrationMeta_SetBytes_InvalidLength(t *testing.T) {
	bytes := utils.PadLeft([]byte{1, 2, 3}, 130)
	meta := MassMigrationMeta{}
	err := meta.SetBytes(bytes)
	require.NotNil(t, err)
	require.ErrorIs(t, err, ErrInvalidLength)
}
