package admin

import (
	"context"

	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/Worldcoin/hubble-commander/models/stored"
	"github.com/sirupsen/logrus"
)

func (a *API) InspectStateID(ctx context.Context, stateID uint32) error {
	// print the current nonce and balance (from the state tree)
	batchedState, err := a.storage.StateTree.Leaf(stateID)
	if err != nil {
		return err
	}

	logrus.WithFields(logrus.Fields{
		"StateID": stateID,
		"Nonce":   batchedState.Nonce.Uint64(),
		"Balance": batchedState.Balance.Uint64(),
	}).Info("batched state")

	// print the first 10 and last 10 txns from the mempool
	allMempoolTxns, err := a.storage.GetAllMempoolTransactions()
	if err != nil {
		return err
	}

	ourTxns := make([]stored.PendingTx, 0, len(allMempoolTxns))
	for i := range allMempoolTxns {
		txn := allMempoolTxns[i]
		if txn.FromStateID != stateID {
			continue
		}

		ourTxns = append(ourTxns, txn)
	}

	for i := 0; i < len(ourTxns); i++ {
		if i >= 10 && i < len(ourTxns)-10 {
			continue
		}

		txn := ourTxns[i]
		logrus.WithFields(logrus.Fields{
			"Nonce":  txn.Nonce.Uint64(),
			"Amount": txn.Amount.Uint64(),
			"Type":   txtype.TransactionTypes[txn.TxType],
			"Hash":   txn.Hash.Hex(),
		}).Info("mempool txn")
	}

	// print the pending state
	pendingNonce, err := a.storage.GetPendingNonce(stateID)
	if err != nil {
		return err
	}

	pendingBalance, err := a.storage.GetPendingBalance(stateID)
	if err != nil {
		return err
	}

	logrus.WithFields(logrus.Fields{
		"StateID": stateID,
		"Nonce":   pendingNonce.Uint64(),
		"Balance": pendingBalance.Uint64(),
	}).Info("pending state")

	return nil
}
