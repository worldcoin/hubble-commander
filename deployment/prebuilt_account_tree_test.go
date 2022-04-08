package deployment

import (
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/Worldcoin/hubble-commander/contracts/accountregistry"
	"github.com/Worldcoin/hubble-commander/contracts/defaultaccountregistry"
	"github.com/Worldcoin/hubble-commander/eth/deployer"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/testutils/simulator"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// I. The tests:

type PrebuiltAccountTreeTestSuite struct {
	*require.Assertions
	suite.Suite

	sim     *simulator.Simulator
	chooser *common.Address
}

func (s *PrebuiltAccountTreeTestSuite) TestEmptyRootsMatch() {
	// all three methods of building empty trees better result in the same root hash
	defaultAccountRegistry := s.deployEmptyRegistry()
	tree := NewTree(32)
	accountRegistry := s.deployAccountRegistry(tree)

	defaultAccountRegistryHash := s.getContractRootHash(defaultAccountRegistry.Root)
	s.Equal(tree.Root(), defaultAccountRegistryHash)

	accountRegistryHash := s.getContractRootHash(accountRegistry.Root)
	s.Equal(tree.Root(), accountRegistryHash)
}

func (s *PrebuiltAccountTreeTestSuite) TestAddAccountsDefaultRegistry() {
	// if adding some accounts gives us the correct root then all the intermediate
	// subtrees were initialized correctly

	accountRegistry := s.deployEmptyRegistry()
	tree := NewTree(32)

	randomKeys := randomPublicKeys(129)

	s.registerPublicKeys(accountRegistry.Register, randomKeys)

	for i := range randomKeys {
		key := randomKeys[i]
		tree.RegisterAccount(key)
	}

	populatedRootHash := s.getContractRootHash(accountRegistry.Root)
	s.Equal(tree.Root(), populatedRootHash)
}

func (s *PrebuiltAccountTreeTestSuite) TestAddAccountsPrebuiltRegistry() {
	// if we can add more accounts and get the correct root than all the intermediate
	// subtrees were initialized correctly

	// 1. come up with an initial state to deploy our account registry with

	accountRegistry := s.deployEmptyRegistry()
	tree := NewTree(32)

	randomKeys := randomPublicKeys(129)

	s.registerPublicKeys(accountRegistry.Register, randomKeys)

	for i := range randomKeys {
		key := randomKeys[i]
		tree.RegisterAccount(key)
	}

	// 2. add some more accounts to our registry

	newAccountRegistry := s.deployAccountRegistry(tree)

	randomKeys = randomPublicKeys(129)

	s.registerPublicKeys(accountRegistry.Register, randomKeys)
	s.registerPublicKeys(newAccountRegistry.Register, randomKeys)

	for i := range randomKeys {
		key := randomKeys[i]
		tree.RegisterAccount(key)
	}

	oldRegistryRootHash := s.getContractRootHash(accountRegistry.Root)
	s.Equal(tree.Root(), oldRegistryRootHash)

	newRegistryRootHash := s.getContractRootHash(newAccountRegistry.Root)
	s.Equal(tree.Root(), newRegistryRootHash)
}

func (s *PrebuiltAccountTreeTestSuite) TestIncorrectAccountCount() {
	// the contract needs to know which leaf was last added, add an explicit test that
	// we're correctly threading it through

	tree := NewTree(32)

	// 1. build an initial state

	randomKeys := randomPublicKeys(129)
	for i := range randomKeys {
		key := randomKeys[i]
		tree.RegisterAccount(key)
	}

	// 2. deploy our account registry

	seconds, err := time.ParseDuration("5s")
	s.NoError(err)

	leftRoot := tree.LeftRoot()
	_, _, accountRegistry, err := deployer.DeployAccountRegistry(
		s.sim,
		s.chooser,
		seconds,
		&leftRoot,
		tree.AccountCount-1, // intentionally giving it an incorrect count
		(*[31]common.Hash)(tree.Subtrees),
	)
	s.NoError(err)
	s.NotNil(accountRegistry)

	// 3. notice that adding more leaves moves us to an incorrect state

	randomKeys = randomPublicKeys(5)
	s.registerPublicKeys(accountRegistry.Register, randomKeys)
	for i := range randomKeys {
		key := randomKeys[i]
		tree.RegisterAccount(key)
	}

	rootHash := s.getContractRootHash(accountRegistry.Root)
	s.NotEqual(tree.Root(), rootHash)
}

// II. Utilities for running the tests:

func (s *PrebuiltAccountTreeTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *PrebuiltAccountTreeTestSuite) SetupTest() {
	sim, err := simulator.NewSimulator()
	s.NoError(err)
	s.sim = sim

	seconds, err := time.ParseDuration("5s")
	s.NoError(err)

	chooser, _, err := deployer.DeployProofOfAuthority(
		s.sim,
		seconds,
		[]common.Address{s.sim.GetAccount().From},
	)
	s.NoError(err)
	s.chooser = chooser

	s.sim.Commit()
}

func (s *PrebuiltAccountTreeTestSuite) TearDownTest() {
	s.sim.Close()
	s.chooser = nil
}

func TestPrebuiltAccountTreeTestSuite(t *testing.T) {
	suite.Run(t, new(PrebuiltAccountTreeTestSuite))
}

// all the tests are above, below are utilities

func randomPublicKey() *models.PublicKey {
	keyBytes := make([]byte, 128)
	rand.Read(keyBytes) //nolint:gosec

	result := models.PublicKey(*(*[128]byte)(keyBytes))
	return &result
}

func randomPublicKeys(count int) []*models.PublicKey {
	var randomKeys []*models.PublicKey
	for i := 0; i < count; i++ {
		randomKeys = append(randomKeys, randomPublicKey())
	}
	return randomKeys
}

type callHash func(*bind.CallOpts) ([32]byte, error)

func (s *PrebuiltAccountTreeTestSuite) getContractRootHash(f callHash) common.Hash {
	root, err := f(nil)
	s.NoError(err)

	return common.BytesToHash(root[:])
}

type registerKey func(*bind.TransactOpts, [4]*big.Int) (*types.Transaction, error)

func (s *PrebuiltAccountTreeTestSuite) registerPublicKeys(f registerKey, keys []*models.PublicKey) {
	for i := range keys {
		key := keys[i]
		ints := key.BigInts()

		_, err := f(s.sim.GetAccount(), ints)
		s.NoError(err)

		s.sim.Commit()
	}
}

func (s *PrebuiltAccountTreeTestSuite) deployEmptyRegistry() *defaultaccountregistry.DefaultAccountRegistry {
	_, _, accountRegistry, err := defaultaccountregistry.DeployDefaultAccountRegistry(
		s.sim.GetAccount(), s.sim.GetBackend(), *s.chooser,
	)
	s.NoError(err)

	s.sim.Commit()

	return accountRegistry
}

func (s *PrebuiltAccountTreeTestSuite) deployAccountRegistry(tree *Tree) *accountregistry.AccountRegistry {
	seconds, err := time.ParseDuration("5s")
	s.NoError(err)

	leftRoot := tree.LeftRoot()
	_, _, accountRegistry, err := deployer.DeployAccountRegistry(
		s.sim,
		s.chooser,
		seconds,
		&leftRoot,
		tree.AccountCount,
		(*[31]common.Hash)(tree.Subtrees),
	)
	s.NoError(err)
	s.NotNil(accountRegistry)

	return accountRegistry
}
