package api

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
)

func logReceivedTransaction(tx interface{}) {
	if log.IsLevelEnabled(log.DebugLevel) {
		jsonTx, err := json.Marshal(tx)
		if err != nil {
			log.Errorln("Marshaling received transaction failed")
			return
		}
		log.Debugf("API: received new transaction: %s", string(jsonTx))
	}
}

func logDuplicateTransaction(txHash *common.Hash) {
	log.WithField("txHash", txHash).Debug("API: received transaction already exists")
}
