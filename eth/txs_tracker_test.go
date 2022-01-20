package eth

import (
	"testing"
	"time"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type TxsTrackerTestSuite struct {
	*require.Assertions
	suite.Suite
	client  *TestClient
	tracker *TxsTracker
}

func (s *TxsTrackerTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *TxsTrackerTestSuite) SetupTest() {
	client, err := NewTestClient()
	s.NoError(err)
	s.client = client
	s.tracker = NewTxsTracker(client.Blockchain)
}

func (s *TxsTrackerTestSuite) TearDownTest() {
	s.client.Close()
}

func (s *TxsTrackerTestSuite) TestCheckTransaction_SuccessfulTransaction() {
	publicKey := models.PublicKey{1, 2, 3}
	tx, err := s.client.accountRegistry().
		WithGasLimit(500_000).
		Register(publicKey.BigInts())
	s.NoError(err)

	s.tracker.CheckTransaction(tx)

	s.Never(func() bool {
		err = <-s.tracker.Fail()
		return true
	}, time.Second, time.Millisecond*200)
	s.NoError(err)
}

func (s *TxsTrackerTestSuite) TestCheckTransaction_FailedTransaction() {
	publicKey := models.PublicKey{1, 2, 3}
	tx, err := s.client.accountRegistry().
		WithGasLimit(24_000).
		Register(publicKey.BigInts())
	s.NoError(err)

	s.tracker.CheckTransaction(tx)

	s.Eventually(func() bool {
		err = <-s.tracker.Fail()
		return true
	}, time.Second, time.Millisecond*200)
	s.Error(err)
}

func TestTxsTrackerTestSuite(t *testing.T) {
	suite.Run(t, new(TxsTrackerTestSuite))
}
