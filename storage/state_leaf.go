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

func (s *Storage) GetFeeReceiverStateLeaf(pubKeyID uint32, tokenIndex models.Uint256) (*models.StateLeaf, error) {
	stateID, ok := s.feeReceiverStateIDs[tokenIndex.String()]
	if ok {
		return s.GetStateLeaf(stateID)
	}
	stateLeaf, err := s.GetStateLeafByPubKeyIDAndTokenIndex(pubKeyID, tokenIndex)
	if err != nil {
		return nil, err
	}
	s.feeReceiverStateIDs[stateLeaf.TokenIndex.String()] = stateLeaf.StateID
	return stateLeaf, nil
}

func (s *Storage) GetStateLeafByPubKeyIDAndTokenIndex(pubKeyID uint32, tokenIndex models.Uint256) (*models.StateLeaf, error) {
	leaves := make([]models.FlatStateLeaf, 0, 1)
	err := s.Badger.Find(
		&leaves,
		bh.Where("TokenIndex").Eq(tokenIndex).Index("TokenIndex").
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
