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

const transactionReceiptMethod = "TransactionReceipt"

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
	s.tx = newTx()
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
	rp.On(transactionReceiptMethod, mock.Anything, mock.Anything).
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
	rp.On(transactionReceiptMethod, mock.Anything, mock.Anything).
		Return(nilReceipt, ethereum.NotFound).
		Once()

	rp.On(transactionReceiptMethod, mock.Anything, mock.Anything).
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

func withMineTimeout(timeout time.Duration, fn func()) {
	initialMineTimeout := MineTimeout
	MineTimeout = timeout
	defer func() { MineTimeout = initialMineTimeout }()
	fn()
}

func (s *WaitToBeMinedTestSuite) TestWaitToBeMined_EventuallyTimesOut() {
	var nilReceipt *types.Receipt

	rp := new(MockReceiptProvider)
	rp.On(transactionReceiptMethod, mock.Anything, mock.Anything).
		Return(nilReceipt, ethereum.NotFound)

	testTimeout := 50 * time.Millisecond
	expected := time.Now().Add(testTimeout)

	withMineTimeout(testTimeout, func() {
		_, err := WaitToBeMined(rp, s.tx)
		s.ErrorIs(err, ErrWaitToBeMinedTimedOut)
	})

	s.WithinDuration(expected, time.Now(), 20*time.Millisecond)
}

func (s *WaitToBeMinedTestSuite) TestWaitForMultiple_WaitsForAllTransactions() {
	txs := make([]types.Transaction, 2)
	for i := range txs {
		txs[i] = *newTx()
	}
	expectedReceipts := make([]types.Receipt, len(txs))
	expectedReceipts[0] = types.Receipt{BlockNumber: big.NewInt(0)}
	expectedReceipts[1] = types.Receipt{BlockNumber: big.NewInt(1)}

	calls := make([]bool, len(txs))
	rp := new(MockReceiptProvider)
	rp.On(transactionReceiptMethod, mock.Anything, txs[0].Hash()).
		Return(&expectedReceipts[0], nil).
		Run(func(args mock.Arguments) {
			time.Sleep(50 * time.Millisecond)
			calls[0] = true
		})

	rp.On(transactionReceiptMethod, mock.Anything, txs[1].Hash()).
		Return(&expectedReceipts[1], nil).
		Run(func(args mock.Arguments) { calls[1] = true })

	receipts, err := WaitForMultiple(rp, txs)
	s.NoError(err)
	s.Equal([]bool{true, true}, calls)
	s.Contains(receipts, expectedReceipts[0])
	s.Contains(receipts, expectedReceipts[1])
}

func (s *WaitToBeMinedTestSuite) TestWaitForMultiple_FinishesOnTimeout() {
	txs := make([]types.Transaction, 2)
	for i := range txs {
		txs[i] = *newTx()
	}

	var nilReceipt *types.Receipt

	calls := make([]bool, len(txs))
	rp := new(MockReceiptProvider)
	rp.On(transactionReceiptMethod, mock.Anything, txs[0].Hash()).
		Return(nilReceipt, ethereum.NotFound)

	rp.On(transactionReceiptMethod, mock.Anything, txs[1].Hash()).
		Return(s.minedReceipt, nil).
		Run(func(args mock.Arguments) { calls[1] = true })

	testTimeout := 50 * time.Millisecond
	expected := time.Now().Add(testTimeout)

	withMineTimeout(testTimeout, func() {
		receipts, err := WaitForMultiple(rp, txs)
		s.ErrorIs(err, ErrWaitToBeMinedTimedOut)
		s.Nil(receipts)
	})

	s.WithinDuration(expected, time.Now(), 20*time.Millisecond)
	s.True(calls[1])
}

func newTx() *types.Transaction {
	return types.NewTx(&types.DynamicFeeTx{
		ChainID:   big.NewInt(1),
		Nonce:     0,
		To:        ref.Address(utils.RandomAddress()),
		Gas:       123457,
		GasFeeCap: big.NewInt(12_500_000),
		Data:      []byte{},
	})
}

func TestWaitToBeMinedTestSuite(t *testing.T) {
	suite.Run(t, new(WaitToBeMinedTestSuite))
}
