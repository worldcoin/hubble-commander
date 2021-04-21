// +build e2e

package e2e

import (
	"fmt"
	"testing"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/ybbus/jsonrpc/v2"
)

func TestCommanderLocally(t *testing.T) {
	cfg := config.GetConfig()
	endpoint := fmt.Sprintf("http://localhost:%s", cfg.API.Port)
	client := jsonrpc.NewClient(endpoint)
	runE2ETest(t, client)
}
