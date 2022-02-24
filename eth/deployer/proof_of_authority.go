package deployer

import (
	"time"

	"github.com/Worldcoin/hubble-commander/contracts/proofofauthority"
	"github.com/Worldcoin/hubble-commander/eth/chain"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func DeployProofOfAuthority(c chain.Connection, mineTimeout time.Duration) (*common.Address, *proofofauthority.ProofOfAuthority, error) {
	log.Println("Deploying ProofOfAuthority")
	poaAddress, tx, poa, err := proofofauthority.DeployProofOfAuthority(c.GetAccount(), c.GetBackend(), []common.Address{c.GetAccount().From})
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	_, err = chain.WaitToBeMined(c.GetBackend(), mineTimeout, tx)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	return &poaAddress, poa, nil
}
