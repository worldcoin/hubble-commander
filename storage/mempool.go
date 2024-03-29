package storage

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"sort"

	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/Worldcoin/hubble-commander/models/stored"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/Worldcoin/hubble-commander/utils/consts"
	"github.com/dgraph-io/badger/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// TODO: make sure you wrap all errors whenever you talk to badger

// TODO: Why does storage/stored_transaction.go create a separate TranasactionStorage?
//       Should we do the same here?

/// These are used by the API

var pendingStatePrefix = []byte("PendingAccountState")
var pendingPubkeyBalancePrefix = []byte("PendingPubKeyBalance")
var pendingTxPrefix = []byte("PendingTxs")
var migratedPubkeyStatePrefix = []byte("migration:PubKeyPendingState")

var ErrBalanceTooLow = fmt.Errorf("balance too low")

func pendingTxStateIDPrefix(stateID uint32) []byte {
	encodedStateID := make([]byte, 4)
	binary.BigEndian.PutUint32(encodedStateID, stateID)

	return bytes.Join(
		[][]byte{pendingTxPrefix, encodedStateID},
		[]byte(":"),
	)
}

func pendingTxKey(stateID uint32, nonce uint64) []byte {
	prefix := pendingTxStateIDPrefix(stateID)

	encodedNonce := make([]byte, 8)
	binary.BigEndian.PutUint64(encodedNonce, nonce)

	return bytes.Join(
		[][]byte{prefix, encodedNonce},
		[]byte(":"),
	)
}

func pendingStateKey(stateID uint32) []byte {
	encodedStateID := make([]byte, 4)
	binary.BigEndian.PutUint32(encodedStateID, stateID)

	return bytes.Join(
		[][]byte{pendingStatePrefix, encodedStateID},
		[]byte(":"),
	)
}

func pendingPubkeyBalanceKey(pubkey *models.PublicKey) []byte {
	return bytes.Join(
		[][]byte{pendingPubkeyBalancePrefix, pubkey[:]},
		[]byte(":"),
	)
}

func decodePendingPubkeyBalanceKey(keyBytes []byte) *models.PublicKey {
	pubkeyBytes := keyBytes[len(pendingPubkeyBalancePrefix)+1:]
	var pubkey models.PublicKey
	copy(pubkey[:], pubkeyBytes)

	return &pubkey
}

func decodePendingStateKey(keyBytes []byte) uint32 {
	lastFour := keyBytes[len(keyBytes)-4:]
	return binary.BigEndian.Uint32(lastFour)
}

func encodePendingState(nonce, balance models.Uint256) []byte {
	var result bytes.Buffer
	result.Grow(64)
	result.Write(nonce.Bytes())
	result.Write(balance.Bytes())

	return result.Bytes()
}

func decodePendingState(data []byte) (nonce, balance models.Uint256) {
	if len(data) != 64 {
		panic("corrupted database")
	}

	first32, last32 := data[:32], data[32:]
	nonce.SetBytes(first32)
	balance.SetBytes(last32)

	return nonce, balance
}

func (s *Storage) rawLookup(key []byte) (value []byte, err error) {
	err = s.database.Badger.RawUpdate(func(txn *badger.Txn) error {
		item, innerErr := txn.Get(key)
		if innerErr != nil {
			return errors.WithStack(innerErr)
		}

		value, innerErr = item.ValueCopy(nil)
		return errors.WithStack(innerErr)
	})
	if err != nil {
		value = nil
	}

	return value, err
}

func (s *Storage) hasKey(key []byte) (bool, error) {
	_, err := s.rawLookup(key)

	if err == nil {
		return true, nil
	}

	if errors.Is(err, badger.ErrKeyNotFound) {
		return false, nil
	}

	return false, err
}

func (s *Storage) rawSet(key, value []byte) error {
	return s.database.Badger.RawUpdate(func(txn *badger.Txn) error {
		err := txn.Set(key, value)
		if err != nil {
			return errors.WithStack(err)
		}
		return nil
	})
}

func (s *Storage) alreadyRanPubKeyMigration() (bool, error) {
	alreadyMigrated, err := s.hasKey(migratedPubkeyStatePrefix)
	if err != nil {
		return false, err
	}
	return alreadyMigrated, nil
}

func (s *Storage) markRanPubKeyMigration() error {
	return s.rawSet(migratedPubkeyStatePrefix, []byte("true"))
}

