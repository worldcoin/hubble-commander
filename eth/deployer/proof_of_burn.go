package deployer

import (
	"github.com/Worldcoin/hubble-commander/contracts/proofofburn"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func DeployProofOfBurn(c ChainConnection) (*common.Address, *proofofburn.ProofOfBurn, error) {
	log.Println("Deploying ProofOfBurn")
	proofOfBurnAddress, tx, proofOfBurn, err := proofofburn.DeployProofOfBurn(c.GetAccount(), c.GetBackend())
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	c.Commit()
	_, err = WaitToBeMined(c.GetBackend(), tx)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	return &proofOfBurnAddress, proofOfBurn, nil
}
