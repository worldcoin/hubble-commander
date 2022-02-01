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

var defaultTestMineTimeout = 5 * time.Minute

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

func (m *MockReceiptProvider) Commit() {
	// NOOP
}

func (s *WaitToBeMinedTestSuite) TestWaitToBeMined_CallsTransactionReceiptImmediately() {
	var callTime time.Time

	rp := new(MockReceiptProvider)
	rp.On(transactionReceiptMethod, mock.Anything, mock.Anything).
		Return(s.minedReceipt, nil).
		Run(func(args mock.Arguments) { callTime = time.Now() })

	now := time.Now()
	_, err := WaitToBeMined(rp, defaultTestMineTimeout, s.tx)
	s.NoError(err)
	s.WithinDuration(now, callTime, 20*time.Millisecond)
}

func withPollInterval(interval time.Duration, fn func()) {
	initialPollInterval := pollInterval
	pollInterval = interval
	defer func() { pollInterval = initialPollInterval }()
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
		_, err := WaitToBeMined(rp, defaultTestMineTimeout, s.tx)
		s.NoError(err)
	})

	s.WithinDuration(expected, secondCallTime, 20*time.Millisecond)
}

func (s *WaitToBeMinedTestSuite) TestWaitToBeMined_EventuallyTimesOut() {
	var nilReceipt *types.Receipt

	rp := new(MockReceiptProvider)
	rp.On(transactionReceiptMethod, mock.Anything, mock.Anything).
		Return(nilReceipt, ethereum.NotFound)

	testTimeout := 50 * time.Millisecond
	expected := time.Now().Add(testTimeout)

	_, err := WaitToBeMined(rp, testTimeout, s.tx)
	s.ErrorIs(err, ErrWaitToBeMinedTimedOut)

	s.WithinDuration(expected, time.Now(), 20*time.Millisecond)
}

func (s *WaitToBeMinedTestSuite) TestWaitToBeMined_HandlesTransactionReceiptCallTimeouts() {
	var nilReceipt *types.Receipt

	rp := new(MockReceiptProvider)
	rp.On(transactionReceiptMethod, mock.Anything, mock.Anything).
		Return(nilReceipt, context.DeadlineExceeded)

	callTime := time.Now()

	_, err := WaitToBeMined(rp, 500*time.Millisecond, s.tx)
	s.ErrorIs(err, ErrWaitToBeMinedTimedOut)

	s.WithinDuration(callTime, time.Now(), 20*time.Millisecond)
}

func (s *WaitToBeMinedTestSuite) TestWaitForMultipleTxs_WaitsForAllTransactionsAndReturnsReceiptsInOrder() {
	txs := make([]types.Transaction, 2)
	expectedReceipts := make([]types.Receipt, len(txs))
	for i := range txs {
		txs[i] = *newTx()
		expectedReceipts[i] = types.Receipt{BlockNumber: big.NewInt(int64(i))}
	}

	rp := new(MockReceiptProvider)
	rp.On(transactionReceiptMethod, mock.Anything, txs[0].Hash()).
		Return(&expectedReceipts[0], nil).
		Run(func(args mock.Arguments) { time.Sleep(50 * time.Millisecond) })

	rp.On(transactionReceiptMethod, mock.Anything, txs[1].Hash()).
		Return(&expectedReceipts[1], nil)

	receipts, err := WaitForMultipleTxs(rp, defaultTestMineTimeout, txs...)
	s.NoError(err)
	s.Equal(receipts, []types.Receipt{expectedReceipts[0], expectedReceipts[1]})
}

func (s *WaitToBeMinedTestSuite) TestWaitForMultipleTxs_WorksForDuplicatedTransactions() {
	rp := new(MockReceiptProvider)
	rp.On(transactionReceiptMethod, mock.Anything, mock.Anything).
		Return(s.minedReceipt, nil)

	receipts, err := WaitForMultipleTxs(rp, defaultTestMineTimeout, *s.tx, *s.tx)
	s.NoError(err)
	s.Equal(receipts, []types.Receipt{*s.minedReceipt, *s.minedReceipt})
}

func (s *WaitToBeMinedTestSuite) TestWaitForMultipleTxs_FinishesOnTimeout() {
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

	receipts, err := WaitForMultipleTxs(rp, testTimeout, txs...)
	s.ErrorIs(err, ErrWaitToBeMinedTimedOut)
	s.Nil(receipts)

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
