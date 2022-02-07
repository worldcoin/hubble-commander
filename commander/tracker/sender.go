package tracker

import (
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

var ErrIncorrectTxParamsType = errors.New("incorrect tx params type")

type TxsSender struct {
	client       *eth.Client
	requestsChan chan *txRequest
}

func newTxRequestsSender(ethClient *eth.Client) *TxsSender {
	return &TxsSender{
		client:       ethClient,
		requestsChan: make(chan *txRequest, 32),
	}
}

func (t *TxsSender) SubmitTransfersBatch(
	batchID *models.Uint256,
	commitments []models.CommitmentWithTxs,
) (*types.Transaction, error) {
	return t.sendTxRequest(&submitTransfersBatchRequestParams{
		batchID:     batchID,
		commitments: commitments,
	})
}

func (t *TxsSender) SubmitCreate2TransfersBatch(
	batchID *models.Uint256,
	commitments []models.CommitmentWithTxs,
) (*types.Transaction, error) {
	return t.sendTxRequest(&submitCreate2TransfersBatchRequestParams{
		batchID:     batchID,
		commitments: commitments,
	})
}

func (t *TxsSender) SubmitMassMigrationsBatch(
	batchID *models.Uint256,
	commitments []models.CommitmentWithTxs,
) (*types.Transaction, error) {
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

func (t *TxsSender) sendRequest(request *txRequest) (*types.Transaction, error) {
	tx, err := t.sendTransaction(request.params)
	request.responseChannel <- txRequestResponse{Tx: tx, Err: err}
	return tx, err
}

func (t *TxsSender) sendTransaction(params interface{}) (*types.Transaction, error) {
	switch typedParams := params.(type) {
	case *submitTransfersBatchRequestParams:
		return t.client.SubmitTransfersBatch(typedParams.batchID, typedParams.commitments)
	case *submitCreate2TransfersBatchRequestParams:
		return t.client.SubmitCreate2TransfersBatch(typedParams.batchID, typedParams.commitments)
	case *submitMassMigrationsBatchRequestParams:
		return t.client.SubmitMassMigrationsBatch(typedParams.batchID, typedParams.commitments)
	case *withdrawStakeRequestParams:
		return t.client.WithdrawStake(typedParams.batchID)
	}
	return nil, ErrIncorrectTxParamsType
}

func (t *TxsSender) sendTxRequest(params interface{}) (*types.Transaction, error) {
	respChan := make(chan txRequestResponse, 1)
	t.requestsChan <- &txRequest{
		params:          params,
		responseChannel: respChan,
	}
	resp := <-respChan
	return resp.Tx, resp.Err
}
