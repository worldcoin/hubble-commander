package storage

import (
	"github.com/Worldcoin/hubble-commander/models"
	bdg "github.com/dgraph-io/badger/v3"
	"github.com/ethereum/go-ethereum/common"
	bh "github.com/timshannon/badgerhold/v3"
)

func (s *Storage) UpsertStateNode(node *models.StateNode) error {
	return s.Badger.Upsert(node.MerklePath, node)
}

func (s *Storage) BatchUpsertStateNodes(nodes []models.StateNode) (err error) {
	tx, storage, err := s.BeginTransaction(TxOptions{Postgres: true, Badger: true})
	if err != nil {
		return err
	}
	defer tx.Rollback(&err)
	for i := range nodes {
		err = storage.UpsertStateNode(&nodes[i])
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (s *Storage) AddStateNode(node *models.StateNode) error {
	return s.Badger.Insert(node.MerklePath, node)
}

func (s *Storage) GetStateNodeByPath(path *models.MerklePath) (*models.StateNode, error) {
	var node models.StateNode
	err := s.Badger.Get(path, &node)
	if err == bh.ErrNotFound {
		return newZeroStateNode(path), nil
	}
	if err != nil {
		return nil, err
	}
	return &node, nil
}

func (s *Storage) GetStateNodeByHash(dataHash *common.Hash) (*models.StateNode, error) {
	nodes := make([]models.StateNode, 0, 1)
	err := s.Badger.Find(&nodes, bh.Where("DataHash").Eq(dataHash).Index("DataHash"))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, NewNotFoundError("state node")
	}
	return &nodes[0], nil
}

func newZeroStateNode(path *models.MerklePath) *models.StateNode {
	return &models.StateNode{
		MerklePath: *path,
		DataHash:   GetZeroHash(leafDepth - uint(path.Depth)),
	}
}

func (s *Storage) GetStateNodes(paths []models.MerklePath) (nodes []models.StateNode, err error) {
	tx, storage, err := s.BeginTransaction(TxOptions{Badger: true, ReadOnly: true})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(&err)

	nodes = make([]models.StateNode, 0)
	for i := range paths {
		var node *models.StateNode
		node, err = storage.GetStateNodeByPath(&paths[i])
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, *node)
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return nodes, nil
}

func (s *Storage) GetNextAvailableStateID() (*uint32, error) {
	var nextAvailableStateID uint32

	err := s.Badger.BadgerInstance().View(func(txn *bdg.Txn) error {
		opts := bdg.DefaultIteratorOptions
		opts.PrefetchValues = false
		opts.Reverse = true
		it := txn.NewIterator(opts)
		defer it.Close()
		prefix := []byte("bh_FlatStateLeaf")

		seekPrefix := make([]byte, 0, len(prefix)+1)
		seekPrefix = append(seekPrefix, prefix...)
		seekPrefix = append(seekPrefix, 0xFF) // Required to loop backwards

		it.Seek(seekPrefix)
		if it.ValidForPrefix(prefix) {
			lastItem := it.Item()
			lastItemKey := lastItem.Key()
			lastStateID := lastItemKey[len(lastItemKey)-1]

			nextAvailableStateID = uint32(lastStateID) + 1
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &nextAvailableStateID, nil
}
