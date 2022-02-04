package tracker

import (
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/core/types"
)

type TxsSender struct {
	client       *eth.Client
	requestsChan chan *txRequest
}

func newTxRequestsSender(ethClient *eth.Client) *TxsSender {
	return &TxsSender{
		client:       ethClient,
		requestsChan: make(chan *txRequest),
	}
}

func (t *TxsSender) SubmitTransfersBatchRequest(batchID *models.Uint256, commitments []models.CommitmentWithTxs) (*types.Transaction, error) {
	return t.sendTxRequest(&submitTransfersBatchRequestParams{
		batchID:     batchID,
		commitments: commitments,
	})
}

func (t *TxsSender) SubmitCreate2TransfersBatchRequest(batchID *models.Uint256, commitments []models.CommitmentWithTxs) (*types.Transaction, error) {
	return t.sendTxRequest(&submitCreate2TransfersBatchRequestParams{
		batchID:     batchID,
		commitments: commitments,
	})
}

func (t *TxsSender) SubmitMassMigrationsBatchRequest(batchID *models.Uint256, commitments []models.CommitmentWithTxs) (*types.Transaction, error) {
	return t.sendTxRequest(&submitMassMigrationsBatchRequestParams{
		batchID:     batchID,
		commitments: commitments,
	})
}

func (t *TxsSender) WithdrawStakeRequest(batchID *models.Uint256) (*types.Transaction, error) {
	return t.sendTxRequest(&withdrawStakeRequestParams{
		batchID: batchID,
	})
}

func (t *TxsSender) sendTransaction(request *txRequest) (tx *types.Transaction, err error) {
	switch params := request.params.(type) {
	case *submitTransfersBatchRequestParams:
		tx, err = t.client.SubmitTransfersBatch(params.batchID, params.commitments)
	case *submitCreate2TransfersBatchRequestParams:
		tx, err = t.client.SubmitCreate2TransfersBatch(params.batchID, params.commitments)
	case *submitMassMigrationsBatchRequestParams:
		tx, err = t.client.SubmitMassMigrationsBatch(params.batchID, params.commitments)
	case *withdrawStakeRequestParams:
		tx, err = t.client.WithdrawStake(params.batchID)
	}
	request.responseChannel <- txRequestResponse{Tx: tx, Err: err}
	return tx, err
}

func (t *TxsSender) sendTxRequest(params interface{}) (*types.Transaction, error) {
	respChan := make(chan txRequestResponse)
	t.requestsChan <- &txRequest{
		params:          params,
		responseChannel: respChan,
	}
	resp := <-respChan
	return resp.Tx, resp.Err
}
