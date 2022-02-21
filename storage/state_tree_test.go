package storage

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils/merkletree"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var (
	updatedUserState = models.UserState{
		PubKeyID: 1,
		TokenID:  models.MakeUint256(1),
		Balance:  models.MakeUint256(800),
		Nonce:    models.MakeUint256(1),
	}
)

type StateTreeTestSuite struct {
	*require.Assertions
	suite.Suite
	storage   *TestStorage
	leaf      *models.StateLeaf
	treeDepth uint8
}

func (s *StateTreeTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
	s.treeDepth = 32
}

func (s *StateTreeTestSuite) SetupTest() {
	var err error
	s.storage, err = NewTestStorage()
	s.NoError(err)

	state := models.UserState{
		PubKeyID: 1,
		TokenID:  models.MakeUint256(1),
		Balance:  models.MakeUint256(420),
		Nonce:    models.MakeUint256(0),
	}
	leaf, err := NewStateLeaf(0, &state)
	s.NoError(err)
	s.leaf = leaf
}

func (s *StateTreeTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *StateTreeTestSuite) TestLeaf_ReturnsCorrectStruct() {
	leaf, err := NewStateLeaf(0, &models.UserState{
		PubKeyID: 1,
		TokenID:  models.MakeUint256(1),
		Balance:  models.MakeUint256(420),
		Nonce:    models.MakeUint256(0),
	})
	s.NoError(err)

	_, err = s.storage.StateTree.Set(leaf.StateID, &leaf.UserState)
	s.NoError(err)

	actual, err := s.storage.StateTree.Leaf(leaf.StateID)
	s.NoError(err)
	s.Equal(leaf, actual)
}

func (s *StateTreeTestSuite) TestLeaf_NonexistentLeaf() {
	_, err := s.storage.StateTree.Leaf(0)
	s.ErrorIs(err, NewNotFoundError("state leaf"))
}

func (s *StateTreeTestSuite) setStateLeaves(stateIDs ...uint32) {
	for _, stateID := range stateIDs {
		_, err := s.storage.StateTree.Set(stateID, userState1)
		s.NoError(err)
	}
}

func (s *StateTreeTestSuite) TestNextAvailableStateID_EmptyStateTree() {
	stateID, err := s.storage.StateTree.NextAvailableStateID()
	s.NoError(err)
	s.Equal(uint32(0), *stateID)
}

func (s *StateTreeTestSuite) TestNextAvailableStateID_AllLeavesInOrder() {
	s.setStateLeaves(0, 1, 2)
	stateID, err := s.storage.StateTree.NextAvailableStateID()
	s.NoError(err)
	s.Equal(uint32(3), *stateID)
}

func (s *StateTreeTestSuite) TestNextAvailableStateID_BigGap() {
	s.setStateLeaves(0, 12345)

	stateID, err := s.storage.StateTree.NextAvailableStateID()
	s.NoError(err)
	s.Equal(uint32(1), *stateID)
}

func (s *StateTreeTestSuite) TestNextAvailableStateID_ReturnsTheFirstEmptyLeaf() {
	s.setStateLeaves(0, 2, 4)

	stateID, err := s.storage.StateTree.NextAvailableStateID()
	s.NoError(err)
	s.Equal(uint32(1), *stateID)
}

func (s *StateTreeTestSuite) TestNextVacantSubtree_ForSubtreeOfDepth1_AllLeavesInOrder() {
	s.setStateLeaves(0, 1)

	stateID, err := s.storage.StateTree.NextVacantSubtree(1)
	s.NoError(err)
	s.Equal(uint32(2), *stateID)
}

func (s *StateTreeTestSuite) TestNextVacantSubtree_ForSubtreeOfDepth1_ReturnsTheFirstEmptyLeaf() {
	s.setStateLeaves(0, 2, 3, 6)

	stateID, err := s.storage.StateTree.NextVacantSubtree(1)
	s.NoError(err)
	s.Equal(uint32(4), *stateID)
}

func (s *StateTreeTestSuite) TestNextVacantSubtree_ForSubtreeOfDepth1_AlignsTheSlotCorrectly() {
	s.setStateLeaves(0, 2, 6)

	stateID, err := s.storage.StateTree.NextVacantSubtree(1)
	s.NoError(err)
	s.Equal(uint32(4), *stateID)
}

func (s *StateTreeTestSuite) TestNextVacantSubtree_ForSubtreeOfDepth1_IgnoresMisalignedGaps() {
	s.setStateLeaves(0, 2, 5, 8)

	stateID, err := s.storage.StateTree.NextVacantSubtree(1)
	s.NoError(err)
	s.Equal(uint32(6), *stateID)
}

