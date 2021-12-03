package storage

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/stored"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/pkg/errors"
	bh "github.com/timshannon/badgerhold/v4"
)

func (s *StateTree) upsertStateLeaf(leaf *models.StateLeaf) error {
	flatLeaf := stored.MakeStateLeaf(leaf)
	return s.database.Badger.Upsert(leaf.StateID, flatLeaf)
}

func (s *Storage) GetStateLeavesByPublicKey(publicKey *models.PublicKey) (stateLeaves []models.StateLeaf, err error) {
	accounts, err := s.AccountTree.Leaves(publicKey)
	if err != nil {
		return nil, err
	}

	pubKeyIDs := utils.ValueToInterfaceSlice(accounts, "PubKeyID")

	flatStateLeaves := make([]stored.StateLeaf, 0, 1)
	err = s.database.Badger.Find(
		&flatStateLeaves,
		bh.Where("PubKeyID").In(pubKeyIDs...).Index("PubKeyID").SortBy("StateID"),
	)
	if err != nil {
		return nil, err
	}
	if len(flatStateLeaves) == 0 {
		return nil, errors.WithStack(NewNotFoundError("user states"))
	}

	stateLeaves = make([]models.StateLeaf, 0, len(flatStateLeaves))
	for i := range flatStateLeaves {
		stateLeaves = append(stateLeaves, *flatStateLeaves[i].StateLeaf())
	}

	return stateLeaves, nil
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
