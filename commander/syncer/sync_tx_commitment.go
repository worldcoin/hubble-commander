package syncer

import (
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/stored"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
)

var (
	ErrInvalidDataLength = NewDisputableError(Transition, "invalid data length")
	ErrTooManyTxs        = NewDisputableError(Transition, "too many transactions in a commitment")
)

const (
	invalidStateRootMessage      = "invalid commitment post state root"
	invalidWithdrawRootMessage   = "invalid commitment withdraw root"
	invalidTokenID               = "invalid commitment token id"
	mismatchedTotalAmountMessage = "mismatched commitment total amount"
	nonexistentReceiverMessage   = "nonexistent receiver"
)

func (c *TxsContext) syncTxCommitment(commitment encoder.Commitment) error {
	decodedCommitment := commitment.ToDecodedCommitment()
	if len(decodedCommitment.Transactions)%c.Syncer.TxLength() != 0 {
		return ErrInvalidDataLength
	}

	syncedTxs, err := c.Syncer.DeserializeTxs(decodedCommitment.Transactions)
	if err != nil {
		return err
	}

	if uint32(syncedTxs.Txs().Len()) > c.cfg.MaxTxsPerCommitment {
		return ErrTooManyTxs
	}

	appliedTxs, stateProofs, err := c.SyncTxs(syncedTxs, decodedCommitment.FeeReceiver)
	if err != nil {
		return err
	}
	syncedTxs.SetTxs(appliedTxs)

	err = c.verifyStateRoot(decodedCommitment.StateRoot, stateProofs)
	if err != nil {
		return err
	}

	err = c.Syncer.VerifyAmountAndWithdrawRoots(commitment, appliedTxs, stateProofs)
	if err != nil {
		return err
	}

	err = c.Syncer.SetMissingTxsData(commitment, syncedTxs)
	if st.IsNotFoundError(err) {
		return c.createDisputableSignatureError(nonexistentReceiverMessage, syncedTxs.Txs())
	}
	if err != nil {
		return err
	}
	if !c.cfg.DisableSignatures {
		err = c.verifyTxSignature(decodedCommitment, syncedTxs.Txs())
		if err != nil {
			return err
		}
	}

	return c.addTxs(syncedTxs.Txs(), &decodedCommitment.ID)
}

func (c *TxsContext) verifyStateRoot(commitmentPostState common.Hash, proofs []models.StateMerkleProof) error {
	postStateRoot, err := c.storage.StateTree.Root()
	if err != nil {
		return err
	}
	if *postStateRoot != commitmentPostState {
		return NewDisputableErrorWithProofs(Transition, invalidStateRootMessage, proofs)
	}
	return nil
}

func (c *TxsContext) addTxs(txs models.GenericTransactionArray, commitmentID *models.CommitmentID) error {
	if txs.Len() == 0 {
		return nil
	}

	if txs.Len() > 255 {
		panic("Commitments cannot have more than 255 transactions")
	}

	for i := 0; i < txs.Len(); i++ {
		txs.At(i).GetBase().CommitmentSlot = models.NewCommitmentSlot(*commitmentID, uint8(i))

		hashTransfer, err := c.Syncer.HashTx(txs.At(i))
		if err != nil {
			return err
		}
		txs.At(i).GetBase().Hash = *hashTransfer
	}

	// TODO: We have just received a batch which includes transactions. Those
	//       transactions no longer belong in the mempool, the copies we have in our
	//       mempool are now duplicates and should be removed!
	log.WithFields(log.Fields{
		"commitmentID": commitmentID,
	}).Errorf("Unimplemented: synced transactions not removed from mempool.")

	batchedTxs := make([]stored.BatchedTx, 0)
	for i := 0; i < txs.Len(); i++ {
		batchedTxs = append(batchedTxs, *stored.NewBatchedTx(txs.At(i)))
	}

	err := c.storage.MempoolDropBatchedTransactions(batchedTxs)
	if err != nil {
		return err
	}

	// the below is going to add each of these batched transactions and should only fail
	// if for some reason the txs were already batched

	return c.storage.BatchUpsertTransaction(txs)
}