func (s *StateTreeTestSuite) TestNextVacantSubtree_ForSubtreeOfDepth2_IgnoresMisalignedGapsAndAlignsTheSlotCorrectly() {
	s.setStateLeaves(0, 5)

	stateID, err := s.storage.StateTree.NextVacantSubtree(2)
	s.NoError(err)
	s.Equal(uint32(8), *stateID)
}

func (s *StateTreeTestSuite) TestNextVacantSubtree_ForSubtreeOfDepth3_IgnoresMisalignedGapsAndAlignsTheSlotCorrectly() {
	s.setStateLeaves(0, 9)

	stateID, err := s.storage.StateTree.NextVacantSubtree(3)
	s.NoError(err)
	s.Equal(uint32(16), *stateID)
}

func (s *StateTreeTestSuite) TestNextVacantSubtree_ReturnsErrorWhenThereIsNoVacantSubtreeOfGivenDepth() {
	s.setStateLeaves(0)
	stateID, err := s.storage.StateTree.NextVacantSubtree(32)
	s.Nil(stateID)
	s.ErrorIs(err, NewNoVacantSubtreeError(32))

	s.setStateLeaves(uint32(1) << 31)
	stateID, err = s.storage.StateTree.NextVacantSubtree(31)
	s.Nil(stateID)
	s.ErrorIs(err, NewNoVacantSubtreeError(31))
}

func (s *StateTreeTestSuite) TestSet_StoresStateLeafRecord() {
	s.leaf.StateID = 0
	_, err := s.storage.StateTree.Set(0, &s.leaf.UserState)
	s.NoError(err)

	actualLeaf, err := s.storage.StateTree.Leaf(s.leaf.StateID)
	s.NoError(err)
	s.Equal(s.leaf, actualLeaf)
}

func (s *StateTreeTestSuite) TestSet_RootIsDifferentAfterSet() {
	state1 := models.UserState{
		PubKeyID: 1,
		TokenID:  models.MakeUint256(1),
		Balance:  models.MakeUint256(420),
		Nonce:    models.MakeUint256(0),
	}
	state2 := models.UserState{
		PubKeyID: 2,
		TokenID:  models.MakeUint256(5),
		Balance:  models.MakeUint256(100),
		Nonce:    models.MakeUint256(0),
	}

	_, err := s.storage.StateTree.Set(0, &state1)
	s.NoError(err)

	stateRootAfter1, err := s.storage.StateTree.Root()
	s.NoError(err)

	_, err = s.storage.StateTree.Set(0, &state2)
	s.NoError(err)

	stateRootAfter2, err := s.storage.StateTree.Root()
	s.NoError(err)

	s.NotEqual(stateRootAfter1, stateRootAfter2)
}

func (s *StateTreeTestSuite) TestSet_StoresLeafMerkleTreeNodeRecord() {
	_, err := s.storage.StateTree.Set(0, &s.leaf.UserState)
	s.NoError(err)

	expectedNode := &models.MerkleTreeNode{
		MerklePath: models.MerklePath{
			Path:  0,
			Depth: s.treeDepth,
		},
		DataHash: s.leaf.DataHash,
	}

	node, err := s.storage.StateTree.merkleTree.Get(expectedNode.MerklePath)
	s.NoError(err)
	s.Equal(expectedNode, node)
}

func (s *StateTreeTestSuite) TestSet_UpdatesRootMerkleTreeNodeRecord() {
	_, err := s.storage.StateTree.Set(0, &s.leaf.UserState)
	s.NoError(err)

	root, err := s.storage.StateTree.Root()
	s.NoError(err)
	s.Equal(common.HexToHash("0xd8cb702fc833817dccdc3889282af96755b2909274ca2f1a3827a60d11d796eb"), *root)
}

func (s *StateTreeTestSuite) TestSet_CalculatesCorrectRootForLeafOfId1() {
	_, err := s.storage.StateTree.Set(1, &s.leaf.UserState)
	s.NoError(err)

	root, err := s.storage.StateTree.Root()
	s.NoError(err)
	s.Equal(common.HexToHash("0xbec68099063e1499a5144a2d5b41f6a3e005ceac77caef6a171d77573570a000"), *root)
}

