package storage

import (
	"github.com/Worldcoin/hubble-commander/db/badger"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	bh "github.com/timshannon/badgerhold/v3"
)

func (s *Storage) UpsertStateLeaf(leaf *models.StateLeaf) error {
	flatLeaf := models.NewFlatStateLeaf(leaf)
	err := s.Badger.Upsert(leaf.StateID, &flatLeaf)
	if err == bh.ErrKeyExists {
		return nil
	}
	return err
}

func (s *Storage) GetStateLeafByStateID(stateID uint32) (stateLeaf *models.StateLeaf, err error) {
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

// TODO move to state_node, make sure to only iterate over keys (Badger PrefetchValues=false)
func (s *Storage) GetNextAvailableStateID() (*uint32, error) {
	nodes := make([]models.StateNode, 0, 1)
	err := s.Badger.Find(
		&nodes,
		bh.Where("MerklePath").
			MatchFunc(badger.MatchAll). // TODO possibly performance killer
			SortBy("MerklePath.Path").
			Reverse().
			Limit(1),
	)
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return ref.Uint32(0), nil
	}
	stateID := nodes[0].MerklePath.Path + 1
	return &stateID, nil
}

func (s *Storage) GetUserStatesByPublicKey(publicKey *models.PublicKey) (userStates []models.UserStateWithID, err error) {
	accounts, err := s.GetAccounts(publicKey)
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
	if len(leaves) == 0 {
		return nil, NewNotFoundError("user states")
	}

	userStates = make([]models.UserStateWithID, 0, 1)
	for i := range leaves {
		userStates = append(userStates, models.UserStateWithID{
			StateID: leaves[i].StateID,
			UserState: models.UserState{
				PubKeyID:   leaves[i].PubKeyID,
				TokenIndex: leaves[i].TokenIndex,
				Balance:    leaves[i].Balance,
				Nonce:      leaves[i].Nonce,
			},
		})
	}

	return userStates, nil
}

// TODO this method is only used in ApplyFee to get the fee receiver UserState, this can be easily cached
func (s *Storage) GetUserStateByPubKeyIDAndTokenIndex(pubKeyID uint32, tokenIndex models.Uint256) (*models.UserStateWithID, error) {
	tx, storage, err := s.BeginTransaction(TxOptions{Postgres: true, Badger: true, ReadOnly: true})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(&err)

	nodes := make([]models.StateNode, 0)
	err = storage.Badger.Find(&nodes, nil)
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, NewNotFoundError("user state")
	}

	dataHashes := utils.ValueToInterfaceSlice(nodes, "DataHash")

	leaves := make([]models.FlatStateLeaf, 0, 1)
	err = storage.Badger.Find(
		&leaves,
		bh.Where("DataHash").In(dataHashes...).
			And("TokenIndex").Eq(tokenIndex).Index("TokenIndex").
			And("PubKeyID").Eq(pubKeyID).Index("PubKeyID"),
	)
	if err != nil {
		return nil, err
	}
	if len(leaves) == 0 {
		return nil, NewNotFoundError("user state")
	}

	stateNode, err := s.GetStateNodeByHash(&leaves[0].DataHash)
	if err != nil {
		return nil, err
	}
	userState := &models.UserStateWithID{
		StateID: stateNode.MerklePath.Path,
		UserState: models.UserState{
			PubKeyID:   leaves[0].PubKeyID,
			TokenIndex: leaves[0].TokenIndex,
			Balance:    leaves[0].Balance,
			Nonce:      leaves[0].Nonce,
		},
	}

	return userState, nil
}

func (s *Storage) GetUserStateByID(stateID uint32) (*models.UserStateWithID, error) {
	var leaf models.FlatStateLeaf
	err := s.Badger.Get(stateID, &leaf)
	if err == bh.ErrNotFound {
		return nil, NewNotFoundError("user state")
	}
	if err != nil {
		return nil, err
	}

	userState := &models.UserStateWithID{
		StateID: stateID,
		UserState: models.UserState{
			PubKeyID:   leaf.PubKeyID,
			TokenIndex: leaf.TokenIndex,
			Balance:    leaf.Balance,
			Nonce:      leaf.Nonce,
		},
	}

	return userState, nil
}
