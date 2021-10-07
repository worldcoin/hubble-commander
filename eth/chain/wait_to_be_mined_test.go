package chain

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type WaitToBeMinedTestSuite struct {
	*require.Assertions
	suite.Suite
	tx           *types.Transaction
	minedReceipt *types.Receipt
}

func (s *WaitToBeMinedTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *WaitToBeMinedTestSuite) SetupTest() {
	s.tx = types.NewTx(&types.DynamicFeeTx{
		ChainID:   big.NewInt(1),
		Nonce:     0,
		To:        ref.Address(utils.RandomAddress()),
		Gas:       123457,
		GasFeeCap: big.NewInt(12_500_000),
		GasTipCap: big.NewInt(0),
		Data:      []byte{},
	})
	s.minedReceipt = &types.Receipt{BlockNumber: big.NewInt(1234)}
}

type MockReceiptProvider struct {
	mock.Mock
}

func (m *MockReceiptProvider) TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	args := m.Called(ctx, txHash)
	return args.Get(0).(*types.Receipt), args.Error(1)
}

func (s *WaitToBeMinedTestSuite) TestWaitToBeMined_CallsTransactionReceiptImmediately() {
	var callTime time.Time

	rp := new(MockReceiptProvider)
	rp.On("TransactionReceipt", mock.Anything, mock.Anything).
		Return(s.minedReceipt, nil).
		Run(func(args mock.Arguments) { callTime = time.Now() })

	now := time.Now()
	_, err := WaitToBeMined(rp, s.tx)
	s.NoError(err)
	s.WithinDuration(now, callTime, 20*time.Millisecond)
}

func withPollInterval(interval time.Duration, fn func()) {
	initialPollInterval := PollInterval
	PollInterval = interval
	defer func() { PollInterval = initialPollInterval }()
	fn()
}

func (s *WaitToBeMinedTestSuite) TestWaitToBeMined_MakesTheSecondCallAfterInterval() {
	var nilReceipt *types.Receipt
	var secondCallTime time.Time

	rp := new(MockReceiptProvider)
	rp.On("TransactionReceipt", mock.Anything, mock.Anything).
		Return(nilReceipt, ethereum.NotFound).
		Once()

	rp.On("TransactionReceipt", mock.Anything, mock.Anything).
		Return(s.minedReceipt, nil).
		Run(func(args mock.Arguments) { secondCallTime = time.Now() }).
		Once()

	testPollInterval := 50 * time.Millisecond
	expected := time.Now().Add(testPollInterval)

	withPollInterval(testPollInterval, func() {
		_, err := WaitToBeMined(rp, s.tx)
		s.NoError(err)
	})

	s.WithinDuration(expected, secondCallTime, 20*time.Millisecond)
}

func withChainTimeout(timeout time.Duration, fn func()) {
	initialChainTimeout := MineTimeout
	MineTimeout = timeout
	defer func() { MineTimeout = initialChainTimeout }()
	fn()
}

func (s *WaitToBeMinedTestSuite) TestWaitToBeMined_EventuallyTimesOut() {
	var nilReceipt *types.Receipt

	rp := new(MockReceiptProvider)
	rp.On("TransactionReceipt", mock.Anything, mock.Anything).
		Return(nilReceipt, ethereum.NotFound)

	testTimeout := 50 * time.Millisecond
	expected := time.Now().Add(testTimeout)

	withChainTimeout(testTimeout, func() {
		_, err := WaitToBeMined(rp, s.tx)
		s.ErrorIs(err, ErrWaitToBeMinedTimedOut)
	})

	s.WithinDuration(expected, time.Now(), 20*time.Millisecond)
}

func TestWaitToBeMinedTestSuite(t *testing.T) {
	suite.Run(t, new(WaitToBeMinedTestSuite))
}
