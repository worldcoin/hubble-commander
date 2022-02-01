package deployer

import (
	"time"

	"github.com/Worldcoin/hubble-commander/contracts/proofofburn"
	"github.com/Worldcoin/hubble-commander/eth/chain"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func DeployProofOfBurn(c chain.Connection, mineTimeout time.Duration) (*common.Address, *proofofburn.ProofOfBurn, error) {
	log.Println("Deploying ProofOfBurn")
	proofOfBurnAddress, tx, proofOfBurn, err := proofofburn.DeployProofOfBurn(c.GetAccount(), c.GetBackend())
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	_, err = chain.WaitToBeMined(c.GetBackend(), mineTimeout, tx)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	return &proofOfBurnAddress, proofOfBurn, nil
}
