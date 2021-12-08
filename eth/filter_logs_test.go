package eth

import (
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type FilterLogsTestSuite struct {
	*require.Assertions
	suite.Suite
	client *TestClient
}

func (s *FilterLogsTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *FilterLogsTestSuite) SetupTest() {
	client, err := NewTestClient()
	s.NoError(err)
	s.client = client
}

func (s *FilterLogsTestSuite) TearDownTest() {
	s.client.Close()
}

func (s *FilterLogsTestSuite) TestClient_FilterLogs_ErrorHandling() {
	var (
		start uint64 = 0
		end   uint64 = 100
	)

	// failed iterator
	it := mockIterator{fail: errors.New("something bad happened")}

	err := s.client.FilterLogs(s.client.AccountRegistry.BoundContract, "BatchPubkeyRegistered", &bind.FilterOpts{
		Start: start,
		End:   &end,
	}, &it)
	s.Error(err)
}

type mockIterator struct {
	fail error
}

func (t *mockIterator) SetData(_ *bind.BoundContract, _ string, _ chan types.Log, _ ethereum.Subscription) {
	// there is no need to do anything here
}

func (t *mockIterator) Error() error {
	return t.fail
}

func TestFilterLogsTestSuite(t *testing.T) {
	suite.Run(t, new(FilterLogsTestSuite))
}
