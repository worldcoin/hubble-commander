package api

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
)

const txHashField = "txHash"

func logReceivedTransaction(hash common.Hash, tx interface{}) {
	if log.IsLevelEnabled(log.DebugLevel) {
		jsonTx, err := json.Marshal(tx)
		if err != nil {
			log.WithField(txHashField, hash).Errorln("Marshaling received transaction failed")
			return
		}
		log.WithField(txHashField, hash).Debugf("API: received new transaction: %s", string(jsonTx))
	}
}
