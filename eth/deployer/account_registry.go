package deployer

import (
	"time"
	"math/big"

	"github.com/Worldcoin/hubble-commander/contracts/accountregistry"
	"github.com/Worldcoin/hubble-commander/eth/chain"
	"github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func DeployAccountRegistry(c chain.Connection, chooser *common.Address, mineTimeout time.Duration, root *common.Hash, initialAccountCount uint32, subtrees [storage.AccountTreeDepth-1]common.Hash) (
	*common.Address, *uint64, *accountregistry.AccountRegistry, error,
) {
	rootBytes := utils.HashToByteArray(root)

	var accountCountBig big.Int
	accountCountBig.SetUint64(uint64(initialAccountCount))

	var subtreesBytes [storage.AccountTreeDepth-1][32]byte
	for i, hash := range subtrees {
		subtreesBytes[i] = utils.HashToByteArray(&hash)
	}

	log.Println("Deploying AccountRegistry")
	accountRegistryAddress, tx, accountRegistry, err := accountregistry.DeployAccountRegistry(c.GetAccount(), c.GetBackend(), *chooser, rootBytes, &accountCountBig, subtreesBytes)
	if err != nil {
		return nil, nil, nil, errors.WithStack(err)
	}

	txReceipt, err := chain.WaitToBeMined(c.GetBackend(), mineTimeout, tx)
	if err != nil {
		return nil, nil, nil, errors.WithStack(err)
	}

	deploymentBlockNumber := txReceipt.BlockNumber.Uint64()

	return &accountRegistryAddress, &deploymentBlockNumber, accountRegistry, nil
}
