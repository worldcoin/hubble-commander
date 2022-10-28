package storage

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/dgraph-io/badger/v3"
	"github.com/pkg/errors"
)

type PendingState interface {
	GetPendingState(stateID uint32) (*dto.PendingState, error)
	SetPendingState(stateID uint32, state dto.PendingState) error
	HasPendingState(stateID uint32) (bool, error)
	AsMap() (map[uint32]dto.PendingState, error)
}

type StateHelper struct{ ps PendingState }

func (sh *StateHelper) BalanceAdd(stateID uint32, amount *models.Uint256) error {
	state, err := sh.ps.GetPendingState(stateID)
	if err != nil {
		return err
	}

	balance := state.Balance.Add(amount)

	return sh.ps.SetPendingState(stateID, dto.PendingState{Nonce: state.Nonce, Balance: *balance})
}

func (sh *StateHelper) BalanceSub(stateID uint32, amount *models.Uint256) error {
	state, err := sh.ps.GetPendingState(stateID)
	if err != nil {
		return err
	}

	if state.Balance.Cmp(amount) < 0 {
		return errors.WithStack(ErrBalanceTooLow)
	}

	balance := state.Balance.Sub(amount)
	return sh.ps.SetPendingState(stateID, dto.PendingState{Nonce: state.Nonce, Balance: *balance})
}

func (sh *StateHelper) NonceIncr(stateID uint32) error {
	state, err := sh.ps.GetPendingState(stateID)
	if err != nil {
		return err
	}

	one := models.MakeUint256(1)
	nonce := state.Nonce.Add(&one)

	return sh.ps.SetPendingState(stateID, dto.PendingState{Nonce: *nonce, Balance: state.Balance})
}

type BadgerPendingState struct {
	storage *Storage
}

func MakeBadgerPendingState(storage *Storage) *BadgerPendingState {
	return &BadgerPendingState{storage: storage}
}

func (bps *BadgerPendingState) GetPendingState(stateID uint32) (*dto.PendingState, error) {
	key := pendingStateKey(stateID)
	value, err := bps.storage.rawLookup(key)

	if err == nil {
		decodedNonce, decodedBalance := decodePendingState(value)
		return &dto.PendingState{Nonce: decodedNonce, Balance: decodedBalance}, nil
	}

	if !errors.Is(err, badger.ErrKeyNotFound) {
		return nil, err
	}

	// errors.Is(err, badger.ErrKeyNotFound) == true

	state, err := bps.storage.StateTree.Leaf(stateID)
	if err != nil {
		return nil, err
	}

	return &dto.PendingState{Nonce: state.UserState.Nonce, Balance: state.UserState.Balance}, nil
}

func (bps *BadgerPendingState) SetPendingState(stateID uint32, state dto.PendingState) error {
	key := pendingStateKey(stateID)
	value := encodePendingState(state.Nonce, state.Balance)
	return bps.storage.rawSet(key, value)
}

func (bps *BadgerPendingState) HasPendingState(stateID uint32) (bool, error) {
	panic("unimplemented")
}

func (bps *BadgerPendingState) AsMap() (map[uint32]dto.PendingState, error) {
	result := make(map[uint32]dto.PendingState)

	var maxInt uint32 = ^uint32(0)

	// TODO: why is this unsafe? It calls Badger.View which runs fn inside a new read-only txn

	states, err := bps.storage.unsafeGetPendingStates(0, maxInt)
	if err != nil {
		return nil, err
	}

	for i := range states {
		state := states[i]
		result[state.StateID] = dto.PendingState{
			Nonce:   state.Nonce,
			Balance: state.Balance,
		}
	}

	return result, nil
}

type MemoryPendingState struct {
	fallback *Storage
	states   map[uint32]dto.PendingState
}

func MakeMemoryPendingState(fallback *Storage) *MemoryPendingState {
	states := make(map[uint32]dto.PendingState)
	return &MemoryPendingState{
		fallback: fallback,
		states:   states,
	}
}

// will throw an error if this key is not present either here or in the StateTree
func (mps *MemoryPendingState) GetPendingState(stateID uint32) (*dto.PendingState, error) {
	result, present := mps.states[stateID]
	if !present {
		state, err := mps.fallback.StateTree.Leaf(stateID)
		if err != nil {
			return nil, err
		}

		result = dto.PendingState{
			Nonce:   state.Nonce,
			Balance: state.Balance,
		}

		mps.states[stateID] = result
	}

	return &result, nil
}

func (mps *MemoryPendingState) SetPendingState(stateID uint32, state dto.PendingState) error {
	mps.states[stateID] = state
	return nil
}

func (mps *MemoryPendingState) HasPendingState(stateID uint32) (bool, error) {
	panic("unimplemented")
}

func (mps *MemoryPendingState) AsMap() (map[uint32]dto.PendingState, error) {
	result := make(map[uint32]dto.PendingState)
	for k, v := range mps.states {
		result[k] = v
	}

	return result, nil
}
