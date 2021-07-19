package storage

import (
	"github.com/Worldcoin/hubble-commander/db/badger"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils"
	bdg "github.com/dgraph-io/badger/v3"
	bh "github.com/timshannon/badgerhold/v3"
)

func (s *Storage) UpsertStateLeaf(leaf *models.StateLeaf) error {
	flatLeaf := models.MakeFlatStateLeaf(leaf)
	return s.Badger.Upsert(leaf.StateID, flatLeaf)
}

func (s *Storage) GetStateLeaf(stateID uint32) (stateLeaf *models.StateLeaf, err error) {
	var leaf models.FlatStateLeaf
	err = s.Badger.Get(stateID, &leaf)
	if err == bh.ErrNotFound {
		return nil, NewNotFoundError("state leaf")
	}
	if err != nil {
		return nil, err
	}
	return leaf.StateLeaf(), nil
}

func (s *Storage) GetUserStatesByPublicKey(publicKey *models.PublicKey) (userStates []models.UserStateWithID, err error) {
	accounts, err := s.GetAccountLeaves(publicKey)
	if err != nil {
		return nil, err
	}

	pubKeyIDs := utils.ValueToInterfaceSlice(accounts, "PubKeyID")

	leaves := make([]models.FlatStateLeaf, 0, 1)
	err = s.Badger.Find(
		&leaves,
		bh.Where("PubKeyID").In(pubKeyIDs...).Index("PubKeyID"),
	)
	if err != nil {
		return nil, err
	}
	numLeaves := len(leaves)
	if numLeaves == 0 {
		return nil, NewNotFoundError("user states")
	}

	userStates = make([]models.UserStateWithID, 0, numLeaves)
	for i := range leaves {
		leaf := &leaves[i]
		userStates = append(userStates, models.UserStateWithID{
			StateID: leaf.StateID,
			UserState: models.UserState{
				PubKeyID: leaf.PubKeyID,
				TokenID:  leaf.TokenID,
				Balance:  leaf.Balance,
				Nonce:    leaf.Nonce,
			},
		})
	}

	return userStates, nil
}

func (s *Storage) GetFeeReceiverStateLeaf(pubKeyID uint32, tokenID models.Uint256) (*models.StateLeaf, error) {
	stateID, ok := s.feeReceiverStateIDs[tokenID.String()]
	if ok {
		return s.GetStateLeaf(stateID)
	}
	stateLeaf, err := s.GetStateLeafByPubKeyIDAndTokenID(pubKeyID, tokenID)
	if err != nil {
		return nil, err
	}
	s.feeReceiverStateIDs[stateLeaf.TokenID.String()] = stateLeaf.StateID
	return stateLeaf, nil
}

func (s *Storage) GetStateLeafByPubKeyIDAndTokenID(pubKeyID uint32, tokenID models.Uint256) (*models.StateLeaf, error) {
	leaves := make([]models.FlatStateLeaf, 0, 1)
	err := s.Badger.Find(
		&leaves,
		bh.Where("TokenID").Eq(tokenID).
			And("PubKeyID").Eq(pubKeyID).Index("PubKeyID"),
	)
	if err != nil {
		return nil, err
	}
	if len(leaves) == 0 {
		return nil, NewNotFoundError("state leaf")
	}
	return leaves[0].StateLeaf(), nil
}

func (s *Storage) GetNextAvailableStateID() (*uint32, error) {
	nextAvailableStateID := uint32(0)

	err := s.Badger.View(func(txn *bdg.Txn) error {
		opts := bdg.DefaultIteratorOptions
		opts.PrefetchValues = false
		opts.Reverse = true
		it := txn.NewIterator(opts)
		defer it.Close()

		seekPrefix := make([]byte, 0, len(flatStateLeafPrefix)+1)
		seekPrefix = append(seekPrefix, flatStateLeafPrefix...)
		seekPrefix = append(seekPrefix, 0xFF) // Required to loop backwards

		it.Seek(seekPrefix)
		if it.ValidForPrefix(flatStateLeafPrefix) {
			var key uint32
			err := badger.DecodeKey(it.Item().Key(), &key, flatStateLeafPrefix)
			if err != nil {
				return err
			}
			nextAvailableStateID = key + 1
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &nextAvailableStateID, nil
}