// It is possible that the account does not exist in our pending state and also does not
// exist in the state tree. In that case this function will return a storage.NotFoundError
func (s *Storage) getPendingState(stateID uint32) (
	nonce *models.Uint256,
	balance *models.Uint256,
	err error,
) {
	// TODO: not confident this does the right thing, this wants some tests

	key := pendingStateKey(stateID)
	value, err := s.rawLookup(key)

	if err == nil {
		decodedNonce, decodedBalance := decodePendingState(value)
		return &decodedNonce, &decodedBalance, nil
	}

	if errors.Is(err, badger.ErrKeyNotFound) {
		state, innerErr := s.StateTree.Leaf(stateID)
		if innerErr != nil {
			return nil, nil, innerErr
		}

		nonce = &state.UserState.Nonce
		balance = &state.UserState.Balance
		return nonce, balance, nil
	}

	return nil, nil, err
}

func (s *Storage) setPendingPubkeyBalance(pubkey *models.PublicKey, balance *models.Uint256) error {
	key := pendingPubkeyBalanceKey(pubkey)
	return s.rawSet(key, balance.Bytes())
}

func (s *Storage) getPendingPubkeyBalance(pubkey *models.PublicKey) (*models.Uint256, error) {
	key := pendingPubkeyBalanceKey(pubkey)
	value, err := s.rawLookup(key)
	if err != nil {
		return nil, err
	}

	var balance models.Uint256
	balance.SetBytes(value)

	return &balance, nil
}

func (s *Storage) addToPendingPubkeyBalance(pubkey *models.PublicKey, amount *models.Uint256) error {
	balance, err := s.getPendingPubkeyBalance(pubkey)

	if err != nil && errors.Is(err, badger.ErrKeyNotFound) {
		addressableValue := models.MakeUint256(0)
		balance = &addressableValue
	} else if err != nil {
		return err
	}

	newBalance := balance.Add(amount)
	return s.setPendingPubkeyBalance(pubkey, newBalance)
}

// this is public so tests can call it, nobody else should call it.
func (s *Storage) UnsafeSetPendingState(stateID uint32, nonce, balance models.Uint256) error {
	key := pendingStateKey(stateID)
	value := encodePendingState(nonce, balance)
	return s.rawSet(key, value)
}

func (s *Storage) addToPendingBalance(stateID uint32, amount *models.Uint256) error {
	pendingNonce, pendingBalance, err := s.getPendingState(stateID)
	if err != nil {
		return err
	}

	newBalance := pendingBalance.Add(amount)

	return s.UnsafeSetPendingState(stateID, *pendingNonce, *newBalance)
}

func (s *Storage) GetPendingNonce(stateID uint32) (*models.Uint256, error) {
	pendingNonce, _, err := s.getPendingState(stateID)
	return pendingNonce, err
}

func (s *Storage) GetPendingBalance(stateID uint32) (*models.Uint256, error) {
	_, pendingBalance, err := s.getPendingState(stateID)
	return pendingBalance, err
}

func (s *Storage) GetPendingUserState(stateID uint32) (*models.UserState, error) {
	leaf, err := s.StateTree.Leaf(stateID)
	if err != nil {
		return nil, err
	}

	pendingNonce, pendingBalance, err := s.getPendingState(stateID)
	if err != nil {
		return nil, err
	}

	return &models.UserState{
		PubKeyID: leaf.UserState.PubKeyID,
		TokenID:  leaf.UserState.TokenID,
		Balance:  *pendingBalance,
		Nonce:    *pendingNonce,
	}, nil
}

