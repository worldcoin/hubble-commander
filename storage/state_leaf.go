package storage

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils"
	bh "github.com/timshannon/badgerhold/v3"
)

func (s *Storage) UpsertStateLeaf(leaf *models.StateLeaf) error {
	flatLeaf := models.NewFlatStateLeaf(leaf)
	return s.Badger.Upsert(leaf.StateID, &flatLeaf)
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
				PubKeyID:   leaf.PubKeyID,
				TokenIndex: leaf.TokenIndex,
				Balance:    leaf.Balance,
				Nonce:      leaf.Nonce,
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