func (s *StateTreeTestSuite) TestSet_CalculatesCorrectRootForTwoLeaves() {
	_, err := s.storage.StateTree.Set(0, &s.leaf.UserState)
	s.NoError(err)

	state := models.UserState{
		PubKeyID: 2,
		TokenID:  models.MakeUint256(1),
		Balance:  models.MakeUint256(420),
		Nonce:    models.MakeUint256(0),
	}
	_, err = s.storage.StateTree.Set(1, &state)
	s.NoError(err)

	root, err := s.storage.StateTree.Root()
	s.NoError(err)
	s.Equal(common.HexToHash("0x7b1b0382bdffda7f4a6b24d974189c60797b87ce76836de6f18039e1dc73c050"), *root)
}

func (s *StateTreeTestSuite) TestSet_StoresStateUpdateRecord() {
	_, err := s.storage.StateTree.Set(0, &s.leaf.UserState)
	s.NoError(err)

	expectedUpdate := &models.StateUpdate{
		ID:          0,
		CurrentRoot: common.HexToHash("0xd8cb702fc833817dccdc3889282af96755b2909274ca2f1a3827a60d11d796eb"),
		PrevRoot:    merkletree.GetZeroHash(s.treeDepth),
		PrevStateLeaf: models.StateLeaf{
			StateID:  0,
			DataHash: merkletree.GetZeroHash(0),
		},
	}

	update, err := s.storage.StateTree.getStateUpdate(expectedUpdate.ID)
	s.NoError(err)
	s.Equal(expectedUpdate, update)
}

func (s *StateTreeTestSuite) TestSet_UpdateExistingLeafCorrectRootMerkleTreeNode() {
	_, err := s.storage.StateTree.Set(0, &s.leaf.UserState)
	s.NoError(err)

	_, err = s.storage.StateTree.Set(0, &updatedUserState)
	s.NoError(err)

	root, err := s.storage.StateTree.Root()
	s.NoError(err)
	s.Equal(common.HexToHash("0x406515786640be8c51eacf1221f017e7f59e04ef59637a27dcb2b2f054b309bf"), *root)
}

func (s *StateTreeTestSuite) TestSet_UpdateExistingLeafCorrectLeafMerkleTreeNode() {
	_, err := s.storage.StateTree.Set(0, &s.leaf.UserState)
	s.NoError(err)

	leaf, err := NewStateLeaf(0, &updatedUserState)
	s.NoError(err)
	_, err = s.storage.StateTree.Set(0, &updatedUserState)
	s.NoError(err)

	leafPath := models.MerklePath{
		Path:  0,
		Depth: s.treeDepth,
	}

	expectedLeaf := &models.MerkleTreeNode{
		MerklePath: leafPath,
		DataHash:   leaf.DataHash,
	}

	leafNode, err := s.storage.StateTree.merkleTree.Get(models.MerklePath{Path: 0, Depth: s.treeDepth})
	s.NoError(err)
	s.Equal(expectedLeaf, leafNode)
}

func (s *StateTreeTestSuite) TestSet_UpdateExistingLeafNewStateLeafRecord() {
	_, err := s.storage.StateTree.Set(0, &s.leaf.UserState)
	s.NoError(err)

	expectedLeaf, err := NewStateLeaf(0, &updatedUserState)
	s.NoError(err)
	_, err = s.storage.StateTree.Set(0, &updatedUserState)
	s.NoError(err)

	leaf, err := s.storage.StateTree.Leaf(0)
	s.NoError(err)
	s.Equal(expectedLeaf, leaf)
}

func (s *StateTreeTestSuite) TestSet_UpdateExistingLeafAddsStateUpdateRecord() {
	_, err := s.storage.StateTree.Set(0, &s.leaf.UserState)
	s.NoError(err)

	_, err = s.storage.StateTree.Set(0, &updatedUserState)
	s.NoError(err)

	expectedUpdate := &models.StateUpdate{
		ID:            1,
		CurrentRoot:   common.HexToHash("0x406515786640be8c51eacf1221f017e7f59e04ef59637a27dcb2b2f054b309bf"),
		PrevRoot:      common.HexToHash("0xd8cb702fc833817dccdc3889282af96755b2909274ca2f1a3827a60d11d796eb"),
		PrevStateLeaf: *s.leaf,
	}

	update, err := s.storage.StateTree.getStateUpdate(expectedUpdate.ID)
	s.NoError(err)
	s.Equal(expectedUpdate, update)
}

func (s *StateTreeTestSuite) TestSet_ReturnsWitness() {
	witness, err := s.storage.StateTree.Set(0, &s.leaf.UserState)
	s.NoError(err)
	s.Len(witness, int(s.treeDepth))

	node, err := s.storage.StateTree.merkleTree.Get(models.MerklePath{Depth: s.treeDepth, Path: 1})
	s.NoError(err)
	s.Equal(node.DataHash, witness[0])

	node, err = s.storage.StateTree.merkleTree.Get(models.MerklePath{Depth: 1, Path: 1})
	s.NoError(err)
	s.Equal(node.DataHash, witness[31])
}

