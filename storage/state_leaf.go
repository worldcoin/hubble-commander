package storage

import (
	"github.com/Worldcoin/hubble-commander/db/badger"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/common"
	bh "github.com/timshannon/badgerhold/v3"
)

func (s *Storage) AddStateLeaf(leaf *models.StateLeaf) error {
	_, err := s.Postgres.Query(
		s.QB.Insert("state_leaf").
			Values(
				leaf.DataHash,
				leaf.PubKeyID,
				leaf.TokenIndex,
				leaf.Balance,
				leaf.Nonce,
			).
			Suffix("ON CONFLICT DO NOTHING"),
	).Exec()
	if err != nil {
		return err
	}

	flatLeaf := models.NewFlatStateLeaf(leaf)
	err = s.Badger.Insert(leaf.DataHash, &flatLeaf)
	if err == bh.ErrKeyExists {
		return nil
	}
	return err
}

func (s *Storage) GetStateLeafByHash(hash common.Hash) (*models.StateLeaf, error) {
	var leaf models.FlatStateLeaf
	err := s.Badger.Get(hash, &leaf)
	if err == bh.ErrNotFound {
		return nil, NewNotFoundError("state leaf")
	}
	if err != nil {
		return nil, err
	}
	return leaf.StateLeaf(), nil
}

func (s *Storage) GetStateLeafByPath(path *models.MerklePath) (*models.StateLeaf, error) {
	var node models.StateNode
	err := s.Badger.Get(path, &node)
	if err == bh.ErrNotFound {
		return nil, NewNotFoundError("state leaf")
	}
	if err != nil {
		return nil, err
	}
	return s.GetStateLeafByHash(node.DataHash)
}

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

func (s *Storage) GetUserStatesByPublicKey(publicKey *models.PublicKey) ([]models.UserStateWithID, error) {
	accounts, err := s.GetAccounts(publicKey)
	if err != nil {
		return nil, err
	}

	pubKeyIDs := utils.ValueToInterfaceSlice(accounts, "PubKeyID")

	nodes := make([]models.StateNode, 0)
	err = s.Badger.Find(&nodes, nil) // TODO possibly performance killer
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, NewNotFoundError("user states")
	}

	dataHashes := utils.ValueToInterfaceSlice(nodes, "DataHash")

	leaves := make([]models.FlatStateLeaf, 0, 1)
	err = s.Badger.Find(
		&leaves,
		bh.Where("DataHash").In(dataHashes...). // TODO possibly performance killer
							And("PubKeyID").In(pubKeyIDs...).Index("PubKeyID"),
	)
	if err != nil {
		return nil, err
	}
	if len(leaves) == 0 {
		return nil, NewNotFoundError("user states")
	}

	userStates := make([]models.UserStateWithID, 0, 1)
	for i := range leaves {
		stateNode, err := s.GetStateNodeByHash(&leaves[i].DataHash)
		if err != nil {
			return nil, err
		}
		userStates = append(userStates, models.UserStateWithID{
			StateID: stateNode.MerklePath.Path,
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

func (s *Storage) GetUserStateByPubKeyIDAndTokenIndex(pubKeyID uint32, tokenIndex models.Uint256) (*models.UserStateWithID, error) {
	nodes := make([]models.StateNode, 0)
	err := s.Badger.Find(&nodes, nil)
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, NewNotFoundError("user state")
	}

	dataHashes := utils.ValueToInterfaceSlice(nodes, "DataHash")

	leaves := make([]models.FlatStateLeaf, 0, 1)
	err = s.Badger.Find(
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
	path := models.MerklePath{
		Path:  stateID,
		Depth: leafDepth,
	}

	var node models.StateNode
	err := s.Badger.Get(path, &node)
	if err == bh.ErrNotFound {
		return nil, NewNotFoundError("user state")
	}
	if err != nil {
		return nil, err
	}

	var leaf models.FlatStateLeaf
	err = s.Badger.Get(node.DataHash, &leaf)
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
