package eth

import (
	"context"

	"github.com/Worldcoin/hubble-commander/eth/chain"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

type TxsTracker struct {
	fails      chan error
	blockchain chain.Connection
}

func NewTxsTracker(blockchain chain.Connection) *TxsTracker {
	return &TxsTracker{
		fails:      make(chan error),
		blockchain: blockchain,
	}
}

func (t *TxsTracker) CheckTransactionWithReceipt(tx *types.Transaction, receipt *types.Receipt) {
	if receipt.Status == 1 {
		return
	}

	err := getRevertMessage(t.blockchain, tx, receipt)
	if err != nil {
		t.fails <- err
	}
}

func (t *TxsTracker) CheckTransaction(tx *types.Transaction) {
	receipt, err := t.blockchain.GetBackend().TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		t.fails <- errors.WithStack(errors.Wrap(err, "can't get a tx receipt"))
		return
	}
	t.CheckTransactionWithReceipt(tx, receipt)
}

func (t *TxsTracker) Fail() <-chan error {
	return t.fails
}

func (a *AccountManager) RegistrationFail() <-chan error {
	return a.registrationTxsTracker.Fail()
}