// SELECT pubkey, SUM(amount)
// FROM pending_transactions
// WHERE tx_type = 'Create2Transfer'
// GROUP BY pubkey
func (s *Storage) pendingPubkeyBalances() (map[models.PublicKey]*models.Uint256, error) {
	keyBalances := make(map[models.PublicKey]*models.Uint256)

	err := s.forEachMempoolTransaction(func(pendingTx *stored.PendingTx) error {
		if pendingTx.TxType != txtype.Create2Transfer {
			return nil
		}

		c2t := pendingTx.ToCreate2Transfer()
		balance, present := keyBalances[c2t.ToPublicKey]
		if !present {
			keyBalances[c2t.ToPublicKey] = &c2t.Amount
		} else {
			keyBalances[c2t.ToPublicKey] = balance.Add(&c2t.Amount)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return keyBalances, nil
}

func (s *Storage) MigratePubKeyPendingState() error {
	alreadyMigrated, err := s.alreadyRanPubKeyMigration()
	if err != nil {
		return err
	}
	if alreadyMigrated {
		return nil
	}

	keyBalances, err := s.pendingPubkeyBalances()
	if err != nil {
		return err
	}

	for pubKey, balance := range keyBalances {
		err = s.setPendingPubkeyBalance(&pubKey, balance)
		if err != nil {
			return err
		}
	}

	return s.markRanPubKeyMigration()
}

// does not return the full list of states, only the pending c2ts
func (s *Storage) GetPendingC2TState(pubkey *models.PublicKey) (*models.UserState, error) {
	pubKeyID := consts.PendingID

	// TODO: pull this pubKeyId lookup into a function and use it the other places we do this
	accounts, err := s.AccountTree.Leaves(pubkey)
	if err == nil && len(accounts) > 0 {
		pubKeyID = accounts[0].PubKeyID
	}

	balance, err := s.getPendingPubkeyBalance(pubkey)
	if err != nil && errors.Is(err, badger.ErrKeyNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &models.UserState{
		PubKeyID: pubKeyID,
		TokenID:  models.MakeUint256(0),
		Balance:  *balance,
		Nonce:    models.MakeUint256(0),
	}, nil
}

// TODO: this assumes the sender & receiver stateIDs are in the state tree,
//       make sure all callers check for this? Maybe return a nice error?
func (s *Storage) AddMempoolTx(tx models.GenericTransaction) error {
	// TODO: should we check that this txn does not already exist as a batchedTx?
	// TODO: should this accept a stored.PendingTx, so we do not accidentally accept
	//       a tx which already has a CommitmentSlot?

	// The caller should have opened a txn but doing it here just in case.

	// The txn gives badger enough information to lock out conflicting API handlers
	// for us. storage/database.go is even nice enough to automatically retry the
	// fn for us in case of ErrConflict, so `unsafeAddMempoolTx` had better be
	// idempotent.

	return s.ExecuteInReadWriteTransaction(func(txStorage *Storage) error {
		return txStorage.unsafeAddMempoolTx(tx)
	})
}

// - assumes we are currently inside a transaction
// - checks that the txn cleanly applies to the pending state but assumes all other
//   validation has already been done (e.g. the signature check)
func (s *Storage) unsafeAddMempoolTx(tx models.GenericTransaction) error {
	// (I) Validate the txn against the pending state

	fromStateID := tx.GetFromStateID()

	pendingNonce, pendingBalance, err := s.getPendingState(fromStateID)
	if err != nil {
		return err
	}

	txNonce := tx.GetNonce()
	if txNonce.Cmp(pendingNonce) != 0 {
		return errors.WithStack(
			fmt.Errorf(
				// TODO: how do we format this?
				// TODO: how do we support errors.Is ?
				//       you can fmt.Errorf("%v") to wrap the sentinel
				"expected nonce %d, received nonce %d",
				pendingNonce, txNonce.Uint64(),
			),
		)
	}

	txAmount := tx.GetAmount()
	txFee := tx.GetFee()
	txTotal := (&txAmount).Add(&txFee)
	if pendingBalance.Cmp(txTotal) < 0 {
		return errors.WithStack(ErrBalanceTooLow)
	}

	// (II) Update the pending state

	// TODO: also update the pending state of the fee receiver!

	newPendingBalance := pendingBalance.Sub(txTotal)

	one := models.MakeUint256(1)
	newPendingNonce := pendingNonce.Add(&one)

	err = s.UnsafeSetPendingState(fromStateID, *newPendingNonce, *newPendingBalance)
	if err != nil {
		return err
	}

	if tx.Type() == txtype.Transfer {
		toStateID := *tx.GetToStateID() // will not panic, transfers have this
		err = s.addToPendingBalance(toStateID, &txAmount)
		if err != nil {
			return err
		}
	}

	if tx.Type() == txtype.Create2Transfer {
		toPubKey := tx.ToCreate2Transfer().ToPublicKey
		err = s.addToPendingPubkeyBalance(&toPubKey, &txAmount)
		if err != nil {
			return err
		}
	}

	// (III) Add the received transaction to the relevant queue

	pendingTx := stored.NewPendingTx(tx)
	txKey := pendingTxKey(fromStateID, pendingNonce.Uint64())

	_, err = s.rawLookup(txKey)
	if err == nil {
		return errors.WithStack(fmt.Errorf("cannot replace transactions"))
	}
	if !errors.Is(err, badger.ErrKeyNotFound) {
		return err
	}

	return s.rawSet(txKey, pendingTx.Bytes())
}

func (s *Storage) UnsafeInsertPendingTxSkipValidation(pendingTx *stored.PendingTx) error {
	txKey := pendingTxKey(pendingTx.FromStateID, pendingTx.Nonce.Uint64())
	return s.rawSet(txKey, pendingTx.Bytes())
}

/// These are used by the rollup loop:
// TODO: break into a separate file?

type MempoolHeap struct {
	storage *Storage
	txType  txtype.TransactionType
	heap    *TxHeap

	// stateID -> nonce, the nonce of the last tx we added to the heap for this ID
	lastTx map[uint32]uint64

	// the list of txs which have been added to a batch and need to be deleted from
	// the mempool. When Savepoint() is called these are removed from badger. We can
	// not remove them immediately, because it's possible that the rollup loop will
	// fail to fill a commitment and try to put these back into the mempool. When the
	// rollup loop fills a commitment it calls Savepoint() when is when it's finally
	// safe to write to badger. Once Savepoint() has been called these are either
	// going into a batch or the entire batch will fail (in which case the
	// badger tx will be rolled back).
	toBeDeleted [][]byte
}

func (mh *MempoolHeap) nextTxForAccount(stateID uint32) (
	pendingTx *stored.PendingTx,
	err error,
) {
	// this function has to be very careful. The rollup loop (which calls this method)
	// is not allowed to badger.ErrConflict, which happens if we read from a key which
	// a concurrent transaction wrote to before we could commit.

	// We want to check whether the next txn exists, but we are absolutely not allowed
	// to do the following:
	//   nextKey := pendingTxKey(stateID, lastNonce + 1)
	//   _, doesNotExistError := txn.Get(nextKey)
	// This code would add `nextKey` to our read set. If the txn did not exist but
	// an API handler inserts it before we commit then the rollup loop will crash when
	// it tries to commit.

	// Instead, we make an iterator and ask it to scan past our desired key
	// Next() and Key() and ValidForPrefix are safe but Item() and Seek() are
	// dangerous, they add to the read set.

	err = mh.storage.database.Badger.RawUpdate(func(txn *badger.Txn) error {
		iter := txn.NewIterator(db.PrefetchIteratorOpts)
		defer iter.Close()

		lastNonce, wasInMap := mh.lastTx[stateID]
		if wasInMap {
			// SELECT * FROM pendingTx
			//  WHERE stateID = ??
			//    AND nonce > lastNonce
			// ORDER BY nonce ASC LIMIT 1;
			iteratorStartKey := pendingTxKey(stateID, lastNonce)

			// we are allowed to seek to this key because we have already read
			// from the txn at lastNonce, it is already in our heap!
			iter.Seek(iteratorStartKey)
			iter.Next() // this makes it > instead of >=
		} else {
			// this will happen if we have just sent money to an account and want to
			// check whether it now has a txn which can be sent. Here we just want to
			// look for the txn for this stateID with a lowest nonce:
			// SELECT * FROM pendingTx WHERE stateID = ?? ORDER BY nonce ASC LIMIT 1;
			iteratorStartKey := pendingTxStateIDPrefix(stateID)
			iter.Seek(iteratorStartKey)
		}

		if !iter.ValidForPrefix(pendingTxStateIDPrefix(stateID)) {
			// there is no next tx for the given account
			pendingTx = nil
			return nil
		}

		// now that we have confirmed that this tx exists we are allowed to fetch
		// it and add it to our read set. The API handlers will never update a
		// transaction, only insert new ones.

		item := iter.Item()
		innerPendingTx, innerErr := itemToPendingTx(item)
		if innerErr != nil {
			return innerErr
		}

		pendingTx = innerPendingTx
		return nil
	})

	return pendingTx, err
}

func itemToPendingState(item *badger.Item) (nonce, balance *models.Uint256, err error) {
	value, err := item.ValueCopy(nil)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	decodedNonce, decodedBalance := decodePendingState(value)
	return &decodedNonce, &decodedBalance, nil
}

func itemToPendingTx(item *badger.Item) (*stored.PendingTx, error) {
	value, err := item.ValueCopy(nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result stored.PendingTx
	err = result.SetBytes(value)
	if err != nil {
		// TODO: confirm we don't need to attach a stack here
		return nil, err
	}

	return &result, nil
}

func itemToPendingPubkeyBalance(item *badger.Item) (*models.Uint256, error) {
	value, err := item.ValueCopy(nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result models.Uint256
	result.SetBytes(value)

	return &result, nil
}

// TODO: I know the mempool is small but does it really make sense to scan the entire
//       thing every loop?
// TODO: add a test for this method, lint showed me a bug in it
func (s *Storage) FindOldestMempoolTransaction(txType txtype.TransactionType) (*stored.PendingTx, error) {
	pendingTxs, err := s.lowestNoncePendingTxs()
	if err != nil {
		return nil, err
	}

	var oldestTime *models.Timestamp
	var result stored.PendingTx

	foundPendingTx := false

	for i := range pendingTxs {
		tx := pendingTxs[i]

		txTime := tx.ReceiveTime
		if txTime == nil {
			continue
		}
		// TODO: create a helper for this filter, we do it in two different places
		isExecutable, err := s.txIsExecutable(txType, &tx)
		if err != nil {
			return nil, err
		}
		if !isExecutable {
			continue
		}
		if (oldestTime == nil) || txTime.Before(*oldestTime) {
			oldestTime = txTime
			result = tx
			foundPendingTx = true
		}
	}

	if foundPendingTx {
		return &result, nil
	} else {
		return nil, nil
	}
}

// fetches, for each account with a pending tx, the pending tx with the lowest nonce
func (s *Storage) lowestNoncePendingTxs() ([]stored.PendingTx, error) {
	result := make([]stored.PendingTx, 0)

	err := s.database.Badger.RawUpdate(func(txn *badger.Txn) error {
		iter := txn.NewIterator(db.PrefetchIteratorOpts)
		defer iter.Close()

		iter.Seek(pendingTxPrefix)

		for iter.ValidForPrefix(pendingTxPrefix) {
			item := iter.Item()
			pendingTx, innerErr := itemToPendingTx(item)
			if innerErr != nil {
				return innerErr
			}

			result = append(result, *pendingTx)

			// skip over the rest of the txns from this sender
			nextStateID := pendingTx.FromStateID + 1
			nextPrefix := pendingTxStateIDPrefix(nextStateID)
			iter.Seek(nextPrefix)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *Storage) forEachMempoolTransaction(fun func(*stored.PendingTx) error) error {
	return s.database.Badger.View(func(txn *badger.Txn) error {
		iter := txn.NewIterator(db.PrefetchIteratorOpts)
		defer iter.Close()

		iter.Seek(pendingTxPrefix)

		for iter.ValidForPrefix(pendingTxPrefix) {
			item := iter.Item()
			pendingTx, innerErr := itemToPendingTx(item)
			if innerErr != nil {
				return innerErr
			}

			err := fun(pendingTx)
			if err != nil {
				return err
			}

			iter.Next()
		}

		return nil
	})
}

// This assumes that the mempool is relatively small, but we don't do anything to
// guarantee that the mempool remains small.
// TODO: Add some metrics which will warn us if this method (or its callers) start taking
//       too long, and we can add some restrictions on how many pending trasnactions each
//       account is allowed to have... or add an index
func (s *Storage) GetMempoolTransactionByHash(hash common.Hash) (*stored.PendingTx, error) {
	allPendingTxs, err := s.GetAllMempoolTransactions()
	if err != nil {
		return nil, err
	}

	for i := range allPendingTxs {
		if allPendingTxs[i].Hash == hash {
			return &allPendingTxs[i], nil
		}
	}

	return nil, nil
}

func (s *Storage) GetPendingStates(startStateID, pageSize uint32) (
	result []dto.UserStateWithID,
	err error,
) {
	err = s.ExecuteInReadWriteTransaction(func(txStorage *Storage) error {
		innerResult, innerErr := txStorage.unsafeGetPendingStates(startStateID, pageSize)
		result = innerResult
		return innerErr
	})
	return result, err
}

func (s *Storage) unsafeGetPendingStates(startStateID, pageSize uint32) (
	result []dto.UserStateWithID,
	err error,
) {
	result = make([]dto.UserStateWithID, 0)

	err = s.database.Badger.View(func(txn *badger.Txn) error {
		iter := txn.NewIterator(db.PrefetchIteratorOpts)
		defer iter.Close()

		startKey := pendingStateKey(startStateID)
		iter.Seek(startKey)

		for iter.ValidForPrefix(pendingStatePrefix) {
			item := iter.Item()
			pendingNonce, pendingBalance, innerErr := itemToPendingState(item)
			if innerErr != nil {
				return innerErr
			}

			keyBytes := item.Key()
			stateID := decodePendingStateKey(keyBytes)

			batchedState, innerErr := s.StateTree.Leaf(stateID)
			if innerErr != nil {
				return innerErr
			}

			result = append(result, dto.MakeUserStateWithID(
				stateID,
				&models.UserState{
					Nonce:    *pendingNonce,
					Balance:  *pendingBalance,
					PubKeyID: batchedState.PubKeyID,
					TokenID:  batchedState.TokenID,
				},
			))

			if len(result) >= int(pageSize) {
				break
			}

			iter.Next()
		}

		return nil
	})
	return result, err
}

// pageSize==0 means no limit
func (s *Storage) GetPendingPubkeyBalances(startPrefix []byte, pageSize uint32) ([]dto.PubkeyBalance, error) {
	result := make([]dto.PubkeyBalance, 0)

	err := s.database.Badger.View(func(txn *badger.Txn) error {
		iter := txn.NewIterator(db.PrefetchIteratorOpts)
		defer iter.Close()

		seekPrefix := bytes.Join(
			[][]byte{pendingPubkeyBalancePrefix, startPrefix},
			[]byte(":"),
		)
		iter.Seek(seekPrefix)

		for iter.ValidForPrefix(pendingPubkeyBalancePrefix) {
			item := iter.Item()
			pendingBalance, innerErr := itemToPendingPubkeyBalance(item)
			if innerErr != nil {
				return innerErr
			}

			keyBytes := item.Key()
			pubkey := decodePendingPubkeyBalanceKey(keyBytes)

			result = append(result, dto.PubkeyBalance{
				PubKey:  *pubkey,
				Balance: *pendingBalance,
			})

			if pageSize != 0 && len(result) >= int(pageSize) {
				break
			}

			iter.Next()
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *Storage) RecomputePendingState(stateID uint32, mutate bool) (
	result *dto.RecomputePendingState,
	err error,
) {
	err = s.ExecuteInReadWriteTransaction(func(txStorage *Storage) error {
		innerResult, innerErr := txStorage.unsafeRecomputePendingState(stateID, mutate)
		result = innerResult
		return innerErr
	})
	return result, err
}

func (s *Storage) unsafeRecomputePendingState(stateID uint32, mutate bool) (
	*dto.RecomputePendingState,
	error,
) {
	oldNonce, oldBalance, err := s.getPendingState(stateID)
	if err != nil {
		return nil, err
	}

	batchedState, err := s.StateTree.Leaf(stateID)
	if err != nil {
		return nil, err
	}

	newNonce := &batchedState.UserState.Nonce
	newBalance := &batchedState.UserState.Balance

	receivedTxCount, sentTxCount := 0, 0

	// forEachMempoolTransaction was written to avoid conflicting with any other
	// writers. In this case we would like to conflict with other writers! That is achieved
	// later in the method when we call `UnsafeSetPendingState`. And other api coros which add
	// relevant txns to the mempool will also update our pending state, causing a retry.
	err = s.forEachMempoolTransaction(func(pendingTx *stored.PendingTx) error {
		if pendingTx.FromStateID == stateID {
			one := models.MakeUint256(1)
			newNonce = newNonce.Add(&one)

			txTotal := pendingTx.Amount.Add(&pendingTx.Fee)
			newBalance = newBalance.Sub(txTotal)

			sentTxCount += 1
		}

		if pendingTx.TxType != txtype.Transfer {
			return nil
		}

		transfer := pendingTx.ToTransfer()
		if transfer == nil {
			panic("we just confirmed this is a transfer")
		}

		if transfer.ToStateID == stateID {
			newBalance = newBalance.Add(&transfer.Amount)
			receivedTxCount += 1
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	logFields := log.Fields{
		"stateID":         stateID,
		"oldNonce":        oldNonce,
		"oldBalance":      oldBalance,
		"newNonce":        newNonce,
		"newBalance":      newBalance,
		"sentTxCount":     sentTxCount,
		"receivedTxCount": receivedTxCount,
	}

	if mutate {
		err = s.UnsafeSetPendingState(stateID, *newNonce, *newBalance)
		if err != nil {
			return nil, err
		}
		log.WithFields(logFields).Info("computed and updated pending state")
	} else {
		log.WithFields(logFields).Info("computed new pending state")
	}

	return &dto.RecomputePendingState{
		OldNonce:   *oldNonce,
		OldBalance: *oldBalance,
		NewNonce:   *newNonce,
		NewBalance: *newBalance,
	}, nil
}

func (s *Storage) RecomputePendingPubkeyBalances(startPrefix []byte, pageSize uint32) (
	result []dto.PubkeyBalance,
	err error,
) {
	err = s.ExecuteInReadWriteTransaction(func(txStorage *Storage) error {
		innerResult, innerErr := txStorage.unsafeRecomputePendingPubkeyBalances(startPrefix, pageSize)
		result = innerResult
		return innerErr
	})
	return result, err
}

// SELECT
//   pubkey,
//   SUM(amount) AS balance
// FROM mempool_transactions
// WHERE
//   tx_type = 'CREATE2TRANSFER' AND
//   pubkey >= startPrefix
// GROUP BY pubkey
// ORDER BY pubkey
// LIMIT pageSize
func (s *Storage) unsafeRecomputePendingPubkeyBalances(startPrefix []byte, pageSize uint32) ([]dto.PubkeyBalance, error) {
	keyBalances, err := s.pendingPubkeyBalances()
	if err != nil {
		return nil, err
	}

	keys := make([][]byte, 0, len(keyBalances))
	for key := range keyBalances {
		keys = append(keys, key.Bytes())
	}

	sort.Slice(keys, func(i, j int) bool {
		return bytes.Compare(keys[i], keys[j]) < 0
	})

	var result []dto.PubkeyBalance
	if pageSize == 0 {
		result = make([]dto.PubkeyBalance, 0, len(keys))
	} else {
		result = make([]dto.PubkeyBalance, 0, pageSize)
	}

	for _, key := range keys {
		if bytes.Compare(key, startPrefix) == -1 { // if key < startPrefix
			continue
		}

		if pageSize > 0 && uint32(len(result)) >= pageSize {
			break
		}

		var pubkey models.PublicKey
		err = pubkey.SetBytes(key)
		if err != nil {
			// TODO: remove this wrap, SetBytes should attach the stacktrace
			return nil, errors.WithStack(err)
		}

		balance := keyBalances[pubkey]

		result = append(result, dto.PubkeyBalance{
			PubKey:  pubkey,
			Balance: *balance,
		})
	}

	return result, nil
}

func (s *Storage) GetAllMempoolTransactions() ([]stored.PendingTx, error) {
	result := make([]stored.PendingTx, 0)

	err := s.forEachMempoolTransaction(func(pendingTx *stored.PendingTx) error {
		result = append(result, *pendingTx)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *Storage) CountPendingTxsOfType(txType txtype.TransactionType) (uint32, error) {
	result := uint32(0)

	err := s.forEachMempoolTransaction(func(pendingTx *stored.PendingTx) error {
		if pendingTx.TxType == txType {
			result += 1
		}
		return nil
	})
	if err != nil {
		return 0, err
	}

	return result, nil
}

// TODO: should this accept a pointer to a PendingTx?
func (s *Storage) txIsExecutable(txType txtype.TransactionType, tx *stored.PendingTx) (bool, error) {
	if tx == nil {
		panic("txIsExecutable must be given a tx")
	}

	if tx.TxType != txType {
		return false, nil
	}

	currentState, err := s.StateTree.Leaf(tx.FromStateID)
	if err != nil {
		return false, err
	}

	currentNonce := &currentState.UserState.Nonce
	currentBalance := &currentState.UserState.Balance

	if currentNonce.Cmp(&tx.Nonce) != 0 {
		// This method is only ever called on the txn with the lowest nonce that
		// we have for a given stateID. If the nonce is lower than the current
		// nonce then something terribly wrong has happened, the API should not
		// have accepted a duplicate transaction. If the nonce is higher than the
		// current nonce then, again, something terribly wrong has happened, the
		// API should never introduce a nonce gap when accepting transactions.

		// I can imagine strategies for fixing this but I cannot imagine how we
		// might have gotten into this situation so I can't guess at the best way
		// to fix it. Hopefully this error message is scary enough that it will
		// trigger an investigation and a manual patch-up of the mempool.

		// TODO: test that this error message is emitted when appropriate
		log.Errorf("invalid state, cannot processes transactions. stateID=%d", tx.FromStateID)
		return false, nil
	}

	txTotal := tx.Amount.Add(&tx.Fee)
	if currentBalance.Cmp(txTotal) < 0 {
		// If you don't have enough money to pay for your next transaction then
		// we hold onto the transaction until you do. The API would only have
		// accepted this transaction if your pending balance was high enough to
		// pay, so there is guaranteed to be an inbound transfer which will pay
		// for this, eventually.
		return false, nil
	}

	return true, nil
}

// caution: assumes you are running inside a tx
// caution: there can be only one
func (s *Storage) NewMempoolHeap(txType txtype.TransactionType) (*MempoolHeap, error) {
	lastTx := make(map[uint32]uint64)
	mh := &MempoolHeap{
		storage:     s,
		txType:      txType,
		heap:        NewTxHeap(),
		lastTx:      lastTx,
		toBeDeleted: make([][]byte, 0),
	}

	pendingTxs, err := s.lowestNoncePendingTxs()
	if err != nil {
		return nil, err
	}

	for i := range pendingTxs {
		tx := pendingTxs[i]

		isExecutable, err := s.txIsExecutable(txType, &tx)
		if err != nil {
			return nil, err
		}
		if !isExecutable {
			continue
		}

		mh.pushTx(&tx)
	}

	return mh, nil
}

func (mh *MempoolHeap) PeekHighestFeeExecutableTx() models.GenericTransaction {
	fromHeap := mh.heap.Peek()
	if fromHeap != nil {
		return fromHeap.ToGenericTransaction()
	}
	return nil
}

// TODO: use this in our constructor
func (mh *MempoolHeap) pushTx(tx *stored.PendingTx) {
	if tx == nil {
		panic("pushTx must be given a tx")
	}

	if tx.TxType != mh.txType {
		panic("unexpected txType")
	}

	mh.heap.Push(*tx)
	mh.lastTx[tx.FromStateID] = tx.Nonce.Uint64()
}

// Caution: assumes the pendingTx which we are about to drop has been applied to the state
//          in mh.storage
//nolint:gocyclo  // TODO: consider doing what it says
func (mh *MempoolHeap) DropHighestFeeExecutableTx() error {
	pendingTx := mh.heap.Pop()
	if pendingTx == nil {
		// you should have noticed the nil when you called Peek...()
		panic("unreachable")
	}

	mh.scheduleDeletion(pendingTx)

	// Executing this tx has made up to two other txns executable and those should be
	// added to our heap.

	// (I) we are now free to add the next tx from this account to the heap

	nextTxForID, err := mh.nextTxForAccount(pendingTx.FromStateID)
	if err != nil {
		return err
	}

	// TODO: add a test that we don't insert txs with the wrong type here
	if nextTxForID != nil && nextTxForID.TxType == mh.txType {
		mh.pushTx(nextTxForID)
	}

	// (II) if this was a transfer then we might have given funds to an account which
	//      was blocked for lack of funds

	if pendingTx.TxType != txtype.Transfer {
		return nil
	}

	transfer := pendingTx.ToTransfer()
	if transfer == nil {
		panic("we just confirmed this is a transfer")
	}

	toStateID := transfer.ToStateID
	nextTx, err := mh.nextTxForAccount(toStateID)
	if err != nil {
		return err
	}

	if nextTx == nil {
		// easy, this account was not blocked because it has no txn to be blocked
		return nil
	}

	_, wasInHeap := mh.lastTx[nextTx.FromStateID]
	if !wasInHeap {
		// this account has a pendingTx but that tx was never added to the heap,
		// it must have been blocked by something

		isExecutable, innerErr := mh.storage.txIsExecutable(mh.txType, nextTx)
		if innerErr != nil {
			return innerErr
		}

		if isExecutable {
			// good chance it was blocked on our transfer, let's throw it in
			mh.pushTx(nextTx)
		} else { //nolint:staticcheck // lint does not like empty branches, I do
			// it is still not executable, probably we didn't give it enough
			// additional balance, it will have to wait a little longer
		}

		return nil
	}

	// this last case is a little tricky:
	// - we just transferred money to an account which has another transaction
	//   it wants to send
	// - we also know that one of its txns was previously added to the heap
	// - we need to figure out whether that txn is still in the heap

	currentState, err := mh.storage.StateTree.Leaf(nextTx.FromStateID)
	if err != nil {
		return err
	}
	currentNonce := &currentState.UserState.Nonce

	if currentNonce.Cmp(&nextTx.Nonce) != 0 {
		// this account already has a tx which is added to the heap, we don't
		// want to add another one. If we did then Transfer(nonce=X) and
		// Transfer(nonce=X+1) would both be on the heap, and if Transfer(X+1)
		// paid a higher fee then we would try to apply it first and error out.
		return nil
	}

	isExecutable, err := mh.storage.txIsExecutable(mh.txType, nextTx)
	if err != nil {
		return err
	}

	if isExecutable {
		mh.pushTx(nextTx)
		return nil
	}

	return nil
}

func (mh *MempoolHeap) scheduleDeletion(pendingTx *stored.PendingTx) {
	key := pendingTxKey(pendingTx.FromStateID, pendingTx.Nonce.Uint64())
	mh.toBeDeleted = append(mh.toBeDeleted, key)
}

// These methods exist because badger does not support subtransactions

// Write all our changes back to the Storage transaction.
func (mh *MempoolHeap) Savepoint() error {
	err := mh.storage.database.Badger.RawUpdate(func(txn *badger.Txn) error {
		for _, key := range mh.toBeDeleted {
			innerErr := txn.Delete(key)
			if innerErr != nil {
				return innerErr
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	mh.toBeDeleted = make([][]byte, 0)
	return nil
}

/// TODO: where to put this?

type TxHeap struct {
	heap *utils.MutableHeap
}

func NewTxHeap(txs ...stored.PendingTx) *TxHeap {
	less := func(a, b interface{}) bool {
		txA := a.(stored.PendingTx)
		txB := b.(stored.PendingTx)
		return txA.Fee.Cmp(&txB.Fee) > 0
	}

	elements := make([]interface{}, len(txs))
	for i := range txs {
		elements[i] = txs[i]
	}

	return &TxHeap{
		heap: utils.NewMutableHeap(elements, less),
	}
}

func (h *TxHeap) Peek() *stored.PendingTx {
	return h.toPendingTx(h.heap.Peek())
}

//nolint:gocritic // TODO: slightly improve performance by doing what it wants
func (h *TxHeap) Push(tx stored.PendingTx) {
	h.heap.Push(tx)
}

func (h *TxHeap) Pop() *stored.PendingTx {
	return h.toPendingTx(h.heap.Pop())
}

func (h *TxHeap) Size() int {
	return h.heap.Size()
}

func (h *TxHeap) toPendingTx(element interface{}) *stored.PendingTx {
	if tx, ok := element.(stored.PendingTx); ok {
		return &tx
	}
	return nil
}
