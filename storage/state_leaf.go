package storage

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils"
	bh "github.com/timshannon/badgerhold/v3"
)

func (s *StateTree) upsertStateLeaf(leaf *models.StateLeaf) error {
	flatLeaf := models.MakeFlatStateLeaf(leaf)
	return s.storageBase.Database.Badger.Upsert(leaf.StateID, flatLeaf)
}

func (s *Storage) GetUserStatesByPublicKey(publicKey *models.PublicKey) (userStates []models.UserStateWithID, err error) {
	accounts, err := s.AccountTree.Leaves(publicKey)
	if err != nil {
		return nil, err
	}

	pubKeyIDs := utils.ValueToInterfaceSlice(accounts, "PubKeyID")

	leaves := make([]models.FlatStateLeaf, 0, 1)
	err = s.Database.Badger.Find(
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
		return s.StateTree.Leaf(stateID)
	}
	stateLeaf, err := s.StateTree.getLeafByPubKeyIDAndTokenID(pubKeyID, tokenID)
	if err != nil {
		return nil, err
	}
	s.feeReceiverStateIDs[stateLeaf.TokenID.String()] = stateLeaf.StateID
	return stateLeaf, nil
}
