package utils

import (
	"sort"
	"testing"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
)

func TestEventBefore(t *testing.T) {
	first := &types.Log{
		BlockNumber: 1,
		Index:       1,
	}
	second := &types.Log{
		BlockNumber: 1,
		Index:       2,
	}
	third := &types.Log{
		BlockNumber: 2,
		Index:       1,
	}

	require.True(t, EventBefore(first, second))
	require.True(t, EventBefore(first, third))
	require.True(t, EventBefore(second, third))

	require.False(t, EventBefore(second, first))
	require.False(t, EventBefore(third, first))
	require.False(t, EventBefore(third, second))

	firstCopy := *first
	require.False(t, EventBefore(first, &firstCopy))
	require.False(t, EventBefore(&firstCopy, first))
}

func TestEventBefore_SortsWell(t *testing.T) {
	first := &types.Log{
		BlockNumber: 1,
		Index:       1,
	}
	second := &types.Log{
		BlockNumber: 1,
		Index:       2,
	}
	third := &types.Log{
		BlockNumber: 2,
		Index:       1,
	}

	events := []*types.Log{second, third, first}
	sort.Slice(events, func(i, j int) bool {
		return EventBefore(events[i], events[j])
	})

	require.Equal(t, []*types.Log{first, second, third}, events)
}