func (s *StateTreeTestSuite) TestLeavesCount_UpdateCountAfterAddingNewLeaves() {
	s.setStateLeaves(0, 2, 4, 8, 9)

	count, err := s.storage.StateTree.LeavesCount()
	s.NoError(err)
	s.EqualValues(5, count)
}

func (s *StateTreeTestSuite) TestLeavesCount_NoLeaves() {
	count, err := s.storage.StateTree.LeavesCount()
	s.NoError(err)
	s.EqualValues(0, count)
}

func (s *StateTreeTestSuite) TestRevertTo() {
	states := []models.UserState{
		{
			PubKeyID: 1,
			TokenID:  models.MakeUint256(1),
			Balance:  models.MakeUint256(420),
			Nonce:    models.MakeUint256(0),
		},
		{
			PubKeyID: 2,
			TokenID:  models.MakeUint256(5),
			Balance:  models.MakeUint256(100),
			Nonce:    models.MakeUint256(0),
		},
		{
			PubKeyID: 1,
			TokenID:  models.MakeUint256(1),
			Balance:  models.MakeUint256(500),
			Nonce:    models.MakeUint256(0),
		},
	}

	_, err := s.storage.StateTree.Set(0, &states[0])
	s.NoError(err)

	stateRoot, err := s.storage.StateTree.Root()
	s.NoError(err)

	_, err = s.storage.StateTree.Set(1, &states[1])
	s.NoError(err)
	_, err = s.storage.StateTree.Set(0, &states[2])
	s.NoError(err)

	err = s.storage.StateTree.RevertTo(*stateRoot)
	s.NoError(err)

	newStateRoot, err := s.storage.StateTree.Root()
	s.NoError(err)
	s.Equal(stateRoot, newStateRoot)

	leaf, err := s.storage.StateTree.Leaf(0)
	s.NoError(err)
	s.Equal(states[0], leaf.UserState)
}

func (s *StateTreeTestSuite) TestRevertTo_NonexistentRootHash() {
	states := []models.UserState{
		{
			PubKeyID: 1,
			TokenID:  models.MakeUint256(1),
			Balance:  models.MakeUint256(420),
			Nonce:    models.MakeUint256(0),
		},
		{
			PubKeyID: 2,
			TokenID:  models.MakeUint256(5),
			Balance:  models.MakeUint256(100),
			Nonce:    models.MakeUint256(0),
		},
		{
			PubKeyID: 1,
			TokenID:  models.MakeUint256(1),
			Balance:  models.MakeUint256(500),
			Nonce:    models.MakeUint256(0),
		},
	}
	for i := range states {
		_, err := s.storage.StateTree.Set(uint32(i), &states[i])
		s.NoError(err)
	}

	err := s.storage.StateTree.RevertTo(common.Hash{1, 2, 3})
	s.ErrorIs(err, ErrNonexistentState)
}

func (s *StateTreeTestSuite) TestIterateLeaves() {
	state := s.leaf.UserState
	expectedLeaves := []models.StateLeaf{
		{
			StateID:   0,
			DataHash:  s.leaf.DataHash,
			UserState: state,
		},
		{
			StateID:   1,
			DataHash:  s.leaf.DataHash,
			UserState: state,
		},
	}

	for i := range expectedLeaves {
		_, err := s.storage.StateTree.Set(expectedLeaves[i].StateID, &expectedLeaves[i].UserState)
		s.NoError(err)
	}

	leaves := make([]models.StateLeaf, 0, len(expectedLeaves))
	err := s.storage.StateTree.IterateLeaves(func(stateLeaf *models.StateLeaf) error {
		leaves = append(leaves, *stateLeaf)
		return nil
	})
	s.NoError(err)

	s.Len(leaves, len(expectedLeaves))
	s.Equal(expectedLeaves, leaves)
}

func (s *StateTreeTestSuite) TestIterateLeaves_NoLeaves() {
	leaves := make([]models.StateLeaf, 0, 1)
	err := s.storage.StateTree.IterateLeaves(func(stateLeaf *models.StateLeaf) error {
		leaves = append(leaves, *stateLeaf)
		return nil
	})
	s.NoError(err)
	s.Len(leaves, 0)
}

func TestStateTreeTestSuite(t *testing.T) {
	suite.Run(t, new(StateTreeTestSuite))
}
